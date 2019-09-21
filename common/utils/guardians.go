package utils

import (
	"fmt"
	"github.com/kuritka/onho.io/common/log"
)
var logger = log.Log


func FailOnError(err error, msg string){
	if err != nil {
		fail(msg)
	}
}

func DisposeOnError(err error, msg string, dispose func()){
	if err != nil {
		dispose()
		fail(msg)
	}
}


func FailOnEmptyString(str string, msg string ){
	if str == "" {
		fail(msg)
	}
}


func DisposeOnEmptyString(str string, msg string,dispose func() ){
	if str == "" {
		dispose()
		fail(msg)
	}
}

func FailOnNil(entity interface{}, msg string){
	if entity == nil {
		fail(msg)
	}
}

func DisposeOnNil(entity interface{}, msg string,dispose func() ){
	if entity == nil {
		dispose()
		fail(msg)
	}
}

func fail(msg string){
	logger.Fatal().Msgf("%s", msg)
	panic(fmt.Sprintf("%s", msg))
}