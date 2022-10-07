# Hoist

Hoist is a simple tool for transferring large files or directories over a Local Area Network (LAN). 

## Installation

### Debian / Ubuntu 
```
wget https://github.com/aiden-deloryn/Hoist/releases/download/v1.0.0/hoist_1.0.0_amd64.deb \
&& sudo apt install -y ./hoist_1.0.0_amd64.deb
```

### Windows (PowerShell)
```
Invoke-WebRequest https://github.com/aiden-deloryn/Hoist/releases/download/v1.0.0/hoist_1.0.0_amd64.exe -OutFile hoist.exe
```

### Build form source (Linux)
```
git clone https://github.com/aiden-deloryn/Hoist.git
cd Hoist/
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./hoist ./src/main.go
sudo mv ./hoist /usr/local/bin/
```

## Basic usage

Send a directory using hoist:

```
$ hoist send "Project Resources"
Enter a password: 
The target file or directory is ready to send. To download it on another machine, use:
  hoist get 127.0.0.1:8080
```

Download the directory using hoist:

```
$ hoist get 127.0.0.1:8080
Enter password: 
Copying file Project Resources/README.md...
|========100%========| 769/769 bytes (0 MiB/s)
Copying file Project Resources/Summary.pdf...
|========100%========| 738240/738240 bytes (0 MiB/s)
Copying file Project Resources/logo.png...
|========100%========| 230700/230700 bytes (0 MiB/s)
```