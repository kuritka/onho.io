package utils

import (
	"fmt"
	"github.com/kuritka/onho.io/common/log"
)
var logger = log.Log

type Guard struct {}

func FailOnError(err error, msg string){
	if err != nil {
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


func FailOnNil(entity interface{}, msg string){
	if entity == nil {
		logger.Fatal().Msgf("%s", msg)
		panic(fmt.Sprintf("%s", msg))
	}
}