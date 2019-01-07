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

//CommandCheckWeek
const CommandCheckWeek = "checkWeek"

//CheckWeek
func CheckWeek(mongoDBClient database.MongoDBClient, session *discordgo.Session, message *discordgo.MessageCreate, args []string) error {
	year, week, err := utils.GetYearWeek()
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		t := time.Now()
		year, week = t.ISOWeek()
	}
	weeks, err := parseCheckWeekArgs(args)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return err
	}
	if len(weeks) > 0 {
		week = weeks[0]
	}
	entries, err := mongoDBClient.GetDocsFromIndex(year, week)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
	}
	messageString := "Week " + strconv.Itoa(week)
	stringMonday := "\nMonday:"
	stringTuesday := "\nTuesday:"
	stringWednesday := "\nWednesday:"
	stringThursday := "\nThursday:"
	stringFriday := "\nFriday:"
	stringSaturday := "\nSaturday:"
	stringSunday := "\nSunday:"
	for _, entry := range entries {
		s := " " + entry.UserName
		if entry.Monday {
			stringMonday += s
		}
		if entry.Tuesday {
			stringTuesday += s
		}
		if entry.Wednesday {
			stringWednesday += s
		}
		if entry.Thursday {
			stringThursday += s
		}
		if entry.Friday {
			stringFriday += s
		}
		if entry.Saturday {
			stringSaturday += s
		}
		if entry.Sunday {
			stringSunday += s
		}
	}
	messageString += stringMonday
	messageString += stringTuesday
	messageString += stringWednesday
	messageString += stringThursday
	messageString += stringFriday
	messageString += stringSaturday
	messageString += stringSunday
	_, err = session.ChannelMessageSend(message.ChannelID, messageString)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return err
	}
	return nil
}

func parseCheckWeekArgs(weeks []string) ([]int, error) {
	regexpCheckWeek := regexp.MustCompile(`^(\d\d?)$`)
	ret := []int{}

	for _, week := range weeks {
		if regexpCheckWeek.MatchString(week) {
			day, err := strconv.Atoi(week)
			if err != nil {
				logrus.WithFields(utils.Locate()).Error(err.Error())
				return []int{}, err
			}
			ret = append(ret, day)
		}
	}

	return ret, nil
}
