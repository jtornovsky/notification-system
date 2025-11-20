# Real-Time Notification System

A microservices-based notification delivery system with real-time analytics, built as a learning project to demonstrate modern backend architecture and full-stack development skills.

## 🎯 Project Overview

This system provides a scalable notification delivery platform supporting multiple channels (Email, SMS, Push notifications) with real-time tracking, delivery simulation, and analytics dashboard.

**Key Features:**
- Multi-channel notification delivery (Email, SMS, Push)
- Real-time message routing and processing
- Delivery status tracking with MongoDB persistence
- Analytics aggregation with Elasticsearch
- Failure simulation for realistic testing (10% random failure rate)
- Microservices architecture with event-driven communication

---

## 🏗️ Architecture
```
┌─────────────┐
│   Client    │
│  (cURL/UI)  │
└──────┬──────┘
       │ HTTP POST
       ▼
┌─────────────────┐
│  API Gateway    │  (Golang)
│  - REST API     │
│  - Redis Cache  │
└────────┬────────┘
         │ Kafka: notifications
         ▼
┌──────────────────────┐
│ Notification         │  (Node.js)
│ Processor            │
│ - Route by type      │
└─────────┬────────────┘
          │ Kafka: email-notifications
          │        sms-notifications
          │        push-notifications
          ▼
┌──────────────────────┐
│ Delivery Service     │  (Golang)
│ - 3 Handlers         │
│ - MongoDB Storage    │
│ - Simulate Delivery  │
└─────────┬────────────┘
          │ Kafka: delivery-events
          ▼
┌──────────────────────┐
│ Analytics Service    │  (Node.js)
│ - Aggregate Metrics  │
│ - Elasticsearch      │
└─────────┬────────────┘
          │
          ▼
┌──────────────────────┐
│ React Dashboard      │  (React)
│ - Submit Notifs      │
│ - View Analytics     │
└──────────────────────┘
```

---

## 🛠️ Tech Stack

### **Backend Services**
- **Golang** (1.21+) - API Gateway, Delivery Service
- **Node.js** (18+) - Notification Processor, Analytics Service

### **Infrastructure**
- **Apache Kafka** - Event streaming and message queuing
- **MongoDB** - Delivery results persistence
- **Redis** - Caching and rate limiting
- **Elasticsearch** - Analytics data indexing and search

### **Frontend**
- **React** - Analytics dashboard UI

### **DevOps**
- **Docker Compose** - Local development environment

---

## 🚀 Quick Start

### **Prerequisites**
- Docker & Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ (for local development)

### **1. Start Infrastructure**
```bash
cd notification-system
docker-compose up -d
```

This starts:
- Kafka + Zookeeper
- MongoDB
- Redis
- Elasticsearch

### **2. Run Services Locally**

**Terminal 1 - API Gateway:**
```bash
cd api-gateway
go run cmd/main.go
```

**Terminal 2 - Notification Processor:**
```bash
cd notification-processor
npm install
node index.js
```

**Terminal 3 - Delivery Service:**
```bash
cd delivery-service
go run cmd/main.go
```

### **3. Send Test Notification**

**PowerShell:**
```powershell
Invoke-WebRequest -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body '{"user_id":"test-user","type":"EMAIL","recipient":"test@example.com","subject":"Hello","message":"Test notification"}'
```

**Bash/Git Bash:**
```bash
curl -X POST http://localhost:8080/notifications \
  -H "Content-Type: application/json" \
  -d '{"user_id":"test-user","type":"EMAIL","recipient":"test@example.com","subject":"Hello","message":"Test notification"}'
```

---

## 📡 API Documentation

### **POST /notifications**

Submit a new notification for delivery.

**Endpoint:** `http://localhost:8080/notifications`

**Request Body:**
```json
{
  "user_id": "string",
  "type": "EMAIL | SMS | PUSH",
  "recipient": "string",
  "subject": "string (optional, EMAIL only)",
  "message": "string"
}
```

