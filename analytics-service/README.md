# Analytics Service

## 📌 Overview

The Analytics Service consumes delivery events from Kafka, indexes them into Elasticsearch, and provides a searchable repository for notification analytics and metrics.

**Service Type:** Event consumer/indexer  
**Language:** Node.js  
**Primary Function:** Real-time event indexing for analytics

---

## 🏗️ Architecture
```
Kafka: delivery-events topic
         ↓
┌──────────────────────┐
│  Analytics Service   │
│  (Node.js)           │
└──────────┬───────────┘
           │
           └→ Elasticsearch
              - Index: notification-analytics
              - Searchable events
              - Aggregations
```

**Position in System:**
- Consumes from: `delivery-events` topic
- Writes to: Elasticsearch
- Enables: Real-time analytics and search
- Stateless: No in-memory state (all data stored in Elasticsearch)

---

## 🛠️ Tech Stack

- **Runtime:** Node.js 18+
- **Kafka Client:** kafkajs
- **Search Engine:** Elasticsearch (@elastic/elasticsearch v8)
- **Architecture Pattern:** Event Consumer/Indexer

---

## 🚀 Quick Start

### Prerequisites

- Node.js 18 or higher
- Kafka running on `localhost:9092`
- Elasticsearch running on `localhost:9200`
- Topic must exist: `delivery-events`

### Installation
```bash
cd analytics-service

# Install dependencies
npm install

# Run the service
node index.js
```

### Expected Output
```
🚀 Starting Analytics Service...
✅ Elasticsearch index already exists
✅ Connected to Kafka
✅ Subscribed to delivery-events topic
✅ Analytics Service running
📊 Consuming delivery events and indexing to Elasticsearch...
```

---

## 📊 Data Flow

### Input: Kafka Messages

**Topic:** `delivery-events`  
**Format:** JSON

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

### Output: Elasticsearch Documents

**Index:** `notification-analytics`  
**Mapping:**
```javascript
{
  "timestamp": { "type": "date" },
  "type": { "type": "keyword" },
  "status": { "type": "keyword" },
  "delivery_time_ms": { "type": "integer" },
  "recipient": { "type": "keyword" },
  "notification_id": { "type": "keyword" },
  "error_message": { "type": "text" }
}
```

**Field Types Explained:**
- `keyword` - Exact match, aggregations, sorting
- `text` - Full-text search
- `date` - Timestamp queries
- `integer` - Numeric calculations

---

## 🔍 Processing Logic

### Step-by-Step Flow
```
1. Consume message from delivery-events
         ↓
2. Parse JSON
         ↓
3. Extract fields
         ↓
4. Index to Elasticsearch
         ↓
5. Log success
         ↓
6. Commit Kafka offset
```

### Code Overview
```javascript
// Consume from Kafka
await consumer.subscribe({ 
    topic: 'delivery-events',
    fromBeginning: true
});

// For each message
const event = JSON.parse(message.value);

// Index to Elasticsearch
await esClient.index({
    index: 'notification-analytics',
    body: {
        notification_id: event.notification_id,
        type: event.type,
        status: event.status,
        timestamp: event.timestamp,
        delivery_time_ms: event.delivery_time_ms,
        error_message: event.error_message || null
    }
});
```

---

## 🧪 Testing

### Test 1: Verify Service is Running
```bash
node index.js
```

**Expected output:**
```
🚀 Starting Analytics Service...
✅ Elasticsearch index notification-analytics created
✅ Connected to Kafka
✅ Subscribed to delivery-events topic
✅ Analytics Service running
📊 Consuming delivery events and indexing to Elasticsearch...
```

### Test 2: Send Test Notification
```powershell
Invoke-WebRequest -Uri http://localhost:8080/notifications -Method POST -ContentType "application/json" -Body '{"user_id":"test","type":"EMAIL","recipient":"test@example.com","subject":"Test","message":"Analytics test"}'
```

**Analytics Service logs should show:**
```
📩 Received: abc-123 (EMAIL, SENT)
✅ Indexed: abc-123 (EMAIL, SENT)
```

### Test 3: Verify in Elasticsearch

**Count indexed documents:**
```powershell
Invoke-WebRequest -Uri http://localhost:9200/notification-analytics/_count
```

**Expected response:**
```json
{
  "count": 1,
  "_shards": {"total": 1, "successful": 1}
}
```

### Test 4: Query Elasticsearch

**Get all documents:**
```powershell
Invoke-WebRequest -Uri "http://localhost:9200/notification-analytics/_search?pretty" | Select-Object -ExpandProperty Content
```

