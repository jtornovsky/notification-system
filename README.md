# Real-Time Notification System

A full-stack microservices-based notification delivery system with real-time analytics and interactive dashboard, built as a portfolio project to demonstrate modern backend architecture, event-driven design, and full-stack development skills.

## 🎯 Project Overview

This system provides a scalable notification delivery platform supporting multiple channels (Email, SMS, Push notifications) with real-time tracking, delivery simulation, analytics aggregation, and a React-based web dashboard.

**Key Features:**
- Multi-channel notification delivery (Email, SMS, Push)
- Real-time message routing and processing
- Delivery status tracking with MongoDB persistence
- Analytics aggregation with Elasticsearch
- Interactive React dashboard for submission and analytics
- Smart type auto-detection (Email/SMS/Push)
- Failure simulation for realistic testing (10% random failure rate)
- Microservices architecture with event-driven communication
- CORS-enabled API for frontend integration

---

## 🏗️ System Architecture
```
┌──────────────────┐
│ React Dashboard  │  (Port 5173)
│  - Submit Form   │
│  - Analytics UI  │
└────────┬─────────┘
         │ HTTP (CORS)
         ▼
┌─────────────────┐
│  API Gateway    │  (Port 8080 - Golang)
│  - REST API     │
│  - Redis Cache  │
│  - CORS Proxy   │
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
└──────────────────────┘
```

---

## 🛠️ Tech Stack

### **Backend Services**
- **Golang** (1.21+) - API Gateway, Delivery Service
- **Node.js** (18+) - Notification Processor, Analytics Service
- **Gin** - Web framework for API Gateway

### **Infrastructure**
- **Apache Kafka** - Event streaming and message queuing
- **MongoDB** - Delivery results persistence
- **Redis** - Request caching (1-hour TTL)
- **Elasticsearch** - Analytics data indexing and aggregation

### **Frontend**
- **React** (18.3) - Dashboard UI
- **Vite** (7.2) - Build tool and dev server

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
- **Kafka + Zookeeper** (Port 9092)
- **MongoDB** (Port 27017)
- **Redis** (Port 6379)
- **Elasticsearch** (Port 9200)

### **2. Run Backend Services**

**Terminal 1 - API Gateway:**
```bash
cd api-gateway
go run main.go
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
go run main.go
```

**Terminal 4 - Analytics Service:**
```bash
cd analytics-service
npm install
node index.js
```

### **3. Run Frontend Dashboard**

**Terminal 5 - React Dashboard:**
```bash
cd dashboard
npm install
npm run dev
```

**Open browser:** `http://localhost:5173`

### **4. Send Test Notification**

**Option 1: Via Dashboard**
- Open `http://localhost:5173`
- Fill out the form
- Click "Send Notification"

**Option 2: Via Command Line**

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

**Response (201 Created):**
```json
{
  "id": "20251120123456",
  "user_id": "test-user",
  "type": "EMAIL",
  "recipient": "test@example.com",
  "subject": "Hello",
  "message": "Test notification",
  "created_at": "2025-11-20T12:34:56Z"
}
```

---

### **GET /notifications/:id**

Retrieve a cached notification (1-hour TTL).

**Endpoint:** `http://localhost:8080/notifications/{id}`

**Response (200 OK):**
```json
{
  "id": "20251120123456",
  "user_id": "test-user",
  "type": "EMAIL",
  "recipient": "test@example.com",
  "subject": "Hello",
  "message": "Test notification",
  "created_at": "2025-11-20T12:34:56Z"
}
```

---

### **GET /analytics**

Fetch analytics from Elasticsearch (proxied through API Gateway for CORS).

**Endpoint:** `http://localhost:8080/analytics`

**Response (200 OK):**
```json
{
  "took": 5,
  "hits": { "total": { "value": 100 } },
  "aggregations": {
    "by_status": {
      "buckets": [
        { "key": "SENT", "doc_count": 90 },
        { "key": "FAILED", "doc_count": 10 }
      ]
    },
    "by_type": {
      "buckets": [
        { "key": "EMAIL", "doc_count": 60 },
        { "key": "SMS", "doc_count": 25 },
        { "key": "PUSH", "doc_count": 15 }
      ]
    }
  }
}
```

---

### **GET /health**

Check API Gateway health.

**Endpoint:** `http://localhost:8080/health`

