# API Gateway Service

## 📌 Overview

The API Gateway is the entry point for the notification system. It provides a RESTful HTTP API for submitting notifications, caches them in Redis, and publishes them to Kafka for downstream processing.

**Service Type:** REST API  
**Language:** Golang  
**Port:** 8080

---

## 🏗️ Architecture
```
HTTP Client
    ↓
┌─────────────────┐
│   API Gateway   │
│   (Port 8080)   │
└────────┬────────┘
         │
         ├→ Redis (Cache)
         │  - Key: notification:{id}
         │  - TTL: 1 hour
         │
         └→ Kafka (Publish)
            - Topic: notifications
            - Format: JSON
```

**Position in System:**
- First service in the pipeline
- Receives HTTP POST/GET requests
- Validates and caches notifications
- Publishes to Kafka for async processing

---

## 🛠️ Tech Stack

- **Language:** Go 1.21+
- **Web Framework:** Standard library (`net/http`)
- **Cache:** Redis (go-redis/redis)
- **Message Queue:** Kafka (segmentio/kafka-go)
- **UUID Generation:** google/uuid

---

## 🚀 Quick Start

### **Prerequisites**
- Go 1.21 or higher
- Redis running on `localhost:6379`
- Kafka running on `localhost:9092`

### **Installation**
```bash
cd api-gateway

# Install dependencies
go mod download

# Run the service
go run main.go
```

**Expected output:**
```
🚀 Starting API Gateway...
✓ Connected to Redis
✓ Kafka producer initialized
✅ API Gateway running on :8080
```

---

## 📡 API Endpoints

### **GET /health**
- Returns API GW health status
  **Request:**
```http
Invoke-WebRequest -Uri http://localhost:8080/health
```

**Response (Success - 200 OK):**
```json
{
  "StatusCode"        : "200",
  "StatusDescription" : "OK",
  "Content"           : "{\"service\":\"api-gateway\",\"status\":\"healthy\",\"timestamp\":1763372088}",
  "RawContent"        : "HTTP/1.1 200 OK",
  ...
}
```

### **GET /notifications/:id**

Retrieve a **recently submitted** notification from cache.

**⚠️ Cache-Only Endpoint:**
- Returns notifications from the **last 1 hour only**
- Redis cache with 1-hour TTL
- For delivery status, use Delivery Service API
- For analytics/history, use Analytics Service API

**Architecture Note:**
API Gateway only handles submission and short-term caching.
For persistent data:
- Delivery results → Query MongoDB (Delivery Service)
- Analytics/history → Query Elasticsearch (Analytics Service)

**Use Case:** Quick lookup of recently submitted notifications

**Request:**
```http
GET http://localhost:8080/notifications/{notification-id}
```

**Response (Success - 200 OK):**
```json
{
  "id": "abc-123-def-456",
  "user_id": "user-123",
  "type": "EMAIL",
  ...
}
```

**Response (After 1 hour - 404 Not Found):**
```json
{
  "error": "Notification not found"
}
```

### **POST /notifications**

Submit a new notification for delivery.

**Request:**
```http
POST http://localhost:8080/notifications
Content-Type: application/json

{
  "user_id": "string",
  "type": "EMAIL | SMS | PUSH",
  "recipient": "string",
  "subject": "string (optional, EMAIL only)",
  "message": "string"
}
```

**Response (Success - 200 OK):**
```json
{
  "id": "uuid",
  "status": "queued",
  "timestamp": "2025-11-17T08:00:00Z"
}
```

**Response (Error - 400 Bad Request):**
```json
{
  "error": "Invalid notification type. Must be EMAIL, SMS, or PUSH"
}
```

**Response (Error - 500 Internal Server Error):**
```json
{
  "error": "Failed to publish to Kafka: ..."
}
```

---

## 📊 Request Examples

### **Email Notification:**
```bash
curl -X POST http://localhost:8080/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "type": "EMAIL",
    "recipient": "user@example.com",
    "subject": "Welcome!",
    "message": "Thanks for signing up"
  }'
```

**PowerShell:**
```powershell
Invoke-WebRequest -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body '{"user_id":"user-123","type":"EMAIL","recipient":"user@example.com","subject":"Welcome!","message":"Thanks for signing up"}'
```

