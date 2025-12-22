package main

import (
	"flag"
	"fmt"
	"linkedin-automation/auth"
	"linkedin-automation/browser"
	"linkedin-automation/config"
	"linkedin-automation/connect"
	"linkedin-automation/logger"
	"linkedin-automation/messaging"
	"linkedin-automation/search"
	"linkedin-automation/storage"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Parse command line flags
	configPath := flag.String("config", "config/config.yaml", "Path to config file")
	action := flag.String("action", "full", "Action to perform: login, search, connect, message, full")
	flag.Parse()

	// Initialize logger
	log := logger.InitLogger()
	log.Info("🚀 LinkedIn Automation PoC Starting...")

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Info("✅ Configuration loaded successfully")

	// Validate credentials
	if cfg.LinkedIn.Email == "" || cfg.LinkedIn.Password == "" {
		log.Fatal("❌ LinkedIn credentials not set. Check .env file")
	}

	// Initialize database
	db, err := storage.InitDB(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	log.Info("✅ Database initialized")

	// Initialize browser
	browserManager := browser.NewBrowserManager(cfg, log)
	page, cleanup, err := browserManager.Launch()
	if err != nil {
		log.Fatalf("Failed to launch browser: %v", err)
	}
	defer cleanup()
	log.Info("✅ Browser launched with stealth mode")

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Warn("⚠️ Interrupt received, shutting down gracefully...")
		cleanup()
		os.Exit(0)
	}()

	// Initialize session manager
	sessionMgr := auth.NewSessionManager(cfg, log, db)

	// Perform login
	loginMgr := auth.NewLoginManager(cfg, log, sessionMgr)
	if err := loginMgr.Login(page); err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	log.Info("✅ Successfully logged in to LinkedIn")

	// Execute requested action
	switch *action {
	case "login":
		log.Info("✅ Login complete. Exiting...")
		return

	case "search":
		executeSearch(cfg, log, db, page)

	case "connect":
		executeSearch(cfg, log, db, page)
		executeConnect(cfg, log, db, page)

	case "message":
		executeMessaging(cfg, log, db, page)

	case "full":
		executeSearch(cfg, log, db, page)
		executeConnect(cfg, log, db, page)
		executeMessaging(cfg, log, db, page)

	default:
		log.Fatalf("Unknown action: %s", *action)
	}

	log.Info("🎉 Automation completed successfully!")
}

// executeSearch performs people search and stores results
func executeSearch(cfg *config.Config, log *logger.Logger, db *storage.DB, page *browser.Page) {
	log.Info("🔍 Starting people search...")

	searcher := search.NewPeopleSearcher(cfg, log, db)
	profiles, err := searcher.Search(page, cfg.Search.Keywords[0])
	if err != nil {
		log.Errorf("Search failed: %v", err)
		return
	}

	log.Infof("✅ Found %d profiles", len(profiles))
	for i, profile := range profiles {
		log.Debugf("[%d] %s - %s", i+1, profile.Name, profile.ProfileURL)
	}
}

// executeConnect sends connection requests
func executeConnect(cfg *config.Config, log *logger.Logger, db *storage.DB, page *browser.Page) {
	log.Info("🤝 Starting connection requests...")

	connector := connect.NewConnector(cfg, log, db)
	if err := connector.SendConnectionRequests(page); err != nil {
		log.Errorf("Connection requests failed: %v", err)
		return
	}

	log.Info("✅ Connection requests completed")
}

// executeMessaging sends follow-up messages
func executeMessaging(cfg *config.Config, log *logger.Logger, db *storage.DB, page *browser.Page) {
	log.Info("💬 Starting messaging campaign...")

	messenger := messaging.NewMessenger(cfg, log, db)
	if err := messenger.SendMessages(page); err != nil {
		log.Errorf("Messaging failed: %v", err)
		return
	}

	log.Info("✅ Messaging campaign completed")
}

// printBanner displays ASCII art banner
func printBanner() {
	banner := `
╔═══════════════════════════════════════════════════════════════╗
║                                                               ║
║        LinkedIn Automation PoC - Educational Use Only         ║
║                                                               ║
║  ⚠️  WARNING: This violates LinkedIn Terms of Service        ║
║  ⚠️  DO NOT use on real accounts                             ║
║  ⚠️  Risk of permanent ban                                   ║
║                                                               ║
╚═══════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
	time.Sleep(2 * time.Second)
}

func init() {
	printBanner()
}
