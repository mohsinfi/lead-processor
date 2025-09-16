# Mock API Server Setup

Simple setup for the mock API server needed to test the lead processor.

## Quick Start

```bash
# Install dependencies
npm install

# Start the server
node server.js
```

You should see:
```
🚀 Mock API Server running on http://localhost:3030
```

## Test the Server

```bash
# Health check
curl http://localhost:3030/api/health

# Lookup a lead
curl "http://localhost:3030/api/leads/lookup?email=alice@example.com"
```

## What the Server Does

- Simulates a lead management API
- Randomly injects rate limiting (429) and server errors (500) to test error handling
- Stores leads in memory
- Runs on port 3030

## API Endpoints

- `GET /api/leads/lookup?email={email}` - Lookup lead by email
- `POST /api/leads/create` - Create new lead  
- `POST /api/leads/update` - Update existing lead
- `GET /api/health` - Health check