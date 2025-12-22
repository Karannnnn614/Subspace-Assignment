package messaging

import (
	"fmt"
	"linkedin-automation/config"
	"linkedin-automation/connect"
	"linkedin-automation/logger"
	"linkedin-automation/models"
	"linkedin-automation/search"
	"linkedin-automation/stealth"
	"linkedin-automation/storage"
	"time"

	"github.com/go-rod/rod"
)

// Messenger handles LinkedIn messaging
type Messenger struct {
	config        *config.Config
	log           *logger.Logger
	db            *storage.DB
	limitsChecker *connect.LimitsChecker
	templateMgr   *TemplateManager
}

// NewMessenger creates a new messenger
func NewMessenger(cfg *config.Config, log *logger.Logger, db *storage.DB) *Messenger {
	return &Messenger{
		config:        cfg,
		log:           log,
		db:            db,
		limitsChecker: connect.NewLimitsChecker(cfg, log, db),
		templateMgr:   NewTemplateManager(cfg),
	}
}

// SendMessages sends follow-up messages to accepted connections
func (m *Messenger) SendMessages(page *rod.Page) error {
	// Check if we can send messages
	if !m.limitsChecker.CanSendMessage() {
		m.log.Warn("Daily message limit reached or outside business hours")
		return nil
	}

	// Get profiles that accepted connection but haven't received message
	profiles, err := m.db.GetProfilesNeedingMessage(m.config.Limits.MaxMessagesPerDay)
	if err != nil {
		return fmt.Errorf("failed to get profiles: %w", err)
	}

	if len(profiles) == 0 {
		m.log.Info("No profiles need messaging")
		return nil
	}

	m.log.Info(fmt.Sprintf("Found %d profiles to message", len(profiles)))

	// Navigate to messaging page
	if err := page.Navigate(m.config.LinkedIn.BaseURL + "/messaging"); err != nil {
		return fmt.Errorf("failed to navigate to messaging: %w", err)
	}

	stealth.RandomDelay(2000, 4000)
	page.MustWaitLoad()

	successCount := 0
	for i, profile := range profiles {
		// Check limits
		if !m.limitsChecker.CanSendMessage() {
			m.log.Warn("Daily limit reached during processing")
			break
		}

		m.log.Info(fmt.Sprintf("[%d/%d] Messaging: %s", i+1, len(profiles), profile.Name))

		// Send message
		if err := m.sendMessage(page, profile); err != nil {
			m.log.Failure("Failed to send message", err, map[string]interface{}{
				"profile": profile.Name,
			})
			continue
		}

		successCount++
		m.log.Success("Message sent", map[string]interface{}{
			"profile": profile.Name,
		})

		// Record message
		m.limitsChecker.RecordMessage()

		// Random delay between messages
		if i < len(profiles)-1 {
			delay := stealth.RandomInterval(m.config.Limits.MinDelaySeconds*2, m.config.Limits.MaxDelaySeconds*2)
			m.log.Debug(fmt.Sprintf("Waiting %v before next message...", delay))
			time.Sleep(delay)
		}
	}

	m.log.Success(fmt.Sprintf("Sent %d messages", successCount), nil)
	return nil
}

// sendMessage sends a single message to a profile
func (m *Messenger) sendMessage(page *rod.Page, profile *models.Profile) error {
	// Search for conversation
	searchBox, err := page.Timeout(10 * time.Second).Element("input[placeholder*='Search messages']")
	if err != nil {
		return fmt.Errorf("message search box not found: %w", err)
	}

	// Type profile name to search
	parser := search.NewParser()
	searchName := parser.ExtractFirstName(profile.Name)
	
	if err := stealth.HumanType(searchBox, searchName, m.log); err != nil {
		return fmt.Errorf("failed to type in search: %w", err)
	}

	stealth.RandomDelay(2000, 3000)

	// Find and click on conversation
	conversation, err := page.Timeout(5 * time.Second).Element("li.msg-conversation-listitem")
	if err != nil {
		return fmt.Errorf("conversation not found: %w", err)
	}

	if err := stealth.ClickElement(page, conversation, m.log); err != nil {
		return fmt.Errorf("failed to click conversation: %w", err)
	}

	stealth.RandomDelay(1500, 3000)

	// Find message input box
	messageBox, err := page.Timeout(10 * time.Second).Element("div.msg-form__contenteditable")
	if err != nil {
		return fmt.Errorf("message box not found: %w", err)
	}

	// Generate message from template
	message := m.templateMgr.GenerateMessage("default", profile)

	// Type message
	if err := stealth.TypeWithPauses(messageBox, message, m.log); err != nil {
		return fmt.Errorf("failed to type message: %w", err)
	}

	stealth.RandomDelay(1000, 2000)

	// Find and click Send button
	sendBtn, err := page.Timeout(5 * time.Second).Element("button.msg-form__send-button")
	if err != nil {
		return fmt.Errorf("send button not found: %w", err)
	}

	if err := stealth.ClickElement(page, sendBtn, m.log); err != nil {
		return fmt.Errorf("failed to click send: %w", err)
	}

	stealth.RandomDelay(1000, 2000)

	// Update database
	if err := m.db.MarkMessageSent(profile.ProfileURL); err != nil {
		m.log.Warn("Failed to update database")
	}

	return nil
}

// SendDirectMessage sends a message directly from profile page
func (m *Messenger) SendDirectMessage(page *rod.Page, profile *models.Profile, message string) error {
	// Navigate to profile
	if err := page.Navigate(profile.ProfileURL); err != nil {
		return fmt.Errorf("failed to navigate: %w", err)
	}

	stealth.RandomDelay(2000, 4000)
	page.MustWaitLoad()

	// Find Message button
	messageBtn, err := m.findMessageButton(page)
	if err != nil {
		return fmt.Errorf("message button not found: %w", err)
	}

	// Click Message button
	if err := stealth.ClickElement(page, messageBtn, m.log); err != nil {
		return fmt.Errorf("failed to click message button: %w", err)
	}

	stealth.RandomDelay(1500, 2500)

	// Find message modal input
	messageInput, err := page.Timeout(5 * time.Second).Element("div.msg-form__contenteditable")
	if err != nil {
		return fmt.Errorf("message input not found: %w", err)
	}

	// Type message
	if err := stealth.TypeWithPauses(messageInput, message, m.log); err != nil {
		return fmt.Errorf("failed to type message: %w", err)
	}

	stealth.RandomDelay(1000, 2000)

	// Send
	sendBtn, err := page.Timeout(5 * time.Second).Element("button[type='submit']")
	if err != nil {
		return fmt.Errorf("send button not found: %w", err)
	}

	if err := stealth.ClickElement(page, sendBtn, m.log); err != nil {
		return fmt.Errorf("failed to send: %w", err)
	}

	return nil
}

// findMessageButton finds the Message button on profile page
func (m *Messenger) findMessageButton(page *rod.Page) (*rod.Element, error) {
	selectors := []string{
		"button[aria-label*='Message']",
		"button:has-text('Message')",
		"a[href*='/messaging/thread/']",
	}

	for _, selector := range selectors {
		btn, err := page.Timeout(5 * time.Second).Element(selector)
		if err == nil {
			return btn, nil
		}
	}

	return nil, fmt.Errorf("message button not found")
}
