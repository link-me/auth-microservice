$ErrorActionPreference = 'Stop'

# Defaults
if (-not $env:JWT_SECRET) { $env:JWT_SECRET = 'dev-secret' }
if (-not $env:ACCESS_TTL) { $env:ACCESS_TTL = '900s' }
if (-not $env:REFRESH_TTL) { $env:REFRESH_TTL = '720h' }
if (-not $env:REDIS_URL)  { $env:REDIS_URL  = 'redis://localhost:6379' }
if (-not $env:ADDR)       { $env:ADDR       = ':8080' }

Write-Host "Starting Auth Microservice with:" -ForegroundColor Cyan
Write-Host " JWT_SECRET   = $env:JWT_SECRET"
Write-Host " ACCESS_TTL   = $env:ACCESS_TTL"
Write-Host " REFRESH_TTL  = $env:REFRESH_TTL"
Write-Host " REDIS_URL    = $env:REDIS_URL"
Write-Host " ADDR         = $env:ADDR"

go mod tidy
go run ./src