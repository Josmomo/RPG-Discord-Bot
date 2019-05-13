package commands

import (
	"regexp"
	"strconv"
	"time"

	"github.com/Josmomo/RPG-Discord-Bot/database"
	"github.com/Josmomo/RPG-Discord-Bot/utils"
	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
)

//CommandAddNextWeek
const CommandAddNextWeek = "addNextWeek"

//AddNextWeek
func AddNextWeek(mongoDBClient database.MongoDBClient, session *discordgo.Session, message *discordgo.MessageCreate, args []string) error {
	weekdays, err := parseAddArgs(args)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return err
	}
	year, week, err := utils.GetYearWeekNext()
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		t := time.Now()
		year, week = t.ISOWeek()
	}
	entry, err := mongoDBClient.GetDocFromIndex(message.Author.Mention(), year, week)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		entry.UserID = message.Author.ID
		entry.UserName = message.Author.Username
		entry.Year = year
		entry.Week = week
	}
	if utils.ContainsInt(weekdays, 1) {
		entry.Monday = true
	}
	if utils.ContainsInt(weekdays, 2) {
		entry.Tuesday = true
	}
	if utils.ContainsInt(weekdays, 3) {
		entry.Wednesday = true
	}
	if utils.ContainsInt(weekdays, 4) {
		entry.Thursday = true
	}
	if utils.ContainsInt(weekdays, 5) {
		entry.Friday = true
	}
	if utils.ContainsInt(weekdays, 6) {
		entry.Saturday = true
	}
	if utils.ContainsInt(weekdays, 7) {
		entry.Sunday = true
	}

	err = mongoDBClient.UpsertWeek(entry)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return err
	}

	messageString := message.Author.Mention() + " can now play these days for week " + strconv.Itoa(entry.Week) + ":"
	if entry.Monday {
		messageString += "\n" + CheckMark + " Monday"
	}
	if entry.Tuesday {
		messageString += "\n" + CheckMark + " Tuesday"
	}
	if entry.Wednesday {
		messageString += "\n" + CheckMark + " Wednesday"
	}
	if entry.Thursday {
		messageString += "\n" + CheckMark + " Thursday"
	}
	if entry.Friday {
		messageString += "\n" + CheckMark + " Friday"
	}
	if entry.Saturday {
		messageString += "\n" + CheckMark + " Saturday"
	}
	if entry.Sunday {
		messageString += "\n" + CheckMark + " Sunday"
	}
	channelID := message.ChannelID
	userCannelID, err := session.UserChannelCreate(message.Author.ID)
	if err == nil {
		channelID = userCannelID.ID
	}

	botMessage, err := session.ChannelMessageSend(channelID, messageString)
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

func parseAddNextWeekArgs(adds []string) ([]int, error) {
	regexpAdd := regexp.MustCompile(`^[1234567]$`)
	ret := []int{}

	for _, add := range adds {
		if regexpAdd.MatchString(add) {
			day, err := strconv.Atoi(add)
			if err != nil {
				logrus.WithFields(utils.Locate()).Error(err.Error())
				return []int{}, err
			}
			ret = append(ret, day)
		}
	}

	return ret, nil
}
