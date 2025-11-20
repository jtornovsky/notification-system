# Delivery Service

## 📌 Overview

The Delivery Service simulates notification delivery across multiple channels (Email, SMS, Push). It consumes type-specific notifications from Kafka, simulates realistic delivery with random failures, persists results to MongoDB, and publishes delivery events for analytics.

**Service Type:** Multi-handler consumer/processor  
**Language:** Golang  
**Primary Function:** Notification delivery simulation and tracking

---

## 🏗️ Architecture
```
Kafka Topics (Input):
├─ email-notifications
├─ sms-notifications  
└─ push-notifications
         ↓
┌──────────────────────────┐
│   Delivery Service       │
│   ┌────────────────┐     │
│   │ Email Handler  │     │
│   │ SMS Handler    │     │
│   │ Push Handler   │     │
│   └────────────────┘     │
└──────────┬───────────────┘
           │
           ├→ MongoDB (delivery_results)
           │  - Persistent delivery records
           │
           └→ Kafka (delivery-events)
              - Analytics events
```

**Position in System:**
- Consumes from: Type-specific topics
- Writes to: MongoDB
- Publishes to: `delivery-events` topic
- 3 concurrent handlers running in parallel

---

## 🛠️ Tech Stack

- **Language:** Go 1.21+
- **Kafka Client:** segmentio/kafka-go
- **Database:** MongoDB (go.mongodb.org/mongo-driver)
- **Concurrency:** Goroutines (3 concurrent handlers)
- **Architecture Pattern:** Strategy Pattern for delivery simulation

---

## 🚀 Quick Start

### **Prerequisites**
- Go 1.21 or higher
- Kafka running on `localhost:9092`
- MongoDB running on `localhost:27017` (with auth: `admin:password`)
- Topics must exist:
    - `email-notifications` (input)
    - `sms-notifications` (input)
    - `push-notifications` (input)
    - `delivery-events` (output)

### **Installation**
```bash
cd delivery-service

# Install dependencies
go mod download

# Run the service
go run main.go
```

**Expected output:**
```
🚀 Starting Delivery Service...
✓ Connected to MongoDB
✓ Email handler initialized for topic: email-notifications
✓ SMS handler initialized for topic: sms-notifications
✓ Push handler initialized for topic: push-notifications
✅ All handlers started successfully
🚀 [Email] Handler started, waiting for messages...
🚀 [SMS] Handler started, waiting for messages...
🚀 [Push] Handler started, waiting for messages...
```

---

## 📊 Delivery Handlers

### **Email Handler**

**Consumes from:** `email-notifications`  
**Delivery simulation:**
- Delay: 50-200ms (network latency)
- Failure rate: 10% random
- Error: "SMTP connection timeout"

**Example log:**
```
📩 [Email] Processing notification: abc-123 for test@example.com
📧 Email sent to test@example.com (took 125ms)
✓ [Email] Saved delivery result to MongoDB
✓ [Email] Published delivery event to Kafka
```

---

### **SMS Handler**

**Consumes from:** `sms-notifications`  
**Delivery simulation:**
- Delay: 30-100ms (faster than email)
- Failure rate: 10% random
- Error: "carrier gateway unreachable"

**Example log:**
```
📩 [SMS] Processing notification: def-456 for +1234567890
📱 SMS sent to +1234567890 (took 65ms)
✓ [SMS] Saved delivery result to MongoDB
✓ [SMS] Published delivery event to Kafka
```

---

### **Push Handler**

**Consumes from:** `push-notifications`  
**Delivery simulation:**
- Delay: 20-80ms (fastest)
- Failure rate: 10% random
- Error: "device token invalid"

**Example log:**
```
📩 [Push] Processing notification: ghi-789 for device-token-123
🔔 Push sent to device-token-123 (took 42ms)
✓ [Push] Saved delivery result to MongoDB
✓ [Push] Published delivery event to Kafka
```

---

## 🗄️ Data Storage

### **MongoDB - Delivery Results**

**Database:** `notifications`  
**Collection:** `delivery_results`

**Document Schema:**
```javascript
{
  _id: ObjectId("..."),
  notification_id: "uuid",              // Original notification ID
  type: "EMAIL | SMS | PUSH",           // Delivery channel
  recipient: "string",                  // Email/phone/device token
  status: "SENT | FAILED",              // Delivery outcome
  timestamp: ISODate("..."),            // When delivery was attempted
  delivery_time_ms: Long(125),          // How long it took
  error_message: "string (optional)"    // Error if failed
}
```

**Indexes (recommended):**
```javascript
db.delivery_results.createIndex({ notification_id: 1 })
db.delivery_results.createIndex({ status: 1 })
db.delivery_results.createIndex({ type: 1 })
db.delivery_results.createIndex({ timestamp: -1 })
```

---

### **Kafka - Delivery Events**

**Topic:** `delivery-events`  
**Format:** JSON  
**Purpose:** Feed analytics service

