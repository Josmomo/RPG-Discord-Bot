package commands

import (
	"regexp"
	"strconv"

	"github.com/Josmomo/RPG-Discord-Bot/database"
	"github.com/Josmomo/RPG-Discord-Bot/utils"
	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
)

//CommandAdd
const CommandAdd = "add"

//Add
func Add(mongoDBClient database.MongoDBClient, session *discordgo.Session, message *discordgo.MessageCreate, args []string) error {
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
	if contains(weekdays, 1) {
		entry.Monday = true
	}
	if contains(weekdays, 2) {
		entry.Tuesday = true
	}
	if contains(weekdays, 3) {
		entry.Wednesday = true
	}
	if contains(weekdays, 4) {
		entry.Thursday = true
	}
	if contains(weekdays, 5) {
		entry.Friday = true
	}
	if contains(weekdays, 6) {
		entry.Saturday = true
	}
	if contains(weekdays, 7) {
		entry.Sunday = true
	}

	err = mongoDBClient.UpsertWeek(entry)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return err
	}

	return nil
}

func parseAddArgs(adds []string) ([]int, error) {
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

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