**Response (200 OK):**
```json
{
  "status": "ok"
}
```

---

## 📊 Services Overview

### **1. API Gateway** (Port 8080)
- **Technology:** Golang + Gin framework
- **Purpose:** HTTP REST API, caching, Kafka publishing, Elasticsearch proxy
- **Features:**
    - Accepts notification requests via REST API
    - Caches notifications in Redis (1 hour TTL)
    - Publishes to Kafka `notifications` topic
    - CORS-enabled for frontend access
    - Proxies analytics requests to Elasticsearch
- **[Full Documentation →](api-gateway/README.md)**

---

### **2. Notification Processor**
- **Technology:** Node.js + KafkaJS
- **Purpose:** Route notifications by type
- **Features:**
    - Consumes from `notifications` topic
    - Routes to type-specific topics:
        - `email-notifications`
        - `sms-notifications`
        - `push-notifications`
    - Pure routing logic (no transformation)
- **[Full Documentation →](notification-processor/README.md)**

---

### **3. Delivery Service**
- **Technology:** Golang + Kafka + MongoDB
- **Purpose:** Simulate delivery, persist results
- **Features:**
    - 3 concurrent handlers (Email, SMS, Push)
    - Simulates delivery with realistic delays:
        - Email: 50-200ms
        - SMS: 30-100ms
        - Push: 20-80ms
    - 10% random failure rate for testing
    - Saves results to MongoDB
    - Publishes delivery events to `delivery-events` topic
    - Strategy pattern implementation
- **[Full Documentation →](delivery-service/README.md)**

---

### **4. Analytics Service**
- **Technology:** Node.js + Elasticsearch
- **Purpose:** Aggregate metrics, index to Elasticsearch
- **Features:**
    - Consumes from `delivery-events` topic
    - Indexes delivery results to Elasticsearch
    - Enables aggregations by status and type
    - Real-time analytics updates
- **[Full Documentation →](analytics-service/README.md)**

---

### **5. React Dashboard** (Port 5173)
- **Technology:** React 18 + Vite 7
- **Purpose:** Web UI for notification submission and analytics
- **Features:**
    - **Smart notification form:**
        - Auto-detects type from recipient format
        - Email: Contains `@` and `.`
        - SMS: 10-15 digits with optional `+()- `
        - Push: Everything else
    - **Real-time analytics display:**
        - Auto-loads on page mount
        - Manual refresh button
        - Shows aggregated metrics
    - **Responsive two-column layout**
    - **Form validation and error handling**
- **[Full Documentation →](dashboard/README.md)**

---

## 🗄️ Data Storage

### **MongoDB - Delivery Results**
**Database:** `notifications`  
**Collection:** `delivery_results`

**Document Schema:**
```javascript
{
  _id: ObjectId("..."),
  notification_id: "20251120123456",
  type: "EMAIL | SMS | PUSH",
  recipient: "string",
  status: "SENT | FAILED",
  timestamp: ISODate("2025-11-20T12:00:00Z"),
  delivery_time_ms: 125,
  error_message: "string (if failed)"
}
```

**Indexes:**
- `notification_id` (unique)
- `type`
- `status`
- `timestamp`

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

---

### **Redis - Notification Cache**
**Key Pattern:** `notification:{id}`  
**TTL:** 1 hour (3600 seconds)  
**Value:** JSON string of notification

**Example:**
```
Key: notification:20251120123456
Value: {"id":"20251120123456","user_id":"test",...}
Expiry: 3600 seconds
```

---

### **Elasticsearch - Analytics Index**
**Index Name:** `notification-analytics`  
**Purpose:** Aggregated metrics for dashboard

**Document Schema:**
```json
{
  "notification_id": "20251120123456",
  "type": "EMAIL",
  "status": "SENT",
  "timestamp": "2025-11-20T12:00:00Z",
  "delivery_time_ms": 125
}
```

**Common Queries:**
```bash
# Count documents
curl http://localhost:9200/notification-analytics/_count

# Aggregate by status
curl -X POST http://localhost:9200/notification-analytics/_search \
  -H "Content-Type: application/json" \
  -d '{"size":0,"aggs":{"by_status":{"terms":{"field":"status"}}}}'

# Aggregate by type
curl -X POST http://localhost:9200/notification-analytics/_search \
  -H "Content-Type: application/json" \
  -d '{"size":0,"aggs":{"by_type":{"terms":{"field":"type"}}}}'
```

