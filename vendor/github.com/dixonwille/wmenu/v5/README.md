# WMenu[![Build Status](https://travis-ci.org/dixonwille/wmenu.svg?branch=master)](https://travis-ci.org/dixonwille/wmenu) [![codecov](https://codecov.io/gh/dixonwille/wmenu/branch/master/graph/badge.svg)](https://codecov.io/gh/dixonwille/wmenu)

Package wmenu creates menus for cli programs. It uses wlog for its interface
with the command line. It uses os.Stdin, os.Stdout, and os.Stderr with
concurrency by default. wmenu allows you to change the color of the different
parts of the menu. This package also creates it's own error structure so you can
type assert if you need to. wmenu will validate all responses before calling any function. It will also figure out which function should be called so you don't have to.

[![Watch example](https://asciinema.org/a/4lv3ash3ubtnsclindvzdf320.png)](https://asciinema.org/a/4lv3ash3ubtnsclindvzdf320)

## Import

### Post Go1.11

```go
import "github.com/dixonwille/wmenu/v5"
```

### Pre Go1.11

I try and keep up with my tags. To use the version and stable it is recommended to use `govendor` or another vendoring tool that allows you to build your project for specific tags.

```sh
govendor fetch github.com/dixonwille/wmenu@v4
```

The above will grab the latest v4 at that time and mark it. It will then be stable for you to use.

I will try to support as many versions as possable but please be patient.

### V1.0.0 - Major Release [![Go Report Card](https://goreportcard.com/badge/gopkg.in/dixonwille/wmenu.v1)](https://goreportcard.com/report/gopkg.in/dixonwille/wmenu.v1) [![GoDoc](https://godoc.org/https://godoc.org/gopkg.in/dixonwille/wmenu.v1?status.svg)](https://godoc.org/gopkg.in/dixonwille/wmenu.v1)

### V2.0.0 - Allowing an interface to be passed in for options [![Go Report Card](https://goreportcard.com/badge/gopkg.in/dixonwille/wmenu.v2)](https://goreportcard.com/report/gopkg.in/dixonwille/wmenu.v2) [![GoDoc](https://godoc.org/https://godoc.org/gopkg.in/dixonwille/wmenu.v2?status.svg)](https://godoc.org/gopkg.in/dixonwille/wmenu.v2)

### V3.0.0 - Pass in the option to that option's function [![Go Report Card](https://goreportcard.com/badge/gopkg.in/dixonwille/wmenu.v3)](https://goreportcard.com/report/gopkg.in/dixonwille/wmenu.v3) [![GoDoc](https://godoc.org/https://godoc.org/gopkg.in/dixonwille/wmenu.v3?status.svg)](https://godoc.org/gopkg.in/dixonwille/wmenu.v3)

### V4.0.0 - Now have an Action that supports multiple options [![Go Report Card](https://goreportcard.com/badge/gopkg.in/dixonwille/wmenu.v4)](https://goreportcard.com/report/gopkg.in/dixonwille/wmenu.v4) [![GoDoc](https://godoc.org/https://godoc.org/gopkg.in/dixonwille/wmenu.v4?status.svg)](https://godoc.org/gopkg.in/dixonwille/wmenu.v4)

### v5.0.0 - Support Go Mods

https://pkg.go.dev/github.com/dixonwille/wmenu/v5

## Features

- Force single selection
- Allow multiple selection
- Change the delimiter
- Change the color of different parts of the menu
- Easily see which option(s) are default
- Change the symbol used for default option(s)
- Ask simple yes and no questions
- Validate all responses before calling any functions
- With yes and no can accept:
  - yes, Yes, YES, y, Y
  - no, No, NO, n, N
- Figure out which Action should be called (Options, Default, or Multiple Action)
- Re-ask question if invalid response up to a certain number of times
- Can change max number of times to ask before failing output
- Change reader and writer
- Clear the screen whenever the menu is brought up
- Has its own error structure so you can type assert menu errors

### V2 - Adds these Features

- Allowing any interface to be passed through for the options.

### V3 - Adds these Features

- Pass the option chosen to that options function

### V4 - Adds these Features

- Have one function for both single and multiple select. Allowing the user to an easier way of handeling the request.

### v5 - Support Go Mods

- No other change except you should import with the following now

```go
import "github.com/dixonwille/wmenu/v5"
```

## Usage

This is a simple use of the package. (**NOTE: THIS IS A V4 SAMPLE**)

```go
menu := wmenu.NewMenu("What is your favorite food?")
menu.Action(func (opts []wmenu.Opt) error {fmt.Printf(opts[0].Text + " is your favorite food."); return nil})
menu.Option("Pizza", nil, true, nil)
menu.Option("Ice Cream", nil, false, nil)
menu.Option("Tacos", nil, false, func(opt wmenu.Opt) error {
  fmt.Printf("Tacos are great")
  return nil
})
err := menu.Run()
if err != nil{
  log.Fatal(err)
}
```

The output would look like this:

```
1) *Pizza
2) Ice Cream
3) Tacos
What is your favorite food?
```

If the user just presses `[Enter]` then the option(s) with the `*` will be selected. This indicates that it is a default function. If they choose `1` then they would see `Ice Cream is your favorite food.`. This used the Action's function because the option selected didn't have a function along with it. But if they choose `2` they would see `Tacos are great`. That option did have a function with it which take precedence over Action.

You can you also use:

```go
menu.AllowMultiple()
```

This will allow the user to select multiple options. The default delimiter is a `[space]`, but can be changed by using:

```go
menu.SetSeperator("some string")
```

Another feature is the ability to ask yes or no questions.

```go
menu.IsYesNo(0)
```

This will remove any options previously added options and hide the ones used for the menu. It will simply just ask yes or no. Menu will parse and validate the response for you. This option will always call the Action's function and pass in the option that was selected.

## V3+ - Release

Allows the user to pass anything for the value so it can be retrieved later in the function. The following is to show case the power of this.

> The following was written in V3 but the concept holds for V4. V4 just changed `actFunc` to be `func([]wmenu.Opt) error` instead.

```go
type NameEntity struct {
  FirstName string
  LastName  string
}

optFunc := func(opt wmenu.Opt) error {
  fmt.Println("Option 1 was chosen.")
  return nil
}
actFunc := func(opt wmenu.Opt) error {
  name, ok := opt.Value.(NameEntity)
  if !ok {
    log.Fatal("Could not cast option's value to NameEntity")
  }
  fmt.Printf("%s has an id of %d.\n", opt.Text, opt.ID)
  fmt.Printf("Hello, %s %s.\n", name.FirstName, name.LastName)
  return nil
}
menu := NewMenu("Choose an option.")
menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
menu.Action(actFunc)
menu.Option("Option 1", NameEntity{"Bill", "Bob"}, true, optFunc)
menu.Option("Option 2", NameEntity{"John", "Doe"}, false, nil)
menu.Option("Option 3", NameEntity{"Jane", "Doe"}, false, nil)
err := menu.Run()
if err != nil {
  log.Fatal(err)
}
```

The immediate output would be:

```
Output:
1) *Option 1
2) Option 2
3) Option 3
Choose an option.
```

Now if the user pushes `[ENTER]` the output would be `Options 0 was chosen.`. But now if either option 1 or 2 were chosen it would cast the options value to a NameEntity allowing the function to be able to gather both the first name and last name of the NameEntity. If you want though you can just pass in `nil` as the value or even a string (`"hello"`) since both of these implement the empty interface required by value. Just make sure to cast the values so you can use them appropriately.

## Further Reading

This whole package has been documented and has a few examples in:

- [godocs V1](https://godoc.org/gopkg.in/dixonwille/wmenu.v1)
- [godocs V2](https://godoc.org/gopkg.in/dixonwille/wmenu.v2)
- [godocs V3](https://godoc.org/gopkg.in/dixonwille/wmenu.v3)
- [godocs V4](https://godoc.org/gopkg.in/dixonwille/wmenu.v4)
- [pkg.go.dev V5](https://pkg.go.dev/github.com/dixonwille/wmenu/v5)

You should read the docs to find all functions and structures at your finger tips.
