# API Gateway Service

## 📌 Overview

The API Gateway is the entry point for the notification system. It provides a RESTful HTTP API for submitting notifications, retrieving analytics, caches notifications in Redis, and publishes them to Kafka for downstream processing.

**Service Type:** REST API  
**Language:** Golang  
**Framework:** Gin  
**Port:** 8080

---

## 🏗️ Architecture
```
HTTP Client (React Dashboard)
    ↓
┌─────────────────────┐
│   API Gateway       │
│   (Port 8080)       │
│   + CORS Enabled    │
└──────────┬──────────┘
           │
           ├→ Redis (Cache)
           │  - Key: notification:{id}
           │  - TTL: 1 hour
           │
           ├→ Kafka (Publish)
           │  - Topic: notifications
           │  - Format: JSON
           │
           └→ Elasticsearch (Proxy)
              - Analytics queries
              - CORS bypass
```

**Position in System:**
- First service in the pipeline
- Receives HTTP POST/GET requests
- Validates and caches notifications
- Publishes to Kafka for async processing
- Proxies analytics requests to Elasticsearch

---

## 🛠️ Tech Stack

- **Language:** Go 1.21+
- **Web Framework:** Gin (github.com/gin-gonic/gin)
- **Cache:** Redis (github.com/redis/go-redis/v9)
- **Message Queue:** Kafka (github.com/segmentio/kafka-go)

---

## 🚀 Quick Start

### **Prerequisites**
- Go 1.21 or higher
- Redis running on `localhost:6379`
- Kafka running on `localhost:9092`
- Elasticsearch running on `localhost:9200`

### **Installation**
```bash
cd api-gateway

# Install dependencies
go mod tidy

# Run the service
go run cmd/main.go
```

**Expected output:**
```
[GIN-debug] Listening and serving HTTP on :8080
API Gateway listening on :8080
```

---

## 📡 API Endpoints

### **GET /health**
Returns API Gateway health status.

**Request:**
```http
GET http://localhost:8080/health
```

**Response (200 OK):**
```json
{
  "status": "ok"
}
```

---

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

**Response (201 Created):**
```json
{
  "id": "20251117123456",
  "user_id": "user-123",
  "type": "EMAIL",
  "recipient": "user@example.com",
  "subject": "Welcome",
  "message": "Thanks for signing up",
  "created_at": "2025-11-17T12:34:56Z"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "validation error message"
}
```

---

### **GET /notifications/:id**

Retrieve a **recently submitted** notification from cache.

**⚠️ Cache-Only Endpoint:**
- Returns notifications from the **last 1 hour only**
- Redis cache with 1-hour TTL
- For delivery status, use Delivery Service API
- For analytics/history, use Analytics Service API

**Request:**
```http
GET http://localhost:8080/notifications/20251117123456
```

**Response (200 OK):**
```json
{
  "id": "20251117123456",
  "user_id": "user-123",
  "type": "EMAIL",
  "recipient": "user@example.com",
  "subject": "Welcome",
  "message": "Thanks for signing up",
  "created_at": "2025-11-17T12:34:56Z"
}
```

**Response (404 Not Found):**
```json
{
  "error": "Notification not found"
}
```

---

### **POST /analytics** ⭐ NEW

Proxy endpoint for fetching analytics from Elasticsearch.

**Purpose:** Provides CORS-enabled access to Elasticsearch analytics for frontend applications.

**Request:**
```http
POST http://localhost:8080/analytics
Content-Type: application/json
```

**Response (200 OK):**
```json
{
  "took": 5,
  "hits": { "total": { "value": 55 } },
  "aggregations": {
    "by_status": {
      "buckets": [
        { "key": "SENT", "doc_count": 51 },
        { "key": "FAILED", "doc_count": 4 }
      ]
    },
    "by_type": {
      "buckets": [
        { "key": "EMAIL", "doc_count": 40 },
        { "key": "SMS", "doc_count": 8 },
        { "key": "PUSH", "doc_count": 7 }
      ]
    }
  }
}
```

**Why This Endpoint?**
- Elasticsearch doesn't natively support CORS
- Frontend (React) can't directly call Elasticsearch
- API Gateway acts as a proxy with CORS enabled

---

## 🌐 CORS Configuration

**CORS is enabled for all endpoints** to allow frontend applications (React dashboard on localhost:5173) to make requests.

**Allowed Origins:** `*` (all origins)  
**Allowed Methods:** `GET, POST, OPTIONS`  
**Allowed Headers:** `Content-Type`

**Implementation:**
```go
router.Use(func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    
    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(http.StatusOK)
        return
    }
    
    c.Next()
})
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

### **Fetch Analytics:**
```bash
curl -X POST http://localhost:8080/analytics \
  -H "Content-Type: application/json"
