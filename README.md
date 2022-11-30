# Continuous Seed
<p style="text-align: center">
  <img src="https://github.com/cloudbit-interactive/continuous-seed/blob/main/seed.png?raw=true" height="100" title="hover text">
</p>
<p style="text-align: center">
	Light tool for continuous delivery CI/CD without external dependencies.
</p>

# Yaml File Example

```yaml
vars:
  name: SeedApp
  environment: staging
  port: 4000
  branch: main
  projectDir: /path/to/project/
  os: ${os} # replaced with: windows, linux, mac, etc...
  arch: ${arch} #	replaced with: amd64, arm64, etc...
  date: ${date} # replaced with: yyyy-mm-dd
  datetime: ${datetime} # replaced with: yyyy-mm-dd 00:00:00
tasks:
  my-first-task: &my-first-task-anchor
		echo: device=${os} # output: device=mac
jobs:
	- cmd: echo 'Hello Seed' # output: Hello Seed 
  - execute-first-task: *my-first-task-anchor
	-	other-job: 
```

# Generate Binaries
```
GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o bin/seed-windows-amd64.exe src/main.go;
GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/seed-linux-amd64 src/main.go; 
GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o bin/seed-mac-amd64 src/main.go;
GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o bin/seed-mac-arm64 src/main.go;
lipo -create -output bin/seed-mac-universal bin/seed-mac-amd64 bin/seed-mac-arm64;
rm bin/seed-mac-amd64; 
rm bin/seed-mac-arm64;
```
