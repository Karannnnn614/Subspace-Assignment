package stealth

import (
	"math/rand"
	"time"
)

// RandomDelay introduces a random delay between min and max milliseconds
func RandomDelay(minMs, maxMs int) {
	if minMs >= maxMs {
		time.Sleep(time.Duration(minMs) * time.Millisecond)
		return
	}

	rand.Seed(time.Now().UnixNano())
	delay := minMs + rand.Intn(maxMs-minMs)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// ThinkDelay simulates human thinking time
func ThinkDelay() {
	rand.Seed(time.Now().UnixNano())
	// Human think time: 1-5 seconds
	delay := 1000 + rand.Intn(4000)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// ShortPause simulates a short pause
func ShortPause() {
	RandomDelay(500, 1500)
}

// MediumPause simulates a medium pause
func MediumPause() {
	RandomDelay(2000, 4000)
}

// LongPause simulates a long pause
func LongPause() {
	RandomDelay(5000, 10000)
}

// ReadingDelay simulates time spent reading content
func ReadingDelay(wordCount int) {
	// Average reading speed: 200-250 words per minute
	// That's about 240-300ms per word
	rand.Seed(time.Now().UnixNano())
	
	msPerWord := 240 + rand.Intn(60)
	totalMs := wordCount * msPerWord
	
	// Add some variance
	variance := rand.Intn(totalMs/4) - totalMs/8
	totalMs += variance
	
	// Cap at reasonable limits
	if totalMs < 1000 {
		totalMs = 1000
	}
	if totalMs > 30000 {
		totalMs = 30000
	}
	
	time.Sleep(time.Duration(totalMs) * time.Millisecond)
}

// ExponentialBackoff implements exponential backoff for retries
func ExponentialBackoff(attempt int, maxDelay int) {
	rand.Seed(time.Now().UnixNano())
	
	// Calculate exponential delay: 2^attempt * base
	base := 1000 // 1 second base
	delay := base * (1 << uint(attempt))
	
	// Add jitter (±25%)
	jitter := rand.Intn(delay/2) - delay/4
	delay += jitter
	
	// Cap at max delay
	if delay > maxDelay {
		delay = maxDelay
	}
	
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// BusinessHoursDelay waits until business hours if outside them
func BusinessHoursDelay(startHour, endHour int) {
	now := time.Now()
	currentHour := now.Hour()
	
	// If within business hours, no delay
	if currentHour >= startHour && currentHour < endHour {
		return
	}
	
	// Calculate time until next business hour
	var nextBusinessHour time.Time
	if currentHour < startHour {
		// Wait until start of business hours today
		nextBusinessHour = time.Date(now.Year(), now.Month(), now.Day(), startHour, 0, 0, 0, now.Location())
	} else {
		// Wait until start of business hours tomorrow
		tomorrow := now.AddDate(0, 0, 1)
		nextBusinessHour = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), startHour, 0, 0, 0, now.Location())
	}
	
	waitDuration := nextBusinessHour.Sub(now)
	time.Sleep(waitDuration)
}

// RandomInterval returns a random duration within a range
func RandomInterval(minSeconds, maxSeconds int) time.Duration {
	rand.Seed(time.Now().UnixNano())
	seconds := minSeconds + rand.Intn(maxSeconds-minSeconds)
	return time.Duration(seconds) * time.Second
}
