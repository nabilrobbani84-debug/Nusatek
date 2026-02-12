param (
    [string]$Command
)

switch ($Command) {
    "run" { go run cmd/api/main.go }
    "build" { go build -o bin/nusatek-backend.exe cmd/api/main.go }
    "test" { go test -v ./... }
    "clean" { if (Test-Path bin) { Remove-Item bin -Recurse -Force } }
    "docker-up" { docker-compose up -d }
    "docker-down" { docker-compose down }
    Default { 
        Write-Host "Usage: .\tasks.ps1 [command]"
        Write-Host "Commands:"
        Write-Host "  run         - Run the application"
        Write-Host "  build       - Build the binary"
        Write-Host "  test        - Run unit tests"
        Write-Host "  docker-up   - Start database and queue"
        Write-Host "  docker-down - Stop database and queue"
    }
}
