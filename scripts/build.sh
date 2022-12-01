#!/bin/bash

set -e

VERSION=$1
VERSION_REGEX="^[0-9]+[.][0-9]+[.][0-9]+$"

# Check prerequisites
if ! test -f "go.mod"; then
    >&2 echo "Error: This script must be run from the project root directory."
    exit 1
fi

if [[ ! $VERSION =~ $VERSION_REGEX ]]; then
    >&2 echo "Error: Version '${VERSION}' is not valid."
    exit 1
fi

if ! grep -q $VERSION "./src/values/values.go"; then
    >&2 echo "Error: Version '${VERSION}' does not match ./src/values/values.go"
    exit 1
fi

if ! grep -q $VERSION "./README.md"; then
    >&2 echo "Error: Please update README.md for version '${VERSION}'"
    exit 1
fi

# Create the required directory structure to create the debian packages
mkdir -p ./debpkgs/hoist_${VERSION}_amd64/DEBIAN
mkdir -p ./debpkgs/hoist_${VERSION}_amd64/usr/bin

mkdir -p ./debpkgs/hoist_${VERSION}_arm64/DEBIAN
mkdir -p ./debpkgs/hoist_${VERSION}_arm64/usr/bin

# Create a bin folder for Windows executables
mkdir -p ./bin


# Build the binary from source
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./debpkgs/hoist_${VERSION}_amd64/usr/bin/hoist ./src/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./debpkgs/hoist_${VERSION}_arm64/usr/bin/hoist ./src/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/hoist_${VERSION}_amd64.exe ./src/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o ./bin/hoist_${VERSION}_arm64.exe ./src/main.go

# Generate a metadata file for the debian package
cat >./debpkgs/hoist_${VERSION}_amd64/DEBIAN/control <<EOL
Package: hoist
Version: ${VERSION}
Architecture: amd64
Essential: no
Priority: optional
Depends:
Maintainer: Aiden De Loryn
Description: Hoist is a simple tool for transferring large files or directories over a Local Area Network (LAN).
EOL

cat >./debpkgs/hoist_${VERSION}_arm64/DEBIAN/control <<EOL
Package: hoist
Version: ${VERSION}
Architecture: arm64
Essential: no
Priority: optional
Depends:
Maintainer: Aiden De Loryn
Description: Hoist is a simple tool for transferring large files or directories over a Local Area Network (LAN).
EOL

# Build the debian package
dpkg-deb --build ./debpkgs/hoist_${VERSION}_amd64
dpkg-deb --build ./debpkgs/hoist_${VERSION}_arm64

# Tag the release using git
git tag -a v${VERSION} -m "Release version ${VERSION}"
git push origin v${VERSION}
git log --oneline --decorate --graph -n 5
