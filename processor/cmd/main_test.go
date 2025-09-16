package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCommand_Initialization(t *testing.T) {
	t.Run("root command should be properly configured", func(t *testing.T) {
		// Arrange
		expectedUse := "lead-processor"
		expectedShort := "Lead ingestion automation tool"

		// Act
		actualUse := rootCmd.Use
		actualShort := rootCmd.Short

		// Assert
		assert.Equal(t, expectedUse, actualUse)
		assert.Equal(t, expectedShort, actualShort)
	})

	t.Run("root command should have required flags", func(t *testing.T) {
		// Arrange & Act
		apiURLFlag := rootCmd.PersistentFlags().Lookup("api-url")

		// Assert
		assert.NotNil(t, apiURLFlag, "api-url flag should exist")
		assert.Equal(t, "http://localhost:3030", apiURLFlag.DefValue)
	})
}
