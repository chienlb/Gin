#!/usr/bin/env pwsh

# Script to run the application locally with environment configuration

$ErrorActionPreference = "Stop"

# Load environment variables from local.env
if (Test-Path "configs/local.env") {
    Write-Host "Loading environment variables from configs/local.env..." -ForegroundColor Green
    Get-Content "configs/local.env" | ForEach-Object {
        if ($_ -match "^([^=]+)=(.*)$") {
            $key = $matches[1]
            $value = $matches[2]
            [Environment]::SetEnvironmentVariable($key, $value)
            Write-Host "  Set $key" -ForegroundColor Cyan
        }
    }
} else {
    Write-Host "Warning: configs/local.env not found. Using default values." -ForegroundColor Yellow
}

# Check if Go is installed
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "Error: Go is not installed or not in PATH" -ForegroundColor Red
    exit 1
}

Write-Host "Starting Gin application..." -ForegroundColor Green

# Build and run
go run ./cmd/api

if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Application failed to start" -ForegroundColor Red
    exit 1
}
