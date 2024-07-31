package wlog

import "sync"

// ConcurrentUI is a wrapper for UI that makes the UI thread safe.
type ConcurrentUI struct {
	UI UI
	l  sync.Mutex
}

// Log calls UI.Log to write.
// This is a thread safe function.
func (ui *ConcurrentUI) Log(message string) {
	ui.l.Lock()
	defer ui.l.Unlock()
	ui.UI.Log(message)
}

// Output calls UI.Output to write.
// This is a thread safe function.
func (ui *ConcurrentUI) Output(message string) {
	ui.l.Lock()
	defer ui.l.Unlock()
	ui.UI.Output(message)
}

// Success calls UI.Success to write.
// Useful when you want separate colors or prefixes.
// This is a thread safe function.
func (ui *ConcurrentUI) Success(message string) {
	ui.l.Lock()
	defer ui.l.Unlock()
	ui.UI.Success(message)
}

// Info calls UI.Info to write.
// Useful when you want separate colors or prefixes.
// This is a thread safe function.
func (ui *ConcurrentUI) Info(message string) {
	ui.l.Lock()
	defer ui.l.Unlock()
	ui.UI.Info(message)
}

// Error calls UI.Error to write.
// This is a thread safe function.
func (ui *ConcurrentUI) Error(message string) {
	ui.l.Lock()
	defer ui.l.Unlock()
	ui.UI.Error(message)
}

// Warn calls UI.Warn to write.
// Useful when you want separate colors or prefixes.
// This is a thread safe function.
func (ui *ConcurrentUI) Warn(message string) {
	ui.l.Lock()
	defer ui.l.Unlock()
	ui.UI.Warn(message)
}

// Running calls UI.Running to write.
// Useful when you want separate colors or prefixes.
// This is a thread safe function.
func (ui *ConcurrentUI) Running(message string) {
	ui.l.Lock()
	defer ui.l.Unlock()
	ui.UI.Running(message)
}

// Ask will call UI.Ask with message then wait for UI.Ask to return a response and/or error.
// It will clean the response by removing any carriage returns and new lines that if finds.
//Then it will trim the message using the trim variable.
//Use and empty string to specify you do not want to trim.
// If a message is not used ("") then it will not prompt user before waiting on a response.
// This is a thread safe function.
func (ui *ConcurrentUI) Ask(message, trim string) (string, error) {
	ui.l.Lock()
	defer ui.l.Unlock()
	res, err := ui.UI.Ask(message, trim)
	return res, err
}
