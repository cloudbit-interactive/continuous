# Continuous Seed
<p align="center">
  <img src="https://github.com/cloudbit-interactive/continuous-seed/blob/main/seed.png?raw=true" height="100" title="hover text">
</p>
<p align="center">
Light tool for continuous delivery CI/CD without extra dependencies.
</p>

# Generate Binaries
```
GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o bin/seed-win-amd64.exe src/main.go;
GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/seed-linux-amd64 src/main.go; 
GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o bin/seed-mac-amd64 src/main.go;
GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o bin/seed-mac-arm64 src/main.go;
lipo -create -output bin/seed-mac-universal bin/seed-mac-amd64 bin/seed-mac-arm64;
rm bin/seed-mac-amd64; 
rm bin/seed-mac-arm64;
```