### Test 5: Aggregate by Status
```powershell
$body = @"
{
  "size": 0,
  "aggs": {
    "by_status": {
      "terms": { "field": "status" }
    }
  }
}
"@

Invoke-WebRequest -Uri "http://localhost:9200/notification-analytics/_search" -Method POST -ContentType "application/json" -Body $body | Select-Object -ExpandProperty Content
```

**Expected output:**
```json
{
  "aggregations": {
    "by_status": {
      "buckets": [
        {"key": "SENT", "doc_count": 45},
        {"key": "FAILED", "doc_count": 5}
      ]
    }
  }
}
```

### Test 6: Aggregate by Type
```powershell
$body = @"
{
  "size": 0,
  "aggs": {
    "by_type": {
      "terms": { "field": "type" }
    }
  }
}
"@

Invoke-WebRequest -Uri "http://localhost:9200/notification-analytics/_search" -Method POST -ContentType "application/json" -Body $body | Select-Object -ExpandProperty Content
```

### Test 7: Average Delivery Time
```powershell
$body = @"
{
  "size": 0,
  "aggs": {
    "avg_delivery_time": {
      "avg": { "field": "delivery_time_ms" }
    }
  }
}
"@

Invoke-WebRequest -Uri "http://localhost:9200/notification-analytics/_search" -Method POST -ContentType "application/json" -Body $body | Select-Object -ExpandProperty Content
```

---

## 🔧 Configuration

**All configuration is in `index.js`:**
```javascript
const KAFKA_BROKERS = ['localhost:9092'];
const ELASTICSEARCH_NODE = 'http://localhost:9200';
const INDEX_NAME = 'notification-analytics';
```

**Kafka Settings:**
```javascript
const consumer = kafka.consumer({ 
    groupId: 'analytics-service-group'
});

await consumer.subscribe({ 
    topic: 'delivery-events',
    fromBeginning: true
});
```

**Key Setting:** `fromBeginning: true`
- Processes ALL messages from topic start
- Important for analytics (want complete history)
- Without this: only processes new messages

---

## 📂 Project Structure
```
analytics-service/
├── package.json
├── index.js
└── README.md
```

**Single-file architecture:**
- Simple indexing logic (~200 lines)
- No complex business rules
- Easy to understand and maintain

---

## 🐛 Troubleshooting

### Error: "ResponseError: index_not_found_exception"

**Create index manually:**
```bash
curl -X PUT "localhost:9200/notification-analytics" -H 'Content-Type: application/json' -d'
{
  "mappings": {
    "properties": {
      "timestamp": { "type": "date" },
      "type": { "type": "keyword" },
      "status": { "type": "keyword" },
      "delivery_time_ms": { "type": "integer" },
      "recipient": { "type": "keyword" },
      "notification_id": { "type": "keyword" },
      "error_message": { "type": "text" }
    }
  }
}'
```

### Error: "KafkaJSNumberOfRetriesExceeded"

**Check Kafka:**
```bash
docker-compose ps
docker-compose logs notification_kafka
```

**Verify topic exists:**
```bash
docker exec -it notification_kafka kafka-topics --list --bootstrap-server localhost:9092
```

**Create topic if missing:**
```bash
docker exec -it notification_kafka kafka-topics --create --bootstrap-server localhost:9092 --topic delivery-events --partitions 1 --replication-factor 1
```

### Error: "connection refused" (Elasticsearch)

**Check Elasticsearch:**
```bash
docker-compose ps
docker-compose logs notification_elasticsearch
curl http://localhost:9200
```

### No messages being indexed

**Check Delivery Service is running and verify messages in topic:**
```bash
docker exec -it notification_kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic delivery-events \
  --from-beginning
```

---

## 🔍 Monitoring

### Check Consumer Group Status
```bash
docker exec -it notification_kafka kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --describe \
  --group analytics-service-group
```

**Output shows:**
- Current offset
- Log end offset
- Lag (how far behind)

### Check Elasticsearch Index

**Document count:**
```bash
curl http://localhost:9200/notification-analytics/_count
```

**Index stats:**
```bash
curl http://localhost:9200/notification-analytics/_stats?pretty
```

**Mappings:**
```bash
curl http://localhost:9200/notification-analytics/_mapping?pretty
```

---

## 📈 Query Examples

### 1. Recent Deliveries (Last 10)
```bash
curl -X POST "localhost:9200/notification-analytics/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "size": 10,
  "sort": [{ "timestamp": "desc" }]
}'
```

