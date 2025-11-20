import { useState, useEffect } from 'react'
import './App.css'

function App() {
    const [userId, setUserId] = useState('')
    const [type, setType] = useState('')
    const [recipient, setRecipient] = useState('')
    const [subject, setSubject] = useState('')
    const [message, setMessage] = useState('')
    const [analytics, setAnalytics] = useState(null)

    useEffect(() => {
        fetchAnalytics();  // Auto-call on mount
    }, []);  // run once on load

    function detectTypeAndSetRecipient(value) {
        setRecipient(value);

        // If empty, clear type
        if (value.length === 0) {
            setType('');
        }
        // Check for email (has @ and domain with dot)
        else if (value.includes('@') && value.includes('.') && value.indexOf('@') < value.lastIndexOf('.')) {
            setType('EMAIL');
        }
        // Check for phone (digits with optional +, (), -, spaces)
        else if (/^[\d\+\(\)\-\s]{10,20}$/.test(value)) {
            const digitCount = value.replace(/\D/g, '').length;
            if (digitCount >= 10 && digitCount <= 15) {
                setType('SMS');
            } else {
                setType('PUSH');
            }
        }
        // Anything else is PUSH
        else if (value.length > 0) {
            setType('PUSH');
        }
    }

    function handleSubmit() {
        const notification = {
            user_id: userId,
            type: type,
            recipient: recipient,
            subject: subject,
            message: message
        };

        console.log('Sending notification:', notification);

        // TODO: Send to API Gateway
    }

    async function fetchAnalytics() {
        try {
            const response = await fetch('http://localhost:9200/notification-analytics/_search', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    size: 0,
                    aggs: {
                        by_status: {
                            terms: { field: 'status' }
                        },
                        by_type: {
                            terms: { field: 'type' }
                        }
                    }
                })
            });

            const data = await response.json();
            console.log('Analytics data:', data);
            setAnalytics(data);
        } catch (error) {
            console.error('Error fetching analytics:', error);
        }
    }

    return (
        <div className="app">
            <header>
                <h1>Notification System Dashboard</h1>
            </header>

            <main>
                <section>
                    <h2>Send Notification</h2>
                    <div className="form">
                        <div>
                            <label>User ID:</label>
                            <input
                                type="text"
                                value={userId}
                                onChange={(e) => setUserId(e.target.value)}
                                placeholder="Enter user ID"
                            />
                        </div>

                        <div>
                            <label>Type (auto-detected):</label>
                            <input
                                type="text"
                                value={type}
                                placeholder="Auto-detected from recipient"
                                readOnly
                                style={{ backgroundColor: '#f0f0f0', cursor: 'not-allowed' }}
                            />
                        </div>

                        <div>
                            <label>Recipient:</label>
                            <input
                                type="text"
                                value={recipient}
                                onChange={(e) => detectTypeAndSetRecipient(e.target.value)}
                                placeholder="email@example.com / phone / device_token"
                            />
                        </div>

                        <div>
                            <label>Subject:</label>
                            <input
                                type="text"
                                value={subject}
                                onChange={(e) => setSubject(e.target.value)}
                                placeholder="Notification subject"
                            />
                        </div>

                        <div>
                            <label>Message:</label>
                            <textarea
                                value={message}
                                onChange={(e) => setMessage(e.target.value)}
                                placeholder="Notification message"
                                rows="4"
                            />
                        </div>

                        <button onClick={handleSubmit}>Send Notification</button>
                    </div>
                </section>

                <section>
                    <h2>Analytics</h2>
                    <div className="form">
                        <div>
                            <label>Statistics:</label>
                            <textarea
                                value={analytics ? JSON.stringify(analytics, null, 2) : ''}
                                placeholder="Analytics data appear here..."
                                readOnly
                                rows="10"
                                style={{ backgroundColor: '#f0f0f0', cursor: 'not-allowed' }}
                            />
                        </div>

                        <button onClick={fetchAnalytics}>Refresh Analytics</button>
                    </div>
                </section>
            </main>
        </div>
    )
}

export default App