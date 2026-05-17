# Sports Analytics Platform - Deployment Guide

## Development Setup

### Quick Start (Local)

1. **Start with Docker Compose:**
```bash
# Linux/Mac
./start.sh

# Windows
start.bat
```

2. **Access the application:**
   - Frontend: http://localhost
   - Check services are running: `docker-compose ps`

3. **Stop services:**
```bash
# Linux/Mac
./stop.sh

# Windows
stop.bat
```

## Production Deployment

### Prerequisites
- Docker & Docker Compose
- Nginx (reverse proxy)
- SSL certificates (Let's Encrypt recommended)
- Domain name
- Server with at least 4GB RAM, 20GB disk

### Step 1: Prepare Server

```bash
# Update system
sudo apt-get update && sudo apt-get upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker
```

### Step 2: Clone Repository

```bash
git clone <your-repo-url> /opt/sport_app
cd /opt/sport_app
```

### Step 3: Configure Environment

Create `.env` files for each service:

**services/auth_service/.env:**
```
PORT=8081
DATABASE_ADDR=postgres:5432
DB_NAME=sport_db
DB_USER=postgres
DB_PASS=<strong_password_here>
JWT_SECRET=<generate_random_secret>
```

**services/core_api/.env:**
```
PORT=8080
DATABASE_ADDR=postgres:5432
DB_NAME=sport_db
DB_USER=postgres
DB_PASS=<strong_password_here>
REDIS_ADDR=redis:6379
ELASTICSEARCH_ADDR=http://elasticsearch:9200
ANALYTICS_SERVICE_ADDR=analytics_service:50051
```

**services/analytics_service/.env:**
```
CLICKHOUSE_HOST=clickhouse
CLICKHOUSE_PORT=9000
```

Generate secure JWT secret:
```bash
openssl rand -base64 32
```

### Step 4: Update docker-compose.yml

Modify services to use volume mounts for data persistence:

```yaml
volumes:
  postgres_data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /data/postgres
  clickhouse_data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /data/clickhouse
  elasticsearch_data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /data/elasticsearch
  redis_data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /data/redis
```

Create directories:
```bash
sudo mkdir -p /data/{postgres,clickhouse,elasticsearch,redis}
sudo chown -R $USER:$USER /data
```

### Step 5: Configure Nginx Reverse Proxy

**Create `/etc/nginx/sites-available/sports-analytics`:**

```nginx
upstream app_backend {
    server 127.0.0.1:80;
}

server {
    listen 80;
    listen [::]:80;
    server_name yourdomain.com www.yourdomain.com;

    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name yourdomain.com www.yourdomain.com;

    # SSL certificates
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Proxy to Docker
    location / {
        proxy_pass http://app_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_buffering off;
    }

    # Cache static assets
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Logging
    access_log /var/log/nginx/sports_access.log;
    error_log /var/log/nginx/sports_error.log;
}
```

Enable site:
```bash
sudo ln -s /etc/nginx/sites-available/sports-analytics /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### Step 6: Set Up SSL Certificates

```bash
sudo apt-get install certbot python3-certbot-nginx
sudo certbot certonly --nginx -d yourdomain.com -d www.yourdomain.com
```

### Step 7: Start Services

```bash
cd /opt/sport_app
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

### Step 8: Verify Deployment

```bash
# Check frontend
curl https://yourdomain.com

# Check API
curl https://yourdomain.com/api/search?q=test

# Check databases
docker-compose exec postgres psql -U postgres -d sport_db -c "\dt"
```

## Monitoring & Maintenance

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f core_api

# Last 100 lines
docker-compose logs --tail=100 frontend
```

### Health Checks

```bash
# Check all containers
docker-compose ps

# Inspect container
docker inspect <container_name>

# Check disk usage
docker system df
```

### Database Backups

**PostgreSQL backup:**
```bash
docker-compose exec postgres pg_dump -U postgres sport_db > backup_$(date +%Y%m%d).sql
```

**PostgreSQL restore:**
```bash
docker-compose exec -T postgres psql -U postgres sport_db < backup_20231215.sql
```

**ClickHouse backup:**
```bash
docker-compose exec clickhouse clickhouse-client --query "SELECT * FROM system.tables" > tables_backup.sql
```

### SSL Certificate Renewal

```bash
# Set up auto-renewal
sudo certbot renew --dry-run

# Manual renewal
sudo certbot renew --force-renewal
```

## Scaling & Performance

### Load Testing

```bash
# Install Apache Bench
sudo apt-get install apache2-utils

# Test frontend
ab -n 1000 -c 10 https://yourdomain.com/

# Test API
ab -n 1000 -c 10 https://yourdomain.com/api/search?q=test
```

### Database Optimization

**PostgreSQL:**
```bash
docker-compose exec postgres psql -U postgres -d sport_db -c "ANALYZE;"
docker-compose exec postgres psql -U postgres -d sport_db -c "REINDEX DATABASE sport_db;"
```

**Redis:**
```bash
docker-compose exec redis redis-cli INFO stats
docker-compose exec redis redis-cli DBSIZE
```

### Resource Limits

Update `docker-compose.yml`:

```yaml
services:
  core_api:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 256M
```

## Troubleshooting

### Service Won't Start
```bash
# Check logs
docker-compose logs <service_name>

# Rebuild
docker-compose build --no-cache <service_name>

# Clean and restart
docker-compose down -v
docker-compose up -d
```

### Database Connection Issues
```bash
# Check PostgreSQL
docker-compose exec postgres pg_isready

# Check connection
docker-compose exec core_api nc -zv postgres 5432
```

### High Memory Usage
```bash
# Check memory usage
docker stats

# Cleanup
docker system prune -a --volumes
```

### API Slow Response
```bash
# Check Redis cache
docker-compose exec redis redis-cli ping

# Monitor database
docker-compose exec postgres psql -U postgres -d sport_db -c "SELECT * FROM pg_stat_statements LIMIT 10;"
```

## Security Checklist

- [ ] Change all default passwords
- [ ] Enable SSL/TLS certificates
- [ ] Set up firewalls (UFW)
- [ ] Configure fail2ban for DDoS protection
- [ ] Regular security updates
- [ ] Enable database backups
- [ ] Monitor access logs
- [ ] Set up rate limiting in Nginx
- [ ] Use strong JWT secrets
- [ ] Enable CORS restrictions

## Upgrade Instructions

```bash
# Backup data
docker-compose exec postgres pg_dump -U postgres sport_db > backup.sql

# Pull latest code
git pull origin main

# Rebuild images
docker-compose build --no-cache

# Restart services
docker-compose up -d

# Verify
docker-compose ps
```

## Support

For issues and questions:
1. Check logs: `docker-compose logs <service>`
2. Review environment variables
3. Verify network connectivity
4. Check resource availability
5. Open an issue on GitHub

---

Happy deploying! 🚀
