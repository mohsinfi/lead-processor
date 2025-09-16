package csv

import (
	"code/internal/models"
	"encoding/csv"
	"os"
)

// CSVReader handles reading and parsing CSV files
type CSVReader struct{}

// NewCSVReader creates a new CSV reader
func NewCSVReader() *CSVReader {
	return &CSVReader{}
}

// ReadLeads reads leads from a CSV file
func (r *CSVReader) ReadLeads(filePath string) ([]*models.Lead, error) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create CSV reader
	csvReader := csv.NewReader(file)

	// Read all records
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Skip header row and convert records to leads
	var leads []*models.Lead
	for i, record := range records {
		if i == 0 { // Skip header
			continue
		}

		if len(record) >= 4 {
			lead := models.NewLead(record[0], record[1], record[2], record[3])
			leads = append(leads, lead)
		}
	}

	return leads, nil
}
