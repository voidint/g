//Package wmenu creates menus for cli programs.
//It uses wlog for it's interface with the command line.
//It uses os.Stdin, os.Stdout, and os.Stderr with concurrency by default.
//wmenu allows you to change the color of the different parts of the menu.
//This package also creates it's own error structure so you can type assert if you need to.
//wmenu will validate all responses before calling any function.
//It will also figure out which function should be called so you don't have to.
package wmenu

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/dixonwille/wlog/v3"
	"github.com/mattn/go-isatty"
)

//DefaultYN is used to specify what the default answer is to a yes/no question.
type DefaultYN int

const (
	//DefY defaults yes/no question to use yes.
	DefY DefaultYN = iota + 1
	//DefN defaults yes/no question to use no.
	DefN
)

var (
	noColor = os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))
)

//Menu is used to display options to a user.
//A user can then select options and Menu will validate the response and perform the correct action.
type Menu struct {
	question       string
	function       func([]Opt) error
	options        []Opt
	ui             wlog.UI
	multiSeparator string
	allowMultiple  bool
	loopOnInvalid  bool
	clear          bool
	tries          int
	defIcon        string
	isYN           bool
	ynDef          DefaultYN
	padOptionID    bool
	initialIndex   int
}

//NewMenu creates a menu with a wlog.UI as the writer.
func NewMenu(question string) *Menu {
	//Create a default ui to use for menu
	var ui wlog.UI
	ui = wlog.New(os.Stdin, os.Stdout, os.Stderr)
	ui = wlog.AddConcurrent(ui)

	return &Menu{
		question:       question,
		function:       nil,
		options:        nil,
		ui:             ui,
		multiSeparator: " ",
		allowMultiple:  false,
		loopOnInvalid:  false,
		clear:          false,
		tries:          3,
		defIcon:        "*",
		isYN:           false,
		ynDef:          0,
		initialIndex:   1,
	}
}

//AddColor will change the color of the menu items.
//optionColor changes the color of the options.
//questionColor changes the color of the questions.
//errorColor changes the color of the question.
//Use wlog.None if you do not want to change the color.
func (m *Menu) AddColor(optionColor, questionColor, responseColor, errorColor wlog.Color) {
	if !noColor {
		m.ui = wlog.AddColor(questionColor, errorColor, wlog.None, wlog.None, optionColor, responseColor, wlog.None, wlog.None, wlog.None, m.ui)
	}
}

// PadOptionID will pad the option IDs when printing, so they all right-align.
func (m *Menu) PadOptionID() {
	m.padOptionID = true
}

//ClearOnMenuRun will clear the screen when a menu is ran.
//This is checked when LoopOnInvalid is activated.
//Meaning if an error occurred then it will clear the screen before asking again.
func (m *Menu) ClearOnMenuRun() {
	m.clear = true
}

//SetSeparator sets the separator to use when multiple options are valid responses.
//Default value is a space.
func (m *Menu) SetSeparator(sep string) {
	m.multiSeparator = sep
}

//SetTries sets the number of tries on the loop before failing out.
//Default is 3.
//Negative values act like 0.
func (m *Menu) SetTries(i int) {
	m.tries = i
}

//LoopOnInvalid is used if an invalid option was given then it will prompt the user again.
func (m *Menu) LoopOnInvalid() {
	m.loopOnInvalid = true
}

//SetDefaultIcon sets the icon used to identify which options will be selected by default
func (m *Menu) SetDefaultIcon(icon string) {
	m.defIcon = icon
}

//IsYesNo sets the menu to a yes/no state.
//Does not show options but does ask question.
//Will also parse the answer to allow for all variants of yes/no (IE Y yes No ...)
//Both will call the Action function you specified.
// Opt{ID: 1, Text: "y"} for yes and Opt{ID: 2, Text: "n"} for no will be passed to the function.
func (m *Menu) IsYesNo(def DefaultYN) {
	m.isYN = true
	m.ynDef = def
}

