# Sports Analytics Platform

A complete sports analytics platform with a modern frontend, authentication, and comprehensive backend services.

## 🎯 Features

- **Modern Frontend** - Clean, responsive HTML/CSS/JavaScript interface
- **User Authentication** - Register and login with JWT tokens
- **Player Statistics** - Detailed player profiles with team history and fixtures
- **Team Analytics** - Team standings, player rosters, manager information
- **Tournament Management** - Full tournament standings and statistics
- **Match Details** - Complete match statistics including player and team performance
- **Search** - Global search across players, teams, tournaments, and managers
- **Favorites** - Save your favorite players, teams, tournaments, and managers
- **Analytics Service** - gRPC-based analytics with ClickHouse backend
- **Real-time Data** - Kafka-based event streaming with data collection

## 📋 Project Structure

```
sport_app/
├── services/
│   ├── frontend/          # Nginx + HTML/CSS/JS frontend
│   ├── auth_service/      # JWT authentication (Go)
│   ├── core_api/          # Main REST API (Go)
│   ├── analytics_service/ # Analytics service (Python gRPC)
│   └── data_collector/    # Data collection service (Go)
├── proto/                 # Protocol Buffer definitions
└── docker-compose.yml     # Docker Compose orchestration
```

## 🚀 Quick Start

### Prerequisites
- Docker
- Docker Compose

### Installation & Running

1. **Clone the repository:**
```bash
cd sport_app
```

2. **Start all services:**
```bash
docker-compose up -d
```

3. **Access the application:**
- Frontend: http://localhost
- API Documentation:
  - Core API: http://localhost:8080/api
  - Auth Service: http://localhost:8081
  - Analytics Service: http://localhost:8082

### Services Overview

| Service | Port | Description |
|---------|------|-------------|
| Frontend (Nginx) | 80 | Web interface |
| Core API | 8080 | Main REST API |
| Auth Service | 8081 | Authentication service |
| Analytics Service | 50051 (gRPC), 8082 (HTTP) | Analytics & statistics |
| Data Collector | - | Background data collection |
| PostgreSQL | 5432 | SQL database |
| ClickHouse | 8123, 9000 | Analytics database |
| Redis | 6379 | Caching layer |
| Elasticsearch | 9200 | Search engine |
| Kafka | 9092 | Message streaming |

## 📱 Frontend Pages

### Public Pages
- **Home** (`/`) - Main dashboard with recent matches and tournaments
- **Search** (`/search.html`) - Global search functionality
- **Player** (`/player.html?id=<id>`) - Player profile with stats
- **Team** (`/team.html?id=<id>`) - Team information and standings
- **Tournament** (`/tournament.html?id=<id>`) - Tournament table and statistics
- **Manager** (`/manager.html?id=<id>`) - Manager profile
- **Match** (`/match.html?id=<id>`) - Match details and statistics

### Protected Features (Requires Login)
- Save favorite players, teams, tournaments, and managers
- View favorites section on home page
- Persistent user preferences

## 🔌 API Endpoints

### Authentication
```
POST   /register          - User registration
POST   /login             - User login (returns JWT token)
GET    /me                - Get current user info
```

### Players
```
GET    /api/player/{id}/details       - Player info
GET    /api/player/{id}/stats         - Player statistics
GET    /api/player/{id}/fixtures      - Player match history
GET    /api/player/{id}/teams         - Teams history
GET    /api/player/favorite           - Get favorite players (protected)
POST   /api/player/{id}/favorite      - Add to favorites (protected)
DELETE /api/player/{id}/favorite      - Remove from favorites (protected)
```

### Teams
```
GET    /api/team/{id}/details         - Team info
GET    /api/team/{id}/stats           - Team statistics
GET    /api/team/{id}/next_game       - Next match
GET    /api/team/{id}/standings       - Tournament standings
GET    /api/team/{id}/players         - Team roster
GET    /api/team/{id}/fixtures        - Match history
GET    /api/team/{id}/manager         - Team manager
GET    /api/team/{id}/players_stats   - Players with stats
GET    /api/team/favorite             - Get favorite teams (protected)
POST   /api/team/{id}/favorite        - Add to favorites (protected)
DELETE /api/team/{id}/favorite        - Remove from favorites (protected)
```

