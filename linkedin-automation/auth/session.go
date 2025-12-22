package auth

import (
	"encoding/json"
	"fmt"
	"linkedin-automation/config"
	"linkedin-automation/logger"
	"linkedin-automation/storage"
	"os"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// SessionManager handles session persistence
type SessionManager struct {
	config      *config.Config
	log         *logger.Logger
	db          *storage.DB
	sessionFile string
}

// Session represents stored session data
type Session struct {
	Cookies   []*proto.NetworkCookie `json:"cookies"`
	UserAgent string                 `json:"user_agent"`
	Timestamp time.Time              `json:"timestamp"`
}

// NewSessionManager creates a new session manager
func NewSessionManager(cfg *config.Config, log *logger.Logger, db *storage.DB) *SessionManager {
	sessionFile := filepath.Join(cfg.Browser.UserDataDir, "session.json")
	return &SessionManager{
		config:      cfg,
		log:         log,
		db:          db,
		sessionFile: sessionFile,
	}
}

// HasValidSession checks if a valid session exists
func (sm *SessionManager) HasValidSession() bool {
	if _, err := os.Stat(sm.sessionFile); os.IsNotExist(err) {
		return false
	}

	session, err := sm.loadSessionFromFile()
	if err != nil {
		sm.log.Debug("Failed to load session: " + err.Error())
		return false
	}

	// Check if session is older than 7 days
	if time.Since(session.Timestamp) > 7*24*time.Hour {
		sm.log.Debug("Session expired (older than 7 days)")
		return false
	}

	return len(session.Cookies) > 0
}

// SaveSession saves current browser session
func (sm *SessionManager) SaveSession(page *rod.Page) error {
	// Get all cookies
	cookies, err := page.Cookies([]string{})
	if err != nil {
		return fmt.Errorf("failed to get cookies: %w", err)
	}

	// Get user agent
	userAgent, err := page.Eval(`() => navigator.userAgent`)
	if err != nil {
		return fmt.Errorf("failed to get user agent: %w", err)
	}

	session := Session{
		Cookies:   cookies,
		UserAgent: userAgent.Value.String(),
		Timestamp: time.Now(),
	}

	// Save to file
	if err := sm.saveSessionToFile(&session); err != nil {
		return fmt.Errorf("failed to save session to file: %w", err)
	}

	sm.log.Success("Session saved", map[string]interface{}{
		"cookies": len(cookies),
		"file":    sm.sessionFile,
	})

	return nil
}

// RestoreSession restores a saved session
func (sm *SessionManager) RestoreSession(page *rod.Page) error {
	session, err := sm.loadSessionFromFile()
	if err != nil {
		return fmt.Errorf("failed to load session: %w", err)
	}

	// Set cookies
	if err := page.SetCookies(session.Cookies); err != nil {
		return fmt.Errorf("failed to set cookies: %w", err)
	}

	// Navigate to LinkedIn to activate session
	if err := page.Navigate(sm.config.LinkedIn.BaseURL + "/feed"); err != nil {
		return fmt.Errorf("failed to navigate: %w", err)
	}

	page.MustWaitLoad()

	sm.log.Success("Session restored", map[string]interface{}{
		"cookies": len(session.Cookies),
	})

	return nil
}

// ClearSession removes saved session
func (sm *SessionManager) ClearSession() error {
	if err := os.Remove(sm.sessionFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove session file: %w", err)
	}
	sm.log.Info("Session cleared")
	return nil
}

// saveSessionToFile saves session to JSON file
func (sm *SessionManager) saveSessionToFile(session *Session) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(sm.sessionFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Write to file
	if err := os.WriteFile(sm.sessionFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// loadSessionFromFile loads session from JSON file
func (sm *SessionManager) loadSessionFromFile() (*Session, error) {
	data, err := os.ReadFile(sm.sessionFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}
