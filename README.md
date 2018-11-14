wfcd
====

[![Build Status](https://travis-ci.org/btcsuite/btcd.png?branch=master)](https://travis-ci.org/btcsuite/btcd)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/btcsuite/btcd)

wfcd is an alternative full node wificoin implementation written in Go (golang).


## Requirements

[Go](http://golang.org) 1.8 or newer.

## Installation

#### Windows - MSI Available

https://github.com/wificoin-project/wfc/releases

#### Linux/BSD/MacOSX/POSIX - Build from Source

- Install Go according to the installation instructions here:
  http://golang.org/doc/install

- Mirror set

```
glide mirror set https://golang.org/x/crypto https://github.com/golang/crypto --vcs git
```

- Ensure Go was installed properly and is a supported version:

```bash
$ go version
$ go env GOROOT GOPATH
```

NOTE: The `GOROOT` and `GOPATH` above must not be the same path.  It is
recommended that `GOPATH` is set to a directory in your home directory such as
`~/goprojects` to avoid write permission issues.  It is also recommended to add
`$GOPATH/bin` to your `PATH` at this point.

- Run the following commands to obtain btcd, all dependencies, and install it:

```bash
$ go get -u github.com/Masterminds/glide
$ git clone https://github.com/wificoin-project/wfcd $GOPATH/src/github.com/wificoin-project/wfcd
$ cd $GOPATH/src/github.com/wificoin-project/wfcd
$ glide install
$ go install . ./cmd/...
```

- wfcd (and utilities) will now be installed in ```$GOPATH/bin```.  If you did
  not already add the bin directory to your system path during Go installation,
  we recommend you do so now.

## Updating

#### Windows

Install a newer MSI

#### Linux/BSD/MacOSX/POSIX - Build from Source

- Run the following commands to update btcd, all dependencies, and install it:

```bash
$ cd $GOPATH/src/github.com/wificoin-project/wfcd
$ git pull && glide install
$ go install . ./cmd/...
```

## Getting Started

wfcd has several configuration options available to tweak how it runs, but all
of the basic operations described in the intro section work with zero
configuration.

#### Windows (Installed from MSI)



#### Linux/BSD/POSIX/Source

```bash
$ ./wfcd
```


## Issue Tracker


## Documentation

The documentation is a work-in-progress.  

## GPG Verification Key


## License

wfcd is licensed under the [copyfree](http://copyfree.org) ISC License.
