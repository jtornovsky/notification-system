# Notification System Dashboard

## 📌 Overview

A React-based web dashboard for the notification system. Provides a user interface to submit notifications and view real-time analytics from Elasticsearch.

**Framework:** React 18 + Vite  
**Port:** 5173 (dev server)  
**Language:** JavaScript

---

## 🏗️ Architecture
```
React Dashboard (localhost:5173)
    ↓ HTTP Requests
API Gateway (localhost:8080)
    ↓
Backend Services
```

**Features:**
- Submit notifications (EMAIL, SMS, PUSH)
- Auto-detect notification type based on recipient
- View real-time analytics from Elasticsearch
- Responsive two-column layout
- Form validation and error handling

---

## 🛠️ Tech Stack

- **React:** 18.3.1
- **Vite:** 7.2.2 (Build tool & dev server)
- **Build Tool:** ESBuild (via Vite)
- **CSS:** Custom styles with Flexbox/Grid

---

## 🚀 Quick Start

### **Prerequisites**
- Node.js 18+ and npm
- API Gateway running on `localhost:8080`
- Backend services running (for end-to-end testing)

### **Installation**
```bash
cd dashboard

# Install dependencies
npm install

# Start dev server
npm run dev
```

**Expected output:**
```
VITE v7.2.2  ready in 571 ms
➜  Local:   http://localhost:5173/
➜  Network: use --host to expose
```

**Open browser:** `http://localhost:5173`

---

## 📱 Features

### 1. **Send Notification Form**

**Fields:**
- **User ID** - Identifier for the user
- **Type (auto-detected)** - EMAIL, SMS, or PUSH
- **Recipient** - Email, phone number, or device token
- **Subject** - Optional subject line (for emails)
- **Message** - Notification content

**Smart Type Detection:**
- Contains `@` and `.` → **EMAIL**
- 10-15 digits with `+()- ` → **SMS**
- Anything else → **PUSH**

**Example:**
- Type `test@example.com` → Auto-detects as EMAIL
- Type `+1 (555) 123-4567` → Auto-detects as SMS
- Type `device_token_abc123` → Auto-detects as PUSH

### 2. **Analytics Display**

**Shows:**
- Total notifications delivered
- Breakdown by status (SENT, FAILED)
- Breakdown by type (EMAIL, SMS, PUSH)

**Updates:**
- Auto-loads on page load
- Manual refresh via "Refresh Analytics" button
- Auto-refreshes after sending notification

---

## 🎨 UI Components

### **Main Layout**
```
┌─────────────────────────────────────────┐
│     Notification System Dashboard      │
└─────────────────────────────────────────┘

┌──────────────────┬──────────────────────┐
│  Send            │  Analytics           │
│  Notification    │                      │
│                  │  [Statistics Box]    │
│  [Form Fields]   │                      │
│                  │  [Refresh Button]    │
│  [Send Button]   │                      │
└──────────────────┴──────────────────────┘
```

### **Color Scheme**
- **Primary Blue:** `#3498db` (buttons, headers)
- **Dark Blue:** `#2c3e50` (header background)
- **Light Gray:** `#f5f5f5` (page background)
- **White:** `#ffffff` (card backgrounds)

---

## 📊 User Flow

### **Sending a Notification**

