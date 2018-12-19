package commands

import (
	"regexp"
	"strconv"

	"github.com/Josmomo/RPG-Discord-Bot/database"
	"github.com/Josmomo/RPG-Discord-Bot/utils"
	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
)

//CommandRemove
const CommandRemove = "remove"

//Remove
func Remove(mongoDBClient database.MongoDBClient, session *discordgo.Session, message *discordgo.MessageCreate, args []string) error {
	weekdays, err := parseAddArgs(args)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return err
	}
	year, week, err := utils.GetYearWeek()
	entry, err := mongoDBClient.GetDocFromIndex(message.Author.Mention(), year, week)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		entry.UserID = message.Author.Mention()
		entry.Year = year
		entry.Week = week
		//return err
	}
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

	return nil
}

func parseRemoveArgs(adds []string) ([]int, error) {
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
