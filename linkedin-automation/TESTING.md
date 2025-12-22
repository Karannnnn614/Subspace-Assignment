## Testing Instructions

### 1. Start the Test Site Server

Open a PowerShell terminal and run:

```powershell
cd 'c:\Users\karan\OneDrive\Creative Cloud Files\Desktop\Subspace Assignment\Test Site'
python -m http.server 8080
```

Keep this terminal open. You should see: `Serving HTTP on :: port 8080`

### 2. Test Site in Browser

Open your web browser and go to: http://localhost:8080

You should see the MockedIn login page.

### 3. Run the Bot

In another PowerShell terminal:

```powershell
cd "c:\Users\karan\OneDrive\Creative Cloud Files\Desktop\Subspace Assignment\linkedin-automation"
.\linkedin-bot.exe -action login
```

### What to Expect

1. ✅ Chrome browser window will open (not headless)
2. ✅ Navigate to http://localhost:8080
3. ✅ Type email slowly with human-like behavior
4. ✅ Type password slowly
5. ✅ Click "Sign in" button
6. ✅ Redirect to feed.html page
7. ✅ Check if automation detected or stealth working

### Common Issues

**Issue**: "Failed to launch browser: Opening in existing browser session"
**Fix**: Close all Chrome windows first:

```powershell
Get-Process chrome -ErrorAction SilentlyContinue | Stop-Process -Force
```

**Issue**: "sign-in button not found: context deadline exceeded"
**Fix**: The page is loading too slowly. Check that test server is running on port 8080.

**Issue**: Opening real LinkedIn instead of localhost
**Fix**: Check .env file has `LINKEDIN_BASE_URL=http://localhost:8080`

### Success Indicators

- Chrome opens to localhost:8080 (not linkedin.com)
- You see typing happen character-by-character
- Page redirects to feed.html after login
- Feed page shows "✓ Stealth Mode Active" (green) instead of "⚠️ Automation Detected" (red)

### Logs

Check logs for details:

```powershell
Get-Content logs\automation.log -Tail 30
```
