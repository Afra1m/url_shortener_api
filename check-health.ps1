# check-health.ps1
Write-Host "=== Checking Shortener Health ===" -ForegroundColor Yellow

Write-Host "`nMethod 1: Invoke-RestMethod" -ForegroundColor Cyan
try {
    $result = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get
    Write-Host "Success: " -NoNewline -ForegroundColor Green
    $result
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`nMethod 2: Invoke-WebRequest" -ForegroundColor Cyan
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -Method Get
    Write-Host "Status Code: $($response.StatusCode)" -ForegroundColor Green
    Write-Host "Content: $($response.Content)" -ForegroundColor Green
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`nMethod 3: Using curl.exe" -ForegroundColor Cyan
try {
    $curlResult = curl.exe -s http://localhost:8080/health
    Write-Host "Result: $curlResult" -ForegroundColor Green
} catch {
    Write-Host "Curl not available or failed" -ForegroundColor Yellow
}

Write-Host "`n=== Also checking Analytics ===" -ForegroundColor Yellow
try {
    $analyticsHealth = Invoke-RestMethod -Uri "http://localhost:8081/health" -Method Get
    Write-Host "Analytics: $($analyticsHealth.status)" -ForegroundColor Green
} catch {
    Write-Host "Analytics check failed" -ForegroundColor Red
}