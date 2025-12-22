# 🎓 LinkedIn Automation PoC - Educational Use Only

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Rod](https://img.shields.io/badge/Rod-Browser_Automation-4B8BBE?style=flat)
![License](https://img.shields.io/badge/License-Educational-yellow?style=flat)
![Status](https://img.shields.io/badge/Status-Proof_of_Concept-orange?style=flat)

</div>

---

## ⚠️ **CRITICAL LEGAL DISCLAIMER**

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃                                                                   ┃
┃  ⛔ THIS SOFTWARE IS FOR EDUCATIONAL AND EVALUATION PURPOSES ONLY  ┃
┃                                                                     ┃
┃  ❌ Automating LinkedIn VIOLATES their Terms of Service            ┃
┃  ❌ DO NOT use this on real LinkedIn accounts                      ┃
┃  ❌ Risk of PERMANENT ACCOUNT BAN                                  ┃
┃  ❌ Potential LEGAL CONSEQUENCES                                   ┃
┃                                                                     ┃
┃  ✅ Use ONLY for learning browser automation techniques            ┃
┃  ✅ Use ONLY in sandboxed test environments                        ┃
┃  ✅ Use ONLY to understand anti-detection methods                  ┃
┃                                                                     ┃
┃  By using this software, you accept ALL responsibility and risk.   ┃
┃                                                                     ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

**⚖️ Legal Notice:**

- This project is a technical proof-of-concept demonstrating browser automation and stealth techniques
- Using this on LinkedIn or any production website violates their Terms of Service
- The author assumes NO responsibility for misuse
- Intended audience: Security researchers, automation engineers, students

---

## 📖 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [Stealth Techniques](#stealth-techniques)
- [Tech Stack](#tech-stack)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [How It Works](#how-it-works)
- [Demo](#demo)
- [Development](#development)
- [FAQ](#faq)
- [Contributing](#contributing)
- [License](#license)

---

## 🎯 Overview

This project is a **sophisticated proof-of-concept** demonstrating advanced browser automation techniques using Go and the Rod library. It showcases:

- **Production-grade Go architecture** with clean separation of concerns
- **Advanced anti-detection (stealth) engineering** to simulate human behavior
- **Realistic browser fingerprint masking**
- **Human-like interaction patterns** (mouse movement, typing, scrolling)
- **Robust state management** with SQLite persistence
- **Rate limiting and business hours simulation**

### What This Is NOT

❌ A production-ready LinkedIn automation tool  
❌ A tool for spam or unsolicited outreach  
❌ A scraping service for commercial use

### What This IS

✅ An educational demonstration of browser automation  
✅ A showcase of stealth programming techniques  
✅ A reference implementation for clean Go architecture  
✅ A learning resource for anti-bot detection methods

---

## 🏗️ Architecture

### High-Level Design

```
┌─────────────────────────────────────────────────────────────┐
│                         Main Entry                          │
│                      (cmd/main.go)                          │
└────────────────┬────────────────────────────────────────────┘
                 │
      ┌──────────┴──────────┐
      │                     │
      ▼                     ▼
┌─────────────┐      ┌─────────────┐
│   Config    │      │   Logger    │
│   System    │      │   System    │
└─────────────┘      └─────────────┘
      │                     │
      └──────────┬──────────┘
                 │
                 ▼
      ┌──────────────────┐
      │  Browser Manager │
      │  (Rod + Stealth) │
      └──────────────────┘
                 │
      ┌──────────┴──────────────────────┐
      │                                  │
      ▼                                  ▼
┌────────────┐                    ┌────────────┐
│    Auth    │                    │  Storage   │
│   Module   │                    │  (SQLite)  │
└────────────┘                    └────────────┘
      │                                  │
      └──────────┬───────────────────────┘
                 │
      ┌──────────┴──────────┐
      │                     │
      ▼                     ▼
┌─────────────┐      ┌─────────────┐
│   Search    │      │   Connect   │
│   Module    │      │   Module    │
└─────────────┘      └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │  Messaging  │
                    │   Module    │
                    └─────────────┘
```

### Module Breakdown

| Module         | Responsibility                                                 |
| -------------- | -------------------------------------------------------------- |
| **browser/**   | Browser initialization, stealth injection, fingerprint masking |
| **auth/**      | Login handling, session management, cookie persistence         |
| **stealth/**   | Human behavior simulation (mouse, typing, scrolling, timing)   |
| **search/**    | People search, profile extraction, pagination                  |
| **connect/**   | Connection requests, rate limiting, business hours             |
| **messaging/** | Message sending, template management                           |
| **storage/**   | SQLite database, state persistence, daily limits               |
| **logger/**    | Structured logging, contextual fields                          |
| **config/**    | YAML + environment variable configuration                      |

---

## ✨ Features

### Core Functionality

- ✅ **Automated Login** with credential validation and 2FA detection
- ✅ **Session Persistence** - Resume without re-login (cookie storage)
- ✅ **People Search** - Keyword-based search with pagination
- ✅ **Profile Scraping** - Extract name, headline, company, location
- ✅ **Connection Requests** - Send personalized connection notes
- ✅ **LinkedIn Messaging** - Follow-up messages with template support
- ✅ **Deduplication** - Avoid contacting same profile multiple times

### Advanced Features

- ✅ **Daily Rate Limits** - Configurable connection/message caps
- ✅ **Business Hours Enforcement** - Only operate during 9-5 workdays
- ✅ **Exponential Backoff** - Intelligent retry logic
- ✅ **Activity Scheduling** - Simulate breaks and realistic timing
- ✅ **Database Persistence** - Resume from last state after restart
- ✅ **Comprehensive Logging** - Debug, Info, Warn, Error levels with context

---

## 🥷 Stealth Techniques

This project implements **15+ anti-detection techniques** to simulate human behavior:

### 1️⃣ **Browser Fingerprint Masking**

```go
// Removes automation signatures
- Custom User-Agent injection
- navigator.webdriver = undefined
- Randomized viewport dimensions
- Realistic plugin/language arrays
- Hardware concurrency spoofing
```

### 2️⃣ **Human-like Mouse Movement**

```go
// Bézier curve-based mouse paths
- Non-linear trajectories
- Variable speed (ease-in/ease-out)
- Natural overshoot + micro-corrections
- Random idle movements
```

**Visualization:**

```
Start ●─────────────╮
                    │  (Bezier curve)
                    ╰────────● Target
     (NOT straight line!)
```

### 3️⃣ **Realistic Typing Simulation**

```go
- Variable keystroke delays (120-300ms)
- Occasional typos + backspace correction (5% chance)
- Burst typing followed by pauses
- Word-boundary delays
```

### 4️⃣ **Natural Scrolling Behavior**

```go
- Multi-chunk scrolling with pauses
- Occasional scroll-up (re-reading)
- Variable scroll speeds
- Random overshoot/correction
```

### 5️⃣ **Randomized Timing**

```go
- Think delays (1-5 seconds)
- Reading delays based on word count
- Random action intervals
- Exponential backoff on errors
```

### 6️⃣ **Rate Limiting & Quotas**

```go
- Daily connection limits (default: 20)
- Daily message limits (default: 10)
- Cooldown periods between actions
- Business hours enforcement (9 AM - 5 PM)
```

### 7️⃣ **Additional Stealth Layers**

- ✅ **Hover before click** - Cursor hovers over elements before clicking
- ✅ **Random scrolling** - Scroll up/down unpredictably
- ✅ **Reading pauses** - Simulate content consumption
- ✅ **Weekend detection** - Pause automation on Sat/Sun
- ✅ **Session rotation** - Cookie-based session reuse

---

## 🛠️ Tech Stack

| Component              | Technology                                          |
| ---------------------- | --------------------------------------------------- |
| **Language**           | Go 1.21+                                            |
| **Browser Automation** | [Rod](https://github.com/go-rod/rod)                |
| **Stealth**            | [go-rod/stealth](https://github.com/go-rod/stealth) |
| **Database**           | SQLite 3                                            |
| **Configuration**      | YAML + dotenv                                       |
| **Logging**            | Logrus (structured JSON)                            |
| **OS Support**         | Windows, macOS, Linux                               |

---

## 📦 Installation

### Prerequisites

- **Go 1.21 or higher** ([Download](https://go.dev/dl/))
- **Git** ([Download](https://git-scm.com/downloads))

### Step 1: Clone Repository

```bash
git clone https://github.com/yourusername/linkedin-automation.git
cd linkedin-automation
```

### Step 2: Install Dependencies

```bash
go mod download
```

### Step 3: Set Up Environment

```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your credentials (NEVER commit this file!)
# Use a TEST ACCOUNT ONLY - not your real LinkedIn!
```

### Step 4: Create Required Directories

```bash
mkdir -p data logs browser-data downloads
```

---

## ⚙️ Configuration

### Environment Variables (.env)

```env
# LinkedIn Credentials (USE TEST ACCOUNT ONLY!)
LINKEDIN_EMAIL=test-account@example.com
LINKEDIN_PASSWORD=test-password-123

# Database
DATABASE_PATH=./data/linkedin.db

# Logging
LOG_LEVEL=info           # debug, info, warn, error
LOG_FILE=./logs/automation.log

# Browser
HEADLESS=false           # Set to true to hide browser window
BROWSER_TIMEOUT=30       # Seconds

# Rate Limiting
MAX_CONNECTIONS_PER_DAY=20
MAX_MESSAGES_PER_DAY=10
MIN_DELAY_SECONDS=5
MAX_DELAY_SECONDS=15

# Search
SEARCH_KEYWORDS=Software Engineer,Product Manager
MAX_SEARCH_PAGES=3

# Messaging
MESSAGE_TEMPLATE=Hi {name}, I'd love to connect!

# Stealth
ENABLE_STEALTH=true
RANDOMIZE_VIEWPORT=true
```

### YAML Configuration (config/config.yaml)

All settings can be overridden via environment variables. See [config/config.yaml](config/config.yaml) for full options.

---

## 🚀 Usage

### Basic Commands

```bash
# Login only (test authentication)
go run cmd/main.go -action=login

# Search for profiles
go run cmd/main.go -action=search

# Send connection requests
go run cmd/main.go -action=connect

# Send messages to accepted connections
go run cmd/main.go -action=message

# Full automation (search + connect + message)
go run cmd/main.go -action=full
```

### Build Binary

```bash
# Compile executable
go build -o linkedin-bot cmd/main.go

# Run binary
./linkedin-bot -action=full
```

### Custom Config

```bash
go run cmd/main.go -config=/path/to/custom/config.yaml -action=full
```

---

## 📁 Project Structure

```
linkedin-automation/
│
├── cmd/
│   └── main.go                 # Entry point
│
├── config/
│   ├── config.go               # Config loader + validation
│   └── config.yaml             # YAML configuration
│
├── browser/
│   └── browser.go              # Browser initialization + stealth injection
│
├── auth/
│   ├── login.go                # Login logic
│   └── session.go              # Session persistence
│
├── stealth/
│   ├── mouse.go                # Bézier curve mouse movement
│   ├── typing.go               # Human-like typing simulation
│   ├── timing.go               # Randomized delays
│   ├── fingerprint.go          # Browser fingerprint masking
│   └── scroll.go               # Natural scrolling behavior
│
├── search/
│   ├── people_search.go        # People search + pagination
│   └── parser.go               # HTML parsing utilities
│
├── connect/
│   ├── connect.go              # Connection request handler
│   └── limits.go               # Rate limiting + business hours
│
├── messaging/
│   ├── messages.go             # Message sending logic
│   └── templates.go            # Message template engine
│
├── storage/
│   └── sqlite.go               # Database operations
│
├── logger/
│   └── logger.go               # Structured logging
│
├── .env.example                # Example environment variables
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
└── README.md                   # This file
```

---

## 🔍 How It Works

### Workflow Diagram

```
┌──────────────┐
│ 1. Load Config│
└───────┬──────┘
        │
        ▼
┌─────────────────────┐
│ 2. Initialize Browser│
│    (with stealth)    │
└──────────┬──────────┘
           │
           ▼
    ┌──────────────┐
    │ 3. Login     │ ──► Session exists? ──► Restore cookies
    └──────┬───────┘              │
           │                      No
           ▼                       │
    ┌──────────────┐              ▼
    │ 4. Search    │       Perform login
    │   Profiles   │       (human-like)
    └──────┬───────┘
           │
           ▼
    ┌──────────────────┐
    │ 5. Send          │ ──► Check daily limits
    │   Connection     │ ──► Enforce business hours
    │   Requests       │ ──► Random delays
    └──────┬───────────┘
           │
           ▼
    ┌──────────────────┐
    │ 6. Send Messages │ ──► Personalized templates
    │   to Accepted    │ ──► Track in database
    └──────┬───────────┘
           │
           ▼
    ┌──────────────────┐
    │ 7. Update Database│
    │   Save State      │
    └───────────────────┘
```

### Database Schema

```sql
-- Profiles discovered during search
CREATE TABLE profiles (
    id INTEGER PRIMARY KEY,
    name TEXT,
    profile_url TEXT UNIQUE,
    headline TEXT,
    company TEXT,
    location TEXT,
    status TEXT,              -- 'discovered', 'connection_sent', 'message_sent'
    contacted_at TIMESTAMP,
    message_sent_at TIMESTAMP,
    discovered_at TIMESTAMP
);

-- Rate limiting counters
CREATE TABLE daily_limits (
    id INTEGER PRIMARY KEY,
    date TEXT UNIQUE,
    connection_count INTEGER,
    message_count INTEGER
);

-- Activity log
CREATE TABLE actions (
    id INTEGER PRIMARY KEY,
    profile_url TEXT,
    action_type TEXT,         -- 'search', 'connect', 'message'
    status TEXT,              -- 'success', 'failed'
    performed_at TIMESTAMP
);
```

---

## 🎥 Demo

> **Note:** No demo video is provided as this is educational code only. Do NOT run against live LinkedIn accounts.

**To test locally (safely):**

1. Create a dummy LinkedIn account (not your real one!)
2. Use `HEADLESS=false` to watch browser automation
3. Monitor logs in real-time: `tail -f logs/automation.log`

---

## 👨‍💻 Development

### Run Tests

```bash
go test ./...
```

### Format Code

```bash
go fmt ./...
```

### Lint Code

```bash
golangci-lint run
```

### Add New Stealth Technique

1. Create function in appropriate `stealth/*.go` file
2. Document technique in README
3. Add unit tests
4. Update config if needed

---

## ❓ FAQ

### Q: Is this safe to use?

**A:** NO - not on real LinkedIn accounts. This violates LinkedIn's ToS. Use only for educational purposes in test environments.

### Q: Will I get banned?

**A:** YES - if you use this on real accounts. LinkedIn has sophisticated bot detection.

### Q: Can I use this for my business?

**A:** NO - this is illegal and unethical. Use LinkedIn's official Sales Navigator or Recruiter tools.

### Q: What is the point of this project?

**A:** To demonstrate advanced browser automation and anti-detection techniques for educational purposes.

### Q: How do I contribute?

**A:** Fork the repo, make improvements to stealth techniques or code quality, submit a PR.

### Q: Does this actually work?

**A:** The techniques are realistic, but LinkedIn's detection evolves constantly. This is a snapshot-in-time PoC.

---

## 🤝 Contributing

Contributions that improve code quality, stealth realism, or educational value are welcome!

### Guidelines

1. **Do NOT** add features that encourage production use
2. **Do** improve anti-detection techniques
3. **Do** add tests and documentation
4. **Do** follow Go best practices

### How to Contribute

```bash
# Fork the repo
git checkout -b feature/your-feature-name

# Make changes
git commit -m "feat: add XYZ stealth technique"

# Push and create PR
git push origin feature/your-feature-name
```

---

## 📜 License

**MIT License** - For educational use only.

```
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND.
THE AUTHOR IS NOT RESPONSIBLE FOR ANY MISUSE OR DAMAGES.
```

---

## 🙏 Acknowledgments

- [Rod](https://github.com/go-rod/rod) - Excellent Go browser automation library
- [go-rod/stealth](https://github.com/go-rod/stealth) - Anti-detection utilities
- Puppeteer Stealth Plugin - Inspiration for fingerprint masking

---

## 📞 Contact

**For educational inquiries only:**

- GitHub Issues: [Create an issue](https://github.com/yourusername/linkedin-automation/issues)
- Educational use questions welcome
- Commercial support requests will be ignored

---

<div align="center">

**⚠️ REMEMBER: This is educational code. DO NOT use on production systems. ⚠️**

Made with ❤️ for the automation engineering community

</div>
