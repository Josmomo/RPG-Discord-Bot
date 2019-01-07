package commands

import (
	"time"

	"github.com/Josmomo/RPG-Discord-Bot/database"
	"github.com/Josmomo/RPG-Discord-Bot/utils"
	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
)

//CommandHelp
const CommandHelp = "help"

//Help
func Help(mongoDBClient database.MongoDBClient, session *discordgo.Session, message *discordgo.MessageCreate, args []string) error {
	messageString := "```" + `Schedule your adventure!

	add <1-7>*
		Add days to current week
	addNextWeek <1-7>*
		Add days to next week
	remove <1-7>*
		Remove days from current week
	removeNextWeek <1-7>*
		Remove days from next week
	checkWeek <1-53>?
		Check who can play on any given week, current week if nothing is specified
	roll <1-99D1-999>?
		Roll a dice of your choice, rolls a D20 if nothing is specified
	help
		Shows help message
` + "```"
	botMessage, err := session.ChannelMessageSend(message.ChannelID, messageString)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return err
	}
	defer func() {
		session.ChannelMessageDelete(message.ChannelID, message.ID)
		time.Sleep(time.Second * 60)
		session.ChannelMessageDelete(botMessage.ChannelID, botMessage.ID)
	}()

	return nil
}
