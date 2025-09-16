package processor

import (
	"code/internal/models"
)

// LeadProcessor handles the business logic for processing leads
type LeadProcessor struct {
	apiClient APIClient
}

// APIClient interface for API operations
type APIClient interface {
	LookupLead(email string) (*LookupResponse, error)
	CreateLead(lead *models.Lead) (*models.Lead, error)
	UpdateLead(lead *models.Lead) (*models.Lead, error)
}

// LookupResponse represents the response from lookup API
type LookupResponse struct {
	Found bool
	Lead  *models.Lead
}

// ProcessResult represents the result of processing a lead
type ProcessResult struct {
	Action      string
	Lead        *models.Lead
	CreatedLead *models.Lead
	UpdatedLead *models.Lead
	Error       error
}

// NewLeadProcessor creates a new lead processor
func NewLeadProcessor(apiClient APIClient) *LeadProcessor {
	return &LeadProcessor{
		apiClient: apiClient,
	}
}

// ProcessLead processes a single lead according to business rules
func (p *LeadProcessor) ProcessLead(lead *models.Lead) (*ProcessResult, error) {
	// Validate the lead first
	if err := lead.Validate(); err != nil {
		return &ProcessResult{
			Action: "VALIDATION_ERROR",
			Lead:   lead,
			Error:  err,
		}, nil
	}

	// Look up existing lead by email
	lookupResp, err := p.apiClient.LookupLead(lead.Email)
	if err != nil {
		return &ProcessResult{
			Action: "API_ERROR",
			Lead:   lead,
			Error:  err,
		}, nil
	}

	// If lead not found, create new lead
	if !lookupResp.Found {
		createdLead, err := p.apiClient.CreateLead(lead)
		if err != nil {
			return &ProcessResult{
				Action: "CREATE_ERROR",
				Lead:   lead,
				Error:  err,
			}, nil
		}

		return &ProcessResult{
			Action:      "CREATE",
			Lead:        lead,
			CreatedLead: createdLead,
		}, nil
	}

	// Lead found - check if data differs
	existingLead := lookupResp.Lead
	if lead.IsEqual(existingLead) {
		// Data is identical, skip
		return &ProcessResult{
			Action: "SKIP",
			Lead:   lead,
		}, nil
	}

	// Data differs, update the lead
	updatedLead, err := p.apiClient.UpdateLead(lead)
	if err != nil {
		return &ProcessResult{
			Action: "UPDATE_ERROR",
			Lead:   lead,
			Error:  err,
		}, nil
	}

	return &ProcessResult{
		Action:      "UPDATE",
		Lead:        lead,
		UpdatedLead: updatedLead,
	}, nil
}