### **SMS Notification:**
```bash
curl -X POST http://localhost:8080/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-456",
    "type": "SMS",
    "recipient": "+1234567890",
    "message": "Your verification code is 123456"
  }'
```

### **Push Notification:**
```bash
curl -X POST http://localhost:8080/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-789",
    "type": "PUSH",
    "recipient": "device-token-abc123",
    "message": "You have a new message"
  }'
```

---

## 🗄️ Data Storage

### **Redis Cache**

**Purpose:** Cache notifications for quick retrieval

**Key Pattern:** `notification:{uuid}`  
**TTL:** 1 hour (3600 seconds)  
**Value:** JSON string of notification

**Example:**
```
Key: notification:abc-123-def-456
Value: {"id":"abc-123-def-456","user_id":"user-123","type":"EMAIL",...}
Expiry: 3600 seconds
```

### **Kafka Topic**

**Topic Name:** `notifications`  
**Partitions:** 1 (default)  
**Format:** JSON  
**Key:** notification ID (for partitioning)

**Message Structure:**
```json
{
  "id": "uuid",
  "user_id": "string",
  "type": "EMAIL|SMS|PUSH",
  "recipient": "string",
  "subject": "string (optional)",
  "message": "string",
  "timestamp": "ISO8601 datetime"
}
```

---

## 🧪 Testing

### **Test the Service is Running:**
```bash
# Should return 404 (no GET endpoint, only POST)
curl http://localhost:8080/notifications
```

### **Send Test Notification:**
```bash
curl -X POST http://localhost:8080/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "test-user",
    "type": "EMAIL",
    "recipient": "test@example.com",
    "subject": "Test",
    "message": "Hello World"
  }'
```

### **Verify in Redis:**
```bash
docker exec -it notification_redis redis-cli

# Find notification keys
KEYS notification:*

# Get a specific notification
GET notification:{id}
```

### **Verify in Kafka:**
```bash
docker exec -it notification_kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic notifications \
  --from-beginning
```

---

## 🔧 Configuration

**All configuration is in `main.go`:**
```go
const (
    serverPort    = ":8080"
    redisAddr     = "localhost:6379"
    kafkaBroker   = "localhost:9092"
    kafkaTopic    = "notifications"
    cacheDuration = 1 * time.Hour
)
```

**To change:**
1. Edit the constants in `main.go`
2. Rebuild and restart the service

**Future improvement:** Use environment variables for configuration

---

## 📂 Project Structure
```
api-gateway/
├── main.go              # Application entry point
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
└── README.md               # This file
```

---

## 🐛 Troubleshooting

### **Error: "connection refused" (Redis)**
- Ensure Redis is running: `docker-compose ps`
- Check Redis port: `docker-compose logs notification_redis`
- Verify connection: `redis-cli ping`

### **Error: "connection refused" (Kafka)**
- Ensure Kafka is running: `docker-compose ps`
- Check Kafka logs: `docker-compose logs notification_kafka`
- Verify topic exists: `kafka-topics --list --bootstrap-server localhost:9092`

### **Error: "bind: address already in use"**
- Port 8080 is already taken
- Find process: `netstat -ano | findstr :8080` (Windows)
- Kill process or change port in code

---

## 🔄 Flow Diagram
```
1. Client sends POST /notifications
         ↓
2. API Gateway receives request
         ↓
3. Generate UUID for notification
         ↓
4. Add timestamp
         ↓
5. Cache in Redis (1 hour TTL)
         ↓
6. Publish to Kafka (notifications topic)
         ↓
7. Return response to client
         ↓
8. Notification Processor picks up message
```

---

## 📈 Metrics & Monitoring

**Current Status:** No metrics implemented

**Future Improvements:**
- Request counter
- Error rate tracking
- Response time histogram
- Cache hit/miss ratio

---

## 🚀 Future Enhancements

- [ ] Authentication/Authorization (API keys)
- [ ] Rate limiting per user
- [ ] Request validation middleware
- [ ] Metrics endpoint (`GET /metrics`)
- [ ] Configuration via environment variables
- [ ] Structured logging with levels
- [ ] OpenAPI/Swagger documentation

---

## 📝 Notes

- No authentication implemented (open API)
- No rate limiting (unlimited requests)
- Redis cache is optional (system works if Redis is down)
- Kafka publish is synchronous (waits for ack)

---

**Last Updated:** November 17, 2025