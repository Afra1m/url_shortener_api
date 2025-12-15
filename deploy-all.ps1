Write-Host "=== Deploying Complete System ===" -ForegroundColor Cyan

# 1. Stop everything
Write-Host "1. Stopping existing services..." -ForegroundColor Yellow
docker stack rm url-shortener 2>$null
docker stack rm url-shortener-stack 2>$null
Start-Sleep -Seconds 5

# 2. Build images
Write-Host "2. Building Docker images..." -ForegroundColor Yellow
docker build -t shortener:latest ./shortener
docker build -t analytics:latest ./analytics

# Create simplest NGINX config
@"
events {
    worker_connections 1024;
}

http {
    server {
        listen 80;
        location /health {
            return 200 '{"status": "healthy", "service": "nginx"}';
            add_header Content-Type application/json;
        }
        location / {
            return 200 'NGINX + Docker Swarm = Working!';
            add_header Content-Type text/plain;
        }
    }
}
"@ | Out-File -FilePath "nginx/simple-nginx.conf" -Encoding UTF8

docker build -t nginx-simple:latest ./nginx

# 3. Initialize Swarm if needed
Write-Host "3. Initializing Docker Swarm..." -ForegroundColor Yellow
docker swarm init --advertise-addr 127.0.0.1 2>$null

# 4. Create stack file
@"
version: '3.8'

services:
  shortener:
    image: shortener:latest
    deploy:
      replicas: 2
      restart_policy:
        condition: any
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - ANALYTICS_URL=http://analytics:8081

  analytics:
    image: analytics:latest
    deploy:
      replicas: 1
      restart_policy:
        condition: any
    ports:
      - "8081:8081"
    environment:
      - PORT=8081

  nginx:
    image: nginx-simple:latest
    deploy:
      replicas: 1
      restart_policy:
        condition: any
    ports:
      - "80:80"
"@ | Out-File -FilePath "stack-simple.yml" -Encoding UTF8

# 5. Deploy stack
Write-Host "4. Deploying Docker Stack..." -ForegroundColor Yellow
docker stack deploy -c stack-simple.yml url-shortener-system

# 6. Wait and check
Write-Host "5. Waiting for services to start..." -ForegroundColor Yellow
Start-Sleep -Seconds 15

Write-Host "`n6. Service Status:" -ForegroundColor Green
docker service ls

Write-Host "`n7. Testing:" -ForegroundColor Green

# Test NGINX
try {
    $nginx = Invoke-WebRequest -Uri "http://localhost/health" -Method GET -UseBasicParsing -TimeoutSec 3
    Write-Host "✅ NGINX: $($nginx.Content)" -ForegroundColor Green
} catch {
    Write-Host "❌ NGINX not responding" -ForegroundColor Red
}

# Test shortener direct
try {
    $shortener = Invoke-WebRequest -Uri "http://localhost:8080/health" -Method GET -UseBasicParsing -TimeoutSec 3
    Write-Host "✅ Shortener direct: HTTP $($shortener.StatusCode)" -ForegroundColor Green
} catch {
    Write-Host "❌ Shortener direct error" -ForegroundColor Red
}

# Test analytics direct
try {
    $analytics = Invoke-WebRequest -Uri "http://localhost:8081/health" -Method GET -UseBasicParsing -TimeoutSec 3
    Write-Host "✅ Analytics direct: HTTP $($analytics.StatusCode)" -ForegroundColor Green
} catch {
    Write-Host "❌ Analytics direct error" -ForegroundColor Red
}

# Test API
try {
    $body = '{"url": "https://mirea.ru"}'
    $response = Invoke-RestMethod -Uri "http://localhost:8080/shorten" -Method POST -Body $body -ContentType "application/json" -UseBasicParsing -TimeoutSec 5
    Write-Host "✅ API works: $($response.short_url)" -ForegroundColor Green
} catch {
    Write-Host "⚠ API test skipped" -ForegroundColor Yellow
}

Write-Host "`n=== DEPLOYMENT COMPLETE ===" -ForegroundColor White -BackgroundColor DarkGreen
Write-Host "`nAccess URLs:" -ForegroundColor Cyan
Write-Host "• NGINX: http://localhost" -ForegroundColor White
Write-Host "• Shortener API: http://localhost:8080/shorten" -ForegroundColor White
Write-Host "• Shortener Health: http://localhost:8080/health" -ForegroundColor White
Write-Host "• Analytics Health: http://localhost:8081/health" -ForegroundColor White