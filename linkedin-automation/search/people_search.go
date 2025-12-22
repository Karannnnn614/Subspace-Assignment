package search

import (
	"fmt"
	"linkedin-automation/config"
	"linkedin-automation/logger"
	"linkedin-automation/stealth"
	"linkedin-automation/storage"
	"net/url"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

// Profile represents a LinkedIn profile
type Profile struct {
	Name        string
	Headline    string
	ProfileURL  string
	Company     string
	Location    string
	ConnectionDegree string
	SearchKeyword    string
	DiscoveredAt     time.Time
}

// PeopleSearcher handles LinkedIn people search
type PeopleSearcher struct {
	config *config.Config
	log    *logger.Logger
	db     *storage.DB
}

// NewPeopleSearcher creates a new people searcher
func NewPeopleSearcher(cfg *config.Config, log *logger.Logger, db *storage.DB) *PeopleSearcher {
	return &PeopleSearcher{
		config: cfg,
		log:    log,
		db:     db,
	}
}

// Search performs a people search and returns found profiles
func (ps *PeopleSearcher) Search(page *rod.Page, keyword string) ([]*Profile, error) {
	ps.log.Info(fmt.Sprintf("Searching for: %s", keyword))

	// Build search URL
	searchURL := ps.buildSearchURL(keyword)
	ps.log.Debug(fmt.Sprintf("Search URL: %s", searchURL))

	// Navigate to search page
	if err := page.Navigate(searchURL); err != nil {
		return nil, fmt.Errorf("failed to navigate to search: %w", err)
	}

	// Wait for results to load
	stealth.RandomDelay(2000, 4000)
	page.MustWaitLoad()

	// Collect profiles from multiple pages
	var allProfiles []*Profile

	for pageNum := 1; pageNum <= ps.config.Search.MaxSearchPages; pageNum++ {
		ps.log.Info(fmt.Sprintf("Scraping page %d/%d", pageNum, ps.config.Search.MaxSearchPages))

		// Random scroll to simulate reading
		stealth.ScrollWithPauses(page, ps.log)

		// Extract profiles from current page
		profiles, err := ps.extractProfiles(page, keyword)
		if err != nil {
			ps.log.Warn(fmt.Sprintf("Failed to extract profiles from page %d: %v", pageNum, err))
			continue
		}

		ps.log.Info(fmt.Sprintf("Found %d profiles on page %d", len(profiles), pageNum))
		allProfiles = append(allProfiles, profiles...)

		// Save profiles to database
		for _, profile := range profiles {
			if err := ps.db.SaveProfile(profile); err != nil {
				ps.log.Warn(fmt.Sprintf("Failed to save profile: %v", err))
			}
		}

		// Navigate to next page if not last
		if pageNum < ps.config.Search.MaxSearchPages {
			if err := ps.goToNextPage(page, pageNum); err != nil {
				ps.log.Warn(fmt.Sprintf("Failed to navigate to next page: %v", err))
				break
			}
		}
	}

	ps.log.Success(fmt.Sprintf("Search completed: %d total profiles found", len(allProfiles)), map[string]interface{}{
		"keyword": keyword,
		"pages":   ps.config.Search.MaxSearchPages,
	})

	return allProfiles, nil
}

// buildSearchURL constructs the LinkedIn people search URL
func (ps *PeopleSearcher) buildSearchURL(keyword string) string {
	baseURL := ps.config.LinkedIn.BaseURL + "/search/results/people/"
	params := url.Values{}
	params.Add("keywords", keyword)
	params.Add("origin", "GLOBAL_SEARCH_HEADER")
	
	return baseURL + "?" + params.Encode()
}

// extractProfiles extracts profile data from search results page
func (ps *PeopleSearcher) extractProfiles(page *rod.Page, keyword string) ([]*Profile, error) {
	var profiles []*Profile

	// Wait for search results container
	resultsContainer, err := page.Timeout(10 * time.Second).Element(".search-results-container")
	if err != nil {
		return nil, fmt.Errorf("search results container not found: %w", err)
	}

	// Find all profile cards
	profileCards, err := resultsContainer.Elements("li.reusable-search__result-container")
	if err != nil {
		return nil, fmt.Errorf("no profile cards found: %w", err)
	}

	ps.log.Debug(fmt.Sprintf("Found %d profile cards", len(profileCards)))

	// Extract data from each card
	for _, card := range profileCards {
		profile, err := ps.extractProfileFromCard(card, keyword)
		if err != nil {
			ps.log.Debug(fmt.Sprintf("Failed to extract profile: %v", err))
			continue
		}

		// Deduplicate
		if !ps.isDuplicate(profiles, profile) {
			profiles = append(profiles, profile)
		}

		// Small delay between processing cards
		stealth.RandomDelay(100, 300)
	}

	return profiles, nil
}

// extractProfileFromCard extracts profile data from a single card element
func (ps *PeopleSearcher) extractProfileFromCard(card *rod.Element, keyword string) (*Profile, error) {
	profile := &Profile{
		SearchKeyword: keyword,
		DiscoveredAt:  time.Now(),
	}

	// Extract name
	nameEl, err := card.Element("span.entity-result__title-text a span[aria-hidden='true']")
	if err == nil {
		profile.Name = strings.TrimSpace(nameEl.MustText())
	}

	// Extract profile URL
	linkEl, err := card.Element("a.app-aware-link")
	if err == nil {
		href := linkEl.MustProperty("href").String()
		profile.ProfileURL = strings.Split(href, "?")[0] // Remove query params
	}

	// Extract headline
	headlineEl, err := card.Element("div.entity-result__primary-subtitle")
	if err == nil {
		profile.Headline = strings.TrimSpace(headlineEl.MustText())
	}

	// Extract company/location
	secondaryEl, err := card.Element("div.entity-result__secondary-subtitle")
	if err == nil {
		profile.Location = strings.TrimSpace(secondaryEl.MustText())
	}

	// Validate profile has minimum required data
	if profile.Name == "" || profile.ProfileURL == "" {
		return nil, fmt.Errorf("incomplete profile data")
	}

	return profile, nil
}

// isDuplicate checks if profile already exists in the list
func (ps *PeopleSearcher) isDuplicate(profiles []*Profile, newProfile *Profile) bool {
	for _, p := range profiles {
		if p.ProfileURL == newProfile.ProfileURL {
			return true
		}
	}
	return false
}

// goToNextPage navigates to the next search results page
func (ps *PeopleSearcher) goToNextPage(page *rod.Page, currentPage int) error {
	ps.log.Debug(fmt.Sprintf("Navigating to page %d", currentPage+1))

	// Human-like delay before pagination
	stealth.ThinkDelay()

	// Scroll to pagination
	stealth.ScrollToBottom(page, ps.log)
	stealth.RandomDelay(1000, 2000)

	// Find next button
	nextButton, err := page.Timeout(5 * time.Second).Element("button[aria-label='Next']")
	if err != nil {
		return fmt.Errorf("next button not found: %w", err)
	}

	// Check if button is disabled
	disabled := nextButton.MustProperty("disabled").Bool()
	if disabled {
		return fmt.Errorf("next button is disabled (last page)")
	}

	// Click next with human-like behavior
	if err := stealth.ClickElement(page, nextButton, ps.log); err != nil {
		return fmt.Errorf("failed to click next button: %w", err)
	}

	// Wait for new results to load
	stealth.RandomDelay(3000, 5000)
	page.MustWaitLoad()

	return nil
}

// SearchMultipleKeywords searches for multiple keywords
func (ps *PeopleSearcher) SearchMultipleKeywords(page *rod.Page) ([]*Profile, error) {
	var allProfiles []*Profile

	for i, keyword := range ps.config.Search.Keywords {
		ps.log.Info(fmt.Sprintf("Searching keyword %d/%d: %s", i+1, len(ps.config.Search.Keywords), keyword))

		profiles, err := ps.Search(page, keyword)
		if err != nil {
			ps.log.Warn(fmt.Sprintf("Search failed for keyword '%s': %v", keyword, err))
			continue
		}

		allProfiles = append(allProfiles, profiles...)

		// Delay between keyword searches
		if i < len(ps.config.Search.Keywords)-1 {
			ps.log.Debug("Pausing between keyword searches...")
			stealth.LongPause()
		}
	}

	return allProfiles, nil
}
