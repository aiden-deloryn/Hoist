# Hoist

Hoist is a simple tool for transferring large files or directories over a Local Area Network (LAN). 

## Installation

### Debian / Ubuntu 
```
VERSION='1.3.0' ARCH='amd64' sh -c 'wget https://github.com/aiden-deloryn/Hoist/releases/download/v${VERSION}/hoist_${VERSION}_${ARCH}.deb \
&& sudo apt install -y ./hoist_${VERSION}_${ARCH}.deb'
```

### Windows (PowerShell)
```
Invoke-WebRequest https://raw.githubusercontent.com/aiden-deloryn/Hoist/v1.3.0/scripts/install.ps1 -OutFile install.ps1
Set-ExecutionPolicy Bypass -Scope Process -Force; .\install.ps1 -version "1.3.0" -arch "amd64"
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
$ hoist send "./Pictures/Funny Cat Photos"
Enter a password: 
The target file or directory is ready to send. To download it on another machine, use:
  hoist get 127.0.0.1:47478
```

Download the directory using hoist:

```
$ hoist get 127.0.0.1:47478
Enter password: 
Copying file Funny Cat Photos/Cat's wearing hats/Cats-Wearing-Hats-social.jpg...
|========100%========| 270249/270249 bytes (0 MiB/s)
Copying file Funny Cat Photos/Cat's wearing hats/cat-wearing-hat-12408151.jpg...
|========100%========| 112350/112350 bytes (0 MiB/s)
Copying file Funny Cat Photos/Cat's wearing hats/cat_wearing_hat_with_ears.jpg...
|========100%========| 32741/32741 bytes (0 MiB/s)
Copying file Funny Cat Photos/cat_looking_shocked.jpeg...
|========100%========| 5871/5871 bytes (0 MiB/s)
Copying file Funny Cat Photos/cat_on_skateboard.jpeg...
|========100%========| 6020/6020 bytes (0 MiB/s)
Copying file Funny Cat Photos/cat_wearing_glasses.jpeg...
|========100%========| 8364/8364 bytes (0 MiB/s)
Copying file Funny Cat Photos/grumpy-cat-meme-of-not-enjoying-a-morning-at-all.jpeg...
|========100%========| 102247/102247 bytes (0 MiB/s)
```