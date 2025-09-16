package csv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSVReader_ReadLeads(t *testing.T) {
	t.Run("reads valid CSV with all required fields", func(t *testing.T) {
		// Arrange
		reader := NewCSVReader()
		filePath := "../../testdata/leads.csv"

		// Act
		leads, err := reader.ReadLeads(filePath)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, leads)
		assert.Greater(t, len(leads), 0)

		// Check first lead
		firstLead := leads[0]
		assert.Equal(t, "Alice Johnson", firstLead.Name)
		assert.Equal(t, "alice@example.com", firstLead.Email)
		assert.Equal(t, "Acme Inc", firstLead.Company)
		assert.Equal(t, "LinkedIn", firstLead.Source)
		assert.NotEmpty(t, firstLead.ID)
		assert.NotZero(t, firstLead.CreatedAt)
	})

	t.Run("handles CSV with missing fields gracefully", func(t *testing.T) {
		// Arrange
		reader := NewCSVReader()
		filePath := "../../testdata/leads_missing_fields.csv"

		// Act
		leads, err := reader.ReadLeads(filePath)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, leads)
	})

	t.Run("handles empty CSV file", func(t *testing.T) {
		// Arrange
		reader := NewCSVReader()
		filePath := "../../testdata/empty_leads.csv"

		// Act
		leads, err := reader.ReadLeads(filePath)

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, leads)
	})
}
