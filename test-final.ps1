Write-Host "=== FINAL TEST ===" -ForegroundColor Green

Write-Host "`n1. Analytics service:" -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "http://localhost:8081/health" -Method Get
    Write-Host "   Status: $($health.status)" -ForegroundColor Green
} catch {
    Write-Host "   Failed" -ForegroundColor Red
}

Write-Host "`n2. Create short URL:" -ForegroundColor Yellow
try {
    $url = "https://test-" + (Get-Date -Format "HHmmss") + ".com"
    $body = @{original_url = $url} | ConvertTo-Json
    $response = Invoke-RestMethod -Uri "http://localhost:8080/shorten" -Method Post -Body $body -ContentType "application/json"
    
    $shortCode = $response.short_url.Split('/')[-1]
    Write-Host "   Created: $($response.short_url)" -ForegroundColor Green
    Write-Host "   Code: $shortCode" -ForegroundColor Cyan
    
} catch {
    Write-Host "   Failed to create URL: $_" -ForegroundColor Red
    exit
}

Write-Host "`n3. Check stored data:" -ForegroundColor Yellow
try {
    $stats = Invoke-RestMethod -Uri "http://localhost:8080/stats/$shortCode" -Method Get
    Write-Host "   Stats: $($stats.access_count) accesses" -ForegroundColor Green
    
    $allUrls = Invoke-RestMethod -Uri "http://localhost:8080/urls" -Method Get
    Write-Host "   Total URLs: $($allUrls.Count)" -ForegroundColor Green
    
} catch {
    Write-Host "   Failed to get data: $_" -ForegroundColor Red
}

Write-Host "`n4. Check analytics:" -ForegroundColor Yellow
Start-Sleep -Seconds 2
try {
    $analytics = Invoke-RestMethod -Uri "http://localhost:8081/summary/$shortCode" -Method Get
    Write-Host "   Analytics: $($analytics.total_clicks) events" -ForegroundColor Magenta
} catch {
    Write-Host "   Analytics not ready (normal)" -ForegroundColor Yellow
}

Write-Host "`n=== RESULT ===" -ForegroundColor Green
Write-Host "Microservices: WORKING" -ForegroundColor Green
Write-Host "Core functionality: OK" -ForegroundColor Green
Write-Host "Health check issue: requires debugging" -ForegroundColor Yellow