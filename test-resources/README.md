# AI-Assisted Live Coding Challenge

**Duration**: 75 minutes  
**Position**: Mid-Senior Software Engineer  
**Focus**: Engineering principles, TDD practices, AI-assisted development

## Objective

Build a lead ingestion automation solution that demonstrates your ability to write quality code using AI assistance while following TDD practices.

## Challenge Requirements

### Core Functionality

Build a CLI application that:
1. Reads a CSV file containing lead data
2. Processes each lead by calling external APIs
3. Implements business logic for create/update decisions
4. Provides comprehensive logging and error handling

### Technical Requirements

#### 1. Test-Driven Development

- Write tests FIRST before implementing functionality

#### 2. AI-Assisted Development

Use your preferred AI coding assistant (Cursor, Windsurf, Claude, etc.) to help with:
- Exploring the problem
- Plan how to approach specific problems
- Generate test cases for specific scenarios
- Implement the specific solution in code 
- Review the specific implementation

#### 3. Engineering Principles 

- **Clean Code**: Readable, maintainable, well-structured
- **SOLID Principles**: Proper separation of concerns
- **Error Handling**: Graceful handling of failures and invalid data
- **Logging**: Structured logging with appropriate levels

### Business Logic

For each lead in the CSV:
1. **Lookup** existing lead by email
2. **Decision Logic**:
   - If found AND data differs → UPDATE
   - If found AND data identical → SKIP (log only)
   - If not found → CREATE
   - If API error → RETRY with exponential backoff
3. **Logging**: Record all actions with timestamps and outcomes

### Data Validation

Implement validation for:
- Email format validation
- Required fields presence
- Company name normalization
- Source value allowlist
- Duplicate email detection within CSV

### Error Scenarios

- Network timeouts
- Invalid CSV format
- Missing required fields
- API rate limiting (429 responses)
- Malformed API responses
- File not found
- Permission errors

## Setup Instructions

### 1. Start Mock API Server

```bash
cd ai-interview
npm install
npm run start-api
```

[Mock API - Setup Instructions and test examples](https://gist.github.com/fernando-indebted/56cae53b225467eab9797ad60e772854#file-setup-md)

The API will run on `http://localhost:3001`

## Expected Deliverables

1. **Working CLI Application** with proper argument parsing
2. **Comprehensive Test Suite** with multiple test scenarios
3. **Updated README** with usage instructions
4. **Clean, documented code** following best practices
5. **Demonstration** of AI-assisted development process

## Technology Choices

Choose your preferred:
- **Language**: TypeScript/Node.js, Python, Go, Java, C#, Rust
- **Testing Framework**: Jest, Pytest, Go Test, JUnit, xUnit, etc.
- **HTTP Client**: Axios, Fetch, Requests, etc.
- **CLI Framework**: Commander.js, Click, Cobra, etc.

## Testing Your Implementation

Once you've built your application, test it with the provided `leads.csv` file:

```bash
# Example command (adjust based on your implementation)
./your-app --file leads.csv --api-url http://localhost:3001
```

Expected behavior:
1. `alice@example.com` - Should UPDATE (exists in sample data)
2. `bob@startup.com` - Should UPDATE (exists in sample data)  
3. `charlie@peanuts.com` - Should CREATE (new lead)
4. `diana@wonderwoman.com` - Should CREATE (new lead)
5. `invalid-email` - Should log validation error
6. Duplicate `alice@example.com` - Should skip or handle appropriately

## Tips for Success

- **Start with tests** - Let TDD guide your design
- **Use AI effectively** - Leverage it for planning, boilerplate, test cases, code and code review
- **Document as you go** - Use AI to help with documentation
- **Think out loud** - Explain your process during the interview

## Troubleshooting

### API Server Won't Start

- Check if port 3001 is already in use: `lsof -i :3001`
- Kill existing process: `kill -9 <PID>`
- Try different port: `PORT=3002 npm run start-api`

### Debugging Tips

1. Check the API server logs for request details
2. Use the `/api/health` endpoint to verify server status
3. Use the `/api/leads` endpoint to see current state
4. Test individual API calls with curl before implementing
5. Monitor response times and handle timeouts appropriately

### Dependencies Not Installing

- Clear npm cache: `npm cache clean --force`
- Delete node_modules and reinstall: `rm -rf node_modules && npm install`
- Check Node.js version: `node --version`

### AI Assistant Issues

- Ensure you have active internet connection
- Check AI assistant authentication
- Try restarting your IDE
- Use specific, detailed prompts

### Getting Help

During the interview:
- Ask clarifying questions about requirements
- Explain your thinking process
- Demonstrate how you use AI assistance
- Show your problem-solving approach
- Don't hesitate to ask for guidance

Good luck! 🚀