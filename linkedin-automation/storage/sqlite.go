package storage

import (
	"database/sql"
	"fmt"
	"linkedin-automation/search"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps SQLite database connection
type DB struct {
	conn *sql.DB
}

// InitDB initializes the database and creates tables
func InitDB(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}

	// Create tables
	if err := db.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// createTables creates necessary database tables
func (db *DB) createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS profiles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		profile_url TEXT UNIQUE NOT NULL,
		headline TEXT,
		company TEXT,
		location TEXT,
		connection_degree TEXT,
		search_keyword TEXT,
		discovered_at TIMESTAMP,
		contacted_at TIMESTAMP,
		message_sent_at TIMESTAMP,
		status TEXT DEFAULT 'discovered',
		notes TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS actions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_url TEXT,
		action_type TEXT NOT NULL,
		status TEXT,
		error_message TEXT,
		performed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (profile_url) REFERENCES profiles(profile_url)
	);

	CREATE TABLE IF NOT EXISTS daily_limits (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT UNIQUE NOT NULL,
		connection_count INTEGER DEFAULT 0,
		message_count INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS session_data (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT UNIQUE NOT NULL,
		value TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_profile_url ON profiles(profile_url);
	CREATE INDEX IF NOT EXISTS idx_profile_status ON profiles(status);
	CREATE INDEX IF NOT EXISTS idx_action_date ON actions(performed_at);
	CREATE INDEX IF NOT EXISTS idx_daily_limits_date ON daily_limits(date);
	`

	_, err := db.conn.Exec(schema)
	return err
}

// SaveProfile saves or updates a profile
func (db *DB) SaveProfile(profile *search.Profile) error {
	query := `
	INSERT INTO profiles (name, profile_url, headline, company, location, connection_degree, search_keyword, discovered_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(profile_url) DO UPDATE SET
		name = excluded.name,
		headline = excluded.headline,
		company = excluded.company,
		location = excluded.location,
		updated_at = CURRENT_TIMESTAMP
	`

	_, err := db.conn.Exec(query,
		profile.Name,
		profile.ProfileURL,
		profile.Headline,
		profile.Company,
		profile.Location,
		profile.ConnectionDegree,
		profile.SearchKeyword,
		profile.DiscoveredAt,
	)

	return err
}

// GetUncontactedProfiles retrieves profiles that haven't been contacted yet
func (db *DB) GetUncontactedProfiles(limit int) ([]*search.Profile, error) {
	query := `
	SELECT name, profile_url, headline, company, location, connection_degree, search_keyword, discovered_at
	FROM profiles
	WHERE contacted_at IS NULL AND status = 'discovered'
	ORDER BY discovered_at DESC
	LIMIT ?
	`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return db.scanProfiles(rows)
}

// GetProfilesNeedingMessage retrieves profiles that need follow-up messages
func (db *DB) GetProfilesNeedingMessage(limit int) ([]*search.Profile, error) {
	query := `
	SELECT name, profile_url, headline, company, location, connection_degree, search_keyword, discovered_at
	FROM profiles
	WHERE contacted_at IS NOT NULL 
		AND message_sent_at IS NULL 
		AND status = 'connection_sent'
	ORDER BY contacted_at ASC
	LIMIT ?
	`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return db.scanProfiles(rows)
}

// MarkProfileAsContacted marks a profile as contacted
func (db *DB) MarkProfileAsContacted(profileURL, status string) error {
	query := `
	UPDATE profiles
	SET contacted_at = CURRENT_TIMESTAMP,
		status = ?,
		updated_at = CURRENT_TIMESTAMP
	WHERE profile_url = ?
	`

	_, err := db.conn.Exec(query, status, profileURL)
	return err
}

// MarkMessageSent marks that a message was sent to a profile
func (db *DB) MarkMessageSent(profileURL string) error {
	query := `
	UPDATE profiles
	SET message_sent_at = CURRENT_TIMESTAMP,
		status = 'message_sent',
		updated_at = CURRENT_TIMESTAMP
	WHERE profile_url = ?
	`

	_, err := db.conn.Exec(query, profileURL)
	return err
}

// RecordAction records an action performed
func (db *DB) RecordAction(profileURL, actionType, status, errorMsg string) error {
	query := `
	INSERT INTO actions (profile_url, action_type, status, error_message)
	VALUES (?, ?, ?, ?)
	`

	_, err := db.conn.Exec(query, profileURL, actionType, status, errorMsg)
	return err
}

// GetDailyConnectionCount gets today's connection count
func (db *DB) GetDailyConnectionCount() (int, error) {
	today := time.Now().Format("2006-01-02")
	
	query := `
	SELECT COALESCE(connection_count, 0)
	FROM daily_limits
	WHERE date = ?
	`

	var count int
	err := db.conn.QueryRow(query, today).Scan(&count)
	if err == sql.ErrNoRows {
		return 0, nil
	}

	return count, err
}

// GetDailyMessageCount gets today's message count
func (db *DB) GetDailyMessageCount() (int, error) {
	today := time.Now().Format("2006-01-02")
	
	query := `
	SELECT COALESCE(message_count, 0)
	FROM daily_limits
	WHERE date = ?
	`

	var count int
	err := db.conn.QueryRow(query, today).Scan(&count)
	if err == sql.ErrNoRows {
		return 0, nil
	}

	return count, err
}

// IncrementDailyConnectionCount increments today's connection count
func (db *DB) IncrementDailyConnectionCount() error {
	today := time.Now().Format("2006-01-02")
	
	query := `
	INSERT INTO daily_limits (date, connection_count, message_count)
	VALUES (?, 1, 0)
	ON CONFLICT(date) DO UPDATE SET
		connection_count = connection_count + 1
	`

	_, err := db.conn.Exec(query, today)
	return err
}

// IncrementDailyMessageCount increments today's message count
func (db *DB) IncrementDailyMessageCount() error {
	today := time.Now().Format("2006-01-02")
	
	query := `
	INSERT INTO daily_limits (date, connection_count, message_count)
	VALUES (?, 0, 1)
	ON CONFLICT(date) DO UPDATE SET
		message_count = message_count + 1
	`

	_, err := db.conn.Exec(query, today)
	return err
}

// GetStats returns database statistics
func (db *DB) GetStats() (map[string]int, error) {
	stats := make(map[string]int)

	// Total profiles
	var totalProfiles int
	db.conn.QueryRow("SELECT COUNT(*) FROM profiles").Scan(&totalProfiles)
	stats["total_profiles"] = totalProfiles

	// Contacted profiles
	var contacted int
	db.conn.QueryRow("SELECT COUNT(*) FROM profiles WHERE contacted_at IS NOT NULL").Scan(&contacted)
	stats["contacted"] = contacted

	// Messages sent
	var messagesSent int
	db.conn.QueryRow("SELECT COUNT(*) FROM profiles WHERE message_sent_at IS NOT NULL").Scan(&messagesSent)
	stats["messages_sent"] = messagesSent

	// Today's connections
	todayConnections, _ := db.GetDailyConnectionCount()
	stats["today_connections"] = todayConnections

	// Today's messages
	todayMessages, _ := db.GetDailyMessageCount()
	stats["today_messages"] = todayMessages

	return stats, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// scanProfiles scans rows into Profile structs
func (db *DB) scanProfiles(rows *sql.Rows) ([]*search.Profile, error) {
	var profiles []*search.Profile

	for rows.Next() {
		var p search.Profile
		var discoveredAt sql.NullTime

		err := rows.Scan(
			&p.Name,
			&p.ProfileURL,
			&p.Headline,
			&p.Company,
			&p.Location,
			&p.ConnectionDegree,
			&p.SearchKeyword,
			&discoveredAt,
		)

		if err != nil {
			return nil, err
		}

		if discoveredAt.Valid {
			p.DiscoveredAt = discoveredAt.Time
		}

		profiles = append(profiles, &p)
	}

	return profiles, rows.Err()
}
