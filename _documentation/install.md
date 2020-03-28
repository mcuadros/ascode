---
title: 'Quick Install'
weight: 1 
---

AsCode is written in [Go](https://golang.org/) with support for multiple platforms. 

The latest release can be found at [GitHub Releases.](https://github.com/mcuadros/ascode/releases), currently provides pre-built binaries for the following:

- Linux
- macOS (Darwin)
- Windows

## Binary (Cross-platform)

Download the appropriate version for your platform from [GitHub Releases.](https://github.com/mcuadros/ascode/releases). Once downloaded, the binary can be run from anywhere. You don’t need to install it into a global location. 

Ideally, you should install it somewhere in your PATH for easy use. `/usr/local/bin` is the most probable location.

### Linux 
```sh
wget https://github.com/mcuadros/ascode/releases/download/{{< param "version" >}}/ascode-{{< param "version" >}}-linux-amd64.tar.gz
tar -xvzf ascode-{{< param "version" >}}-linux-amd64.tar.gz
mv ascode /usr/local/bin/
```

### macOS (Darwin)
```sh
wget https://github.com/mcuadros/ascode/releases/download/{{< param "version" >}}/ascode-{{< param "version" >}}-darwin-amd64.tar.gz
tar -xvzf ascode-{{< param "version" >}}-darwin-amd64.tar.gz
mv ascode /usr/local/bin/
```

## Source 

### Prerequisite Tools 

- [Git](https://git-scm.com/)
- [Go](https://golang.org/) (at least Go 1.12)

### Clone from GitHub

AsCode uses the [Go Modules](https://github.com/golang/go/wiki/Modules), so e easiest way to get started is to clone AsCode in a directory outside of the `$GOPATH`, as in the following example:

```sh
git clone https://github.com/mcuadros/ascode.git $HOME/ascode-src
cd $HOME/ascode-src
go install ./... 
```

