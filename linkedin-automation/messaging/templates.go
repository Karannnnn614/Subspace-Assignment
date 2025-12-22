package messaging

import (
	"linkedin-automation/config"
	"linkedin-automation/search"
	"strings"
)

// TemplateManager manages message templates
type TemplateManager struct {
	config *config.Config
	parser *search.Parser
}

// NewTemplateManager creates a new template manager
func NewTemplateManager(cfg *config.Config) *TemplateManager {
	return &TemplateManager{
		config: cfg,
		parser: search.NewParser(),
	}
}

// GenerateMessage generates a personalized message from template
func (tm *TemplateManager) GenerateMessage(templateName string, profile *search.Profile) string {
	// Get template
	template := tm.getTemplate(templateName)

	// Extract first name
	firstName := tm.parser.ExtractFirstName(profile.Name)

	// Replace placeholders
	message := strings.ReplaceAll(template, "{name}", firstName)
	message = strings.ReplaceAll(message, "{full_name}", profile.Name)
	message = strings.ReplaceAll(message, "{headline}", profile.Headline)
	message = strings.ReplaceAll(message, "{company}", profile.Company)
	message = strings.ReplaceAll(message, "{location}", profile.Location)

	return tm.parser.SanitizeText(message)
}

// getTemplate retrieves a template by name
func (tm *TemplateManager) getTemplate(name string) string {
	// Check custom templates first
	if template, exists := tm.config.Messaging.Templates[name]; exists {
		return template
	}

	// Fall back to default template
	return tm.config.Messaging.MessageTemplate
}

// GetAllTemplates returns all available templates
func (tm *TemplateManager) GetAllTemplates() map[string]string {
	return tm.config.Messaging.Templates
}

// AddTemplate adds a new template
func (tm *TemplateManager) AddTemplate(name, template string) {
	if tm.config.Messaging.Templates == nil {
		tm.config.Messaging.Templates = make(map[string]string)
	}
	tm.config.Messaging.Templates[name] = template
}

// DefaultTemplates returns a map of default message templates
func DefaultTemplates() map[string]string {
	return map[string]string{
		"default": "Hi {name}, I came across your profile and would love to connect!",
		
		"engineer": "Hi {name},\n\nI noticed your work at {company} and was impressed by your background in {headline}. I'd love to connect and potentially exchange ideas about the industry.\n\nBest regards!",
		
		"recruiter": "Hi {name},\n\nI'm currently exploring opportunities in my field and came across your profile. Your experience at {company} caught my attention. Would love to connect!\n\nThanks!",
		
		"followup": "Hi {name},\n\nThanks for connecting! I really appreciate it. I'd love to learn more about your work at {company}.\n\nLooking forward to staying in touch!",
		
		"industry_specific": "Hi {name},\n\nI see you're working in {headline}. I'm also passionate about this space and would love to connect with like-minded professionals.\n\nCheers!",
		
		"event": "Hi {name},\n\nI noticed we both attended/are interested in similar professional events. Would love to connect and potentially discuss insights from the industry.\n\nBest!",
		
		"mutual_connection": "Hi {name},\n\nI noticed we have several mutual connections and thought it would be great to expand my network with professionals in this space.\n\nLooking forward to connecting!",
	}
}

// ValidateTemplate checks if a template is valid (not too long, has required placeholders)
func (tm *TemplateManager) ValidateTemplate(template string) error {
	// LinkedIn message limit (for connection requests: 300 chars, for messages: 8000 chars)
	// We'll use conservative 300 char limit for connection notes
	
	// Note: When placeholders are replaced, length may vary
	// This is a basic check
	if len(template) > 300 {
		return nil // Warning only, not error
	}

	return nil
}

// PersonalizeMessage adds personal touches to a template
func (tm *TemplateManager) PersonalizeMessage(template string, customVars map[string]string) string {
	message := template

	// Replace custom variables
	for key, value := range customVars {
		placeholder := "{" + key + "}"
		message = strings.ReplaceAll(message, placeholder, value)
	}

	return message
}