// InitialIndex sets the starting number for the displayed list. It defaults to 1.
func (m *Menu) InitialIndex(index int) {
	m.initialIndex = index
}

//Option adds an option to the menu for the user to select from.
//value is an empty interface that can be used to pass anything through to the function.
//title is the string the user will select
//isDefault is whether this option is a default option (IE when no options are selected).
//function is what is called when only this option is selected.
//If function is nil then it will default to the menu's Action.
func (m *Menu) Option(title string, value interface{}, isDefault bool, function func(Opt) error) {
	option := newOption(len(m.options), title, value, isDefault, function)
	m.options = append(m.options, *option)
}

//Action adds a default action to use in certain scenarios.
//If the selected option (by default or user selected) does not have a function applied to it this will be called.
//If there are no default options and no option was selected this will be called with an option that has an ID of -1.
func (m *Menu) Action(function func([]Opt) error) {
	m.function = function
}

//AllowMultiple will tell the menu to allow multiple selections.
//The menu will fail if this is not called and mulple selections were selected.
func (m *Menu) AllowMultiple() {
	m.allowMultiple = true
}

//ChangeReaderWriter changes where the menu listens and writes to.
//reader is where user input is collected.
//writer and errorWriter is where the menu should write to.
func (m *Menu) ChangeReaderWriter(reader io.Reader, writer, errorWriter io.Writer) {
	var ui wlog.UI
	ui = wlog.New(reader, writer, errorWriter)
	m.ui = ui
}

//Run is used to execute the menu.
//It will print to options and question to the screen.
//It will only clear the screen if ClearOnMenuRun is activated.
//This will validate all responses.
//Errors are of type MenuError.
func (m *Menu) Run() error {
	if m.clear {
		Clear()
	}
	valid := false
	var options []Opt
	//Loop and on error check if loopOnInvalid is enabled.
	//If it is Clear the screen and write error.
	//Then ask again
	for !valid {
		//step 1 print options to screen
		m.print()
		//step 2 ask question, get and validate response
		opt, err := m.ask()
		if err != nil {
			m.tries = m.tries - 1
			if !IsMenuErr(err) {
				err = newMenuError(err, "", m.triesLeft())
			}
			if m.loopOnInvalid && m.tries > 0 {
				if m.clear {
					Clear()
				}
				m.ui.Error(err.Error())
			} else {
				return err
			}
		} else {
			options = opt
			valid = true
		}
	}
	//step 3 call appropriate action with the responses
	return m.callAppropriate(options)
}

func (m *Menu) callAppropriate(options []Opt) (err error) {
	if len(options) == 0 {
		return m.callAppropriateNoOptions()
	}
	if len(options) == 1 && options[0].function != nil {
		return options[0].function(options[0])
	}
	return m.function(options)
}

func (m *Menu) callAppropriateNoOptions() (err error) {
	options := m.getDefault()
	if len(options) == 0 {
		return m.function([]Opt{{ID: -1}})
	}
	if len(options) == 1 && options[0].function != nil {
		return options[0].function(options[0])
	}
	return m.function(options)
}

//hide options when this is a yes or no
func (m *Menu) print() {
	if !m.isYN {
		outputFormat := "%d) %s%s"
		if m.padOptionID {
			padding := len(strconv.Itoa(len(m.options)))
			outputFormat = "%" + strconv.Itoa(padding) + "d) %s%s"
		}
		for _, opt := range m.options {
			icon := m.defIcon
			if !opt.isDefault {
				icon = ""
			}
			m.ui.Output(fmt.Sprintf(outputFormat, opt.ID+m.initialIndex, icon, opt.Text))
		}
	} else {
		//TODO Allow user to specify what to use as value for YN options
		m.options = []Opt{}
		m.Option("y", "yes", m.ynDef == DefY, nil)
		m.Option("n", "no", m.ynDef == DefN, nil)
	}
}

