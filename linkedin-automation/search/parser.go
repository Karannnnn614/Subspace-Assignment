package search

import (
	"regexp"
	"strings"
)

// Parser provides utility functions for parsing LinkedIn data
type Parser struct{}

// NewParser creates a new parser
func NewParser() *Parser {
	return &Parser{}
}

// CleanText removes extra whitespace and newlines from text
func (p *Parser) CleanText(text string) string {
	// Replace multiple spaces with single space
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	
	// Trim whitespace
	text = strings.TrimSpace(text)
	
	return text
}

// ExtractProfileID extracts profile ID from LinkedIn URL
func (p *Parser) ExtractProfileID(profileURL string) string {
	// LinkedIn profile URLs typically: https://www.linkedin.com/in/profile-id/
	re := regexp.MustCompile(`/in/([^/\?]+)`)
	matches := re.FindStringSubmatch(profileURL)
	
	if len(matches) > 1 {
		return matches[1]
	}
	
	return ""
}

// IsValidProfileURL checks if URL is a valid LinkedIn profile URL
func (p *Parser) IsValidProfileURL(url string) bool {
	return strings.Contains(url, "linkedin.com/in/")
}

// ExtractCompanyFromHeadline attempts to extract company name from headline
func (p *Parser) ExtractCompanyFromHeadline(headline string) string {
	// Common patterns: "Position at Company" or "Position @ Company"
	patterns := []string{
		` at (.+)`,
		` @ (.+)`,
		` - (.+)`,
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(headline)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}
	
	return ""
}

// NormalizeURL removes query parameters and fragments from URL
func (p *Parser) NormalizeURL(url string) string {
	// Remove everything after ? or #
	url = regexp.MustCompile(`[?#].*`).ReplaceAllString(url, "")
	
	// Remove trailing slash
	url = strings.TrimRight(url, "/")
	
	return url
}

// ParseConnectionDegree extracts connection degree (1st, 2nd, 3rd)
func (p *Parser) ParseConnectionDegree(text string) string {
	re := regexp.MustCompile(`(\d+)(st|nd|rd|th)`)
	matches := re.FindStringSubmatch(text)
	
	if len(matches) > 0 {
		return matches[0]
	}
	
	return ""
}

// ExtractFirstName gets first name from full name
func (p *Parser) ExtractFirstName(fullName string) string {
	parts := strings.Fields(fullName)
	if len(parts) > 0 {
		return parts[0]
	}
	return fullName
}

// SanitizeText removes special characters and normalizes text
func (p *Parser) SanitizeText(text string) string {
	// Remove non-printable characters
	text = regexp.MustCompile(`[^\x20-\x7E]+`).ReplaceAllString(text, "")
	
	// Clean and normalize
	text = p.CleanText(text)
	
	return text
}
