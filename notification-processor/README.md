# Notification Processor Service

## 📌 Overview

The Notification Processor is a lightweight routing service that consumes notifications from Kafka and routes them to type-specific topics based on notification type (EMAIL, SMS, PUSH).

**Service Type:** Kafka Consumer/Producer (Router)  
**Language:** Node.js  
**Primary Function:** Message routing by type

---

## 🏗️ Architecture
```
Kafka: notifications topic
         ↓
┌──────────────────────┐
│ Notification         │
│ Processor            │
│ (Node.js)            │
└──────────┬───────────┘
           │
           ├→ email-notifications (EMAIL)
           ├→ sms-notifications (SMS)
           └→ push-notifications (PUSH)
```

**Position in System:**
- Consumes from: `notifications` topic
- Routes to: Type-specific topics
- Pure routing logic (no business logic)
- Stateless (no database)

---

## 🛠️ Tech Stack

- **Runtime:** Node.js 18+
- **Kafka Client:** kafkajs
- **Architecture Pattern:** Consumer/Producer (Router)

---

## 🚀 Quick Start

### **Prerequisites**
- Node.js 18 or higher
- Kafka running on `localhost:9092`
- Topics must exist:
    - `notifications` (input)
    - `email-notifications` (output)
    - `sms-notifications` (output)
    - `push-notifications` (output)

### **Installation**
```bash
cd notification-processor

# Install dependencies
npm install

# Run the service
node index.js
```

**Expected output:**
```
🚀 Starting Notification Processor...
✓ Connected to Kafka
✓ Subscribed to notifications topic
✅ Notification Processor is running
📨 Waiting for notifications...
```

---

## 📊 Routing Logic

### **Input Topic:** `notifications`

**Message Format:**
```json
{
  "id": "uuid",
  "user_id": "string",
  "type": "EMAIL | SMS | PUSH",
  "recipient": "string",
  "subject": "string (optional)",
  "message": "string",
  "timestamp": "ISO8601"
}
```

### **Output Topics:**

| **Input Type** | **Output Topic** | **Use Case** |
|----------------|------------------|--------------|
| `EMAIL` | `email-notifications` | Email delivery |
| `SMS` | `sms-notifications` | SMS delivery |
| `PUSH` | `push-notifications` | Push notifications |

### **Routing Decision:**
```javascript
switch (notification.type) {
  case 'EMAIL':
    targetTopic = 'email-notifications';
    break;
  case 'SMS':
    targetTopic = 'sms-notifications';
    break;
  case 'PUSH':
    targetTopic = 'push-notifications';
    break;
  default:
    log.error('Unknown notification type');
}
```

---

## 🔄 Message Flow
```
1. Consume from 'notifications' topic
         ↓
2. Parse JSON message
         ↓
3. Extract notification.type
         ↓
4. Determine target topic based on type
         ↓
5. Publish to target topic
         ↓
6. Log routing action
         ↓
7. Commit offset (mark as processed)
```

---

## 🧪 Testing

### **Test 1: Verify Service is Running**

**Start the service:**
```bash
node index.js
```

**You should see:**
```
🚀 Starting Notification Processor...
✓ Connected to Kafka
✓ Subscribed to notifications topic
✅ Notification Processor is running
📨 Waiting for notifications...
```

---

### **Test 2: Send Test Notification**

**In another terminal, send a notification via API Gateway:**
```powershell
Invoke-WebRequest -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body '{"user_id":"test","type":"EMAIL","recipient":"test@example.com","subject":"Test","message":"Hello"}'
```

**Processor logs should show:**
```
📝 Parsed notification: {
  id: 'abc-123',
  type: 'EMAIL',
  recipient: 'test@example.com'
}
📨 Routing EMAIL notification to email-notifications
✓ Successfully routed notification abc-123
```

---

### **Test 3: Verify Messages in Output Topics**

**Check email-notifications topic:**
```bash
docker exec -it notification_kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic email-notifications \
  --from-beginning
```

**You should see the routed message in JSON format**

---

### **Test 4: Test All Types**
```powershell
# Send EMAIL
iwr -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body '{"user_id":"u1","type":"EMAIL","recipient":"test@example.com","subject":"Test","message":"Email"}'

# Send SMS
iwr -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body '{"user_id":"u2","type":"SMS","recipient":"+1234567890","message":"SMS test"}'

# Send PUSH
iwr -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body '{"user_id":"u3","type":"PUSH","recipient":"device-token-123","message":"Push test"}'
```

