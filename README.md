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