---

## 🧪 Testing

### **End-to-End Test via Dashboard**

1. **Open Dashboard:** `http://localhost:5173`
2. **Submit Notification:**
    - User ID: `test123`
    - Recipient: `test@example.com` (auto-detects EMAIL)
    - Subject: `Test Subject`
    - Message: `Hello World`
3. **Click "Send Notification"**
4. **Verify Success Alert**
5. **Check Analytics:** Click "Refresh Analytics"

---

### **End-to-End Test via Command Line**

**1. Send 10 notifications:**
```powershell
1..10 | ForEach-Object {
    Invoke-WebRequest -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body "{`"user_id`":`"user-$_`",`"type`":`"EMAIL`",`"recipient`":`"test$_@example.com`",`"subject`":`"Test`",`"message`":`"Message $_`"}"
}
```

**2. Check MongoDB results:**
```bash
docker exec -it <mongodb-container> mongosh -u admin -p password

use notifications
db.delivery_results.countDocuments()
db.delivery_results.aggregate([
  { $group: { _id: "$status", count: { $sum: 1 } } }
])
```

**3. Check Elasticsearch:**
```bash
curl http://localhost:9200/notification-analytics/_count
```

**4. View in Dashboard:**
- Open `http://localhost:5173`
- Click "Refresh Analytics"
- See updated counts

---

### **Test Different Notification Types**

**Via Dashboard:**
- Type `test@example.com` → Auto-detects EMAIL
- Type `+1 (555) 123-4567` → Auto-detects SMS
- Type `device_token_abc` → Auto-detects PUSH

**Via Command Line:**
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

### **Service Health Checks**
- **API Gateway:** `http://localhost:8080/health`
- **React Dashboard:** `http://localhost:5173`
- **Elasticsearch:** `curl http://localhost:9200/_cluster/health`

### **Infrastructure Status**
```bash
# Check all services
docker-compose ps

# View Kafka topics
docker exec -it <kafka-container> kafka-topics --list --bootstrap-server localhost:9092

# Check MongoDB
docker exec -it <mongodb-container> mongosh -u admin -p password

# Check Redis
docker exec -it <redis-container> redis-cli
```

### **View Logs**
```bash
# Infrastructure logs
docker-compose logs -f kafka
docker-compose logs -f mongodb
docker-compose logs -f elasticsearch

# Service logs (when running locally)
# Check terminal windows where services are running
```

---

## 🏗️ Project Structure
```
notification-system/
├── api-gateway/              # Golang REST API + CORS
│   ├── main.go
│   ├── go.mod
│   └── README.md
├── notification-processor/   # Node.js routing service
│   ├── index.js
│   ├── package.json
│   └── README.md
├── delivery-service/         # Golang delivery handlers
│   ├── main.go
│   ├── internal/
│   │   ├── handlers/
│   │   ├── models/
│   │   └── mongo/
│   ├── go.mod
│   └── README.md
├── analytics-service/        # Node.js + Elasticsearch
│   ├── index.js
│   ├── package.json
│   └── README.md
├── dashboard/                # React + Vite frontend
│   ├── src/
│   │   ├── App.jsx
│   │   ├── App.css
│   │   └── main.jsx
│   ├── package.json
│   └── README.md
├── docker-compose.yml        # Infrastructure setup
└── README.md                 # This file
```

---

## 🎓 Learning Objectives

This project demonstrates:

### **Backend Development**
- ✅ Microservices architecture design
- ✅ Event-driven design with Apache Kafka
- ✅ Golang backend development (Gin framework)
- ✅ Node.js async programming (async/await)
- ✅ RESTful API design
- ✅ CORS configuration for frontend integration

### **Data Management**
- ✅ MongoDB document storage and indexing
- ✅ Redis caching strategies
- ✅ Elasticsearch indexing and aggregations
- ✅ Data consistency across services

### **Frontend Development**
- ✅ React 18 with Hooks (useState, useEffect)
- ✅ Vite build tooling
- ✅ Controlled components and form handling
- ✅ HTTP client integration (fetch API)
- ✅ Responsive CSS (Flexbox/Grid)

### **DevOps & Architecture**
- ✅ Docker containerization
- ✅ Concurrent programming (goroutines, async handlers)
- ✅ Message queue integration
- ✅ Service orchestration

