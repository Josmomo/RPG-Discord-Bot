package main

import (
	"github.com/Josmomo/RPG-Discord-Bot/client"
	"github.com/Josmomo/RPG-Discord-Bot/utils"
	"github.com/Sirupsen/logrus"
)

func main() {

	bot := client.CreateBot()
	bot.Run()

	logrus.WithFields(utils.Locate()).Info("Main waiting")
	<-make(chan struct{})
}
