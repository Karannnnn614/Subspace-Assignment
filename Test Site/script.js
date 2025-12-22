// Mock LinkedIn Login Behavior
document.getElementById('loginForm').addEventListener('submit', function(e) {
    e.preventDefault();
    
    const email = document.getElementById('session_key').value;
    const password = document.getElementById('session_password').value;
    const messageDiv = document.getElementById('message');
    
    // Simulate processing delay
    messageDiv.textContent = 'Signing in...';
    messageDiv.className = 'message';
    messageDiv.style.display = 'block';
    
    setTimeout(() => {
        // Accept any credentials for testing
        if (email && password) {
            messageDiv.textContent = '✓ Login successful! Redirecting to feed...';
            messageDiv.className = 'message success';
            
            // Redirect to feed page after 1 second
            setTimeout(() => {
                window.location.href = 'feed.html';
            }, 1000);
        } else {
            messageDiv.textContent = '✗ Please enter both email and password';
            messageDiv.className = 'message error';
        }
    }, 500);
});

// Log automation detection attempts (for testing stealth features)
if (navigator.webdriver) {
    console.warn('⚠️ WebDriver detected!');
}

// Detect automation
Object.defineProperty(window, '_automationDetected', {
    get: () => {
        return navigator.webdriver === true;
    }
});

console.log('Test page loaded. Automation detected:', navigator.webdriver);
