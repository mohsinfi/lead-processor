#!/usr/bin/env node

/**
 * Mock API Server for AI Coding Interview Challenge
 * 
 * This server simulates a lead management API with the following features:
 * - CRUD operations for leads
 * - Rate limiting simulation
 * - Error handling
 * - Realistic response times
 * - In-memory data storage
 * 
 * Endpoints:
 * - GET /api/leads/lookup?email={email} - Lookup lead by email
 * - POST /api/leads/create - Create new lead
 * - POST /api/leads/update - Update existing lead
 * 
 * The server randomly injects failures to test error handling:
 * - Rate limiting (429) - 10% chance
 * - Server errors (500) - 5% chance
 * - Network delays - Random 100-1000ms
 */

const express = require('express');
const cors = require('cors');

const app = express();
const PORT = process.env.PORT || 3030;

// Middleware
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// In-memory data store
const leads = new Map();
let nextId = 1;

// Pre-populate with some sample data
const sampleLeads = [
  {
    id: "1",
    name: "Alice Johnson",
    email: "alice@example.com",
    company: "Acme Inc",
    source: "LinkedIn",
    createdAt: "2024-01-01T00:00:00Z"
  },
  {
    id: "2", 
    name: "Bob Smith",
    email: "bob@startup.com",
    company: "Startup Co",
    source: "Webinar",
    createdAt: "2024-01-02T00:00:00Z"
  }
];

// Initialize sample data
sampleLeads.forEach(lead => {
  leads.set(lead.email, lead);
  nextId = Math.max(nextId, parseInt(lead.id) + 1);
});

// Utility functions
const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));

const getRandomDelay = () => Math.floor(Math.random() * 900) + 100; // 100-1000ms

const shouldSimulateRateLimit = () => Math.random() < 0.1; // 10% chance

const shouldSimulateServerError = () => Math.random() < 0.05; // 5% chance

const validateEmail = (email) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

const validateLeadData = (data) => {
  const errors = {};
  
  if (!data.name || data.name.trim().length === 0) {
    errors.name = "Name is required";
  }
  
  if (!data.email || !validateEmail(data.email)) {
    errors.email = "Valid email is required";
  }
  
  if (!data.company || data.company.trim().length === 0) {
    errors.company = "Company is required";
  }
  
  if (!data.source || data.source.trim().length === 0) {
    errors.source = "Source is required";
  }
  
  const validSources = ["LinkedIn", "Webinar", "Conference", "Referral", "Twitter", "Website"];
  if (data.source && !validSources.includes(data.source)) {
    errors.source = `Source must be one of: ${validSources.join(", ")}`;
  }
  
  return {
    isValid: Object.keys(errors).length === 0,
    errors
  };
};

// Request logging middleware
app.use((req, res, next) => {
  const timestamp = new Date().toISOString();
  console.log(`[${timestamp}] ${req.method} ${req.path} - ${JSON.stringify(req.query)}`);
  next();
});

// Routes

/**
 * GET /api/leads/lookup
 * 
 * Lookup a lead by email address
 * 
 * Query Parameters:
 * - email (required): Email address to lookup
 * 
 * Responses:
 * - 200: Lead found or not found
 * - 400: Invalid email parameter
 * - 429: Rate limit exceeded
 * - 500: Server error
 */
app.get('/api/leads/lookup', async (req, res) => {
  // Simulate network delay
  await delay(getRandomDelay());
  
  // Simulate rate limiting
  if (shouldSimulateRateLimit()) {
    return res.status(429).json({
      error: "Rate limit exceeded",
      retryAfter: 5
    });
  }
  
  // Simulate server errors
  if (shouldSimulateServerError()) {
    return res.status(500).json({
      error: "Internal server error"
    });
  }
  
  const { email } = req.query;
  
  if (!email) {
    return res.status(400).json({
      error: "Email parameter is required"
    });
  }
  
  if (!validateEmail(email)) {
    return res.status(400).json({
      error: "Invalid email format"
    });
  }
  
  const lead = leads.get(email.toLowerCase());
  
  if (lead) {
    res.json({
      found: true,
      lead
    });
  } else {
    res.json({
      found: false
    });
  }
});

/**
 * POST /api/leads/create
 * 
 * Create a new lead
 * 
 * Request Body:
 * - name (required): Lead's full name
 * - email (required): Lead's email address
 * - company (required): Lead's company name
 * - source (required): Lead source (LinkedIn, Webinar, etc.)
 * 
 * Responses:
 * - 201: Lead created successfully
 * - 400: Validation error
 * - 409: Lead already exists
 * - 429: Rate limit exceeded
 * - 500: Server error
 */