**Message Structure:**
```json
{
  "notification_id": "uuid",
  "type": "EMAIL | SMS | PUSH",
  "recipient": "string",
  "status": "SENT | FAILED",
  "timestamp": "ISO8601",
  "delivery_time_ms": 125,
  "error_message": "string (optional)"
}
```

---

## 🧪 Testing

### **Test 1: Single Delivery (Email)**

**Send notification:**
```powershell
Invoke-WebRequest -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body '{"user_id":"test","type":"EMAIL","recipient":"test@example.com","subject":"Test","message":"Hello"}'
```

**Watch Delivery Service logs:**
```
📩 [Email] Processing notification: abc-123 for test@example.com
📧 Email sent to test@example.com (took 158ms)
✓ [Email] Saved delivery result to MongoDB
✓ [Email] Published delivery event to Kafka
```

**Verify in MongoDB:**
```javascript
use notifications
db.delivery_results.find().sort({timestamp: -1}).limit(1).pretty()
```

---

### **Test 2: All Three Types**
```powershell
# Email
iwr -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body '{"user_id":"u1","type":"EMAIL","recipient":"test@example.com","subject":"Test","message":"Email"}'

# SMS
iwr -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body '{"user_id":"u2","type":"SMS","recipient":"+1234567890","message":"SMS"}'

# Push
iwr -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body '{"user_id":"u3","type":"PUSH","recipient":"device-token-123","message":"Push"}'
```

**Verify all types in MongoDB:**
```javascript
db.delivery_results.aggregate([
  { $group: { _id: "$type", count: { $sum: 1 } } }
])
```

---

### **Test 3: Verify Failures**

**Send 20 notifications to see random failures:**
```powershell
1..20 | ForEach-Object {
    iwr -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body "{`"user_id`":`"user-$_`",`"type`":`"EMAIL`",`"recipient`":`"test$_@example.com`",`"subject`":`"Test`",`"message`":`"Message $_`"}"
}
```

**Check failure rate:**
```javascript
db.delivery_results.countDocuments({status: "SENT"})
db.delivery_results.countDocuments({status: "FAILED"})

// Should be roughly 90% success, 10% failure
```

---

### **Test 4: Performance (Concurrent Processing)**

**Send 50 notifications rapidly:**
```powershell
1..50 | ForEach-Object {
    Start-Job -ScriptBlock {
        param($i)
        Invoke-WebRequest -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body "{`"user_id`":`"user-$i`",`"type`":`"EMAIL`",`"recipient`":`"test$i@example.com`",`"message`":`"Test`"}"
    } -ArgumentList $_
}

Get-Job | Wait-Job | Remove-Job
```

**All 3 handlers process concurrently - should complete quickly**

---

## 🔧 Configuration

**MongoDB Connection (`main.go`):**
```go
mongoURI := "mongodb://admin:password@localhost:27017"
mongoDatabase := "notifications"
mongoCollection := "delivery_results"
```

**Kafka Brokers:**
```go
kafkaBrokers := []string{"localhost:9092"}
```

**Delivery Configuration (`internal/handlers/simulators.go`):**
```go
var emailConfig = DeliveryConfig{
    Name:         "Email",
    MinDelayMs:   50,      // Minimum delivery time
    MaxDelayMs:   200,     // Maximum delivery time
    FailureRate:  0.10,    // 10% failure rate
    ErrorMessage: "SMTP connection timeout",
    SuccessEmoji: "📧",
}
```

**To change failure rate:**
```go
FailureRate: 0.20,  // 20% failure rate
```

---

## 📂 Project Structure
```
delivery-service/
├── main.go                      # Application entry point
├── internal/
│   ├── handlers/
│   │   ├── handler.go               # Generic handler (Strategy Pattern)
│   │   └── simulators.go            # Delivery simulation functions
│   ├── models/
│   │   └── notification.go          # Data models
│   └── mongo/
│       └── client.go                # MongoDB client
├── go.mod                           # Go module definition
├── go.sum                           # Dependency checksums
└── README.md                        # This file
```

---

## 🎨 Design Patterns

### **Strategy Pattern**

**Context:** `Handler` struct  
**Strategy:** `DeliverySimulator` function type  
**Concrete Strategies:** `SimulateEmailDelivery`, `SimulateSmsDelivery`, `SimulatePushDelivery`

**Implementation:**
```go
// Strategy interface (function type)
type DeliverySimulator func(models.Notification) (string, int64, error)

// Context holds strategy
type Handler struct {
    simulator DeliverySimulator  // Strategy injected here
}

// Context uses strategy
func (h *Handler) processMessage(...) {
    status, time, err := h.simulator(notification)  // Execute strategy
}

