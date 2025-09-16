# Lead Processor

A simple CLI tool that processes lead data from CSV files and manages them via external APIs.

## Features

- CSV file processing with validation
- API integration for lead management  
- CREATE/UPDATE/SKIP logic based on data comparison
- Rate limiting handling with exponential backoff
- Structured logging with configurable levels

## Prerequisites

- Go 1.24.5 or higher
- Node.js (for mock API server)

## Quick Start

### 1. Start Mock API Server

```bash
# Navigate to test-resources directory
cd test-resources

# Install dependencies
npm install

# Start API server
node server.js
```

You should see:
```
🚀 Mock API Server running on http://localhost:3030
```

### 2. Run the Application

```bash
# Navigate to processor directory
cd processor

# Install dependencies
go mod tidy

# Run directly
go run . process ../test-resources/leads.csv

# Or build and run
go build -o lead-processor .
./lead-processor process ../test-resources/leads.csv
```

## Usage

```bash
# Process leads from CSV file
go run . process ../test-resources/leads.csv

# With custom API URL
go run . process ../test-resources/leads.csv --api-url http://localhost:3030

# Show help
go run . --help
```

## CSV Format

The CSV file should have these columns (header row required):
```csv
Name,Email,Company,Source
Alice Johnson,alice@example.com,Acme Inc,LinkedIn
Bob Smith,bob@startup.com,Startup Co,Webinar
Charlie Brown,charlie@peanuts.com,Peanuts Corp,Conference
Diana Prince,diana@wonderwoman.com,Justice League,Referral
Invalid User,invalid-email,Test Company,LinkedIn
```

**Valid sources:** LinkedIn, Website, Conference, Referral, Webinar, Twitter

**Sample file:** `../test-resources/leads.csv` (contains 10 test leads including validation errors)

## Project Structure

```
processor/
├── cmd/main.go              # CLI commands
├── internal/
│   ├── api/client.go        # API communication
│   ├── csv/reader.go        # CSV reading
│   ├── models/lead.go       # Data models
│   └── processor/processor.go # Business logic
├── testdata/                # Sample CSV files
└── main.go                  # Entry point
```

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...
```

## File Locations

- **Application:** `processor/` directory
- **Mock API Server:** `test-resources/` directory  
- **Sample CSV:** `test-resources/leads.csv`
- **API Server:** Runs on `http://localhost:3030`

## Business Logic

1. **Validation** - Validates email format and required fields
2. **Lookup** - Checks if lead exists in API by email
3. **Decision:**
   - **CREATE** - If lead not found → Create new lead
   - **UPDATE** - If lead found and data differs → Update existing lead
   - **SKIP** - If lead found and data identical → Skip processing
   - **ERROR** - If validation fails → Log validation error

## Error Handling

- Network timeouts
- API rate limiting (429) with exponential backoff
- Invalid CSV format
- Missing required fields
- Malformed API responses