**Response:**
```json
{
  "id": "uuid",
  "status": "queued",
  "timestamp": "2025-11-16T12:00:00Z"
}
```

**Example - Send Email:**
```json
{
  "user_id": "user-123",
  "type": "EMAIL",
  "recipient": "user@example.com",
  "subject": "Welcome!",
  "message": "Thanks for signing up"
}
```

**Example - Send SMS:**
```json
{
  "user_id": "user-456",
  "type": "SMS",
  "recipient": "+1234567890",
  "message": "Your verification code is 123456"
}
```

**Example - Send Push:**
```json
{
  "user_id": "user-789",
  "type": "PUSH",
  "recipient": "device-token-abc123",
  "message": "You have a new message"
}
```

---

## 📊 Services

### **1. API Gateway** (Port 8080)
- **Technology:** Golang
- **Purpose:** HTTP REST API, caching, Kafka publishing
- **Features:**
    - Accepts notification requests
    - Caches in Redis (1 hour TTL)
    - Publishes to Kafka `notifications` topic

### **2. Notification Processor**
- **Technology:** Node.js
- **Purpose:** Route notifications by type
- **Features:**
    - Consumes from `notifications` topic
    - Routes to type-specific topics:
        - `email-notifications`
        - `sms-notifications`
        - `push-notifications`

### **3. Delivery Service**
- **Technology:** Golang
- **Purpose:** Simulate delivery, persist results
- **Features:**
    - 3 concurrent handlers (Email, SMS, Push)
    - Simulates delivery with realistic delays:
        - Email: 50-200ms
        - SMS: 30-100ms
        - Push: 20-80ms
    - 10% random failure rate for testing
    - Saves results to MongoDB
    - Publishes events to `delivery-events` topic

### **4. Analytics Service** *(Work in Progress)*
- **Technology:** Node.js
- **Purpose:** Aggregate metrics, index to Elasticsearch
- **Planned Features:**
    - Success/failure rates
    - Average delivery times
    - Time-based trends

### **5. React Dashboard** *(To Do)*
- **Technology:** React
- **Purpose:** Web UI for notifications and analytics
- **Planned Features:**
    - Submit notifications form
    - Real-time delivery status
    - Analytics charts and graphs

---

## 🗄️ Data Storage

### **MongoDB - Delivery Results**
**Database:** `notifications`  
**Collection:** `delivery_results`

**Document Schema:**
```javascript
{
  _id: ObjectId("..."),
  notification_id: "uuid",
  type: "EMAIL | SMS | PUSH",
  recipient: "string",
  status: "SENT | FAILED",
  timestamp: ISODate("2025-11-16T12:00:00Z"),
  delivery_time_ms: 125,
  error_message: "string (if failed)"
}
```

**Query Examples:**
```javascript
// Count by status
db.delivery_results.countDocuments({status: "SENT"})
db.delivery_results.countDocuments({status: "FAILED"})

// Aggregate by type
db.delivery_results.aggregate([
  { $group: { _id: "$type", count: { $sum: 1 } } }
])

// Find failures
db.delivery_results.find({status: "FAILED"}).pretty()
```

### **Redis - Notification Cache**
**Key Pattern:** `notification:{id}`  
**TTL:** 1 hour  
**Value:** JSON string of notification

### **Elasticsearch - Analytics**
**Index:** `notification-analytics`  
**Purpose:** Aggregated metrics for dashboard

---

## 🧪 Testing

### **End-to-End Test**

**1. Send 10 notifications:**
```powershell
1..10 | ForEach-Object {
    Invoke-WebRequest -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body "{`"user_id`":`"user-$_`",`"type`":`"EMAIL`",`"recipient`":`"test$_@example.com`",`"subject`":`"Test`",`"message`":`"Message $_`"}"
}
```

**2. Check MongoDB results:**
```bash
docker exec -it notification_mongodb mongosh -u admin -p password

