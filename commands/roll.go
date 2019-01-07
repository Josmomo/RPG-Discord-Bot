package commands

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/Josmomo/RPG-Discord-Bot/utils"
	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
)

//CommandRoll
const CommandRoll = "roll"

//Roll makes a dice roll from input and writes a message in channel
func Roll(session *discordgo.Session, message *discordgo.MessageCreate, args []string) error {
	rolls, err := parseRollArgs(args)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return err
	}
	messageString := ""
	noRolls := len(rolls)
	if noRolls > 1 {
		for _, roll := range rolls[:noRolls-1] {
			messageString += " " + strconv.Itoa(roll) + ","
		}
		messageString += " " + strconv.Itoa(rolls[len(rolls)-1])
	} else if noRolls == 1 {
		messageString = strconv.Itoa(rolls[0])
	} else {
		messageString = "No dices to roll"
	}

	_, err = session.ChannelMessageSend(message.ChannelID, message.Author.Mention()+":game_die:"+messageString)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return err
	}
	return nil
}

func parseRollArgs(rolls []string) ([]int, error) {
	regexpRoll := regexp.MustCompile(`^\d\d?D\d+$`)
	ret := []int{}

	if len(rolls) == 0 {
		rolls = append(rolls, "1D20")
	}

	for _, roll := range rolls {
		if regexpRoll.MatchString(roll) {
			s := strings.Split(strings.TrimSpace(roll), "D")
			numberOfThrows, err := strconv.Atoi(s[0])
			if err != nil {
				logrus.WithFields(utils.Locate()).Error(err.Error())
				return []int{}, err
			}
			diceType, err := strconv.Atoi(s[1])
			if err != nil {
				logrus.WithFields(utils.Locate()).Error(err.Error())
				return []int{}, err
			}
			ret = append(ret, numberOfThrows*(rand.Intn(diceType)+1))
		}
	}

	return ret, nil
}