### Tournaments
```
GET    /api/tournament/{id}/details             - Tournament info
GET    /api/tournament/{id}/table               - Standings table
GET    /api/tournament/{id}/stats/teams         - Teams statistics
GET    /api/tournament/{id}/stats/players       - Players statistics
GET    /api/tournament/{id}/table/graph         - Table graph data
GET    /api/tournament/{id}/fixtures            - Tournament matches
GET    /api/tournament/favorite                 - Get favorites (protected)
POST   /api/tournament/{id}/favorite            - Add to favorites (protected)
DELETE /api/tournament/{id}/favorite            - Remove from favorites (protected)
```

### Managers
```
GET    /api/manager/{id}/details      - Manager info
GET    /api/manager/{id}/stats        - Manager statistics
GET    /api/manager/{id}/teams        - Teams managed
GET    /api/manager/{id}/fixtures     - Match history
GET    /api/manager/favorite          - Get favorites (protected)
POST   /api/manager/{id}/favorite     - Add to favorites (protected)
DELETE /api/manager/{id}/favorite     - Remove from favorites (protected)
```

### Fixtures (Matches)
```
GET    /api/fixture/{id}/details      - Match details
GET    /api/fixture/{id}/stats/players - Player statistics
GET    /api/fixture/{id}/stats/goalies - Goalie statistics
GET    /api/fixture/{id}/stats/teams  - Team match statistics
GET    /api/fixture?date=YYYY-MM-DD  - Matches by date
```

### Search
```
GET    /api/search?q=<query>          - Global search
```

### Analytics
```
GET    /api/analytics/...             - Various analytics endpoints
```

## 🔐 Authentication

The platform uses JWT (JSON Web Tokens) for authentication:

1. **Register**: POST to `/register` with `username` and `password`
2. **Login**: POST to `/login` with credentials to receive a token
3. **Token Usage**: Include `Authorization: Bearer <token>` header in requests to protected endpoints

## 🛠️ Development

### Building Individual Services

**Frontend:**
```bash
cd services/frontend
docker build -t sports-frontend .
```

**Auth Service:**
```bash
cd services/auth_service
docker build -t sports-auth .
```

**Core API:**
```bash
cd services/core_api
docker build -t sports-api .
```

### Environment Variables

Create `.env` files in each service directory:

**auth_service/.env:**
```
PORT=8081
DATABASE_ADDR=postgres:5432
DB_NAME=sport_db
DB_USER=postgres
DB_PASS=postgres
JWT_SECRET=your_secret_key
```

**core_api/.env:**
```
PORT=8080
DATABASE_ADDR=postgres:5432
DB_NAME=sport_db
DB_USER=postgres
DB_PASS=postgres
REDIS_ADDR=redis:6379
ELASTICSEARCH_ADDR=http://elasticsearch:9200
ANALYTICS_SERVICE_ADDR=analytics_service:50051
```

## 📊 Frontend Technology Stack

- **HTML5** - Semantic markup
- **CSS3** - Modern responsive design with grid and flexbox
- **Vanilla JavaScript** - No dependencies, lightweight
- **Nginx** - Web server and reverse proxy
- **Docker** - Containerization

## 🎨 Design Features

- **Responsive Design** - Works on desktop, tablet, and mobile
- **Modern UI** - Clean, minimalist interface
- **Fast Loading** - Optimized assets and lazy loading
- **Accessible** - Semantic HTML and ARIA labels
- **Dark-friendly** - Light color scheme suitable for all backgrounds

## 🔄 Data Flow

1. **Frontend** makes HTTP requests to Nginx
2. **Nginx** routes to appropriate backend services
3. **Auth Service** handles user authentication and JWT tokens
4. **Core API** provides REST endpoints for data
5. **Analytics Service** (gRPC) serves advanced analytics
6. **Data Collector** continuously updates data from external sources
7. **Databases** store and retrieve data (PostgreSQL, ClickHouse)
8. **Cache** (Redis) improves performance
9. **Search Engine** (Elasticsearch) provides full-text search

## 📝 License

MIT License

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📧 Support

For support, please open an issue in the repository.

---

**Happy analyzing! ⚽📊**
