package client

import (
	"regexp"
	"strings"

	"github.com/Josmomo/RPG-Discord-Bot/commands"
	"github.com/Josmomo/RPG-Discord-Bot/constants"
	"github.com/Josmomo/RPG-Discord-Bot/database"
	"github.com/Josmomo/RPG-Discord-Bot/utils"
	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
)

//Bot struct
type Bot struct {
	ID            string
	session       *discordgo.Session
	user          *discordgo.User
	closeChannel  chan bool
	mongoDBClient database.MongoDBClient
}

//CreateBot create and returns a new Bot
func CreateBot() *Bot {
	bot := new(Bot)
	return bot
}

//Run starts the bot
func (bot *Bot) Run() {
	var err error

	// Create MongoDBClient
	mongoDBClient, err := database.CreateNewClient()
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return
	}
	bot.mongoDBClient = mongoDBClient
	//defer bot.mongoDBClient.Session.Close()

	// Create a new Discord session
	bot.session, err = discordgo.New("Bot " + constants.Token)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return
	}

	// Create a new Discord user
	bot.user, err = bot.session.User("@me")
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return
	}

	bot.ID = bot.user.ID
	bot.session.AddHandler(bot.commandHandler)
	bot.session.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "A friendly helpful bot!")
		if err != nil {
			logrus.WithFields(utils.Locate()).Error(err.Error())
		}
		for _, guild := range discord.State.Guilds {
			logrus.WithFields(utils.Locate()).Info("RPG-Discord-Bot started on server " + guild.Name)
			channels, _ := discord.GuildChannels(guild.ID)
			for _, channel := range channels {
				logrus.WithFields(utils.Locate()).Info("RPG-Discord-Bot started on server " + channel.Name)
			}
		}
	})

	err = bot.session.Open()
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return
	}
	defer bot.session.Close()
	<-bot.closeChannel
}

func (bot *Bot) commandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	logrus.WithFields(utils.Locate()).Info("Author=", message.Author)
	logrus.WithFields(utils.Locate()).Info("Content=" + message.Content)

	user := message.Author
	channelID := message.ChannelID
	if user.ID == bot.ID || user.Bot || channelID != constants.BotChannelID {
		// Do nothing, a bot wrote this message
		return
	}

	command, args, err := parseMessage(message.Content)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return
	}

	switch command {
	case commands.CommandRoll:
		err := commands.Roll(session, message, args)
		if err != nil {
			return
		}
	case commands.CommandAdd:
		err := commands.Add(bot.mongoDBClient, session, message, args)
		if err != nil {
			return
		}
	default:
		logrus.WithFields(utils.Locate()).Info("Command not recognized")
	}
}

//Stop closes the bot
func (bot *Bot) Stop() {
	bot.closeChannel <- true
}

// Extract command and args
func parseMessage(message string) (string, []string, error) {
	command := ""
	args := []string{}
	regexpCommand := regexp.MustCompile("^" + constants.Prefix + `\S*`)
	if regexpCommand.MatchString(message) {
		splitMessage := strings.Split(strings.TrimSpace(message), " ")
		command = regexpCommand.FindString(splitMessage[0])[len(constants.Prefix):]
		args = splitMessage[1:]
	}
	return command, args, nil
}
