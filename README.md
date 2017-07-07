# dew
Domain Error Wrapper, combining causing error, domain error, custom message, and stacktrace in Go

## Example
```
package main

import (
	"errors"
	"fmt"
	"github.com/kristofferingemansson/dew"
)

func main() {
	var i Asd = &dd{}

	err := Package1(i)
	if err == nil {
		return
	}

	switch dew.Domain(err) {
	case ErrPackage1Lust:
		fmt.Println(err.Error())
	case ErrPackage1Gluttony:
		fmt.Println(err.Error())
	default:
		fmt.Println("Unknown error")
	}

	trace := dew.StackTrace(err)
	for _, r := range trace {
		fmt.Println(r)
	}
}

// Some package
var (
	ErrPackage1Lust     = errors.New("Lust")
	ErrPackage1Gluttony = errors.New("Gluttony")
)

func Package1(a Asd) error {
	err := a.Sdf()
	if err != nil {
		switch dew.Domain(err) {
		case ErrAsdGreed:
			return dew.New(ErrPackage1Gluttony, err, "Omnomnom")
		case ErrAsdSloth:
			return dew.New(ErrPackage1Lust, err, "Grrr")
		default:
			return err
		}
	}
	return nil
}

// Some interface-package
var (
	ErrAsdGreed = errors.New("Greed")
	ErrAsdSloth = errors.New("Sloth")
)

type Asd interface {
	Sdf() error
}

// Some package implementing above interface-package
type dd struct{}

func (d dd) Sdf() error {
	err := errors.New("Some random error from some external package")
	if err != nil {
		return dew.New(ErrAsdGreed, err, "All the stuff!")
	}

	return nil
}
```