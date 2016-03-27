tapet
=========
Requires >= OS X 10.9

Tested and working on OS X 10.11 El Capitan.

## Installation

### Short version

If you already have Go installed and configured:

```> go get github.com/joecheatham/tapet```

### Less Short Version

This method doesn't require a GOPATH.

- Install Go
- Clone this repository
- Run `go build -o tapet` in to build binary
- Run newly created binary from the command line

## Usage

#### Getting image of the day, extracting colors, and setting as background image:
`> tapet`

#### Where are my images stored?
Look in `~/.tapet/`.