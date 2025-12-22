# MockedIn - Test Site for LinkedIn Automation

This is a local test website that mimics LinkedIn's structure for safe automation testing.

## Features

- ✅ Login page with same element IDs as LinkedIn (`#session_key`, `#session_password`)
- ✅ Feed page with navigation bar
- ✅ Automation detection display
- ✅ No real authentication required
- ✅ Safe for testing bot behavior

## How to Use

### Option 1: Simple Python HTTP Server

```bash
cd test-site
python -m http.server 8080
```

Then visit: `http://localhost:8080`

### Option 2: Node.js HTTP Server

```bash
cd test-site
npx http-server -p 8080
```

Then visit: `http://localhost:8080`

### Option 3: VS Code Live Server Extension

1. Install "Live Server" extension in VS Code
2. Right-click on `index.html`
3. Select "Open with Live Server"

## Testing Your Bot

Update your `.env` file to point to the test site:

```env
# For local testing
LINKEDIN_BASE_URL=http://localhost:8080
LINKEDIN_EMAIL=test@example.com
LINKEDIN_PASSWORD=testpassword123
```

The test site will accept **any** credentials, so you can use dummy values.

## What to Test

1. **Login Flow** - Bot should fill email/password and click submit
2. **Navigation Detection** - Check if bot detects successful login
3. **Stealth Features** - Feed page shows if `navigator.webdriver` is hidden
4. **Element Interaction** - Test clicking, typing, scrolling
5. **Session Persistence** - Test cookie saving/loading

## Pages

- `index.html` - Login page
- `feed.html` - Feed/dashboard page (shown after login)

## Automation Detection

The feed page displays:
- Whether automation is detected (`navigator.webdriver`)
- Current user agent
- Green = Stealth working, Red = Automation detected

---

**Safe Testing Environment - No Real LinkedIn Interaction!**
