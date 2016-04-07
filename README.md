tapet
=========
Requires >= OS X 10.9

Tested and working on OS X 10.11 El Capitan.

## Installation

### Short version

If you already have Go installed and configured:

```> go get github.com/joecheatham/tapet```

### Less Short Version

If you have homebrew installed, first make sure you've tapped my repository:

```> brew tap joecheatham/homebrew-formulae```

Then you can install with the following:

```> brew install tapet```

### Long Version

This method doesn't require a GOPATH.

- Install Go
- Clone this repository
- Run `go build` to build binary
- Run newly created binary from the command line

## Usage

#### Getting image of the day, extracting colors, and setting as background image:
`> tapet`

#### Where are my images stored?
Look in `~/.tapet/`.
