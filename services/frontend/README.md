# Sports Analytics Frontend

A lightweight, modern sports analytics web interface built with vanilla HTML, CSS, and JavaScript.

## 📁 File Structure

```
frontend/
├── index.html              # Main homepage
├── search.html             # Search page
├── player.html             # Player profile page
├── team.html               # Team information page
├── tournament.html         # Tournament standings page
├── manager.html            # Manager profile page
├── match.html              # Match details page
├── styles.css              # Global styles
├── api.js                  # API client library
├── app.js                  # Main application logic
├── search.js               # Search page logic
├── player.js               # Player page logic
├── team.js                 # Team page logic
├── tournament.js           # Tournament page logic
├── manager.js              # Manager page logic
├── match.js                # Match page logic
├── Dockerfile              # Docker configuration
├── nginx.conf              # Nginx server configuration
└── .dockerignore            # Docker ignore patterns
```

## 🎯 Features

### Core Pages

1. **Index (Home)** - Main dashboard with:
   - Hero section with search
   - Recent matches
   - Popular tournaments
   - User favorites (when authenticated)

2. **Search** - Global search for:
   - Players
   - Teams
   - Tournaments
   - Managers

3. **Player Profile** - Shows:
   - Player details (name, position, nation)
   - Career statistics
   - Teams history
   - Match fixtures
   - Favorite button (if authenticated)

4. **Team Profile** - Shows:
   - Team information
   - Team statistics
   - Next upcoming game
   - Tournament standing
   - Team roster
   - Manager information
   - Recent fixtures
   - Favorite button (if authenticated)

5. **Tournament** - Shows:
   - Tournament standings table
   - Team statistics
   - Top players
   - All fixtures
   - Favorite button (if authenticated)

6. **Manager Profile** - Shows:
   - Manager details
   - Career statistics
   - Managed teams
   - Match history
   - Favorite button (if authenticated)

7. **Match Details** - Shows:
   - Match score and status
   - Team statistics comparison
   - Player performance
   - Goalie statistics
   - Detailed match stats

## 🔐 Authentication

The frontend includes built-in authentication with:

- **Registration** - Create new user account
- **Login** - Authenticate with JWT token
- **Session Management** - Token stored in localStorage
- **Protected Endpoints** - JWT token sent with protected requests
- **Auto-logout** - Redirects to login if token expires (401)

### Auth Flow

1. User clicks "Login" button
2. Auth panel appears with login/register tabs
3. User submits credentials
4. JWT token received and stored
5. User email displayed in navbar
6. Token included in all API requests

## 🎨 Design & UX

### Design Principles

- **Minimalist** - Clean, uncluttered interface
- **Responsive** - Works on desktop, tablet, mobile
- **Fast** - No heavy frameworks, pure vanilla JS
- **Accessible** - Semantic HTML, proper ARIA labels
- **Dark-friendly** - Light color scheme

### Color Scheme

- Primary: #3498db (Blue)
- Secondary: #2c3e50 (Dark Blue-Gray)
- Accent: #e74c3c (Red)
- Background: #f5f7fa (Light Gray)
- Text: #2c3e50 (Dark Gray)

### Typography

- Font Family: System fonts (Apple System, Segoe UI, Arial)
- Heading: Bold, 1.8-2.5rem
- Body: 1rem
- Small: 0.9rem

## 🔌 API Integration

### Key API Modules (in api.js)

```javascript
// Token Management
TokenManager.setToken(token)     // Store JWT
TokenManager.getToken()          // Retrieve JWT
TokenManager.removeToken()       // Clear token
TokenManager.hasToken()          // Check if authenticated

// Authentication
authAPI.register(username, password)
authAPI.login(username, password)
authAPI.getMe()

// Players
playerAPI.getDetails(id)
playerAPI.getStats(id)
playerAPI.getFixtures(id)
playerAPI.getTeams(id)
playerAPI.getFavorites()
playerAPI.addFavorite(id)
playerAPI.removeFavorite(id)

// Teams, Tournaments, Managers, Fixtures
// Similar structure for each resource type
```

