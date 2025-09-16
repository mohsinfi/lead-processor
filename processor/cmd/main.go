package cmd

import (
	"code/internal/api"
	"code/internal/csv"
	"code/internal/models"
	"code/internal/processor"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var currentLogLevel LogLevel = INFO

// initLogger initializes the logger with the specified level
func initLogger(level string) {
	currentLogLevel = parseLogLevel(level)
	log.SetFlags(0) // Remove default timestamp, we'll add our own
}

// parseLogLevel converts string to LogLevel
func parseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn", "warning":
		return WARN
	case "error":
		return ERROR
	default:
		return INFO
	}
}

// logMessage logs a message with structured format
func logMessage(level LogLevel, levelStr, msg string, fields ...interface{}) {
	if level < currentLogLevel {
		return
	}

	timestamp := time.Now().Format("2006-01-02T15:04:05Z")
	logEntry := fmt.Sprintf("[%s] %s: %s", timestamp, levelStr, msg)

	// Add fields if provided
	if len(fields) > 0 {
		logEntry += " |"
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				logEntry += fmt.Sprintf(" %v=%v", fields[i], fields[i+1])
			}
		}
	}

	log.Println(logEntry)
}

// LogDebug logs a debug message
func LogDebug(msg string, fields ...interface{}) {
	logMessage(DEBUG, "DEBUG", msg, fields...)
}

// LogInfo logs an info message
func LogInfo(msg string, fields ...interface{}) {
	logMessage(INFO, "INFO", msg, fields...)
}

// LogWarn logs a warning message
func LogWarn(msg string, fields ...interface{}) {
	logMessage(WARN, "WARN", msg, fields...)
}

// LogError logs an error message
func LogError(msg string, err error, fields ...interface{}) {
	allFields := append(fields, "error", err.Error())
	logMessage(ERROR, "ERROR", msg, allFields...)
}

// APIClientAdapter adapts the api.APIClient to the processor.APIClient interface
type APIClientAdapter struct {
	client *api.APIClient
}

func (a *APIClientAdapter) LookupLead(email string) (*processor.LookupResponse, error) {
	resp, err := a.client.LookupLead(email)
	if err != nil {
		return nil, err
	}

	return &processor.LookupResponse{
		Found: resp.Found,
		Lead:  convertAPIToProcessorLead(resp.Lead),
	}, nil
}

func (a *APIClientAdapter) CreateLead(lead *models.Lead) (*models.Lead, error) {
	return a.client.CreateLead(lead)
}

func (a *APIClientAdapter) UpdateLead(lead *models.Lead) (*models.Lead, error) {
	return a.client.UpdateLead(lead)
}

func convertAPIToProcessorLead(apiLead *api.Lead) *models.Lead {
	if apiLead == nil {
		return nil
	}

	return &models.Lead{
		ID:        apiLead.ID,
		Name:      apiLead.Name,
		Email:     apiLead.Email,
		Company:   apiLead.Company,
		Source:    apiLead.Source,
		CreatedAt: apiLead.CreatedAt,
	}
}

var rootCmd = &cobra.Command{
	Use:   "lead-processor",
	Short: "Lead ingestion automation tool",
	Long:  `A CLI tool for processing lead data from CSV files and managing them via external APIs.`,
}

// Execute runs the CLI application
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add global flags here
	rootCmd.PersistentFlags().StringP("api-url", "u", "http://localhost:3030", "API base URL")
}

var processCmd = &cobra.Command{
	Use:   "process [file]",
	Short: "Process leads from a CSV file",
	Long:  `Process leads from a CSV file and manage them via external APIs.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runProcessCommand,
}

func init() {
	rootCmd.AddCommand(processCmd)
}

func runProcessCommand(cmd *cobra.Command, args []string) error {
	// Get flags
	apiURL, _ := cmd.Flags().GetString("api-url")

	// Initialize structured logging with default level
	initLogger("info")

	// Get CSV file path
	csvFile := args[0]

	LogInfo("Starting lead processing", "csvFile", csvFile, "apiURL", apiURL)

	fmt.Printf("Processing leads from: %s\n", csvFile)
	fmt.Printf("API URL: %s\n", apiURL)

	// Initialize components
	apiClient := api.NewAPIClient(apiURL)
	csvReader := csv.NewCSVReader()

	// Create adapter to make API client compatible with processor interface
	apiAdapter := &APIClientAdapter{client: apiClient}
	leadProcessor := processor.NewLeadProcessor(apiAdapter)

	// Read leads from CSV
	LogInfo("Reading leads from CSV file")
	fmt.Println("Reading leads from CSV file...")
	leads, err := csvReader.ReadLeads(csvFile)
	if err != nil {
		LogError("Failed to read CSV file", err, "csvFile", csvFile)
		return fmt.Errorf("failed to read CSV file: %w", err)
	}

	LogInfo("CSV file read successfully", "leadCount", len(leads))
	fmt.Printf("Found %d leads to process\n", len(leads))

	// Process each lead
	createCount := 0
	updateCount := 0
	skipCount := 0
	errorCount := 0

	for i, lead := range leads {
		LogInfo("Processing lead", "progress", fmt.Sprintf("%d/%d", i+1, len(leads)), "name", lead.Name, "email", lead.Email)
		fmt.Printf("Processing lead %d/%d: %s (%s)\n", i+1, len(leads), lead.Name, lead.Email)

		result, err := leadProcessor.ProcessLead(lead)
		if err != nil {
			LogError("Lead processing failed", err, "name", lead.Name, "email", lead.Email)
			fmt.Printf("  Error: %v\n", err)
			errorCount++
			continue
		}

		switch result.Action {
		case "CREATE":
			LogInfo("Lead created successfully", "name", lead.Name, "email", lead.Email)
			fmt.Printf("  ✓ Created new lead\n")
			createCount++
		case "UPDATE":
			LogInfo("Lead updated successfully", "name", lead.Name, "email", lead.Email)
			fmt.Printf("  ✓ Updated existing lead\n")
			updateCount++
		case "SKIP":
			LogInfo("Lead skipped (no changes needed)", "name", lead.Name, "email", lead.Email)
			fmt.Printf("  - Skipped (no changes needed)\n")
			skipCount++
		case "VALIDATION_ERROR":
			LogWarn("Lead validation failed", "name", lead.Name, "email", lead.Email, "error", result.Error.Error())
			fmt.Printf("  ✗ Validation error: %v\n", result.Error)
			errorCount++
		case "API_ERROR":
			LogError("API error during lead processing", result.Error, "name", lead.Name, "email", lead.Email)
			fmt.Printf("  ✗ API error: %v\n", result.Error)
			errorCount++
		default:
			LogWarn("Unknown action result", "action", result.Action, "name", lead.Name, "email", lead.Email)
			fmt.Printf("  ? Unknown action: %s\n", result.Action)
			errorCount++
		}
	}

	// Log and print summary
	LogInfo("Processing completed", "totalLeads", len(leads), "created", createCount, "updated", updateCount, "skipped", skipCount, "errors", errorCount)

	fmt.Println("\n=== Processing Summary ===")
	fmt.Printf("Total leads: %d\n", len(leads))
	fmt.Printf("Created: %d\n", createCount)
	fmt.Printf("Updated: %d\n", updateCount)
	fmt.Printf("Skipped: %d\n", skipCount)
	fmt.Printf("Errors: %d\n", errorCount)

	return nil
}

func init() {
	cobra.OnInitialize(func() {
		// Initialize logging here
		fmt.Println("Lead Processor CLI initialized")
	})
}
