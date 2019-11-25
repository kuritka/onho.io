package utils

import "fmt"

func Keypress(){
	var a string
	_,err := fmt.Scanln(&a)
	FailFastOnError(err)
}
