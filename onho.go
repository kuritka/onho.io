package main

import (
	"github.com/kuritka/onho.io/cmd"
	"github.com/kuritka/onho.io/common/log"
	"os/user"
)

var logger   = log.Log

func main() {

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Current User
	logger.Info().Msg("Hi " + user.Name + " (id: " + user.Uid + ")")
	logger.Info().Msg("Username: " + user.Username)
	logger.Info().Msg("Home Dir: " + user.HomeDir)


	cmd.Execute()
}


