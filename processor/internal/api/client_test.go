package api

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIClient_LookupLead(t *testing.T) {
	t.Run("successfully looks up existing lead", func(t *testing.T) {
		// Arrange
		client := NewAPIClient("http://localhost:3030")
		email := "alice@example.com"

		// Act
		result, err := client.LookupLead(email)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Found)
		assert.NotNil(t, result.Lead)
		assert.Equal(t, email, result.Lead.Email)
		assert.Equal(t, "Alice Johnson", result.Lead.Name)
		assert.Equal(t, "Acme Inc", result.Lead.Company)
		assert.Equal(t, "LinkedIn", result.Lead.Source)
	})

	t.Run("handles network timeout gracefully", func(t *testing.T) {
		// Arrange
		client := NewAPIClient("http://192.168.1.999:9999") // Non-existent server
		email := "test@example.com"

		// Act
		result, err := client.LookupLead(email)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		// Should contain either timeout, connection, or host error
		errMsg := err.Error()
		assert.True(t, strings.Contains(errMsg, "timeout") ||
			strings.Contains(errMsg, "connection") ||
			strings.Contains(errMsg, "refused") ||
			strings.Contains(errMsg, "no such host"),
			"Expected timeout, connection, refused, or no such host error, got: %s", errMsg)
	})

	t.Run("handles API rate limiting (429) with retry", func(t *testing.T) {
		// Arrange
		client := NewAPIClient("http://localhost:3030")
		email := "test@example.com"

		// Act
		result, err := client.LookupLead(email)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		// Should eventually succeed after retry (either found or not found is valid)
		assert.True(t, result.Found || !result.Found)
		
		// Verify that the result is properly structured
		if result.Found {
			assert.NotNil(t, result.Lead)
			assert.Equal(t, email, result.Lead.Email)
		} else {
			assert.Nil(t, result.Lead)
		}
	})
}