use notifications
db.delivery_results.countDocuments()
db.delivery_results.aggregate([
  { $group: { _id: "$status", count: { $sum: 1 } } }
])
```

**3. Check Kafka events:**
```bash
docker exec -it notification_kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic delivery-events --from-beginning
```

### **Test Different Types**

**Email, SMS, and Push:**
```powershell
# Email
iwr -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body '{"user_id":"u1","type":"EMAIL","recipient":"test@example.com","subject":"Hi","message":"Email test"}'

# SMS
iwr -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body '{"user_id":"u2","type":"SMS","recipient":"+1234567890","message":"SMS test"}'

# Push
iwr -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body '{"user_id":"u3","type":"PUSH","recipient":"device-token-123","message":"Push test"}'
```

---

## 📈 Monitoring

### **Service Health**
- **API Gateway:** `http://localhost:8080/health`
- **Kafka Topics:** `docker exec -it notification_kafka kafka-topics --list --bootstrap-server localhost:9092`
- **MongoDB:** `docker exec -it notification_mongodb mongosh -u admin -p password`
- **Redis:** `docker exec -it notification_redis redis-cli`
- **Elasticsearch:** `curl http://localhost:9200/_cluster/health`

### **View Logs**
```bash
# Infrastructure logs
docker-compose logs -f kafka
docker-compose logs -f mongodb

# Service logs (when running locally)
# Check terminal windows where services are running
```

---

## 🏗️ Development

### **Project Structure**
```
notification-system/
├── api-gateway/              # Golang REST API
│   ├── cmd/main.go
│   └── internal/
├── notification-processor/   # Node.js routing service
│   └── index.js
├── delivery-service/         # Golang delivery handlers
│   ├── cmd/main.go
│   └── internal/
│       ├── handlers/
│       ├── models/
│       └── mongo/
├── analytics-service/        # Node.js analytics (WIP)
├── dashboard/                # React UI (TODO)
└── docker-compose.yml        # Infrastructure setup
```

### **Adding a New Notification Type**

**1. Add simulator in `delivery-service/internal/handlers/simulators.go`:**
```go
var webhookConfig = DeliveryConfig{
    Name:         "Webhook",
    MinDelayMs:   100,
    MaxDelayMs:   500,
    FailureRate:  0.10,
    ErrorMessage: "HTTP request failed",
    SuccessEmoji: "🔗",
}

func SimulateWebhookDelivery(notification models.Notification) (string, int64, error) {
    return simulateDelivery(notification, webhookConfig)
}
```

**2. Create Kafka topic:**
```bash
docker exec -it notification_kafka kafka-topics --create --bootstrap-server localhost:9092 --topic webhook-notifications --partitions 1 --replication-factor 1
```

**3. Add handler in `main.go`:**
```go
webhookHandler, err := handlers.NewHandler(
    "Webhook",
    kafkaBrokers,
    "webhook-notifications",
    "delivery-service-webhook",
    mongoClient,
    handlers.SimulateWebhookDelivery,
)
```

---

## 🎓 Learning Objectives

This project demonstrates:
- ✅ Microservices architecture
- ✅ Event-driven design with Kafka
- ✅ Golang backend development
- ✅ Node.js async programming
- ✅ MongoDB document storage
- ✅ Redis caching
- ✅ Elasticsearch indexing
- ✅ Docker containerization
- ✅ Concurrent programming (goroutines, async/await)
- ✅ Clean code principles (DRY, SOLID)
- ✅ Strategy pattern implementation

---

## 📝 Current Status

- ✅ **API Gateway** - Complete
- ✅ **Notification Processor** - Complete
- ✅ **Delivery Service** - Complete (Email, SMS, Push)
- 🔄 **Analytics Service** - In Progress
- ⏳ **React Dashboard** - To Do

---

## 🤝 Contributing

This is a personal learning project. Feel free to fork and experiment!

---

## 📄 License

MIT License - Free to use for learning and portfolio purposes.

---

## 👤 Author

**Jonah Tornovsky**

---

**Last Updated:** November 16, 2025