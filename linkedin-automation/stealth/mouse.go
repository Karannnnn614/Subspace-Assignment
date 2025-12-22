package stealth

import (
	"fmt"
	"linkedin-automation/logger"
	"math"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// Point represents a 2D coordinate
type Point struct {
	X float64
	Y float64
}

// HumanMouseMove moves mouse in a human-like Bézier curve path
func HumanMouseMove(page *rod.Page, element *rod.Element, log *logger.Logger) error {
	if element == nil {
		return fmt.Errorf("element is nil")
	}

	// Get element center position
	box, err := element.Box()
	if err != nil {
		return fmt.Errorf("failed to get element box: %w", err)
	}

	targetX := box.X + box.Width/2
	targetY := box.Y + box.Height/2

	// Get current mouse position (simulate from random start)
	rand.Seed(time.Now().UnixNano())
	startX := float64(rand.Intn(100))
	startY := float64(rand.Intn(100))

	// Generate Bézier curve points
	points := generateBezierCurve(
		Point{X: startX, Y: startY},
		Point{X: targetX, Y: targetY},
		20, // number of steps
	)

	// Move mouse along curve with variable speed
	for i, point := range points {
		// Move mouse to point
		page.Mouse.Move(point.X, point.Y, 1)

		// Variable delay - slower at start/end, faster in middle
		delay := calculateMouseDelay(i, len(points))
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	// Small overshoot and correction for realism
	overshootX := targetX + float64(rand.Intn(5)-2)
	overshootY := targetY + float64(rand.Intn(5)-2)
	page.Mouse.Move(overshootX, overshootY, 1)
	time.Sleep(50 * time.Millisecond)

	// Correct back to target
	page.Mouse.Move(targetX, targetY, 1)
	time.Sleep(30 * time.Millisecond)

	log.Stealth("Human mouse movement", map[string]interface{}{
		"target_x": targetX,
		"target_y": targetY,
		"steps":    len(points),
	})

	return nil
}

// generateBezierCurve generates points along a Bézier curve
func generateBezierCurve(start, end Point, steps int) []Point {
	rand.Seed(time.Now().UnixNano())
	
	// Generate random control points for natural curve
	cp1 := Point{
		X: start.X + (end.X-start.X)*0.25 + float64(rand.Intn(100)-50),
		Y: start.Y + (end.Y-start.Y)*0.25 + float64(rand.Intn(100)-50),
	}
	cp2 := Point{
		X: start.X + (end.X-start.X)*0.75 + float64(rand.Intn(100)-50),
		Y: start.Y + (end.Y-start.Y)*0.75 + float64(rand.Intn(100)-50),
	}

	points := make([]Point, steps)
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)
		points[i] = cubicBezier(start, cp1, cp2, end, t)
	}

	return points
}

// cubicBezier calculates a point on a cubic Bézier curve
func cubicBezier(p0, p1, p2, p3 Point, t float64) Point {
	u := 1 - t
	tt := t * t
	uu := u * u
	uuu := uu * u
	ttt := tt * t

	x := uuu*p0.X + 3*uu*t*p1.X + 3*u*tt*p2.X + ttt*p3.X
	y := uuu*p0.Y + 3*uu*t*p1.Y + 3*u*tt*p2.Y + ttt*p3.Y

	return Point{X: x, Y: y}
}

// calculateMouseDelay calculates delay based on position in path
func calculateMouseDelay(step, totalSteps int) int {
	// Slower at start and end, faster in middle (ease-in-out)
	progress := float64(step) / float64(totalSteps)
	
	// Use sine wave for smooth acceleration/deceleration
	factor := math.Sin(progress * math.Pi)
	
	minDelay := 5
	maxDelay := 20
	
	delay := minDelay + int(factor*float64(maxDelay-minDelay))
	
	// Add random jitter
	jitter := rand.Intn(5) - 2
	delay += jitter
	
	if delay < 1 {
		delay = 1
	}
	
	return delay
}

// HoverElement hovers over an element before clicking
func HoverElement(page *rod.Page, element *rod.Element, log *logger.Logger) error {
	if err := HumanMouseMove(page, element, log); err != nil {
		return err
	}

	// Hover for random duration
	hoverDuration := time.Duration(rand.Intn(1000)+500) * time.Millisecond
	time.Sleep(hoverDuration)

	log.Stealth("Hovered element", map[string]interface{}{
		"duration_ms": hoverDuration.Milliseconds(),
	})

	return nil
}

// IdleMouseMovement simulates random mouse movement
func IdleMouseMovement(page *rod.Page, log *logger.Logger) {
	rand.Seed(time.Now().UnixNano())
	
	// Random position within viewport
	x := float64(rand.Intn(800) + 100)
	y := float64(rand.Intn(600) + 100)
	
	// Move in small increments
	steps := rand.Intn(5) + 3
	for i := 0; i < steps; i++ {
		dx := float64(rand.Intn(50) - 25)
		dy := float64(rand.Intn(50) - 25)
		page.Mouse.Move(x+dx, y+dy, 1)
		time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
	}

	log.Stealth("Idle mouse movement", map[string]interface{}{
		"steps": steps,
	})
}

// ClickElement clicks an element with human-like behavior
func ClickElement(page *rod.Page, element *rod.Element, log *logger.Logger) error {
	// Hover before clicking
	if err := HoverElement(page, element, log); err != nil {
		return err
	}

	// Random delay before click
	RandomDelay(200, 600)

	// Click
	if err := element.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click element: %w", err)
	}

	// Small delay after click
	RandomDelay(300, 800)

	log.Stealth("Clicked element", nil)
	return nil
}
