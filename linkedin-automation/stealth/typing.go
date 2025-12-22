package stealth

import (
	"linkedin-automation/logger"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

// HumanType simulates human-like typing with variable speed and occasional typos
func HumanType(element *rod.Element, text string, log *logger.Logger) error {
	rand.Seed(time.Now().UnixNano())

	// Focus on element
	if err := element.Focus(); err != nil {
		return err
	}

	// Random delay before typing
	RandomDelay(300, 800)

	// Type each character with variable delay
	for i, char := range text {
		// Simulate occasional typo (5% chance)
		if rand.Float64() < 0.05 && i < len(text)-1 {
			// Type wrong character
			wrongChar := getRandomChar()
			element.MustInput(string(wrongChar))
			time.Sleep(time.Duration(rand.Intn(100)+50) * time.Millisecond)

			// Backspace
			element.MustType(input.Backspace)
			time.Sleep(time.Duration(rand.Intn(150)+100) * time.Millisecond)
		}

		// Type correct character
		element.MustInput(string(char))

		// Variable keystroke delay
		delay := calculateKeystrokeDelay()
		time.Sleep(time.Duration(delay) * time.Millisecond)

		// Occasional longer pause (thinking)
		if rand.Float64() < 0.1 {
			time.Sleep(time.Duration(rand.Intn(500)+300) * time.Millisecond)
		}
	}

	log.Stealth("Human typing completed", map[string]interface{}{
		"length": len(text),
		"typos":  "simulated",
	})

	return nil
}

// calculateKeystrokeDelay calculates realistic keystroke delay
func calculateKeystrokeDelay() int {
	rand.Seed(time.Now().UnixNano())
	
	// Human typing speed varies between 40-80 WPM
	// Average character interval: 120-300ms
	baseDelay := 120
	variance := rand.Intn(180)
	
	// Add random spikes for natural rhythm
	if rand.Float64() < 0.15 {
		variance += rand.Intn(200)
	}
	
	return baseDelay + variance
}

// getRandomChar returns a random character for typo simulation
func getRandomChar() rune {
	chars := []rune("qwertyuiopasdfghjklzxcvbnm")
	rand.Seed(time.Now().UnixNano())
	return chars[rand.Intn(len(chars))]
}

// TypeWithPauses types text with realistic pauses between words
func TypeWithPauses(element *rod.Element, text string, log *logger.Logger) error {
	rand.Seed(time.Now().UnixNano())
	
	words := splitIntoWords(text)
	
	for i, word := range words {
		// Type word
		if err := HumanType(element, word, log); err != nil {
			return err
		}
		
		// Pause between words (not after last word)
		if i < len(words)-1 {
			// Add space
			element.MustInput(" ")
			
			// Longer pause between words
			pause := rand.Intn(300) + 200
			time.Sleep(time.Duration(pause) * time.Millisecond)
		}
	}
	
	return nil
}

// splitIntoWords splits text into words
func splitIntoWords(text string) []string {
	words := []string{}
	currentWord := ""
	
	for _, char := range text {
		if char == ' ' {
			if currentWord != "" {
				words = append(words, currentWord)
				currentWord = ""
			}
		} else {
			currentWord += string(char)
		}
	}
	
	if currentWord != "" {
		words = append(words, currentWord)
	}
	
	return words
}

// PasteText simulates pasting text (occasionally used instead of typing)
func PasteText(element *rod.Element, text string, log *logger.Logger) error {
	// Focus element
	if err := element.Focus(); err != nil {
		return err
	}
	
	// Random delay
	RandomDelay(200, 500)
	
	// Simulate Ctrl+V paste
	element.MustSelectAllText()
	element.MustInput(text)
	
	// Slight delay after paste
	RandomDelay(300, 700)
	
	log.Stealth("Pasted text", map[string]interface{}{
		"length": len(text),
	})
	
	return nil
}