func (m *Menu) ask() ([]Opt, error) {
	if m.isYN {
		if m.ynDef == DefY {
			m.question += " (Y/n)"
		} else {
			m.question += " (y/N)"
		}
	}
	var trim string
	if m.multiSeparator == " " {
		trim = m.multiSeparator
	} else {
		trim = m.multiSeparator + " "
	}
	res, err := m.ui.Ask(m.question, trim)
	if err != nil {
		return nil, err
	}
	//Validate responses
	//Check if no responses are returned and no action to call
	if res == "" {
		//get default options
		opt := m.getDefault()
		if !m.validOptAndFunc(opt) {
			return nil, newMenuError(ErrNoResponse, "", m.triesLeft())
		}
		return nil, nil
	}

	var responses []int
	if !m.isYN {
		responses, err = m.resToInt(res)
		if err != nil {
			return nil, err
		}

		err = m.validateResponses(responses)
		if err != nil {
			return nil, err
		}
	} else {
		responses, err = m.ynResParse(res)
		if err != nil {
			return nil, err
		}
	}

	//Parse responses and return them as options
	var finalOptions []Opt
	for _, response := range responses {
		finalOptions = append(finalOptions, m.options[response-m.initialIndex])
	}

	return finalOptions, nil
}

//Converts the response string to a slice of ints, also validates along the way.
func (m *Menu) resToInt(res string) ([]int, error) {
	resStrings := strings.Split(res, m.multiSeparator)
	//Check if we don't want multiple responses
	if !m.allowMultiple && len(resStrings) > 1 {
		return nil, newMenuError(ErrTooMany, "", m.triesLeft())
	}

	//Convert responses to intigers
	var responses []int
	for _, response := range resStrings {
		//Check if it is an intiger
		response = strings.Trim(response, " ")
		r, err := strconv.Atoi(response)
		if err != nil {
			return nil, newMenuError(ErrInvalid, response, m.triesLeft())
		}
		responses = append(responses, r)
	}
	return responses, nil
}

func (m *Menu) ynResParse(res string) ([]int, error) {
	resStrings := strings.Split(res, m.multiSeparator)
	if len(resStrings) > 1 {
		return nil, newMenuError(ErrTooMany, "", m.triesLeft())
	}
	re := regexp.MustCompile("^\\s*(?:([Yy])(?:es|ES)?|([Nn])(?:o|O)?)\\s*$")
	matches := re.FindStringSubmatch(res)
	if len(matches) < 2 {
		return nil, newMenuError(ErrInvalid, res, m.triesLeft())
	}
	if strings.ToLower(matches[1]) == "y" {
		return []int{int(DefY)}, nil
	}
	return []int{int(DefN)}, nil
}

//Check if response is in the range of options
//If it is make sure it is not duplicated
func (m *Menu) validateResponses(responses []int) error {
	var tmp []int
	for _, response := range responses {
		realIndex := response - m.initialIndex
		if realIndex < 0 || len(m.options) <= realIndex {
			return newMenuError(ErrInvalid, strconv.Itoa(response), m.triesLeft())
		}

		if exist(tmp, response) {
			return newMenuError(ErrDuplicate, strconv.Itoa(response), m.triesLeft())
		}

		tmp = append(tmp, response)
	}
	return nil
}

//Simply checks if number exists in the slice
func exist(slice []int, number int) bool {
	for _, s := range slice {
		if number == s {
			return true
		}
	}
	return false
}

//gets a list of default options
func (m *Menu) getDefault() []Opt {
	var opt []Opt
	for _, o := range m.options {
		if o.isDefault {
			opt = append(opt, o)
		}
	}
	return opt
}

//make sure that there is an action available to be called in certain cases
//returns false if it chould not find an action for the number options available
func (m *Menu) validOptAndFunc(opt []Opt) bool {
	if m.function == nil {
		if len(opt) == 1 && opt[0].function != nil {
			return true
		}
		return false
	}
	return true
}

func (m *Menu) triesLeft() int {
	if m.loopOnInvalid && m.tries > 0 {
		return m.tries
	}
	return 0
}
