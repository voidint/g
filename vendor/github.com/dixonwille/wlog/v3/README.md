# WLog [![Build Status](https://travis-ci.org/dixonwille/wlog.svg?branch=master)](https://travis-ci.org/dixonwille/wlog) [![Go Report Card](https://goreportcard.com/badge/github.com/dixonwille/wlog/v3)](https://goreportcard.com/report/github.com/dixonwille/wlog/v3) [![codecov](https://codecov.io/gh/dixonwille/wlog/branch/master/graph/badge.svg)](https://codecov.io/gh/dixonwille/wlog)

Package wlog creates simple to use UI structure. The UI is used to simply print
to the screen. There a wrappers that will wrap each other to create a good
looking UI. You can add color and prefixes as well as make it thread safe.

## Documentation

https://pkg.go.dev/github.com/dixonwille/wlog/v3

## Installation

WLog can be added to your go module file by running:

```bash
go get github.com/dixonwille/wlog/v3@latest
```

You can them import the library using an import statement:

```go
import "github.com/dixonwille/wlog/v3"
```

## Idea Behind WLog

I used Mitchellh's [CLI](https://github.com/mitchellh/cli) structure and
wrapping for the different structures. It was a clean look and feel. Plus it
was pretty simple to use. But I didn't want all the other cli stuff that came
with the package so I created this.

For color I use DavidDenGCN's
[Go-ColorText](https://github.com/daviddengcn/go-colortext). His color package
allows for color that is available cross-platforms. I made a wrapper with all
possible color combinations with his package. So you only have to import this
package (one less line).

## Example Usage

This example creates a new `wlog.UI` instance, simulates a user providing input and calls UI functions that show output. If you wish to try the example and provide your own user input you can replace the `reader` variable with a reader such as `os.Stdin` which will read from a terminal.

```go
var ui wlog.UI
reader := strings.NewReader("User Input\r\n") //Simulate user typing "User Input" then pressing [enter] when reading from os.Stdin
ui = wlog.New(reader, os.Stdout, os.Stdout)
ui = wlog.AddPrefix("?", wlog.Cross, " ", "", "", "~", wlog.Check, "!", ui)
ui = wlog.AddConcurrent(ui)

ui.Ask("Ask question", "")
ui.Error("Error message")
ui.Info("Info message")
ui.Output("Output message")
ui.Running("Running message")
ui.Success("Success message")
ui.Warn("Warning message")
```

Output:

```
? Ask question
✗ Error message
 Info message
Output message
~ Running message
✓ Success message
! Warning message
```

On Windows it outputs to this (this includes color):

![winss](https://raw.githubusercontent.com/dixonwille/wlog/master/resources/winss.png)

On Mac it outputs to this (this includes color):

![macss](https://raw.githubusercontent.com/dixonwille/wlog/master/resources/macss.png)
