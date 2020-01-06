package utils

import (
	"errors"
	"github.com/kuritka/onho.io/common/log"
	"os"
)
var logger = log.Log


func FailOnError(err error, msg string){
	if err != nil {
		fail(err)
	}
}


func FailFast(msg string){
	fail(errors.New(msg))
}


func NotImplemented(){
	fail(errors.New("function not implemented"))
}


func FailFastOnError(err error){
	if err != nil {
		fail(err)
	}
}

func FailOnLessOrEqualToZero(num int, msg string) {
	if num <= 0 {
		fail(errors.New(msg))
	}
}


func DisposeOnError(err error, msg string, dispose func()){
	if err != nil {
		dispose()
		fail(errors.New(msg))
	}
}


func FailOnEmptyString(str string, msg string ){
	if str == "" {
		fail(errors.New(msg))
	}
}

func FailIfFileNotExist(path string, msg string) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			fail(errors.New(msg))
		}
	}
}


func DisposeOnEmptyString(str string, msg string,dispose func() ){
	if str == "" {
		dispose()
		fail(errors.New(msg))
	}
}

func FailOnNil(entity interface{}, msg string){
	if entity == nil {
		fail(errors.New(msg))
	}
}

func DisposeOnNil(entity interface{}, msg string,dispose func() ){
	if entity == nil {
		dispose()
		fail(errors.New(msg))
	}
}

func fail(err error){
	logger.Panic().Msgf("%+v", err.Error())
}