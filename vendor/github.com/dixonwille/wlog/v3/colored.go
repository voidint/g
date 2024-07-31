package wlog

import "github.com/daviddengcn/go-colortext"

// ColorUI is a wrapper for UI that adds color.
type ColorUI struct {
	LogFGColor      Color
	OutputFGColor   Color
	SuccessFGColor  Color
	InfoFGColor     Color
	ErrorFGColor    Color
	WarnFGColor     Color
	RunningFGColor  Color
	AskFGColor      Color
	ResponseFGColor Color
	LogBGColor      Color
	OutputBGColor   Color
	SuccessBGColor  Color
	InfoBGColor     Color
	ErrorBGColor    Color
	WarnBGColor     Color
	RunningBGColor  Color
	AskBGColor      Color
	ResponseBGColor Color
	UI              UI
}

// Log calls UI.Log to write.
// LogFGColor and LogBGColor are used for color.
func (ui *ColorUI) Log(message string) {
	ct.ChangeColor(ui.LogFGColor.Code, ui.LogFGColor.Bright, ui.LogBGColor.Code, ui.LogBGColor.Bright)
	ui.UI.Log(message)
	ct.ResetColor()
}

// Output calls UI.Output to write.
// OutputFGColor and OutputBGColor are used for color.
func (ui *ColorUI) Output(message string) {
	ct.ChangeColor(ui.OutputFGColor.Code, ui.OutputFGColor.Bright, ui.OutputBGColor.Code, ui.OutputBGColor.Bright)
	ui.UI.Output(message)
	ct.ResetColor()
}

// Success calls UI.Success to write.
// Useful when you want separate colors or prefixes.
// SuccessFGColor and SuccessBGColor are used for color.
func (ui *ColorUI) Success(message string) {
	ct.ChangeColor(ui.SuccessFGColor.Code, ui.SuccessFGColor.Bright, ui.SuccessBGColor.Code, ui.SuccessBGColor.Bright)
	ui.UI.Success(message)
	ct.ResetColor()
}

// Info calls UI.Info to write.
// Useful when you want separate colors or prefixes.
// InfoFGColor and InfoBGColor are used for color.
func (ui *ColorUI) Info(message string) {
	ct.ChangeColor(ui.InfoFGColor.Code, ui.InfoFGColor.Bright, ui.InfoBGColor.Code, ui.InfoBGColor.Bright)
	ui.UI.Info(message)
	ct.ResetColor()
}

// Error calls UI.Error to write.
// ErrorFGColor and ErrorBGColor are used for color.
func (ui *ColorUI) Error(message string) {
	ct.ChangeColor(ui.ErrorFGColor.Code, ui.ErrorFGColor.Bright, ui.ErrorBGColor.Code, ui.ErrorBGColor.Bright)
	ui.UI.Error(message)
	ct.ResetColor()
}

// Warn calls UI.Warn to write.
// Useful when you want separate colors or prefixes.
// WarnFGColor and WarnBGColor are used for color.
func (ui *ColorUI) Warn(message string) {
	ct.ChangeColor(ui.WarnFGColor.Code, ui.WarnFGColor.Bright, ui.WarnBGColor.Code, ui.WarnBGColor.Bright)
	ui.UI.Warn(message)
	ct.ResetColor()
}

// Running calls UI.Running to write.
// Useful when you want separate colors or prefixes.
// RunningFGColor and RunningBGColor are used for color.
func (ui *ColorUI) Running(message string) {
	ct.ChangeColor(ui.RunningFGColor.Code, ui.RunningFGColor.Bright, ui.RunningBGColor.Code, ui.RunningBGColor.Bright)
	ui.UI.Running(message)
	ct.ResetColor()
}

//Ask will call UI.Output with message then wait for UI.Ask to return a response and/or error.
//It will clean the response by removing any carriage returns and new lines that if finds.
//Then it will trim the message using the trim variable.
//Use and empty string to specify you do not want to trim.
//If a message is not used ("") then it will not prompt user before waiting on a response.
//AskFGColor and AskBGColor are used for message color.
//ResponseFGColor and ResponseBGColor are used for response color.
func (ui *ColorUI) Ask(message, trim string) (string, error) {
	if message != "" {
		ct.ChangeColor(ui.AskFGColor.Code, ui.AskFGColor.Bright, ui.AskBGColor.Code, ui.AskBGColor.Bright)
		ui.UI.Output(message)
		ct.ResetColor()
	}
	ct.ChangeColor(ui.ResponseFGColor.Code, ui.ResponseFGColor.Bright, ui.ResponseBGColor.Code, ui.ResponseBGColor.Bright)
	res, err := ui.UI.Ask("", trim)
	ct.ResetColor()
	return res, err
}