**Processor should route each to correct topic:**
```
📨 Routing EMAIL notification to email-notifications
📨 Routing SMS notification to sms-notifications
📨 Routing PUSH notification to push-notifications
```

---

## 🔧 Configuration

**All configuration is in `index.js`:**
```javascript
const kafka = new Kafka({
    clientId: 'notification-processor',
    brokers: ['localhost:9092']
});

const consumer = kafka.consumer({ 
    groupId: 'notification-processor-group' 
});

const producer = kafka.producer();
```

**Key Settings:**
- **Client ID:** `notification-processor`
- **Consumer Group:** `notification-processor-group`
- **Input Topic:** `notifications`
- **Output Topics:** Dynamically determined by type

---

## 📂 Project Structure
```
notification-processor/
├── package.json             # Dependencies
├── index.js                 # Main application (single file)
└── README.md               # This file
```

**Single-file architecture:**
- Simple routing logic (~150 lines)
- No need for complex structure
- Easy to understand and maintain

---

## 🐛 Troubleshooting

### **Error: "KafkaJSNumberOfRetriesExceeded"**
- Kafka broker not reachable
- Check: `docker-compose ps`
- Verify: `docker-compose logs notification_kafka`

### **Error: "This server does not host this topic-partition"**
- Topic doesn't exist
- Create topics:
```bash
docker exec -it notification_kafka kafka-topics --create --bootstrap-server localhost:9092 --topic email-notifications --partitions 1 --replication-factor 1
docker exec -it notification_kafka kafka-topics --create --bootstrap-server localhost:9092 --topic sms-notifications --partitions 1 --replication-factor 1
docker exec -it notification_kafka kafka-topics --create --bootstrap-server localhost:9092 --topic push-notifications --partitions 1 --replication-factor 1
```

### **No messages being processed**
- Check API Gateway is running
- Verify messages in input topic:
```bash
docker exec -it notification_kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic notifications \
  --from-beginning
```

---

## 🔍 Monitoring

### **Check Consumer Group Status**
```bash
docker exec -it notification_kafka kafka-consumer-groups --bootstrap-server localhost:9092 --describe --group notification-processor-group
```

**Output shows:**
- Current offset
- Log end offset
- Lag (how many messages behind)

### **Verify Topic Messages**
```bash
# Count messages in input topic
docker exec -it notification_kafka kafka-run-class kafka.tools.GetOffsetShell --broker-list localhost:9092 --topic notifications

# Count messages in output topics
docker exec -it notification_kafka kafka-run-class kafka.tools.GetOffsetShell --broker-list localhost:9092 --topic email-notifications
```

---

## 📈 Performance Characteristics

**Throughput:**
- Processes messages as fast as they arrive
- No blocking operations (pure routing)
- Typically <10ms per message

**Scalability:**
- Stateless (can run multiple instances)
- Messages distributed via consumer group
- Limited by Kafka partition count

**Reliability:**
- Automatic reconnection to Kafka
- Offset commits after successful routing
- Failed messages logged but not retried (by design)

---

## 🚀 Future Enhancements

- [ ] Add message validation (schema checking)
- [ ] Dead letter queue for invalid types (separate topic for "bad" messages that can't be processed)
- [ ] Metrics (messages routed per type)
- [ ] Health check endpoint
- [ ] Configuration via environment variables
- [ ] Structured logging

---

## 📝 Design Decisions

### **Why Single File?**
- Simple routing logic
- No complex business rules
- Easy to read and maintain
- ~150 lines total

### **Why No Database?**
- Stateless routing only
- No data persistence needed
- Pure message transformation

### **Why Separate Topics by Type?**
- Different delivery services can scale independently
- Clear separation of concerns
- Easy to add new notification types

---

## 🔄 Error Handling

**Current Strategy:**
- Invalid type → Log error, don't route
- Parse errors → Log error, continue
- Kafka errors → Retry automatically (kafkajs)

**No retry logic for:**
- Unknown notification types
- Malformed JSON

**Rationale:**
- Let upstream (API Gateway) handle validation
- Processor focuses on routing only

---

**Last Updated:** November 17, 2025