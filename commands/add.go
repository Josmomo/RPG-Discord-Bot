package commands

import (
	"os"

	"github.com/Josmomo/RPG-Discord-Bot/utils"
	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
)

//CommandAdd
const CommandAdd = "add"

//Add
func Add(session *discordgo.Session, message *discordgo.MessageCreate, args []string) error {
	file, err := os.Create("result.txt")
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return err
	}
	defer file.Close()

	// Read from file

	return nil
}