### API Call Pattern

```javascript
// All API calls use consistent pattern
async function apiCall(endpoint, options = {}) {
    // Automatically adds JWT token if available
    // Handles 401 responses (token expiry)
    // Returns JSON or raw response
}
```

## 🚀 Running Locally

### Docker

```bash
# Build
docker build -t sports-frontend .

# Run
docker run -p 80:80 sports-frontend

# Access
open http://localhost
```

### Direct with Nginx

```bash
# Start nginx with config
nginx -c $(pwd)/nginx.conf

# Access
open http://localhost
```

## 📱 Mobile Optimization

The frontend is fully responsive with breakpoints:

```css
/* Tablet and below */
@media (max-width: 768px) {
    /* Stacked layout */
    /* Larger touch targets */
    /* Adjusted typography */
}

/* Mobile specific */
@media (max-width: 480px) {
    /* Single column */
    /* Optimized navigation */
}
```

## 🔄 Data Flow

```
User Action (click, type)
    ↓
JavaScript Event Handler
    ↓
API Call (api.js)
    ↓
Nginx Proxy (/:80)
    ↓
Backend Service (core_api, auth_service)
    ↓
Response Processing
    ↓
DOM Update
    ↓
User Sees Result
```

## 🛠️ Development Tips

### Adding a New Page

1. Create `newpage.html` with structure:
```html
<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="styles.css">
</head>
<body>
    <!-- Navigation & Auth Panel (copy from index.html) -->
    <main class="container">
        <!-- Your content -->
    </main>
    <script src="api.js"></script>
    <script src="app.js"></script>
    <script src="newpage.js"></script>
</body>
</html>
```

2. Create `newpage.js` with logic:
```javascript
document.addEventListener('DOMContentLoaded', async () => {
    await checkAuth();
    updateAuthUI();
    // Your initialization code
});
```

3. Add navigation link in navbar:
```html
<a href="/newpage.html" class="nav-link">New Page</a>
```

### Debugging

- **Console**: `F12` or `Ctrl+Shift+I`
- **Network Tab**: Monitor API calls
- **Local Storage**: Check `token` key
- **API Testing**: Use `fetch()` in console

### Performance Tips

- Lazy load images
- Cache DOM queries
- Minimize reflows
- Use CSS Grid/Flexbox
- Optimize asset sizes

## 🔒 Security Considerations

- JWT tokens stored in localStorage (XSS risk - consider httpOnly cookies)
- All API calls go through Nginx proxy
- CORS headers configured on backend
- No sensitive data in localStorage
- Token sent only in Authorization header

## 🚨 Error Handling

The frontend includes:
- API error messages displayed to users
- Graceful fallbacks for failed requests
- Auto-retry for network errors
- Clear loading states
- User-friendly error messages

## 🎯 Features Matrix

| Feature | Public | Authenticated |
|---------|--------|---------------|
| Browse content | ✅ | ✅ |
| Search | ✅ | ✅ |
| View details | ✅ | ✅ |
| Add to favorites | ❌ | ✅ |
| View favorites | ❌ | ✅ |
| User profile | ❌ | ✅ |

## 📚 Resources

- HTML/CSS: [MDN Web Docs](https://developer.mozilla.org/)
- JavaScript: [JavaScript.info](https://javascript.info/)
- Fetch API: [MDN Fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API)
- CSS Grid: [CSS-Tricks Grid Guide](https://css-tricks.com/snippets/css/complete-guide-grid/)

## 🐛 Known Issues & Limitations

- No offline support
- Limited to 20 items per page (add pagination if needed)
- No image caching
- No service worker

## 📝 Future Enhancements

- [ ] Progressive Web App (PWA)
- [ ] Offline support
- [ ] Dark mode toggle
- [ ] Pagination for large datasets
- [ ] Advanced filtering
- [ ] Export data (CSV/PDF)
- [ ] Mobile app (React Native)
- [ ] Real-time updates (WebSocket)

---

**Built with ❤️ using vanilla JavaScript, HTML, and CSS**
