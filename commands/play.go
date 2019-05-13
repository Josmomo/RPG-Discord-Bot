package commands

import (
	"github.com/Josmomo/RPG-Discord-Bot/database"
	"github.com/Josmomo/RPG-Discord-Bot/utils"
	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
)

//CommandPlay
const CommandPlay = "play"

var reactionList = []string{
	"\x30\xE2\x83\xA3",
	"\x31\xE2\x83\xA3",
	"\x32\xE2\x83\xA3",
	"\x33\xE2\x83\xA3",
	"\x34\xE2\x83\xA3",
	"\x35\xE2\x83\xA3",
	"\x36\xE2\x83\xA3",
	"\x37\xE2\x83\xA3",
	//"\x38\xE2\x83\xA3",
	//"\x39\xE2\x83\xA3",
}

//Play
func Play(mongoDBClient database.MongoDBClient, session *discordgo.Session, message *discordgo.MessageCreate, args []string) error {
	messageString := "@everyone, who wants to play something?\n"
	messageString += "\x30\xE2\x83\xA3 Anything goes\n"
	messageString += "\x31\xE2\x83\xA3 Other\n"
	messageString += "\x32\xE2\x83\xA3 Minecraft\n"
	messageString += "\x33\xE2\x83\xA3 Civ6\n"
	messageString += "\x34\xE2\x83\xA3 Vermintide 2\n"
	messageString += "\x35\xE2\x83\xA3 Monster Hunter World\n"
	messageString += "\x36\xE2\x83\xA3 Stardew\n"
	messageString += "\x37\xE2\x83\xA3 Apex Legends\n"

	channelID := message.ChannelID
	messageID := message.ID

	session.ChannelMessageDelete(channelID, messageID)
	botMessage, err := session.ChannelMessageSend(channelID, messageString)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return err
	}

	for _, reaction := range reactionList {
		_ = session.MessageReactionAdd(botMessage.ChannelID, botMessage.ID, reaction)
		if err != nil {
			logrus.WithFields(utils.Locate()).Error(err.Error())
			return err
		}
	}

	return nil
}
