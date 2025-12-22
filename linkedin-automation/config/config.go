package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	LinkedIn LinkedInConfig  `yaml:"linkedin"`
	Database DatabaseConfig  `yaml:"database"`
	Logging  LoggingConfig   `yaml:"logging"`
	Browser  BrowserConfig   `yaml:"browser"`
	Limits   LimitsConfig    `yaml:"limits"`
	Search   SearchConfig    `yaml:"search"`
	Messaging MessagingConfig `yaml:"messaging"`
	Stealth  StealthConfig   `yaml:"stealth"`
}

// LinkedInConfig contains LinkedIn authentication settings
type LinkedInConfig struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	BaseURL  string `yaml:"base_url"`
}

// DatabaseConfig contains database settings
type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level  string `yaml:"level"`
	File   string `yaml:"file"`
	Format string `yaml:"format"`
}

// BrowserConfig contains browser settings
type BrowserConfig struct {
	Headless       bool   `yaml:"headless"`
	Timeout        int    `yaml:"timeout"`
	SlowMotion     int    `yaml:"slow_motion"`
	UserDataDir    string `yaml:"user_data_dir"`
	DownloadPath   string `yaml:"download_path"`
}

// LimitsConfig contains rate limiting settings
type LimitsConfig struct {
	MaxConnectionsPerDay int `yaml:"max_connections_per_day"`
	MaxMessagesPerDay    int `yaml:"max_messages_per_day"`
	MinDelaySeconds      int `yaml:"min_delay_seconds"`
	MaxDelaySeconds      int `yaml:"max_delay_seconds"`
	BusinessHoursOnly    bool `yaml:"business_hours_only"`
	BusinessHoursStart   int `yaml:"business_hours_start"`
	BusinessHoursEnd     int `yaml:"business_hours_end"`
}

// SearchConfig contains search settings
type SearchConfig struct {
	Keywords       []string `yaml:"keywords"`
	MaxSearchPages int      `yaml:"max_search_pages"`
	ResultsPerPage int      `yaml:"results_per_page"`
}

// MessagingConfig contains messaging settings
type MessagingConfig struct {
	MessageTemplate string            `yaml:"message_template"`
	Templates       map[string]string `yaml:"templates"`
	MaxRetries      int               `yaml:"max_retries"`
}

// StealthConfig contains anti-detection settings
type StealthConfig struct {
	Enabled             bool     `yaml:"enabled"`
	RandomizeViewport   bool     `yaml:"randomize_viewport"`
	CustomUserAgent     string   `yaml:"custom_user_agent"`
	MouseMovement       bool     `yaml:"mouse_movement"`
	RealisticTyping     bool     `yaml:"realistic_typing"`
	RandomScrolling     bool     `yaml:"random_scrolling"`
	ViewportSizes       []string `yaml:"viewport_sizes"`
}

// LoadConfig loads configuration from YAML file and overrides with environment variables
func LoadConfig(path string) (*Config, error) {
	// Read YAML file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Override with environment variables
	applyEnvOverrides(&config)

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// applyEnvOverrides overrides configuration with environment variables
func applyEnvOverrides(cfg *Config) {
	if email := os.Getenv("LINKEDIN_EMAIL"); email != "" {
		cfg.LinkedIn.Email = email
	}
	if password := os.Getenv("LINKEDIN_PASSWORD"); password != "" {
		cfg.LinkedIn.Password = password
	}
	if baseURL := os.Getenv("LINKEDIN_BASE_URL"); baseURL != "" {
		cfg.LinkedIn.BaseURL = baseURL
	}
	if dbPath := os.Getenv("DATABASE_PATH"); dbPath != "" {
		cfg.Database.Path = dbPath
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.Logging.Level = logLevel
	}
	if logFile := os.Getenv("LOG_FILE"); logFile != "" {
		cfg.Logging.File = logFile
	}
	if headless := os.Getenv("HEADLESS"); headless != "" {
		cfg.Browser.Headless = headless == "true"
	}
	if timeout := os.Getenv("BROWSER_TIMEOUT"); timeout != "" {
		if val, err := strconv.Atoi(timeout); err == nil {
			cfg.Browser.Timeout = val
		}
	}
	if maxConn := os.Getenv("MAX_CONNECTIONS_PER_DAY"); maxConn != "" {
		if val, err := strconv.Atoi(maxConn); err == nil {
			cfg.Limits.MaxConnectionsPerDay = val
		}
	}
	if maxMsg := os.Getenv("MAX_MESSAGES_PER_DAY"); maxMsg != "" {
		if val, err := strconv.Atoi(maxMsg); err == nil {
			cfg.Limits.MaxMessagesPerDay = val
		}
	}
	if keywords := os.Getenv("SEARCH_KEYWORDS"); keywords != "" {
		cfg.Search.Keywords = strings.Split(keywords, ",")
	}
	if template := os.Getenv("MESSAGE_TEMPLATE"); template != "" {
		cfg.Messaging.MessageTemplate = template
	}
	if enableStealth := os.Getenv("ENABLE_STEALTH"); enableStealth != "" {
		cfg.Stealth.Enabled = enableStealth == "true"
	}
	if userAgent := os.Getenv("CUSTOM_USER_AGENT"); userAgent != "" {
		cfg.Stealth.CustomUserAgent = userAgent
	}
}

// validateConfig validates configuration values
func validateConfig(cfg *Config) error {
	if cfg.LinkedIn.BaseURL == "" {
		return fmt.Errorf("linkedin.base_url is required")
	}
	if cfg.Database.Path == "" {
		return fmt.Errorf("database.path is required")
	}
	if cfg.Limits.MaxConnectionsPerDay <= 0 {
		return fmt.Errorf("limits.max_connections_per_day must be positive")
	}
	if cfg.Limits.MinDelaySeconds < 0 || cfg.Limits.MaxDelaySeconds < cfg.Limits.MinDelaySeconds {
		return fmt.Errorf("invalid delay configuration")
	}
	if cfg.Search.MaxSearchPages <= 0 {
		return fmt.Errorf("search.max_search_pages must be positive")
	}
	return nil
}