1. User opens dashboard (http://localhost:5173/)
2. Fills out form fields
3. Types recipient → Type auto-detects
4. Clicks "Send Notification"
5. Success alert appears
6. Form clears automatically
7. Analytics refreshes

### **Viewing Analytics**

1. Dashboard loads → Auto-fetches analytics
2. User can manually click "Refresh Analytics"
3. Analytics display updates with latest data

---

## 🔧 Configuration

### **API Endpoint**

**Location:** `src/App.jsx`
```javascript
// Notification submission
fetch('http://localhost:8080/notifications', {...})

// Analytics fetch
fetch('http://localhost:8080/analytics', {...})
```

**To change API URL:**
1. Update URLs in `fetchAnalytics()` and `handleSubmit()`
2. Save file (HMR will reload automatically)

### **Styling**

**Location:** `src/App.css`

**Key CSS Variables:**
```css
/* Modify these for theme changes */
header { background-color: #2c3e50; }
button { background-color: #3498db; }
body { background-color: #f5f5f5; }
```

---

## 📂 Project Structure
```
dashboard/
├── public/              # Static assets
├── src/
│   ├── App.jsx         # Main component (form + analytics)
│   ├── App.css         # Styling
│   ├── main.jsx        # React entry point
│   └── assets/         # Images, icons
├── index.html          # Root HTML
├── package.json        # Dependencies
├── vite.config.js      # Vite configuration
└── README.md           # This file
```

---

## 🧪 Testing

### **Manual Testing**

**Test Notification Submission:**
1. Fill form:
    - User ID: `test123`
    - Recipient: `test@example.com`
    - Subject: `Test`
    - Message: `Hello`
2. Click "Send Notification"
3. Verify success alert appears
4. Check browser console for logs

**Test Type Auto-Detection:**
1. Type `test@example.com` in Recipient → Should show EMAIL
2. Clear and type `5551234567` → Should show SMS
3. Clear and type `token_abc` → Should show PUSH

**Test Analytics:**
1. Click "Refresh Analytics" button
2. Verify textarea updates with JSON data
3. Check browser console for logs

### **Browser Console Verification**

**Expected logs:**
```javascript
Sending notification: {user_id: "test123", type: "EMAIL", ...}
Success: {id: "20251117123456", ...}
Analytics data: {took: 5, hits: {...}, aggregations: {...}}
```

---

## 🐛 Troubleshooting

### **Error: "Failed to fetch"**
- **Cause:** API Gateway not running or CORS issue
- **Fix:**
    1. Start API Gateway: `cd api-gateway && go run main.go`
    2. Verify CORS is enabled in API Gateway

### **Type not auto-detecting**
- **Cause:** Recipient field empty or invalid format
- **Fix:**
    - Email must have `@` and `.` with `@` before `.`
    - Phone must be 10-15 digits (can include `+()- `)

### **Analytics showing raw JSON**
- **Status:** Normal behavior (raw data display)
- **Future:** Will be replaced with charts and cards

### **Port 5173 already in use**
- **Cause:** Another dev server running
- **Fix:**
    1. Kill other process
    2. Or change port in `vite.config.js`

### **Form doesn't clear after submit**
- **Cause:** API call failed
- **Fix:** Check browser console for errors

---

## 📜 Scripts
```bash
# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Lint code
npm run lint
```

---

## 🔄 Development Workflow

### **Hot Module Replacement (HMR)**

Vite provides instant updates without full page reload:
1. Edit `src/App.jsx` or `src/App.css`
2. Save file
3. Browser updates instantly

### **Adding New Features**

**Example: Add a notification history list**

1. Add state:
```javascript
const [history, setHistory] = useState([])
```

2. Update on submit:
```javascript
setHistory([...history, notification])
```

3. Display in UI:
```javascript
<ul>
  {history.map(n => <li key={n.id}>{n.message}</li>)}
</ul>
```

---

## 🚀 Future Enhancements

### **Phase 1: Better Analytics Display**
- [ ] Replace JSON textarea with visual cards
- [ ] Add charts (pie chart for types, bar chart for status)
- [ ] Show total count prominently

### **Phase 2: Enhanced UX**
- [ ] Loading spinners during API calls
- [ ] Toast notifications instead of alerts
- [ ] Form validation with error messages
- [ ] Confirmation before sending

### **Phase 3: Advanced Features**
- [ ] Notification history table
- [ ] Search/filter notifications
- [ ] Real-time updates (WebSocket)
- [ ] Dark mode toggle
- [ ] Export analytics to CSV

### **Phase 4: Production Ready**
- [ ] Authentication (login/logout)
- [ ] Environment-based API URLs
- [ ] Error boundary component
- [ ] Retry logic for failed requests
- [ ] Accessibility improvements (ARIA labels)

---

## 🎓 Learning Resources

### **React Concepts Used**

- **`useState`** - Component state management
- **`useEffect`** - Side effects (auto-fetch on mount)
- **`async/await`** - Asynchronous API calls
- **Controlled Components** - Form inputs controlled by state
- **Event Handlers** - `onClick`, `onChange`

### **Key React Patterns**

**State Management:**
```javascript
const [value, setValue] = useState(initial)
```

**Effect Hook:**
```javascript
useEffect(() => {
  // Runs on mount
}, [])
```

**Event Handling:**
```javascript
<input onChange={(e) => setValue(e.target.value)} />
```

---

## 📝 Notes

- **No authentication** - Open access to dashboard
- **No error boundaries** - Errors will break UI
- **No loading states** - Instant transitions
- **Raw JSON analytics** - No data visualization yet
- **No TypeScript** - Pure JavaScript implementation

---

## 🔐 Security Considerations

**Current:**
- No authentication required
- API calls in plain HTTP
- No input sanitization
- CORS allows all origins

**Production Requirements:**
- Add login/authentication
- Use HTTPS for API calls
- Sanitize user inputs
- Implement CSRF protection
- Add rate limiting

---

## 🤝 Contributing

This is a portfolio/learning project. Key areas for improvement:
1. Replace raw JSON with charts
2. Add proper error handling
3. Implement loading states
4. Add TypeScript
5. Write unit tests

---

**Last Updated:** November 20, 2025  
**Built with:** React 18 + Vite 7  
**Part of:** Notification System Microservices Project