```

**PowerShell:**
```powershell
Invoke-WebRequest -Uri http://localhost:8080/analytics -Method POST -ContentType "application/json"
```

---

## 🗄️ Data Storage

### **Redis Cache**

**Purpose:** Cache notifications for quick retrieval

**Key Pattern:** `notification:{id}`  
**TTL:** 1 hour (3600 seconds)  
**Value:** JSON string of notification

**Example:**
```
Key: notification:20251117123456
Value: {"id":"20251117123456","user_id":"user-123","type":"EMAIL",...}
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
  "id": "20251117123456",
  "user_id": "string",
  "type": "EMAIL|SMS|PUSH",
  "recipient": "string",
  "subject": "string (optional)",
  "message": "string",
  "created_at": "ISO8601 datetime"
}
```

---

## 🧪 Testing

### **Test Health Endpoint:**
```bash
curl http://localhost:8080/health
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

### **Test Analytics Endpoint:**
```bash
curl -X POST http://localhost:8080/analytics \
  -H "Content-Type: application/json"
```

### **Verify in Redis:**
```bash
docker exec -it <redis-container-name> redis-cli

# Find notification keys
KEYS notification:*

# Get a specific notification
GET notification:20251117123456
```

### **Verify in Kafka:**
```bash
docker exec -it <kafka-container-name> kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic notifications \
  --from-beginning
```

---

## 🔧 Configuration

**Configuration in `cmd/main.go`:**
```go
// Redis
redisClient = redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

// Kafka
kafkaWriter = &kafka.Writer{
    Addr:  kafka.TCP("localhost:9092"),
    Topic: "notifications",
}

// Server
router.Run(":8080")
```

**To change:**
1. Edit the values in `cmd/main.go`
2. Rebuild and restart: `go run cmd/main.go`

**Future improvement:** Use environment variables for configuration

---

## 📂 Project Structure
```
api-gateway/
├── cmd/
│   └── main.go          # Application entry point
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
└── README.md            # This file
```

---

## 🐛 Troubleshooting

### **Error: "connection refused" (Redis)**
- Ensure Redis is running: `docker-compose ps`
- Check Redis port: `docker ps | grep redis`
- Verify connection: `redis-cli ping`

### **Error: "connection refused" (Kafka)**
- Ensure Kafka is running: `docker-compose ps`
- Check Kafka logs: `docker logs <kafka-container>`
- Verify broker: `localhost:9092`

### **Error: "Failed to fetch analytics"**
- Ensure Elasticsearch is running: `curl http://localhost:9200`
- Check index exists: `curl http://localhost:9200/_cat/indices`
- Verify analytics service is running and indexing data

### **Error: "bind: address already in use"**
- Port 8080 is already taken
- Find process: `netstat -ano | findstr :8080` (Windows) or `lsof -i :8080` (Mac/Linux)
- Kill process or change port in code

### **CORS Errors in Browser**
- Verify API Gateway CORS middleware is enabled
- Check browser console for specific CORS error
- Ensure frontend is making requests to `http://localhost:8080`

---

## 🔄 Flow Diagram
```
1. Client (React) sends POST /notifications
         ↓
2. CORS middleware processes request
         ↓
3. API Gateway receives request
         ↓
4. Generate ID (timestamp-based)
         ↓
5. Cache in Redis (1 hour TTL)
         ↓
6. Publish to Kafka (notifications topic)
         ↓
7. Return response to client
         ↓
8. Notification Processor picks up message

---

Analytics Flow:

1. Client sends POST /analytics
         ↓
2. CORS middleware processes request
         ↓
3. API Gateway proxies to Elasticsearch
         ↓
4. Elasticsearch returns aggregations
         ↓
5. API Gateway returns to client
```

---

## 📈 Metrics & Monitoring

**Current Status:** No metrics implemented

**Future Improvements:**
- Request counter by endpoint
- Error rate tracking
- Response time histogram
- Cache hit/miss ratio
- Kafka publish success/failure rate

---

## 🚀 Future Enhancements

- [ ] Authentication/Authorization (JWT tokens)
- [ ] Rate limiting per user
- [ ] Request validation middleware (gin-validator)
- [ ] Metrics endpoint (`GET /metrics` with Prometheus format)
- [ ] Configuration via environment variables
- [ ] Structured logging (logrus/zap)
- [ ] OpenAPI/Swagger documentation
- [ ] Health checks for dependencies (Redis, Kafka, Elasticsearch)
- [ ] Graceful shutdown
- [ ] Request/Response logging middleware

---

## 📝 Notes

- No authentication implemented (open API)
- No rate limiting (unlimited requests)
- CORS allows all origins (`*`) - tighten for production
- Redis cache is optional (system works if Redis is down)
- Kafka publish is synchronous (waits for ack)
- Analytics endpoint proxies directly to Elasticsearch

---

## 🔐 Security Considerations

**Current:**
- No authentication
- CORS allows all origins
- No input sanitization beyond JSON validation

**Production Requirements:**
- Implement API key authentication
- Restrict CORS to specific domains
- Add input validation and sanitization
- Enable HTTPS/TLS
- Rate limiting per client
- Request size limits

---

**Last Updated:** November 20, 2025