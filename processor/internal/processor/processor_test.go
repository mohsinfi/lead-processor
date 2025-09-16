package processor

import (
	"code/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockAPIClient for testing
type MockAPIClient struct {
	lookupResponse *LookupResponse
	lookupError    error
	createResponse *models.Lead
	createError    error
	updateResponse *models.Lead
	updateError    error
}

func (m *MockAPIClient) LookupLead(email string) (*LookupResponse, error) {
	return m.lookupResponse, m.lookupError
}

func (m *MockAPIClient) CreateLead(lead *models.Lead) (*models.Lead, error) {
	return m.createResponse, m.createError
}

func (m *MockAPIClient) UpdateLead(lead *models.Lead) (*models.Lead, error) {
	return m.updateResponse, m.updateError
}

func TestLeadProcessor_ProcessLead(t *testing.T) {
	t.Run("creates new lead when not found in API", func(t *testing.T) {
		// Arrange
		lead := models.NewLead("John Doe", "john@example.com", "Test Corp", "LinkedIn")
		createdLead := models.NewLead("John Doe", "john@example.com", "Test Corp", "LinkedIn")

		mockAPI := &MockAPIClient{
			lookupResponse: &LookupResponse{Found: false}, // Lead not found
			createResponse: createdLead,
		}

		processor := NewLeadProcessor(mockAPI)

		// Act
		result, err := processor.ProcessLead(lead)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "CREATE", result.Action)
		assert.Equal(t, lead, result.Lead)
		assert.NotNil(t, result.CreatedLead)
		assert.Equal(t, "john@example.com", result.CreatedLead.Email)
	})

	t.Run("updates existing lead when data differs", func(t *testing.T) {
		// Arrange
		newLead := models.NewLead("John Smith", "john@example.com", "New Corp", "Website")
		existingLead := models.NewLead("John Doe", "john@example.com", "Old Corp", "LinkedIn")
		updatedLead := models.NewLead("John Smith", "john@example.com", "New Corp", "Website")

		mockAPI := &MockAPIClient{
			lookupResponse: &LookupResponse{
				Found: true,
				Lead:  existingLead,
			},
			updateResponse: updatedLead,
		}

		processor := NewLeadProcessor(mockAPI)

		// Act
		result, err := processor.ProcessLead(newLead)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "UPDATE", result.Action)
		assert.Equal(t, newLead, result.Lead)
		assert.NotNil(t, result.UpdatedLead)
		assert.Equal(t, "John Smith", result.UpdatedLead.Name)
		assert.Equal(t, "New Corp", result.UpdatedLead.Company)
		assert.Equal(t, "Website", result.UpdatedLead.Source)
	})

	t.Run("skips lead when found and data is identical", func(t *testing.T) {
		// Arrange
		lead := models.NewLead("John Doe", "john@example.com", "Test Corp", "LinkedIn")
		existingLead := models.NewLead("John Doe", "john@example.com", "Test Corp", "LinkedIn")

		mockAPI := &MockAPIClient{
			lookupResponse: &LookupResponse{
				Found: true,
				Lead:  existingLead,
			},
		}

		processor := NewLeadProcessor(mockAPI)

		// Act
		result, err := processor.ProcessLead(lead)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "SKIP", result.Action)
		assert.Equal(t, lead, result.Lead)
		assert.Nil(t, result.CreatedLead)
		assert.Nil(t, result.UpdatedLead)
	})

	t.Run("returns validation error for invalid lead data", func(t *testing.T) {
		// Arrange
		invalidLead := models.NewLead("", "invalid-email", "", "InvalidSource")

		mockAPI := &MockAPIClient{}
		processor := NewLeadProcessor(mockAPI)

		// Act
		result, err := processor.ProcessLead(invalidLead)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "VALIDATION_ERROR", result.Action)
		assert.Equal(t, invalidLead, result.Lead)
		assert.NotNil(t, result.Error)
		assert.Contains(t, result.Error.Error(), "name is required")
		assert.Contains(t, result.Error.Error(), "valid email is required")
		assert.Contains(t, result.Error.Error(), "company is required")
		assert.Contains(t, result.Error.Error(), "source must be one of")
	})

	t.Run("handles duplicate email detection within CSV", func(t *testing.T) {
		// Arrange
		lead1 := models.NewLead("John Doe", "john@example.com", "Test Corp", "LinkedIn")
		lead2 := models.NewLead("Jane Smith", "john@example.com", "Another Corp", "Website")

		mockAPI := &MockAPIClient{
			lookupResponse: &LookupResponse{Found: false}, // Lead not found
			createResponse: lead1,
		}

		processor := NewLeadProcessor(mockAPI)

		// Act - Process first lead
		result1, err1 := processor.ProcessLead(lead1)

		// Act - Process second lead with same email
		result2, err2 := processor.ProcessLead(lead2)

		// Assert - First lead should be created successfully
		assert.NoError(t, err1)
		assert.Equal(t, "CREATE", result1.Action)
		assert.Equal(t, "john@example.com", result1.CreatedLead.Email)

		// Assert - Second lead should also be created (business logic allows duplicates)
		assert.NoError(t, err2)
		assert.Equal(t, "CREATE", result2.Action)
		assert.Equal(t, "john@example.com", result2.CreatedLead.Email)
	})

	t.Run("handles API lookup error gracefully", func(t *testing.T) {
		// Arrange
		lead := models.NewLead("John Doe", "john@example.com", "Test Corp", "LinkedIn")

		mockAPI := &MockAPIClient{
			lookupError: assert.AnError, // Simulate API error
		}

		processor := NewLeadProcessor(mockAPI)

		// Act
		result, err := processor.ProcessLead(lead)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "API_ERROR", result.Action)
		assert.Equal(t, lead, result.Lead)
		assert.NotNil(t, result.Error)
		assert.Equal(t, assert.AnError, result.Error)
	})
}
