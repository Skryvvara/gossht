<div id="top"></div>

<br />
<div align="center">
  <a href="https://github.com/skryvvara/gossht">
      <img src="./.github/assets/gossht.svg" width=124 height=124 style="border-radius: 100%;" alt="G(h)ossht Logo">
  </a>

  <h3 align="center">G(h)ossht - A TUI for ssh written in Go</h3>
</div>

## Description

> ! This is a very early proof of concept

Gossht is a graphical ssh helper written in [Go](https://go.dev) using [tview](https://github.com/rivo/tview).

## Build

You can build gossht using make, all the dependencies are in the vendor directory so not even an internet
connection is required to build gossht.

First, verify that go and make is installed

```sh
go version

# example
#❯ go version
#go version go1.22.5 darwin/arm64

make --version

# example
#❯ make --version
#GNU Make 3.81
#Copyright (C) 2006  Free Software Foundation, Inc.
#This is free software; see the source for copying conditions.
#There is NO warranty; not even for MERCHANTABILITY or FITNESS FOR A
#PARTICULAR PURPOSE.
#
#This program built for i386-apple-darwin11.3.0
```

now either let the makefile automatically determine your platform and architecture

```sh
make build
./bin/<platform>-<arch>

# example
#❯ make build
#rm -rf ./bin
#mkdir -p ./bin
#Building for darwin-arm64
#GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o ./bin/gossht-darwin-arm64 ./cmd
#./bin/gossht-darwin-arm64
```

or to build binaries for all platforms, simply run make

```sh
make

# example
#❯ make
#Building for linux-amd64
#GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/gossht-linux-amd64 ./cmd
#Building for linux-arm64
#GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o ./bin/gossht-linux-arm64 ./cmd
#Building for darwin-amd64
#GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/gossht-darwin-amd64 ./cmd
#Building for darwin-arm64
#GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o ./bin/gossht-darwin-arm64 ./cmd
#Building for windows-amd64
#GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/gossht-windows-amd64.exe ./cmd
#Building for windows-arm64
#GOOS=windows GOARCH=arm64 go build -ldflags "-s -w" -o ./bin/gossht-windows-arm64.exe ./cmd
#Building for freebsd-amd64
#GOOS=freebsd GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/gossht-freebsd-amd64 ./cmd
#Building for freebsd-arm64
#GOOS=freebsd GOARCH=arm64 go build -ldflags "-s -w" -o ./bin/gossht-freebsd-arm64 ./cmd
```

alternatively you can build the binary for the exact platform and architecture you want

```sh
make build-<platform>-<arch>

#example
#❯ make build-windows-amd64
#Building for windows-amd64
#GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X main.Version=v0.0.1-dev" -o ./bin/gossht-windows-amd64.exe ./cmd
```

## License

Gossht is licensed under the [MIT License](https://opensource.org/license/mit).
