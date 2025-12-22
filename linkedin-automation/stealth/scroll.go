package stealth

import (
	"fmt"
	"linkedin-automation/logger"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
)

// HumanScroll scrolls the page in a human-like manner
func HumanScroll(page *rod.Page, targetY int, log *logger.Logger) error {
	rand.Seed(time.Now().UnixNano())

	// Get current scroll position
	currentY := page.MustEval(`window.pageYOffset`).Int()

	// Calculate scroll distance
	distance := targetY - currentY
	if distance == 0 {
		return nil
	}

	// Determine scroll direction
	direction := 1
	if distance < 0 {
		direction = -1
		distance = -distance
	}

	// Scroll in chunks with variable speed
	chunks := rand.Intn(5) + 3 // 3-7 chunks
	chunkSize := distance / chunks

	for i := 0; i < chunks; i++ {
		scrollAmount := chunkSize * direction
		
		// Last chunk: scroll exactly to target
		if i == chunks-1 {
			scrollAmount = targetY - page.MustEval(`window.pageYOffset`).Int()
		}

		// Scroll
		script := fmt.Sprintf(`window.scrollBy(0, %d)`, scrollAmount)
		page.MustEval(script)

		// Variable delay between scrolls
		delay := rand.Intn(300) + 100
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	// Occasionally overshoot and correct
	if rand.Float64() < 0.3 {
		overshoot := (rand.Intn(100) + 50) * direction
		page.MustEval(fmt.Sprintf(`window.scrollBy(0, %d)`, overshoot))
		time.Sleep(200 * time.Millisecond)

		// Correct back
		page.MustEval(fmt.Sprintf(`window.scrollBy(0, %d)`, -overshoot))
		time.Sleep(100 * time.Millisecond)
	}

	log.Stealth("Human scroll", map[string]interface{}{
		"from":   currentY,
		"to":     targetY,
		"chunks": chunks,
	})

	return nil
}

// ScrollToBottom scrolls to the bottom of the page gradually
func ScrollToBottom(page *rod.Page, log *logger.Logger) error {
	rand.Seed(time.Now().UnixNano())

	// Get page height
	pageHeight := page.MustEval(`document.body.scrollHeight`).Int()

	// Scroll in sections
	sections := rand.Intn(5) + 4 // 4-8 sections
	sectionHeight := pageHeight / sections

	for i := 1; i <= sections; i++ {
		targetY := sectionHeight * i
		
		// Scroll to section
		HumanScroll(page, targetY, log)

		// Pause to "read" content
		if rand.Float64() < 0.6 {
			readPause := rand.Intn(2000) + 1000
			time.Sleep(time.Duration(readPause) * time.Millisecond)
		}

		// Occasionally scroll up a bit
		if rand.Float64() < 0.2 {
			scrollUp := rand.Intn(200) + 100
			page.MustEval(fmt.Sprintf(`window.scrollBy(0, -%d)`, scrollUp))
			time.Sleep(500 * time.Millisecond)
		}
	}

	log.Stealth("Scrolled to bottom", map[string]interface{}{
		"sections": sections,
	})

	return nil
}

// RandomScroll performs random scrolling behavior
func RandomScroll(page *rod.Page, log *logger.Logger) {
	rand.Seed(time.Now().UnixNano())

	// Get page height
	pageHeight := page.MustEval(`document.body.scrollHeight`).Int()
	
	// Random target position (not too far)
	currentY := page.MustEval(`window.pageYOffset`).Int()
	maxScroll := 500
	
	targetY := currentY + rand.Intn(maxScroll*2) - maxScroll
	
	// Ensure within bounds
	if targetY < 0 {
		targetY = 0
	}
	if targetY > pageHeight-800 {
		targetY = pageHeight - 800
	}

	HumanScroll(page, targetY, log)

	// Random pause
	RandomDelay(1000, 3000)
}

// ScrollToElement scrolls an element into view naturally
func ScrollToElement(page *rod.Page, element *rod.Element, log *logger.Logger) error {
	if element == nil {
		return fmt.Errorf("element is nil")
	}

	// Get element position
	shape, err := element.Shape()
	if err != nil {
		return fmt.Errorf("failed to get element shape: %w", err)
	}
	box := shape.Box()

	// Calculate target scroll position (center element in viewport)
	viewportHeight := page.MustEval(`window.innerHeight`).Int()
	targetY := int(box.Y) - viewportHeight/2

	if targetY < 0 {
		targetY = 0
	}

	// Scroll to element
	return HumanScroll(page, targetY, log)
}

// ScrollWithPauses scrolls page with reading pauses
func ScrollWithPauses(page *rod.Page, log *logger.Logger) {
	rand.Seed(time.Now().UnixNano())

	pageHeight := page.MustEval(`document.body.scrollHeight`).Int()
	currentY := 0

	for currentY < pageHeight-1000 {
		// Scroll a section
		scrollAmount := rand.Intn(400) + 200
		currentY += scrollAmount
		
		HumanScroll(page, currentY, log)

		// Reading pause
		wordCount := rand.Intn(50) + 20
		ReadingDelay(wordCount)

		// Sometimes scroll up to reread
		if rand.Float64() < 0.15 {
			scrollUp := rand.Intn(150) + 50
			page.MustEval(fmt.Sprintf(`window.scrollBy(0, -%d)`, scrollUp))
			currentY -= scrollUp
			time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
		}
	}

	log.Stealth("Scroll with pauses completed", nil)
}
