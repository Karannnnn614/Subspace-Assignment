package browser

import (
	"fmt"
	"linkedin-automation/config"
	"linkedin-automation/logger"
	"linkedin-automation/stealth"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	rodstealth "github.com/go-rod/stealth"
)

// Page is a type alias for rod.Page
type Page = rod.Page

// BrowserManager handles browser initialization and configuration
type BrowserManager struct {
	config *config.Config
	log    *logger.Logger
}

// NewBrowserManager creates a new browser manager
func NewBrowserManager(cfg *config.Config, log *logger.Logger) *BrowserManager {
	return &BrowserManager{
		config: cfg,
		log:    log,
	}
}

// Launch initializes and returns a new browser page with stealth mode enabled
func (bm *BrowserManager) Launch() (*rod.Page, func(), error) {
	// Initialize launcher
	l := launcher.New().
		Headless(bm.config.Browser.Headless).
		Devtools(false)

	// Set user data directory for session persistence
	if bm.config.Browser.UserDataDir != "" {
		l = l.UserDataDir(bm.config.Browser.UserDataDir)
	}

	// Launch browser
	url, err := l.Launch()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	// Connect to browser
	browser := rod.New().ControlURL(url).MustConnect()

	// Set timeout
	if bm.config.Browser.Timeout > 0 {
		browser = browser.Timeout(time.Duration(bm.config.Browser.Timeout) * time.Second)
	}

	// Create new page
	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		browser.MustClose()
		return nil, nil, fmt.Errorf("failed to create page: %w", err)
	}

	// Apply stealth techniques
	if bm.config.Stealth.Enabled {
		if err := bm.applyStealth(page); err != nil {
			page.MustClose()
			browser.MustClose()
			return nil, nil, fmt.Errorf("failed to apply stealth: %w", err)
		}
		bm.log.Stealth("Applied stealth techniques", map[string]interface{}{
			"user_agent": bm.config.Stealth.CustomUserAgent,
			"viewport":   bm.getRandomViewport(),
		})
	}

	// Cleanup function
	cleanup := func() {
		page.MustClose()
		browser.MustClose()
	}

	return page, cleanup, nil
}

// applyStealth applies anti-detection techniques to the browser
func (bm *BrowserManager) applyStealth(page *rod.Page) error {
	// Apply rod/stealth library
	if err := rodstealth.Apply(page); err != nil {
		return fmt.Errorf("failed to apply rod stealth: %w", err)
	}

	// Set custom user agent
	if bm.config.Stealth.CustomUserAgent != "" {
		if err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: bm.config.Stealth.CustomUserAgent,
		}); err != nil {
			return fmt.Errorf("failed to set user agent: %w", err)
		}
	}

	// Set random viewport
	if bm.config.Stealth.RandomizeViewport {
		viewport := bm.getRandomViewport()
		if err := page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
			Width:  viewport.Width,
			Height: viewport.Height,
		}); err != nil {
			return fmt.Errorf("failed to set viewport: %w", err)
		}
	}

	// Inject additional anti-detection scripts
	if err := bm.injectAntiDetectionScripts(page); err != nil {
		return fmt.Errorf("failed to inject scripts: %w", err)
	}

	return nil
}

// injectAntiDetectionScripts injects JavaScript to mask automation
func (bm *BrowserManager) injectAntiDetectionScripts(page *rod.Page) error {
	scripts := []string{
		// Remove webdriver property
		`Object.defineProperty(navigator, 'webdriver', {get: () => undefined});`,
		
		// Mock plugins
		`Object.defineProperty(navigator, 'plugins', {
			get: () => [1, 2, 3, 4, 5]
		});`,
		
		// Mock languages
		`Object.defineProperty(navigator, 'languages', {
			get: () => ['en-US', 'en']
		});`,
		
		// Chrome runtime
		`window.chrome = {runtime: {}};`,
		
		// Permissions
		`const originalQuery = window.navigator.permissions.query;
		window.navigator.permissions.query = (parameters) => (
			parameters.name === 'notifications' ?
				Promise.resolve({state: Notification.permission}) :
				originalQuery(parameters)
		);`,
	}

	for _, script := range scripts {
		if err := page.Eval(script); err != nil {
			bm.log.Warn(fmt.Sprintf("Failed to inject script: %v", err))
		}
	}

	return nil
}

// getRandomViewport returns a random viewport size from config
func (bm *BrowserManager) getRandomViewport() struct{ Width, Height int } {
	if len(bm.config.Stealth.ViewportSizes) == 0 {
		return struct{ Width, Height int }{1920, 1080}
	}

	rand.Seed(time.Now().UnixNano())
	size := bm.config.Stealth.ViewportSizes[rand.Intn(len(bm.config.Stealth.ViewportSizes))]
	
	parts := strings.Split(size, "x")
	if len(parts) != 2 {
		return struct{ Width, Height int }{1920, 1080}
	}

	width, _ := strconv.Atoi(parts[0])
	height, _ := strconv.Atoi(parts[1])

	return struct{ Width, Height int }{width, height}
}

// WaitForNavigation waits for page navigation with human-like delay
func WaitForNavigation(page *rod.Page, log *logger.Logger) error {
	// Add random delay before navigation
	stealth.RandomDelay(1000, 3000)
	
	wait := page.WaitNavigation(proto.PageLifecycleEventNameNetworkIdle)
	wait()
	
	// Add random delay after navigation
	stealth.RandomDelay(500, 1500)
	
	log.Debug("Navigation completed")
	return nil
}

// ScrollToElement scrolls to an element with human-like behavior
func ScrollToElement(page *rod.Page, element *rod.Element, log *logger.Logger) error {
	if element == nil {
		return fmt.Errorf("element is nil")
	}

	// Get element position
	box, err := element.Box()
	if err != nil {
		return fmt.Errorf("failed to get element box: %w", err)
	}

	// Scroll with human-like behavior
	stealth.HumanScroll(page, int(box.Y), log)
	
	return nil
}
