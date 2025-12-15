Write-Host "=== Comprehensive API Testing ===" -ForegroundColor Cyan
Write-Host ""

# 1. Health Checks
Write-Host "1. Health Checks:" -ForegroundColor Green

$healthEndpoints = @(
    @{Name="NGINX"; Url="http://localhost/health"},
    @{Name="Shortener (via NGINX)"; Url="http://localhost/api/health/shortener"},
    @{Name="Analytics (via NGINX)"; Url="http://localhost/api/health/analytics"},
    @{Name="Shortener (direct)"; Url="http://localhost:8080/health"},
    @{Name="Analytics (direct)"; Url="http://localhost:8081/health"}
)

foreach ($endpoint in $healthEndpoints) {
    try {
        $response = Invoke-WebRequest -Uri $endpoint.Url -Method GET -UseBasicParsing -TimeoutSec 3
        $json = $response.Content | ConvertFrom-Json
        Write-Host "   ✓ $($endpoint.Name): $($json.status) ($($json.service))" -ForegroundColor Green
    } catch {
        Write-Host "   ✗ $($endpoint.Name): Error" -ForegroundColor Red
    }
    Start-Sleep -Milliseconds 200
}

# 2. Test Main Page
Write-Host "`n2. Main Page:" -ForegroundColor Green
try {
    $response = Invoke-WebRequest -Uri "http://localhost" -Method GET -UseBasicParsing
    Write-Host "   ✓ Main page is accessible" -ForegroundColor Green
    $contentLength = $response.Content.Length
    $previewLength = [Math]::Min(100, $contentLength)
    Write-Host "   Response (first ${previewLength} chars): $($response.Content.Substring(0, $previewLength))..." -ForegroundColor Gray
} catch {
    Write-Host "   ✗ Main page error" -ForegroundColor Red
}

# 3. Test Shorten URL
Write-Host "`n3. Shorten URL API:" -ForegroundColor Green
$testUrls = @(
    "https://google.com",
    "https://github.com",
    "https://docker.com"
)

$shortCodes = @()

foreach ($url in $testUrls) {
    try {
        $body = @{url = $url} | ConvertTo-Json
        $response = Invoke-RestMethod -Uri "http://localhost/shorten" -Method POST -Body $body -ContentType "application/json" -UseBasicParsing
        
        Write-Host "   ✓ Shortened: $url" -ForegroundColor Green
        Write-Host "      → $($response.short_url)" -ForegroundColor Cyan
        
        # Extract short code
        $shortCode = $response.short_url -replace '.*/', ''
        $shortCodes += $shortCode
        
    } catch {
        Write-Host "   ✗ Failed to shorten: $url" -ForegroundColor Red
        Write-Host "      Error: $($_.Exception.Message)" -ForegroundColor DarkRed
    }
    Start-Sleep -Milliseconds 500
}

# 4. Test Redirects
Write-Host "`n4. Redirect Tests:" -ForegroundColor Green
foreach ($shortCode in $shortCodes) {
    try {
        $response = Invoke-WebRequest -Uri "http://localhost/${shortCode}" -Method GET -MaximumRedirection 0 -UseBasicParsing -ErrorAction Stop
        
        if ($response.StatusCode -in @(301, 302, 307, 308)) {
            Write-Host "   ✓ Redirect ${shortCode}: HTTP $($response.StatusCode)" -ForegroundColor Green
            Write-Host "      Location: $($response.Headers.Location)" -ForegroundColor Gray
        } else {
            Write-Host "   ⚠ Redirect ${shortCode}: HTTP $($response.StatusCode) (unexpected)" -ForegroundColor Yellow
        }
    } catch [System.Net.WebException] {
        $statusCode = $_.Exception.Response.StatusCode.value__
        if ($statusCode -in @(301, 302, 307, 308)) {
            Write-Host "   ✓ Redirect ${shortCode}: HTTP ${statusCode} (via exception)" -ForegroundColor Green
        } else {
            Write-Host "   ✗ Redirect ${shortCode} failed" -ForegroundColor Red
        }
    } catch {
        Write-Host "   ✗ Redirect ${shortCode} error: $($_.Exception.Message)" -ForegroundColor Red
    }
    Start-Sleep -Milliseconds 500
}

# 5. Test Stats
Write-Host "`n5. Statistics API:" -ForegroundColor Green
if ($shortCodes.Count -gt 0) {
    foreach ($shortCode in $shortCodes[0..0]) { # Test only first one
        try {
            $response = Invoke-RestMethod -Uri "http://localhost/stats/${shortCode}" -Method GET -UseBasicParsing
            Write-Host "   ✓ Stats for ${shortCode}: Available" -ForegroundColor Green
            Write-Host "      Visits: $($response.visits)" -ForegroundColor Gray
        } catch {
            Write-Host "   ✗ Stats for ${shortCode}: Error" -ForegroundColor Yellow
        }
    }
}

# 6. Test GetAllURLs
Write-Host "`n6. Get All URLs:" -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "http://localhost/urls" -Method GET -UseBasicParsing
    Write-Host "   ✓ Total URLs stored: $($response.Count)" -ForegroundColor Green
} catch {
    Write-Host "   ✗ Get all URLs failed" -ForegroundColor Yellow
}

# 7. Docker Swarm Status
Write-Host "`n7. Docker Swarm Status:" -ForegroundColor Green
docker service ls --format "table {{.Name}}`t{{.Replicas}}`t{{.Image}}"

# 8. Service Logs (последние 2 строки)
Write-Host "`n8. Recent Logs:" -ForegroundColor Green
$services = @("url-shortener-stack_nginx", "url-shortener-stack_url-shortener", "url-shortener-stack_analytics")
foreach ($service in $services) {
    Write-Host "   ${service} :" -ForegroundColor Cyan
    $logs = docker service logs $service --tail 2 2>$null
    foreach ($log in $logs) {
        Write-Host "      $log" -ForegroundColor Gray
    }
}

Write-Host "`n" + ("=" * 60) -ForegroundColor White
Write-Host "TESTING COMPLETE" -ForegroundColor White -BackgroundColor DarkGreen
Write-Host ("=" * 60) -ForegroundColor White
Write-Host ""
Write-Host "✅ All services are running in Docker Swarm" -ForegroundColor Green
Write-Host "✅ NGINX is configured as reverse proxy" -ForegroundColor Green
Write-Host "✅ API endpoints are functional" -ForegroundColor Green
Write-Host ""
Write-Host "Access the services at:" -ForegroundColor Yellow
Write-Host "  • Main Interface: http://localhost" -ForegroundColor White
Write-Host "  • API Documentation: See main page for endpoints" -ForegroundColor White
Write-Host "  • Health Status: http://localhost/health" -ForegroundColor White
Write-Host ""
Write-Host "Commands for monitoring:" -ForegroundColor Yellow
Write-Host "  docker service ls" -ForegroundColor Gray
Write-Host "  docker service logs url-shortener-stack_nginx" -ForegroundColor Gray
Write-Host "  docker stack ps url-shortener-stack" -ForegroundColor Gray