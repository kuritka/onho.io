package utils

import (
	"errors"
	"fmt"
	"github.com/kuritka/onho.io/common/log"
)
var logger = log.Log


func Fail(msg string){
	err := errors.New("service error")
	logger.Fatal().Msgf("%s: %s", msg, err)
	panic(fmt.Sprintf("%s: %s", msg, err))
}

func FailOnError(err error, msg string){
	if err != nil {
		logger.Fatal().Msgf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}


func DisposeOnError(err error, msg string, dispose func()){
	if err != nil {
		dispose()
		logger.Fatal().Msgf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}


func FailOnEmptyString(str string, msg string ){
	if str == "" {
		logger.Fatal().Msgf("%s", msg)
		panic(fmt.Sprintf("%s", msg))
	}
}


func DisposeOnEmptyString(str string, msg string,dispose func() ){
	if str == "" {
		dispose()
		logger.Fatal().Msgf("%s", msg)
		panic(fmt.Sprintf("%s", msg))
	}
}

func FailOnNil(entity interface{}, msg string){
	if entity == nil {
		logger.Fatal().Msgf("%s", msg)
		panic(fmt.Sprintf("%s", msg))
	}
}