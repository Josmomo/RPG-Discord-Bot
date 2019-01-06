package utils

import (
	"github.com/Sirupsen/logrus"
	"github.com/beevik/ntp"
)

//GetYearWeek returns the current year and week
func GetYearWeek() (int, int, error) {
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		logrus.WithFields(Locate()).Error(err.Error())
		return 0, 0, err
	}
	year, week := time.ISOWeek()
	return year, week, nil
}

//GetYearWeekNext returns the current year and week plus one week
func GetYearWeekNext() (int, int, error) {
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		logrus.WithFields(Locate()).Error(err.Error())
		return 0, 0, err
	}
	year, week := time.AddDate(0, 0, 7).ISOWeek()
	return year, week, nil
}
