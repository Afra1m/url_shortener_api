Write-Host "=== Simple Test ===" -ForegroundColor Green

# 1. Test health
Write-Host "`nTesting health endpoints..." -ForegroundColor Yellow
try {
    $health1 = Invoke-RestMethod -Uri "http://localhost:8080/health"
    Write-Host "Shortener: OK - $($health1.status)" -ForegroundColor Green
} catch {
    Write-Host "Shortener: FAILED" -ForegroundColor Red
}

try {
    $health2 = Invoke-RestMethod -Uri "http://localhost:8081/health"
    Write-Host "Analytics: OK - $($health2.status)" -ForegroundColor Green
} catch {
    Write-Host "Analytics: FAILED" -ForegroundColor Red
}

# 2. Create one URL
Write-Host "`nCreating short URL..." -ForegroundColor Yellow
try {
    $body = @{original_url = "https://google.com"} | ConvertTo-Json
    $response = Invoke-RestMethod -Uri "http://localhost:8080/shorten" -Method Post -Body $body -ContentType "application/json"
    Write-Host "Created: $($response.short_url)" -ForegroundColor Cyan
    
    $shortCode = $response.short_url.Split('/')[-1]
    Write-Host "Short code: $shortCode" -ForegroundColor Cyan
    
    # 3. Get stats
    $stats = Invoke-RestMethod -Uri "http://localhost:8080/stats/$shortCode" -Method Get
    Write-Host "Stats: $($stats.access_count) accesses" -ForegroundColor Green
    
    # 4. Try analytics
    Start-Sleep -Seconds 1
    try {
        $analytics = Invoke-RestMethod -Uri "http://localhost:8081/summary/$shortCode" -Method Get
        Write-Host "Analytics: $($analytics.total_clicks) clicks" -ForegroundColor Magenta
    } catch {
        Write-Host "Analytics not ready yet (normal for first run)" -ForegroundColor Yellow
    }
    
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
}

Write-Host "`n=== Test Complete ===" -ForegroundColor Green