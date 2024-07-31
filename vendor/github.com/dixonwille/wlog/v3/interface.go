package wlog

// UI simply writes to an io.Writer with a new line appended to each call.
// It also has the ability to ask a question and return a response.
type UI interface {
	// Log writes a timestamped message to the writer
	Log(message string)
	// Output writes a message to the writer
	Output(message string)
	// Success writes a message indicating an success message
	Success(message string)
	// Info writes a message indicating an informational message
	Info(message string)
	// Error writes a message indicating an error
	Error(message string)
	// Warn writes a message indicating a warning
	Warn(message string)
	// Running writes a message indicating a process is running
	Running(message string)
	// Ask writes a message to the writer and reads the user's input
	// Message is written to the writer and the response is trimmed by the trim value
	Ask(message string, trim string) (response string, error error)
}
