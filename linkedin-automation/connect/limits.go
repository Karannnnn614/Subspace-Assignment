package connect

import (
	"linkedin-automation/config"
	"linkedin-automation/logger"
	"linkedin-automation/storage"
	"time"
)

// LimitsChecker enforces rate limits and business hours
type LimitsChecker struct {
	config *config.Config
	log    *logger.Logger
	db     *storage.DB
}

// NewLimitsChecker creates a new limits checker
func NewLimitsChecker(cfg *config.Config, log *logger.Logger, db *storage.DB) *LimitsChecker {
	return &LimitsChecker{
		config: cfg,
		log:    log,
		db:     db,
	}
}

// CanSendConnection checks if we can send a connection request
func (lc *LimitsChecker) CanSendConnection() bool {
	// Check business hours
	if lc.config.Limits.BusinessHoursOnly && !lc.isBusinessHours() {
		lc.log.Warn("Outside business hours. Pausing automation.")
		return false
	}

	// Check daily connection limit
	todayCount, err := lc.db.GetDailyConnectionCount()
	if err != nil {
		lc.log.Warn("Failed to get daily connection count")
		return false
	}

	if todayCount >= lc.config.Limits.MaxConnectionsPerDay {
		lc.log.Warn("Daily connection limit reached")
		return false
	}

	lc.log.Debug("Rate limits OK - can send connection")
	return true
}

// CanSendMessage checks if we can send a message
func (lc *LimitsChecker) CanSendMessage() bool {
	// Check business hours
	if lc.config.Limits.BusinessHoursOnly && !lc.isBusinessHours() {
		lc.log.Warn("Outside business hours. Pausing automation.")
		return false
	}

	// Check daily message limit
	todayCount, err := lc.db.GetDailyMessageCount()
	if err != nil {
		lc.log.Warn("Failed to get daily message count")
		return false
	}

	if todayCount >= lc.config.Limits.MaxMessagesPerDay {
		lc.log.Warn("Daily message limit reached")
		return false
	}

	lc.log.Debug("Rate limits OK - can send message")
	return true
}

// isBusinessHours checks if current time is within business hours
func (lc *LimitsChecker) isBusinessHours() bool {
	now := time.Now()
	currentHour := now.Hour()
	
	// Check if weekend
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		lc.log.Debug("Weekend detected - outside business hours")
		return false
	}

	// Check hour range
	if currentHour < lc.config.Limits.BusinessHoursStart || currentHour >= lc.config.Limits.BusinessHoursEnd {
		lc.log.Debug("Outside business hours range")
		return false
	}

	return true
}

// RecordConnection records a connection attempt
func (lc *LimitsChecker) RecordConnection() {
	if err := lc.db.IncrementDailyConnectionCount(); err != nil {
		lc.log.Warn("Failed to increment connection count")
	}
}

// RecordMessage records a message sent
func (lc *LimitsChecker) RecordMessage() {
	if err := lc.db.IncrementDailyMessageCount(); err != nil {
		lc.log.Warn("Failed to increment message count")
	}
}

// GetRemainingConnections returns how many connections can still be sent today
func (lc *LimitsChecker) GetRemainingConnections() int {
	count, err := lc.db.GetDailyConnectionCount()
	if err != nil {
		return 0
	}

	remaining := lc.config.Limits.MaxConnectionsPerDay - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining
}

// GetRemainingMessages returns how many messages can still be sent today
func (lc *LimitsChecker) GetRemainingMessages() int {
	count, err := lc.db.GetDailyMessageCount()
	if err != nil {
		return 0
	}

	remaining := lc.config.Limits.MaxMessagesPerDay - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining
}

// WaitForBusinessHours blocks until business hours
func (lc *LimitsChecker) WaitForBusinessHours() {
	if !lc.config.Limits.BusinessHoursOnly {
		return
	}

	for !lc.isBusinessHours() {
		now := time.Now()
		
		// Calculate time until next business hour
		var nextBusinessTime time.Time
		
		// If weekend, wait until Monday
		if now.Weekday() == time.Saturday {
			nextBusinessTime = getNextMonday(now, lc.config.Limits.BusinessHoursStart)
		} else if now.Weekday() == time.Sunday {
			nextBusinessTime = getNextMonday(now, lc.config.Limits.BusinessHoursStart)
		} else if now.Hour() < lc.config.Limits.BusinessHoursStart {
			// Before business hours today
			nextBusinessTime = time.Date(now.Year(), now.Month(), now.Day(), 
				lc.config.Limits.BusinessHoursStart, 0, 0, 0, now.Location())
		} else {
			// After business hours today, wait until tomorrow
			tomorrow := now.AddDate(0, 0, 1)
			nextBusinessTime = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
				lc.config.Limits.BusinessHoursStart, 0, 0, 0, now.Location())
		}

		waitDuration := nextBusinessTime.Sub(now)
		lc.log.Info("Waiting for business hours to resume...")
		lc.log.Info("Will resume at: " + nextBusinessTime.Format("2006-01-02 15:04:05"))
		
		time.Sleep(waitDuration)
	}
}

// getNextMonday calculates next Monday at specified hour
func getNextMonday(from time.Time, hour int) time.Time {
	daysUntilMonday := (8 - int(from.Weekday())) % 7
	if daysUntilMonday == 0 {
		daysUntilMonday = 7
	}
	
	nextMonday := from.AddDate(0, 0, daysUntilMonday)
	return time.Date(nextMonday.Year(), nextMonday.Month(), nextMonday.Day(),
		hour, 0, 0, 0, from.Location())
}