### **Software Engineering**
- ✅ Clean code principles (DRY, SOLID)
- ✅ Strategy pattern implementation
- ✅ Error handling and logging
- ✅ API documentation

---

## 🔧 Configuration Reference

### **Ports**
| Service | Port |
|---------|------|
| API Gateway | 8080 |
| React Dashboard | 5173 |
| Kafka | 9092 |
| MongoDB | 27017 |
| Redis | 6379 |
| Elasticsearch | 9200 |

### **Kafka Topics**
| Topic | Producer | Consumer |
|-------|----------|----------|
| `notifications` | API Gateway | Notification Processor |
| `email-notifications` | Notification Processor | Delivery Service (Email) |
| `sms-notifications` | Notification Processor | Delivery Service (SMS) |
| `push-notifications` | Notification Processor | Delivery Service (Push) |
| `delivery-events` | Delivery Service | Analytics Service |

### **Consumer Groups**
| Service | Consumer Group |
|---------|----------------|
| Notification Processor | `notification-processor-group` |
| Delivery Service (Email) | `delivery-service-email` |
| Delivery Service (SMS) | `delivery-service-sms` |
| Delivery Service (Push) | `delivery-service-push` |
| Analytics Service | `analytics-service-group` |

---

## 📝 Current Status

### ✅ **Completed Features**
- ✅ **API Gateway** - REST API with CORS, Redis caching, Kafka publishing
- ✅ **Notification Processor** - Type-based routing
- ✅ **Delivery Service** - Email, SMS, Push handlers with MongoDB persistence
- ✅ **Analytics Service** - Elasticsearch indexing and aggregations
- ✅ **React Dashboard** - Form submission and analytics display
- ✅ **End-to-End Pipeline** - Full flow from submission to analytics

### 🚀 **Future Enhancements**

#### **Phase 1: Enhanced Dashboard**
- [ ] Add pie chart for notification types
- [ ] Add bar chart for success/failure rates
- [ ] Show total notification count prominently
- [ ] Toast notifications instead of alerts

#### **Phase 2: Advanced Features**
- [ ] Notification history table with search/filter
- [ ] Real-time updates (WebSocket integration)
- [ ] Retry logic for failed deliveries
- [ ] Dead Letter Queue (DLQ) for invalid messages
- [ ] User authentication and authorization
- [ ] Rate limiting per user
- [ ] Metrics endpoints (Prometheus format)

#### **Phase 3: Production Readiness**
- [ ] Environment-based configuration
- [ ] Health checks for all dependencies
- [ ] Structured logging (JSON format)
- [ ] OpenAPI/Swagger documentation
- [ ] Kubernetes deployment manifests
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Integration tests
- [ ] Performance benchmarks

---

## 🐛 Troubleshooting

### **Services Won't Start**
```bash
# Check if ports are available
netstat -ano | findstr "8080 9092 27017 6379 9200 5173"

# Restart Docker Compose
docker-compose down
docker-compose up -d

# Check logs
docker-compose logs
```

### **CORS Errors in Dashboard**
- Verify API Gateway is running on port 8080
- Check CORS middleware is enabled in `api-gateway/cmd/main.go`
- Ensure frontend calls `http://localhost:8080` (not 9200 directly)

### **No Analytics Data**
- Verify all services are running
- Send a few test notifications
- Wait 5-10 seconds for processing
- Click "Refresh Analytics" in dashboard
- Check Elasticsearch: `curl http://localhost:9200/notification-analytics/_count`

### **Type Not Auto-Detecting**
- Email must have `@` and `.` with `@` before last `.`
- SMS must be 10-15 digits (can include `+()- ` characters)
- Clear and retype if detection is wrong

---

## 🤝 Contributing

This is a personal learning project built as a portfolio piece. Feel free to fork and experiment!

**If you find this project helpful:**
- ⭐ Star the repository
- 🐛 Report issues
- 💡 Suggest improvements
- 🔀 Submit pull requests

---

## 📄 License

MIT License - Free to use for learning and portfolio purposes.

---

## 👤 Author

**Jonah Tornovsky**  
Backend Developer

**Portfolio Project:** Real-Time Notification System  
**Technologies:** Golang, Node.js, React, Kafka, MongoDB, Redis, Elasticsearch, Docker

---

**Last Updated:** November 20, 2025  
**Project Status:** ✅ Core Features Complete | 🚀 Enhancements in Progress