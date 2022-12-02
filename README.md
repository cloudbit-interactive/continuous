# Continuous Seed
<p align="center">
  <img src="https://github.com/cloudbit-interactive/continuous-seed/blob/main/seed.png?raw=true" height="100" title="hover text">
</p>
<p>
	Light tool for continuous delivery CI/CD that works in the own server/container without external dependencies, runtimes, and almost no consume CPU/RAM.
</p>
<p>
	The original idea of this project is to provide just a seed with internal instructions to create, configure and update itself automatically without any intervention.
</p>
<p>
Check <strong>./example</strong> folder for a more real example
</p>

### Folder structure
```shell
.
└── projectDir
    ├── seed.yaml
    └── seed # the corresponding binary to use (windows, linux or mac)
```

### Current Implementation

<ul>
	<li>Variables</li>
	<li>Dynamic variables</li>
	<li>System variables</li>
	<li>Echo</li>
	<li>Cmd</li>
	<li>Loop</li>
	<li>If</li>
	<li>Kill Port</li>
	<li>Create Folder</li>
	<li>Create File</li>
	<li>Exist</li>
	<li>Delete</li>
</ul>

### Future

<ul>
	<li>For</li>
	<li>Expand If</li>
	<li>Expand Loop</li>
</ul>

### Yaml File Documentation

```yaml
name: Yaml File Example
version: 0.0.1
log: true # check log/ folder: false > only print echo commands, true > print all output in console on log/
vars:
  name: SeedApp
  environment: staging
  port: 4000
  branch: main
  projectDir: /Users/tufikchediak/seed-test
  systemOS: ${os} # replaced with: windows, linux, mac, etc...
  arch: ${arch} # replaced with: amd64, arm64, etc...
  date: ${date} # replaced with: yyyy-mm-dd
  datetime: ${datetime} # replaced with: yyyy-mm-dd 00:00:00
tasks:
  basic-commands: &basic-commands-anchor
    - echo: device=${systemOS} # echo command printing a string with a variable defined  above > output: device=mac
    - echo: ${os} # use a system pre-defined value > output: mac
    - echo: previous output=${output} # print the output returned on the previous line > output: previous output=mac
    - cmd: echo $RANDOM # execute > echo $RANDOM > output: RANDOM_NUMBER
    - var-my-first-variable: ${output} # define a personal variable using var-NAME, in this case is storing the output of the previews command
    - echo: my-first-variable=${my-first-variable} # printing the variable created previously > output: my-first-variable=20208
  other-commands: &other-commands-anchor
    - create-folder: ${projectDir}/folder1/folder2 # create folder
    - create-file: # create a file in a specific folder, if folders doesn't exist it will be auto-created
        file:  ${projectDir}/folder1/folder2/config.js
        content: |
          export const config = {
            name:"${name}",
            environment:"${environment}",
            port:${port}
          }
    - exist: ${projectDir}/folder1/folder2/config.js # check if folder or file exist > output: true, false
    - echo: file_exist=${output} # > output: file exist=true
    - delete: ${projectDir}/tmp/ # delete folder or file
    - kill-port: ${port} # kill process executing in a specific port, for multiple ports use comma separator > 4000, 8080, 80, 443
    - cmd: echo 1 # execute a command. It uses bash -c in linux, mac and powershell in windows
    - cmd: # another way to execute command is passing the app, args, workingDirectory, separator and background params
        app: npm
        args: install
        workingDirectory: ${projectDir}/continuous-seed-test/
        # separator: # args separator, default = ' '
        # background: true # execute command in background to prevent blocking main job
    - echo: ${os}
    - if: # validate the previous output and execute jobs if condition apply
        type: equal # equal, equal! (not equal), contain, contain! (no contain)
        value: mac
        jobs:
          - echo: this is mac
          - cmd: sw_vers -productVersion # get mac version > output: x.x.x
    - loop: # execute a loop each x time in milliseconds
        interval: 5000
        jobs:
          - cmd: echo $RANDOM
    #- stop: true # this command stop the execution
jobs:
  - execute-basic-task: *basic-commands-anchor
  - execute-other-job: *other-commands-anchor
```

### Generate Binaries
```
GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o bin/seed-windows-amd64.exe src/main.go;
GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/seed-linux-amd64 src/main.go; 
GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o bin/seed-mac-amd64 src/main.go;
GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o bin/seed-mac-arm64 src/main.go;
lipo -create -output bin/seed-mac-universal bin/seed-mac-amd64 bin/seed-mac-arm64;
rm bin/seed-mac-amd64; 
rm bin/seed-mac-arm64;
```
