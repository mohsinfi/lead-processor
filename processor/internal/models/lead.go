package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Lead represents a lead in the system
type Lead struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Company   string     `json:"company"`
	Source    string     `json:"source"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

// NewLead creates a new lead with generated ID and timestamp
func NewLead(name, email, company, source string) *Lead {
	return &Lead{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Company:   company,
		Source:    source,
		CreatedAt: time.Now(),
	}
}

// Validate validates the lead data
func (l *Lead) Validate() error {
	var validationErrors []string

	// Validate name
	if strings.TrimSpace(l.Name) == "" {
		validationErrors = append(validationErrors, "name is required")
	}

	// Validate email
	if !isValidEmail(l.Email) {
		validationErrors = append(validationErrors, "valid email is required")
	}

	// Validate company
	if strings.TrimSpace(l.Company) == "" {
		validationErrors = append(validationErrors, "company is required")
	}

	// Validate source
	if !isValidSource(l.Source) {
		validSources := strings.Join(GetValidSources(), ", ")
		validationErrors = append(validationErrors, fmt.Sprintf("source must be one of: %s", validSources))
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(validationErrors, "; "))
	}

	return nil
}

// IsEqual compares two leads for equality (ignoring ID and timestamps)
func (l *Lead) IsEqual(other *Lead) bool {
	if other == nil {
		return false
	}

	return l.Name == other.Name &&
		l.Email == other.Email &&
		l.Company == other.Company &&
		l.Source == other.Source
}

// GetValidSources returns the list of valid source values
func GetValidSources() []string {
	return []string{
		"LinkedIn",
		"Website",
		"Conference",
		"Referral",
		"Webinar",
		"Twitter",
	}
}

// isValidEmail validates email format
func isValidEmail(email string) bool {
	if strings.TrimSpace(email) == "" {
		return false
	}

	// Basic email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// isValidSource checks if the source is in the valid sources list
func isValidSource(source string) bool {
	validSources := GetValidSources()
	for _, validSource := range validSources {
		if source == validSource {
			return true
		}
	}
	return false
}
