package api

import (
	"code/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// APIClient handles communication with the external API
type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

// LookupResponse represents the response from the lookup API
type LookupResponse struct {
	Found bool  `json:"found"`
	Lead  *Lead `json:"lead,omitempty"`
}

// Lead represents a lead from the API response
type Lead struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Company   string    `json:"company"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"createdAt"`
}

// NewAPIClient creates a new API client
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second, // Shorter timeout for testing
		},
	}
}

// LookupLead looks up a lead by email
func (c *APIClient) LookupLead(email string) (*LookupResponse, error) {
	// Build the URL with query parameter
	apiURL := fmt.Sprintf("%s/api/leads/lookup?email=%s", c.baseURL, url.QueryEscape(email))

	// Make HTTP GET request
	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		// Check if it's a timeout error
		if isTimeoutError(err) {
			return nil, fmt.Errorf("request timeout: %w", err)
		}
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode == http.StatusTooManyRequests {
		log.Printf("Rate limit detected for email: %s, status: %d", email, resp.StatusCode)
		// Handle rate limiting with retry
		return c.handleRateLimit(apiURL, email)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Decode JSON response
	var lookupResp LookupResponse
	if err := json.NewDecoder(resp.Body).Decode(&lookupResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &lookupResp, nil
}

// CreateLead creates a new lead
func (c *APIClient) CreateLead(lead *models.Lead) (*models.Lead, error) {
	// TODO: Implement actual HTTP POST request
	// For now, return the lead with a generated ID
	createdLead := &Lead{
		ID:        "generated-id",
		Name:      lead.Name,
		Email:     lead.Email,
		Company:   lead.Company,
		Source:    lead.Source,
		CreatedAt: time.Now(),
	}

	return &models.Lead{
		ID:        createdLead.ID,
		Name:      createdLead.Name,
		Email:     createdLead.Email,
		Company:   createdLead.Company,
		Source:    createdLead.Source,
		CreatedAt: createdLead.CreatedAt,
	}, nil
}

// UpdateLead updates an existing lead
func (c *APIClient) UpdateLead(lead *models.Lead) (*models.Lead, error) {
	// TODO: Implement actual HTTP PUT request
	// For now, return the lead with updated timestamp
	now := time.Now()
	updatedLead := &models.Lead{
		ID:        lead.ID,
		Name:      lead.Name,
		Email:     lead.Email,
		Company:   lead.Company,
		Source:    lead.Source,
		CreatedAt: lead.CreatedAt,
		UpdatedAt: &now,
	}

	return updatedLead, nil
}

// isTimeoutError checks if the error is a timeout error
func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}

	// Check for net.Error timeout
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}

	// Check for context timeout
	if strings.Contains(err.Error(), "timeout") {
		return true
	}

	return false
}

// handleRateLimit handles 429 responses with exponential backoff retry
func (c *APIClient) handleRateLimit(apiURL, email string) (*LookupResponse, error) {
	maxRetries := 3
	baseDelay := 100 * time.Millisecond

	log.Printf("Starting retry with exponential backoff for email: %s, maxRetries: %d, baseDelay: %v", email, maxRetries, baseDelay)

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Calculate exponential backoff delay
		delay := baseDelay * time.Duration(1<<uint(attempt)) // 100ms, 200ms, 400ms

		log.Printf("Retry attempt %d/%d for email: %s, delay: %v", attempt+1, maxRetries, email, delay)

		// Wait before retry
		time.Sleep(delay)

		// Make retry request
		resp, err := c.httpClient.Get(apiURL)
		if err != nil {
			log.Printf("Retry attempt %d failed for email: %s, error: %v", attempt+1, email, err)
			// If it's the last attempt, return the error
			if attempt == maxRetries-1 {
				log.Printf("Max retries exceeded for email: %s, error: %v", email, err)
				return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, err)
			}
			continue
		}
		defer resp.Body.Close()

		// Check if we got a successful response
		if resp.StatusCode == http.StatusOK {
			log.Printf("Retry attempt %d succeeded for email: %s", attempt+1, email)
			var lookupResp LookupResponse
			if err := json.NewDecoder(resp.Body).Decode(&lookupResp); err != nil {
				return nil, fmt.Errorf("failed to decode response: %w", err)
			}
			return &lookupResp, nil
		}

		// If still rate limited and not the last attempt, continue retrying
		if resp.StatusCode == http.StatusTooManyRequests && attempt < maxRetries-1 {
			log.Printf("Still rate limited on attempt %d for email: %s, status: %d", attempt+1, email, resp.StatusCode)
			continue
		}

		// If we get here, it's either the last attempt or a different error
		log.Printf("API returned error after %d retries for email: %s, status: %d", attempt+1, email, resp.StatusCode)
		return nil, fmt.Errorf("API returned status %d after %d retries", resp.StatusCode, attempt+1)
	}

	log.Printf("Max retries exceeded for rate limiting for email: %s", email)
	return nil, fmt.Errorf("max retries exceeded for rate limiting")
}