// Client selects strategy
emailHandler := NewHandler(..., SimulateEmailDelivery)  // Email strategy
smsHandler := NewHandler(..., SimulateSmsDelivery)      // SMS strategy
```

**Benefits:**
- Single handler implementation for all types
- Easy to add new delivery types
- Configuration-driven behavior

---

### **Configuration-Driven Simulation**

**All delivery types use same core logic with different config:**
```go
type DeliveryConfig struct {
    Name         string
    MinDelayMs   int
    MaxDelayMs   int
    FailureRate  float32
    ErrorMessage string
    SuccessEmoji string
}

func simulateDelivery(notification, config) {
    delay := random(config.MinDelayMs, config.MaxDelayMs)
    time.Sleep(delay)
    
    if random() < config.FailureRate {
        return "FAILED", delay, error(config.ErrorMessage)
    }
    
    return "SENT", delay, nil
}
```

---

## 🐛 Troubleshooting

### **Error: "failed to connect to MongoDB"**

**Check MongoDB:**
```bash
docker-compose ps
docker-compose logs notification_mongodb
```

**Test connection:**
```bash
docker exec -it notification_mongodb mongosh -u admin -p password
```

**Verify credentials in code match docker-compose.yml**

---

### **Error: "Unknown Topic Or Partition"**

**Create missing topics:**
```bash
docker exec -it notification_kafka kafka-topics --create --bootstrap-server localhost:9092 --topic email-notifications --partitions 1 --replication-factor 1

docker exec -it notification_kafka kafka-topics --create --bootstrap-server localhost:9092 --topic sms-notifications --partitions 1 --replication-factor 1

docker exec -it notification_kafka kafka-topics --create --bootstrap-server localhost:9092 --topic push-notifications --partitions 1 --replication-factor 1

docker exec -it notification_kafka kafka-topics --create --bootstrap-server localhost:9092 --topic delivery-events --partitions 1 --replication-factor 1
```

---

### **No messages being processed**

**Check upstream services:**
1. API Gateway running?
2. Notification Processor running?
3. Messages in input topics?

**Verify messages exist:**
```bash
docker exec -it notification_kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic email-notifications \
  --from-beginning
```

---

### **Handlers not starting**

**Check logs for:**
- MongoDB connection errors
- Kafka connection errors
- Port conflicts

**Verify all dependencies running:**
```bash
docker-compose ps
```

---

## 📈 Performance Characteristics

**Throughput:**
- Email: ~8-15 msg/sec per handler (limited by delay simulation)
- SMS: ~10-30 msg/sec per handler (faster delays)
- Push: ~12-50 msg/sec per handler (fastest delays)

**Concurrency:**
- 3 handlers run in parallel (goroutines)
- Each handler processes 1 message at a time
- Total: ~30-95 msg/sec across all types

**Scalability:**
- Stateless (can run multiple instances, no in-memory state between requests)
- Kafka consumer group distributes load
- Limited by MongoDB write throughput

---

## 🔍 Monitoring

### **Check Handler Status**
```bash
# Consumer group lag
docker exec -it notification_kafka kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --describe \
  --group delivery-service-email

docker exec -it notification_kafka kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --describe \
  --group delivery-service-sms

docker exec -it notification_kafka kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --describe \
  --group delivery-service-push
```

### **MongoDB Metrics**
```javascript
// Total deliveries
db.delivery_results.countDocuments()

// By status
db.delivery_results.aggregate([
  { $group: { _id: "$status", count: { $sum: 1 } } }
])

// By type
db.delivery_results.aggregate([
  { $group: { _id: "$type", count: { $sum: 1 } } }
])

// Average delivery time by type
db.delivery_results.aggregate([
  { 
    $group: { 
      _id: "$type", 
      avg_time: { $avg: "$delivery_time_ms" } 
    } 
  }
])

// Recent failures
db.delivery_results.find({status: "FAILED"}).sort({timestamp: -1}).limit(10)
```

---

## 🚀 Future Enhancements

- [ ] Retry mechanism for failed deliveries
- [ ] Exponential backoff for transient failures
- [ ] Batch processing for efficiency
- [ ] Real email/SMS/push integration (replace simulation)
- [ ] Metrics endpoint (Prometheus)
- [ ] Health check endpoint
- [ ] Dead letter queue for permanent failures
- [ ] Rate limiting per channel
- [ ] Priority queues (urgent vs normal)
- [ ] Delivery scheduling (send later)

---

## 📝 Key Features

### **Realistic Simulation**

- ✅ Variable delivery times (mimics real network latency)
- ✅ Random failures (mimics real-world issues)
- ✅ Different characteristics per channel
- ✅ Proper error messages

### **Production-Ready Patterns**

- ✅ Graceful shutdown (WaitGroup + Context)
- ✅ Error handling and logging
- ✅ Structured concurrency (goroutines)
- ✅ Clean architecture (Strategy Pattern)
- ✅ Persistent storage (MongoDB)
- ✅ Event publishing (Kafka)

### **Observability**

- ✅ Detailed logging per handler
- ✅ MongoDB for historical queries
- ✅ Kafka events for real-time analytics
- ✅ Consumer group monitoring

---

**Last Updated:** November 17, 2025