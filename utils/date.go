package utils

import (
	"github.com/Sirupsen/logrus"
	"github.com/beevik/ntp"
)

//GetYearWeek return the current year and week
func GetYearWeek() (int, int, error) {
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		logrus.WithFields(Locate()).Error(err.Error())
		return 0, 0, err
	}
	year, week := time.ISOWeek()
	return year, week, nil
}
