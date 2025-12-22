package stealth

import (
	"fmt"
	"linkedin-automation/logger"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
)

// FingerprintConfig contains browser fingerprint settings
type FingerprintConfig struct {
	UserAgent      string
	Platform       string
	Vendor         string
	Language       string
	Timezone       string
	ScreenWidth    int
	ScreenHeight   int
	ColorDepth     int
	HardwareConcurrency int
}

// ApplyFingerprint applies browser fingerprint customizations
func ApplyFingerprint(page *rod.Page, log *logger.Logger) error {
	config := generateRandomFingerprint()

	// Apply user agent
	script := fmt.Sprintf(`
		Object.defineProperty(navigator, 'userAgent', {
			get: () => '%s'
		});
	`, config.UserAgent)
	
	if err := page.Eval(script); err != nil {
		log.Warn("Failed to set user agent fingerprint")
	}

	// Apply platform
	script = fmt.Sprintf(`
		Object.defineProperty(navigator, 'platform', {
			get: () => '%s'
		});
	`, config.Platform)
	
	if err := page.Eval(script); err != nil {
		log.Warn("Failed to set platform fingerprint")
	}

	// Apply vendor
	script = fmt.Sprintf(`
		Object.defineProperty(navigator, 'vendor', {
			get: () => '%s'
		});
	`, config.Vendor)
	
	if err := page.Eval(script); err != nil {
		log.Warn("Failed to set vendor fingerprint")
	}

	// Apply hardware concurrency (CPU cores)
	script = fmt.Sprintf(`
		Object.defineProperty(navigator, 'hardwareConcurrency', {
			get: () => %d
		});
	`, config.HardwareConcurrency)
	
	if err := page.Eval(script); err != nil {
		log.Warn("Failed to set hardware concurrency")
	}

	// Apply screen properties
	script = fmt.Sprintf(`
		Object.defineProperty(screen, 'width', {
			get: () => %d
		});
		Object.defineProperty(screen, 'height', {
			get: () => %d
		});
		Object.defineProperty(screen, 'colorDepth', {
			get: () => %d
		});
	`, config.ScreenWidth, config.ScreenHeight, config.ColorDepth)
	
	if err := page.Eval(script); err != nil {
		log.Warn("Failed to set screen properties")
	}

	// Remove webdriver traces
	antiDetectionScripts := `
		// Remove webdriver property
		delete navigator.__proto__.webdriver;
		
		// Mock permissions
		const originalQuery = window.navigator.permissions.query;
		window.navigator.permissions.query = (parameters) => (
			parameters.name === 'notifications' ?
				Promise.resolve({state: Notification.permission}) :
				originalQuery(parameters)
		);

		// Mock plugins
		Object.defineProperty(navigator, 'plugins', {
			get: () => [
				{name: 'Chrome PDF Plugin'},
				{name: 'Chrome PDF Viewer'},
				{name: 'Native Client'}
			]
		});

		// Mock languages
		Object.defineProperty(navigator, 'languages', {
			get: () => ['en-US', 'en']
		});

		// Add chrome runtime
		if (!window.chrome) {
			window.chrome = {runtime: {}};
		}

		// Mock battery API
		if (navigator.getBattery) {
			navigator.getBattery = () => Promise.resolve({
				charging: true,
				chargingTime: 0,
				dischargingTime: Infinity,
				level: 1
			});
		}
	`
	
	if err := page.Eval(antiDetectionScripts); err != nil {
		log.Warn("Failed to apply anti-detection scripts")
	}

	log.Stealth("Applied browser fingerprint", map[string]interface{}{
		"platform": config.Platform,
		"cores":    config.HardwareConcurrency,
	})

	return nil
}

// generateRandomFingerprint generates a random but realistic browser fingerprint
func generateRandomFingerprint() *FingerprintConfig {
	rand.Seed(time.Now().UnixNano())

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	}

	platforms := []string{"Win32", "MacIntel", "Linux x86_64"}
	vendors := []string{"Google Inc.", "Apple Computer, Inc."}
	
	screenSizes := []struct{ width, height int }{
		{1920, 1080},
		{1366, 768},
		{1536, 864},
		{1440, 900},
		{2560, 1440},
	}

	cores := []int{4, 6, 8, 12, 16}

	selectedUA := userAgents[rand.Intn(len(userAgents))]
	selectedPlatform := platforms[rand.Intn(len(platforms))]
	selectedVendor := vendors[rand.Intn(len(vendors))]
	selectedScreen := screenSizes[rand.Intn(len(screenSizes))]
	selectedCores := cores[rand.Intn(len(cores))]

	return &FingerprintConfig{
		UserAgent:           selectedUA,
		Platform:            selectedPlatform,
		Vendor:              selectedVendor,
		Language:            "en-US",
		Timezone:            "America/New_York",
		ScreenWidth:         selectedScreen.width,
		ScreenHeight:        selectedScreen.height,
		ColorDepth:          24,
		HardwareConcurrency: selectedCores,
	}
}
