# Continuous Seed
Continuous delivery tool CI/CD

```
GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o bin/seed-win-amd64.exe src/main.go;
GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/seed-linux-amd64 src/main.go; 
GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o bin/seed-mac-amd64 src/main.go;
```