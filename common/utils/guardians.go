package utils

import (
	"errors"
	"github.com/kuritka/onho.io/common/log"
)
var logger = log.Log



func FailOnError(err error, msg string){
	if err != nil {
		fail(err)
	}
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
	logger.Fatal().Msgf("%+v", err.Error())
	//panic(fmt.Sprintf("%+v", err))
}