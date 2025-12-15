Write-Host "=== Redeploying with fixed configuration ===" -ForegroundColor Green

# 1. Удаляем старый стек
Write-Host "1. Removing old stack..." -ForegroundColor Yellow
docker stack rm url-shortener-stack
Start-Sleep -Seconds 10

# 2. Пересобираем NGINX с исправленной конфигурацией
Write-Host "2. Rebuilding NGINX image..." -ForegroundColor Yellow
docker build -t nginx-fixed:latest ./nginx

# 3. Проверяем, что образы shortener и analytics существуют
Write-Host "3. Checking existing images..." -ForegroundColor Yellow
docker images | findstr "shortener analytics"

# 4. Если образов нет, собираем их
if (-not (docker images --format "{{.Repository}}" | findstr "shortener")) {
    Write-Host "   Building shortener image..." -ForegroundColor Cyan
    docker build -t shortener:latest ./shortener
}

if (-not (docker images --format "{{.Repository}}" | findstr "analytics")) {
    Write-Host "   Building analytics image..." -ForegroundColor Cyan
    docker build -t analytics:latest ./analytics
}

# 5. Развертываем новый стек
Write-Host "4. Deploying new stack..." -ForegroundColor Yellow
docker stack deploy -c docker-stack.yml url-shortener-stack

# 6. Ждем и проверяем
Write-Host "5. Waiting for services to start..." -ForegroundColor Yellow
Start-Sleep -Seconds 15

# 7. Проверяем статус
Write-Host "6. Checking service status..." -ForegroundColor Green
docker service ls

Write-Host "`n=== Testing ===" -ForegroundColor Green
Write-Host "Testing NGINX health endpoint in 10 seconds..." -ForegroundColor Cyan
Start-Sleep -Seconds 10

try {
    $response = Invoke-WebRequest -Uri "http://localhost/health" -Method GET -ErrorAction SilentlyContinue
    Write-Host "NGINX is healthy: $($response.Content)" -ForegroundColor Green
} catch {
    Write-Host "NGINX is not responding yet. Checking logs..." -ForegroundColor Yellow
    docker service logs url-shortener-stack_nginx --tail 5
}

Write-Host "`nAccess URLs:" -ForegroundColor Cyan
Write-Host "  - Main page: http://localhost" -ForegroundColor White
Write-Host "  - NGINX health: http://localhost/health" -ForegroundColor White
Write-Host "  - Service health: http://localhost/api/health/shortener" -ForegroundColor White