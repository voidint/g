package wlog

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02T15-04-05"
)

// BasicUI simply writes/reads to correct input/output
// It is not thread safe.
// Pretty simple to wrap your own functions around
type BasicUI struct {
	Reader      io.Reader
	Writer      io.Writer
	ErrorWriter io.Writer
}

// Log prefixes to message before writing to Writer.
func (ui *BasicUI) Log(message string) {
	timeString := time.Now().Format(timeFormat)
	message = timeString + ": " + message
	ui.Output(message)
}

// Output simply writes to Writer.
func (ui *BasicUI) Output(message string) {
	fmt.Fprint(ui.Writer, message)
	fmt.Fprint(ui.Writer, "\n")
}

// Success calls Output to write.
// Useful when you want separate colors or prefixes.
func (ui *BasicUI) Success(message string) {
	ui.Output(message)
}

// Info calls Output to write.
// Useful when you want separate colors or prefixes.
func (ui *BasicUI) Info(message string) {
	ui.Output(message)
}

// Error writes message to ErrorWriter.
func (ui *BasicUI) Error(message string) {
	if ui.ErrorWriter != nil {
		fmt.Fprint(ui.ErrorWriter, message)
		fmt.Fprint(ui.ErrorWriter, "\n")
	} else {
		fmt.Fprint(ui.Writer, message)
		fmt.Fprint(ui.Writer, "\n")
	}
}

// Warn calls Error to write.
// Useful when you want separate colors or prefixes.
func (ui *BasicUI) Warn(message string) {
	ui.Error(message)
}

// Running calls Output to write.
// Useful when you want separate colors or prefixes.
func (ui *BasicUI) Running(message string) {
	ui.Output(message)
}

// Ask will call output with message then wait for Reader to print newline (\n).
// If Reader is os.Stdin then that is when ever a user presses [enter].
// It will clean the response by removing any carriage returns and new lines that if finds.
// Then it will trim the response using the trim variable.
// Use an empty string to specify you do not want to trim.
// If the message is left blank ("") then it will not prompt user before waiting on a response.
func (ui *BasicUI) Ask(message, trim string) (string, error) {
	if message != "" {
		ui.Output(message)
	}
	reader := bufio.NewReader(ui.Reader)
	res, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	res = strings.Replace(res, "\r", "", -1) //this will only be useful under windows
	res = strings.Replace(res, "\n", "", -1)
	res = strings.Trim(res, trim)
	return res, err
}