### 2. Failed Deliveries
```bash
curl -X POST "localhost:9200/notification-analytics/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "term": { "status": "FAILED" }
  }
}'
```

### 3. Email Notifications Only
```bash
curl -X POST "localhost:9200/notification-analytics/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "term": { "type": "EMAIL" }
  }
}'
```

### 4. Deliveries in Last Hour
```bash
curl -X POST "localhost:9200/notification-analytics/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "range": {
      "timestamp": {
        "gte": "now-1h"
      }
    }
  }
}'
```

### 5. Success Rate by Type
```bash
curl -X POST "localhost:9200/notification-analytics/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "size": 0,
  "aggs": {
    "by_type": {
      "terms": { "field": "type" },
      "aggs": {
        "by_status": {
          "terms": { "field": "status" }
        }
      }
    }
  }
}'
```

### 6. Average Delivery Time by Type
```bash
curl -X POST "localhost:9200/notification-analytics/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "size": 0,
  "aggs": {
    "by_type": {
      "terms": { "field": "type" },
      "aggs": {
        "avg_time": {
          "avg": { "field": "delivery_time_ms" }
        }
      }
    }
  }
}'
```

### 7. Search Error Messages
```bash
curl -X POST "localhost:9200/notification-analytics/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "match": {
      "error_message": "timeout"
    }
  }
}'
```

---

## 📊 Performance Characteristics

**Throughput:**
- Limited by Elasticsearch indexing speed
- Typically 1,000-10,000 docs/second (single node)
- Batching can improve performance

**Latency:**
- Near real-time indexing (~1 second delay)
- Elasticsearch refresh interval: 1 second default

**Scalability:**
- Service is stateless (no in-memory state)
- Multiple instances can run concurrently
- All instances share Elasticsearch for persistent state
- Kafka consumer group distributes load automatically
- Limited by Elasticsearch cluster capacity

---

## 🚀 Future Enhancements

- [ ] Batch indexing (index multiple docs at once)
- [ ] Data retention policy (delete old documents)
- [ ] Pre-computed aggregations (daily/hourly summaries)
- [ ] Alerting on failure rate thresholds
- [ ] Dashboard API endpoints (REST API for queries)
- [ ] Metrics (Prometheus integration)
- [ ] Health check endpoint
- [ ] Configuration via environment variables
- [ ] Structured logging

---

## 📝 Design Decisions

### Why Single File?
- Simple indexing logic
- No complex transformations
- ~200 lines total
- Consistent with Notification Processor

### Why `fromBeginning: true`?
- Analytics needs complete history
- If service restarts, doesn't miss old events
- Can rebuild index from scratch if needed

### Why Elasticsearch?
- Fast full-text search
- Powerful aggregations
- Time-series data optimization
- Built for analytics workloads

### Why No Database?
- Elasticsearch IS the database
- No need for additional storage
- Direct Kafka → Elasticsearch pipeline

---

## 🔄 Data Retention

**Current:** Unlimited storage (all events kept forever)

**Production Considerations:**

### Option 1: Index Lifecycle Management (ILM)
```javascript
// Auto-delete documents older than 90 days
// Elasticsearch ILM policy
```

### Option 2: Time-based Indexes
```javascript
// Create daily indexes
notification-analytics-2025-11-17
notification-analytics-2025-11-18
// Delete old indexes
```

### Option 3: Manual Cleanup
```bash
# Delete events older than 90 days
curl -X POST "localhost:9200/notification-analytics/_delete_by_query" -d'
{
  "query": {
    "range": {
      "timestamp": {
        "lt": "now-90d"
      }
    }
  }
}'
```

---

## 📊 Use Cases

### 1. Real-time Monitoring
- Track delivery success rates
- Alert on failure spikes
- Monitor delivery times

### 2. Historical Analysis
- Trends over time
- Performance by channel
- Error pattern analysis

### 3. Debugging
- Search for specific notification
- Find all failures for a recipient
- Analyze error messages

### 4. Reporting
- Daily/weekly/monthly stats
- Channel performance comparison
- SLA compliance tracking

---

## 🔗 Integration with Dashboard

**The React Dashboard will query this service's data through Elasticsearch:**
```javascript
// Dashboard queries Elasticsearch directly
fetch('http://localhost:9200/notification-analytics/_search', {
  method: 'POST',
  body: JSON.stringify({
    size: 0,
    aggs: {
      by_type: { terms: { field: 'type' } }
    }
  })
})
```

---

**Last Updated:** November 17, 2025