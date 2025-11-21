const { Kafka } = require('kafkajs');
const winston = require('winston');
const root = require('../proto/github.com/jtornovsky/notification-system/proto/js/notification_pb');
const { Notification, NotificationType } = root.notification;

// ============================================================================
// LOGGER CONFIGURATION
// ============================================================================
const logger = winston.createLogger({
    level: 'info',              // Log level: debug, info, warn, error
    format: winston.format.simple(),  // Simple text format (not JSON)
    transports: [
        new winston.transports.Console()  // Output logs to console (terminal)
    ]
});

logger.info('🚀 Notification Processor starting...');


// ============================================================================
// KAFKA CONFIGURATION
// ============================================================================
// Create Kafka client (connection to Kafka broker)
const kafka = new Kafka({
    clientId: 'notification-processor',  // Name of this application (for Kafka logs)
    brokers: ['localhost:9092']          // Kafka server addresses (array because can have multiple)
});
logger.info('✓ Kafka configured');

// Create producer (for publishing to type-specific topics)
const producer = kafka.producer();
logger.info('✓ Producer created');


// ============================================================================
// KAFKA CONSUMER CREATION
// ============================================================================
// Create a consumer instance
// groupId = consumer group name (multiple consumers with same group share work)
const consumer = kafka.consumer({
    groupId: 'notification-processor-group'
});

logger.info('✓ Consumer created');


// ============================================================================
// MAIN START FUNCTION
// ============================================================================
// Main function to connect and start consuming messages
async function start() {
    try {
        // ----------------------------------------------------------------
        // CONNECT TO KAFKA
        // ----------------------------------------------------------------
        logger.info('Connecting consumer to Kafka...');
        await consumer.connect();
        logger.info('✓ Consumer connected to Kafka');

        logger.info('Connecting producer to Kafka...');
        await producer.connect();
        logger.info('✓ Producer connected to Kafka');

        // ----------------------------------------------------------------
        // SUBSCRIBE TO TOPIC
        // ----------------------------------------------------------------
        logger.info('Subscribing to topic: notifications');

        await consumer.subscribe({
            topic: 'notifications',     // Topic name (created by API Gateway)
            fromBeginning: true         // Read all messages from start (good for testing)
                                        // In production, use false (only new messages)
        });

        logger.info('✓ Subscribed to topic');


        // ----------------------------------------------------------------
        // START CONSUMING MESSAGES
        // ----------------------------------------------------------------
        logger.info('Starting message consumer...');

        // Start the consumer - this runs forever until stopped
        // "eachMessage" is a callback that runs for EVERY message received
        await consumer.run({
            eachMessage: async ({ topic, partition, message }) => {

                logger.info('📩 Received Protobuf message:', {
                    topic: topic,
                    partition: partition,
                    offset: message.offset,
                    sizeBytes: message.value.length
                });

                try {
                    // Decode Protobuf bytes to object
                    const pbNotification = Notification.decode(message.value);

                    // Convert to plain JavaScript object
                    const notification = {
                        id: pbNotification.id,
                        userId: pbNotification.userId,
                        type: pbNotification.type,  // This is a NUMBER (0, 1, 2)
                        recipient: pbNotification.recipient,
                        subject: pbNotification.subject,
                        message: pbNotification.message,
                        createdAt: pbNotification.createdAt,
                        status: pbNotification.status
                    };

                    // Convert type NUMBER to STRING
                    const typeNames = ['EMAIL', 'SMS', 'PUSH'];
                    const typeString = typeNames[notification.type] || 'UNKNOWN';

                    logger.info('📝 Parsed notification:', {
                        id: notification.id,
                        type: typeString,
                        typeNumber: notification.type,
                        recipient: notification.recipient
                    });

                    // Route the notification
                    await routeNotification(notification, typeString, pbNotification);

                } catch (parseError) {
                    logger.error('❌ Failed to parse message:', parseError);
                }
            }
        });

        // This line only reached after consumer.run() starts successfully
        logger.info('✓ Consumer is running and waiting for messages...');

    } catch (error) {
        // If ANY error happens during startup (connect, subscribe, run)
        // Log it and exit the process with error code 1
        logger.error('❌ Error starting consumer:', error);
        process.exit(1);  // Exit program with error status
    }
}


// ============================================================================
// GRACEFUL SHUTDOWN HANDLER
// ============================================================================
// Handle Ctrl+C (SIGINT signal) to shut down gracefully
// When user presses Ctrl+C, this runs instead of abrupt termination
process.on('SIGINT', async () => {
    logger.info('Shutting down...');

    // Disconnect consumer
    await consumer.disconnect();
    logger.info('✓ Consumer disconnected');

    // Disconnect producer
    await producer.disconnect();
    logger.info('✓ Producer disconnected');

    logger.info('✓ Disconnected from Kafka');

    // Exit with success code 0
    process.exit(0);
});

// Function to route notification to type-specific topic
async function routeNotification(notification, typeString, pbNotification) {
    let targetTopic;

    switch (typeString) {
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
            logger.error('❌ Unknown notification type:', typeString);
            return;
    }

    logger.info('📤 Routing to topic:', {
        notificationId: notification.id,
        type: typeString,
        targetTopic: targetTopic
    });

    try {
        // Encode back to Protobuf
        const pbBuffer = Notification.encode(pbNotification).finish();

        await producer.send({
            topic: targetTopic,
            messages: [
                {
                    key: notification.id,
                    value: pbBuffer
                }
            ]
        });

        logger.info('✓ Successfully routed to:', { targetTopic });

    } catch (error) {
        logger.error('❌ Failed to route notification:', error);
    }
}

// ============================================================================
// START THE APPLICATION
// ============================================================================
// Call the start() function to begin consuming messages
// .catch() handles any errors that weren't caught inside start()
start().catch(error => {
    logger.error('❌ Unhandled error:', error);
    process.exit(1);
});
