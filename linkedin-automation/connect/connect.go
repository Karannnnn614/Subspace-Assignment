package connect

import (
	"fmt"
	"linkedin-automation/config"
	"linkedin-automation/logger"
	"linkedin-automation/models"
	"linkedin-automation/search"
	"linkedin-automation/stealth"
	"linkedin-automation/storage"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

// Connector handles sending connection requests
type Connector struct {
	config      *config.Config
	log         *logger.Logger
	db          *storage.DB
	limitsChecker *LimitsChecker
}

// NewConnector creates a new connector
func NewConnector(cfg *config.Config, log *logger.Logger, db *storage.DB) *Connector {
	return &Connector{
		config:        cfg,
		log:           log,
		db:            db,
		limitsChecker: NewLimitsChecker(cfg, log, db),
	}
}

// SendConnectionRequests sends connection requests to profiles
func (c *Connector) SendConnectionRequests(page *rod.Page) error {
	// Check if we've hit daily limit
	if !c.limitsChecker.CanSendConnection() {
		c.log.Warn("Daily connection limit reached. Skipping...")
		return nil
	}

	// Get profiles that haven't been contacted yet
	profiles, err := c.db.GetUncontactedProfiles(c.config.Limits.MaxConnectionsPerDay)
	if err != nil {
		return fmt.Errorf("failed to get profiles: %w", err)
	}

	if len(profiles) == 0 {
		c.log.Info("No profiles to connect with")
		return nil
	}

	c.log.Info(fmt.Sprintf("Found %d profiles to connect with", len(profiles)))

	successCount := 0
	for i, profile := range profiles {
		// Check limits before each connection
		if !c.limitsChecker.CanSendConnection() {
			c.log.Warn("Daily limit reached during processing")
			break
		}

		c.log.Info(fmt.Sprintf("[%d/%d] Connecting with: %s", i+1, len(profiles), profile.Name))

		// Send connection request
		if err := c.sendConnectionRequest(page, profile); err != nil {
			c.log.Failure("Failed to send connection request", err, map[string]interface{}{
				"profile": profile.Name,
				"url":     profile.ProfileURL,
			})
			continue
		}

		successCount++
		c.log.Success("Connection request sent", map[string]interface{}{
			"profile": profile.Name,
		})

		// Record the connection attempt
		c.limitsChecker.RecordConnection()

		// Random delay between connections
		if i < len(profiles)-1 {
			delay := stealth.RandomInterval(c.config.Limits.MinDelaySeconds, c.config.Limits.MaxDelaySeconds)
			c.log.Debug(fmt.Sprintf("Waiting %v before next connection...", delay))
			time.Sleep(delay)
		}
	}

	c.log.Success(fmt.Sprintf("Sent %d connection requests", successCount), nil)
	return nil
}

// sendConnectionRequest sends a single connection request
func (c *Connector) sendConnectionRequest(page *rod.Page, profile *models.Profile) error {
	// Navigate to profile
	if err := page.Navigate(profile.ProfileURL); err != nil {
		return fmt.Errorf("failed to navigate to profile: %w", err)
	}

	// Wait for page load
	stealth.RandomDelay(2000, 4000)
	page.MustWaitLoad()

	// Simulate reading profile
	stealth.ScrollWithPauses(page, c.log)
	stealth.ThinkDelay()

	// Find Connect button
	connectBtn, err := c.findConnectButton(page)
	if err != nil {
		return fmt.Errorf("connect button not found: %w", err)
	}

	// Click Connect button with human-like behavior
	if err := stealth.ClickElement(page, connectBtn, c.log); err != nil {
		return fmt.Errorf("failed to click connect button: %w", err)
	}

	stealth.RandomDelay(1000, 2000)

	// Check if "Add a note" option appears
	addNoteBtn, err := page.Timeout(3 * time.Second).Element("button[aria-label='Add a note']")
	if err == nil {
		// Click "Add a note"
		if err := stealth.ClickElement(page, addNoteBtn, c.log); err == nil {
			stealth.RandomDelay(500, 1000)

			// Find note textarea
			noteTextarea, err := page.Timeout(3 * time.Second).Element("textarea[name='message']")
			if err == nil {
				// Generate personalized message
				message := c.generateMessage(profile)

				// Type message with human-like behavior
				if err := stealth.HumanType(noteTextarea, message, c.log); err != nil {
					c.log.Warn("Failed to type message, sending without note")
				}

				stealth.RandomDelay(1000, 2000)
			}
		}
	}

	// Find and click Send button
	sendBtn, err := c.findSendButton(page)
	if err != nil {
		return fmt.Errorf("send button not found: %w", err)
	}

	// Click Send
	if err := stealth.ClickElement(page, sendBtn, c.log); err != nil {
		return fmt.Errorf("failed to click send button: %w", err)
	}

	stealth.RandomDelay(1000, 2000)

	// Update database
	if err := c.db.MarkProfileAsContacted(profile.ProfileURL, "connection_sent"); err != nil {
		c.log.Warn(fmt.Sprintf("Failed to update database: %v", err))
	}

	return nil
}

// findConnectButton finds the Connect button on profile page
func (c *Connector) findConnectButton(page *rod.Page) (*rod.Element, error) {
	// Try multiple selectors
	selectors := []string{
		"button[aria-label*='Connect']",
		"button.pvs-profile-actions__action:has-text('Connect')",
		"button:has-text('Connect')",
	}

	for _, selector := range selectors {
		btn, err := page.Timeout(5 * time.Second).Element(selector)
		if err == nil {
			return btn, nil
		}
	}

	return nil, fmt.Errorf("connect button not found")
}

// findSendButton finds the Send button in connection modal
func (c *Connector) findSendButton(page *rod.Page) (*rod.Element, error) {
	selectors := []string{
		"button[aria-label='Send now']",
		"button[aria-label='Send invitation']",
		"button:has-text('Send')",
	}

	for _, selector := range selectors {
		btn, err := page.Timeout(5 * time.Second).Element(selector)
		if err == nil {
			return btn, nil
		}
	}

	return nil, fmt.Errorf("send button not found")
}

// generateMessage creates a personalized connection message
func (c *Connector) generateMessage(profile *models.Profile) string {
	template := c.config.Messaging.MessageTemplate
	
	// Extract first name
	parser := search.NewParser()
	firstName := parser.ExtractFirstName(profile.Name)
	
	// Replace placeholders
	message := strings.ReplaceAll(template, "{name}", firstName)
	message = strings.ReplaceAll(message, "{company}", profile.Company)
	message = strings.ReplaceAll(message, "{headline}", profile.Headline)
	
	// Ensure message is within LinkedIn's 300 character limit
	if len(message) > 300 {
		message = message[:297] + "..."
	}
	
	return message
}