app.post('/api/leads/create', async (req, res) => {
  // Simulate network delay
  await delay(getRandomDelay());
  
  // Simulate rate limiting
  if (shouldSimulateRateLimit()) {
    return res.status(429).json({
      error: "Rate limit exceeded",
      retryAfter: 5
    });
  }
  
  // Simulate server errors
  if (shouldSimulateServerError()) {
    return res.status(500).json({
      error: "Internal server error"
    });
  }
  
  const leadData = req.body;
  const validation = validateLeadData(leadData);
  
  if (!validation.isValid) {
    return res.status(400).json({
      error: "Validation failed",
      details: validation.errors
    });
  }
  
  const email = leadData.email.toLowerCase();
  
  // Check if lead already exists
  if (leads.has(email)) {
    return res.status(409).json({
      error: "Lead already exists",
      message: "A lead with this email address already exists"
    });
  }
  
  // Create new lead
  const newLead = {
    id: nextId.toString(),
    name: leadData.name.trim(),
    email: email,
    company: leadData.company.trim(),
    source: leadData.source,
    createdAt: new Date().toISOString()
  };
  
  leads.set(email, newLead);
  nextId++;
  
  res.status(201).json({
    success: true,
    lead: newLead
  });
});

/**
 * POST /api/leads/update
 * 
 * Update an existing lead
 * 
 * Request Body:
 * - email (required): Lead's email address (identifier)
 * - name (optional): New name
 * - company (optional): New company
 * - source (optional): New source
 * 
 * Responses:
 * - 200: Lead updated successfully
 * - 400: Validation error
 * - 404: Lead not found
 * - 429: Rate limit exceeded
 * - 500: Server error
 */
app.post('/api/leads/update', async (req, res) => {
  // Simulate network delay
  await delay(getRandomDelay());
  
  // Simulate rate limiting
  if (shouldSimulateRateLimit()) {
    return res.status(429).json({
      error: "Rate limit exceeded",
      retryAfter: 5
    });
  }
  
  // Simulate server errors
  if (shouldSimulateServerError()) {
    return res.status(500).json({
      error: "Internal server error"
    });
  }
  
  const updateData = req.body;
  
  if (!updateData.email || !validateEmail(updateData.email)) {
    return res.status(400).json({
      error: "Valid email is required"
    });
  }
  
  const email = updateData.email.toLowerCase();
  const existingLead = leads.get(email);
  
  if (!existingLead) {
    return res.status(404).json({
      error: "Lead not found",
      message: "No lead found with the provided email address"
    });
  }
  
  // Validate optional fields if provided
  const fieldsToValidate = {};
  if (updateData.name !== undefined) fieldsToValidate.name = updateData.name;
  if (updateData.company !== undefined) fieldsToValidate.company = updateData.company;
  if (updateData.source !== undefined) fieldsToValidate.source = updateData.source;
  
  // Add email for validation context
  fieldsToValidate.email = updateData.email;
  
  const validation = validateLeadData(fieldsToValidate);
  if (!validation.isValid) {
    return res.status(400).json({
      error: "Validation failed",
      details: validation.errors
    });
  }
  
  // Update lead with provided fields
  const updatedLead = {
    ...existingLead,
    updatedAt: new Date().toISOString()
  };
  
  if (updateData.name !== undefined) {
    updatedLead.name = updateData.name.trim();
  }
  if (updateData.company !== undefined) {
    updatedLead.company = updateData.company.trim();
  }
  if (updateData.source !== undefined) {
    updatedLead.source = updateData.source;
  }
  
  leads.set(email, updatedLead);
  
  res.json({
    success: true,
    lead: updatedLead
  });
});

/**
 * GET /api/health
 * 
 * Health check endpoint
 */
app.get('/api/health', (req, res) => {
  res.json({
    status: 'healthy',
    timestamp: new Date().toISOString(),
    leadsCount: leads.size
  });
});

/**
 * GET /api/leads
 * 
 * List all leads (for debugging purposes)
 */
app.get('/api/leads', (req, res) => {
  const allLeads = Array.from(leads.values());
  res.json({
    leads: allLeads,
    count: allLeads.length
  });
});

// Error handling middleware
app.use((err, req, res, next) => {
  console.error('Unhandled error:', err);
  res.status(500).json({
    error: 'Internal server error',
    message: 'An unexpected error occurred'
  });
});

// 404 handler
app.use((req, res) => {
  res.status(404).json({
    error: 'Not found',
    message: `Endpoint ${req.method} ${req.path} not found`
  });
});

// Start server
app.listen(PORT, () => {
  console.log(`🚀 Mock API Server running on http://localhost:${PORT}`);
  console.log(`📊 Health check: http://localhost:${PORT}/api/health`);
  console.log(`📋 All leads: http://localhost:${PORT}/api/leads`);
  console.log(`\n📚 API Documentation:`);
  console.log(`   GET  /api/leads/lookup?email={email}  - Lookup lead by email`);
  console.log(`   POST /api/leads/create                - Create new lead`);
  console.log(`   POST /api/leads/update                - Update existing lead`);
  console.log(`\n💡 Sample data preloaded:`);
  console.log(`   - alice@example.com (Acme Inc)`);
  console.log(`   - bob@startup.com (Startup Co)`);
  console.log(`\n⚠️  Random failures enabled:`);
  console.log(`   - Rate limiting (429): 10% chance`);
  console.log(`   - Server errors (500): 5% chance`);
  console.log(`   - Network delays: 100-1000ms`);
});

module.exports = app;