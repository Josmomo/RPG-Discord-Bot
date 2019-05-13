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

//CommandRemoveNextWeek
const CommandRemoveNextWeek = "removeNextWeek"

//RemoveNextWeek
func RemoveNextWeek(mongoDBClient database.MongoDBClient, session *discordgo.Session, message *discordgo.MessageCreate, args []string) error {
	weekdays, err := parseRemoveNextWeekArgs(args)
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
	entry.UserName = message.Author.Username
	if utils.ContainsInt(weekdays, 1) {
		entry.Monday = false
	}
	if utils.ContainsInt(weekdays, 2) {
		entry.Tuesday = false
	}
	if utils.ContainsInt(weekdays, 3) {
		entry.Wednesday = false
	}
	if utils.ContainsInt(weekdays, 4) {
		entry.Thursday = false
	}
	if utils.ContainsInt(weekdays, 5) {
		entry.Friday = false
	}
	if utils.ContainsInt(weekdays, 6) {
		entry.Saturday = false
	}
	if utils.ContainsInt(weekdays, 7) {
		entry.Sunday = false
	}

	err = mongoDBClient.UpsertWeek(entry)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return err
	}

	messageString := message.Author.Mention() + " removed some days and can now play these days for week " + strconv.Itoa(entry.Week) + ":"
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

func parseRemoveNextWeekArgs(removes []string) ([]int, error) {
	regexpRemove := regexp.MustCompile(`^[1234567]$`)
	ret := []int{}

	for _, remove := range removes {
		if regexpRemove.MatchString(remove) {
			day, err := strconv.Atoi(remove)
			if err != nil {
				logrus.WithFields(utils.Locate()).Error(err.Error())
				return []int{}, err
			}
			ret = append(ret, day)
		}
	}

	return ret, nil
}
