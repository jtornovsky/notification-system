const { Kafka } = require('kafkajs');
const { Client } = require('@elastic/elasticsearch');
const root = require('../proto/github.com/jtornovsky/notification-system/proto/js/notification_pb');
const { DeliveryEvent, NotificationType, NotificationStatus } = root.notification;

// ============================================
// Configuration
// ============================================

const KAFKA_BROKERS = ['localhost:9092'];
const ELASTICSEARCH_NODE = 'http://localhost:9200';
const INDEX_NAME = 'notification-analytics';

// ============================================
// Elasticsearch Client
// ============================================

const esClient = new Client({
    node: ELASTICSEARCH_NODE
});

// Initialize Elasticsearch index
async function initializeIndex() {
    try {
        const indexExists = await esClient.indices.exists({ index: INDEX_NAME });

        if (!indexExists) {
            console.log('📊 Creating Elasticsearch index...');

            await esClient.indices.create({
                index: INDEX_NAME,
                body: {
                    mappings: {
                        properties: {
                            timestamp: { type: 'date' },
                            type: { type: 'keyword' },
                            status: { type: 'keyword' },
                            delivery_time_ms: { type: 'integer' },
                            recipient: { type: 'keyword' },
                            notification_id: { type: 'keyword' },
                            error_message: { type: 'text' }
                        }
                    }
                }
            });

            console.log('✅ Elasticsearch index ' + INDEX_NAME + ' created');
        } else {
            console.log('✅ Elasticsearch index ' + INDEX_NAME + ' already exists');
        }
    } catch (error) {
        console.error('❌ Failed to initialize Elasticsearch index: ' + INDEX_NAME, error.message);
        throw error;
    }
}

// Index a delivery event into Elasticsearch
async function indexDeliveryEvent(event) {
    try {
        await esClient.index({
            index: INDEX_NAME,
            body: {
                notification_id: event.notification_id,
                type: event.type,
                recipient: event.recipient,
                status: event.status,
                timestamp: event.timestamp,
                delivery_time_ms: event.delivery_time_ms,
                error_message: event.error_message || null
            }
        });

        console.log(`✅ Indexed: ${event.notification_id} (${event.type}, ${event.status})`);
    } catch (error) {
        console.error('❌ Failed to index event:', error.message);
    }
}

// ============================================
// Kafka Consumer
// ============================================

const kafka = new Kafka({
    clientId: 'analytics-service',
    brokers: KAFKA_BROKERS
});

const consumer = kafka.consumer({
    groupId: 'analytics-service-group'
});

// Start consuming delivery events
async function startConsumer() {
    try {
        await consumer.connect();
        console.log('✅ Connected to Kafka');

        await consumer.subscribe({
            topic: 'delivery-events',
            fromBeginning: true
        });
        console.log('✅ Subscribed to delivery-events topic');

        await consumer.run({
            eachMessage: async ({ message }) => {
                try {
                    // Decode Protobuf message
                    const pbEvent = DeliveryEvent.decode(message.value);

                    // Convert to plain object
                    const event = {
                        notification_id: pbEvent.notificationId,
                        type: NotificationType[pbEvent.type] || 'UNKNOWN',
                        recipient: pbEvent.recipient,
                        status: NotificationStatus[pbEvent.status] || 'UNKNOWN',
                        timestamp: new Date(Number(pbEvent.processedAt)),
                        delivery_time_ms: pbEvent.deliveryTimeMs,
                        error_message: pbEvent.errorMessage || null
                    };

                    console.log(`📩 Received Protobuf: ${event.notification_id} (${event.type}, ${event.status})`);

                    await indexDeliveryEvent(event);

                } catch (error) {
                    console.error('❌ Error processing message:', error.message);
                }
            }
        });

    } catch (error) {
        console.error('❌ Failed to start consumer:', error.message);
        throw error;
    }
}

// ============================================
// Main Application
// ============================================

async function main() {
    console.log('🚀 Starting Analytics Service...');

    try {
        // Initialize Elasticsearch
        await initializeIndex();

        // Start Kafka consumer
        await startConsumer();

        console.log('✅ Analytics Service running');
        console.log('📊 Consuming delivery events and indexing to Elasticsearch...');

    } catch (error) {
        console.error('❌ Failed to start Analytics Service:', error);
        process.exit(1);
    }
}

// Graceful shutdown
async function shutdown() {
    console.log('\n📭 Shutting down...');
    await consumer.disconnect();
    console.log('✅ Consumer disconnected');
    process.exit(0);
}

process.on('SIGINT', shutdown);
process.on('SIGTERM', shutdown);

// Start the service
main();