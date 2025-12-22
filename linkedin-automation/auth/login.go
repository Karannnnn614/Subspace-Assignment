package auth

import (
	"fmt"
	"linkedin-automation/config"
	"linkedin-automation/logger"
	"linkedin-automation/stealth"
	"time"

	"github.com/go-rod/rod"
)

// LoginManager handles LinkedIn authentication
type LoginManager struct {
	config     *config.Config
	log        *logger.Logger
	sessionMgr *SessionManager
}

// NewLoginManager creates a new login manager
func NewLoginManager(cfg *config.Config, log *logger.Logger, sessionMgr *SessionManager) *LoginManager {
	return &LoginManager{
		config:     cfg,
		log:        log,
		sessionMgr: sessionMgr,
	}
}

// Login performs LinkedIn login with stealth techniques
func (lm *LoginManager) Login(page *rod.Page) error {
	// Check if session exists
	if lm.sessionMgr.HasValidSession() {
		lm.log.Info("Valid session found, attempting to restore...")
		if err := lm.sessionMgr.RestoreSession(page); err != nil {
			lm.log.Warn("Failed to restore session, performing fresh login")
		} else {
			lm.log.Success("Session restored successfully", nil)
			return nil
		}
	}

	lm.log.Info("Performing fresh login...")

	// Navigate to LinkedIn
	if err := page.Navigate(lm.config.LinkedIn.BaseURL); err != nil {
		return fmt.Errorf("failed to navigate to LinkedIn: %w", err)
	}

	// Wait for page load with random delay
	stealth.RandomDelay(2000, 4000)
	
	// Wait for page load with timeout
	if err := page.Timeout(15 * time.Second).WaitLoad(); err != nil {
		lm.log.Warn("Page load timeout, continuing anyway...")
	}

	// Check if already logged in
	if lm.isLoggedIn(page) {
		lm.log.Success("Already logged in", nil)
		lm.sessionMgr.SaveSession(page)
		return nil
	}

	// Find and fill email field
	emailField, err := page.Timeout(10 * time.Second).Element("#session_key")
	if err != nil {
		return fmt.Errorf("email field not found: %w", err)
	}

	// Human-like typing for email
	if err := stealth.HumanType(emailField, lm.config.LinkedIn.Email, lm.log); err != nil {
		return fmt.Errorf("failed to type email: %w", err)
	}

	stealth.RandomDelay(500, 1500)

	// Find and fill password field
	passwordField, err := page.Timeout(10 * time.Second).Element("#session_password")
	if err != nil {
		return fmt.Errorf("password field not found: %w", err)
	}

	// Human-like typing for password
	if err := stealth.HumanType(passwordField, lm.config.LinkedIn.Password, lm.log); err != nil {
		return fmt.Errorf("failed to type password: %w", err)
	}

	stealth.RandomDelay(1000, 2000)

	// Find and click sign-in button
	signInBtn, err := page.Timeout(10 * time.Second).Element("button[type='submit']")
	if err != nil {
		return fmt.Errorf("sign-in button not found: %w", err)
	}

	// Move mouse to button with human-like movement
	if err := stealth.HumanMouseMove(page, signInBtn, lm.log); err != nil {
		lm.log.Warn("Failed to move mouse humanly, using direct click")
	}

	stealth.RandomDelay(300, 800)

	// Click sign-in button
	signInBtn.MustClick()

	// Wait for navigation
	stealth.RandomDelay(3000, 5000)
	// Gracefully wait for page load
	if err := page.Timeout(10 * time.Second).WaitLoad(); err != nil {
		lm.log.Warn("Page load timeout after login, continuing anyway...")
	}

	// Check for security checkpoint
	if lm.hasSecurityCheckpoint(page) {
		return fmt.Errorf("security checkpoint detected - manual intervention required")
	}

	// Verify login success
	if !lm.isLoggedIn(page) {
		return fmt.Errorf("login failed - credentials may be incorrect")
	}

	// Save session
	if err := lm.sessionMgr.SaveSession(page); err != nil {
		lm.log.Warn("Failed to save session")
	}

	lm.log.Success("Login successful", map[string]interface{}{
		"email": lm.config.LinkedIn.Email,
	})

	return nil
}

// isLoggedIn checks if user is currently logged in
func (lm *LoginManager) isLoggedIn(page *rod.Page) bool {
	// Check for feed page or nav bar
	_, err := page.Timeout(5 * time.Second).Element("nav.global-nav")
	if err == nil {
		return true
	}

	// Alternative check - profile icon
	_, err = page.Timeout(5 * time.Second).Element("img.global-nav__me-photo")
	return err == nil
}

// hasSecurityCheckpoint detects if LinkedIn shows security verification
func (lm *LoginManager) hasSecurityCheckpoint(page *rod.Page) bool {
	// Check for captcha
	_, err := page.Timeout(3 * time.Second).Element("#captcha")
	if err == nil {
		lm.log.Warn("⚠️ CAPTCHA detected!")
		return true
	}

	// Check for 2FA
	_, err = page.Timeout(3 * time.Second).Element("input[name='pin']")
	if err == nil {
		lm.log.Warn("⚠️ Two-Factor Authentication detected!")
		return true
	}

	// Check for email verification
	_, err = page.Timeout(3 * time.Second).Element("input[name='verificationType']")
	if err == nil {
		lm.log.Warn("⚠️ Email verification required!")
		return true
	}

	return false
}

// Logout performs logout
func (lm *LoginManager) Logout(page *rod.Page) error {
	lm.log.Info("Logging out...")

	// Navigate to feed if not there
	page.MustNavigate(lm.config.LinkedIn.BaseURL + "/feed")
	stealth.RandomDelay(2000, 3000)

	// Click profile menu
	profileBtn, err := page.Timeout(10 * time.Second).Element("img.global-nav__me-photo")
	if err != nil {
		return fmt.Errorf("profile button not found: %w", err)
	}

	profileBtn.MustClick()
	stealth.RandomDelay(1000, 2000)

	// Click sign out
	signOutBtn, err := page.Timeout(5 * time.Second).Element("a[href*='logout']")
	if err != nil {
		return fmt.Errorf("sign out button not found: %w", err)
	}

	signOutBtn.MustClick()
	stealth.RandomDelay(2000, 3000)

	lm.log.Success("Logged out successfully", nil)
	return nil
}
