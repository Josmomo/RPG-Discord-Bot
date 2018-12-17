package client

import (
	"github.com/Josmomo/RPG-Discord-Bot/constants"
	"github.com/Josmomo/RPG-Discord-Bot/utils"
	"github.com/Sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
)

type MongoDBClient struct {
	session *mgo.Session
}

type ScheduleWeek struct {
	userID    string `bson:"userID"`
	year      int    `bson:"year"`
	week      int    `bson:"week"`
	monday    bool   `bson:"monday"`
	tuesday   bool   `bson:"tuesday"`
	wednesday bool   `bson:"wednesday"`
	thursday  bool   `bson:"thursday"`
	friday    bool   `bson:"friday"`
	saturday  bool   `bson:"saturday"`
	sunday    bool   `bson:"sunday"`
}

type ScheduleWeekFindQuery struct {
	userID string `bson:"userID"`
	year   int    `bson:"year"`
	week   int    `bson:"week"`
}

//CreateNewClient Creates a new client
func CreateNewClient() (MongoDBClient, error) {
	client := MongoDBClient{}
	session, err := mgo.Dial("cluster0-gvkr3.mongodb.net")
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return client, err
	}
	client.session = session

	return client, nil
}

func (mongoDBClient *MongoDBClient) UpsertWeek(entry ScheduleWeek) error {
	query := ScheduleWeekFindQuery{userID: entry.userID, year: entry.year, week: entry.week}
	_, err := mongoDBClient.session.DB(constants.DataBaseName).C(constants.ScheduleCollection).Upsert(query, entry)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return err
	}
	return nil
}

func (mongoDBClient *MongoDBClient) GetWeek(query ScheduleWeek) ([]ScheduleWeek, error) {
	ret := []ScheduleWeek{}
	err := mongoDBClient.session.DB(constants.DataBaseName).C(constants.ScheduleCollection).Find(query).All(&ret)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return ret, err
	}
	return ret, nil
}
