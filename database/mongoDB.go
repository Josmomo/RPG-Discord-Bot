package database

import (
	"crypto/tls"
	"net"

	"github.com/Josmomo/RPG-Discord-Bot/constants"
	"github.com/Josmomo/RPG-Discord-Bot/utils"
	"github.com/Sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//MongoDBClient ...
type MongoDBClient struct {
	Session *mgo.Session
}

type ScheduleWeek struct {
	UserID    string `bson:"userID"`
	Year      int    `bson:"year"`
	Week      int    `bson:"week"`
	Monday    bool   `bson:"monday"`
	Tuesday   bool   `bson:"tuesday"`
	Wednesday bool   `bson:"wednesday"`
	Thursday  bool   `bson:"thursday"`
	Friday    bool   `bson:"friday"`
	Saturday  bool   `bson:"saturday"`
	Sunday    bool   `bson:"sunday"`
}

type ScheduleWeekFindQuery struct {
	userID string `bson:"userID"`
	year   int    `bson:"year"`
	week   int    `bson:"week"`
}

//CreateNewClient creates a new client
func CreateNewClient() (MongoDBClient, error) {
	client := MongoDBClient{}
	tlsConfig := &tls.Config{}
	dialInfo := &mgo.DialInfo{
		Addrs: []string{"cluster0-shard-00-00-gvkr3.mongodb.net:27017",
			"cluster0-shard-00-01-gvkr3.mongodb.net:27017",
			"cluster0-shard-00-02-gvkr3.mongodb.net:27017"},
		Database: "admin",
		Username: "Josmomo",
		Password: "mQK0RPhnqveD9GntfpVi",
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return client, err
	}
	client.Session = session

	return client, nil
}

func (mongoDBClient *MongoDBClient) UpsertWeek(entry ScheduleWeek) error {
	query := ScheduleWeekFindQuery{userID: entry.UserID, year: entry.Year, week: entry.Week}
	_, err := mongoDBClient.Session.DB(constants.DataBaseName).C(constants.ScheduleCollection).Upsert(query, bson.M{"$set": entry})
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return err
	}
	return nil
}

func (mongoDBClient *MongoDBClient) GetWeek(query ScheduleWeek) ([]ScheduleWeek, error) {
	ret := []ScheduleWeek{}
	err := mongoDBClient.Session.DB(constants.DataBaseName).C(constants.ScheduleCollection).Find(query).All(&ret)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return ret, err
	}
	return ret, nil
}

//GetDocFromIndex returns a single doc matching the index
func (mongoDBClient *MongoDBClient) GetDocFromIndex(uID string, y int, w int) (ScheduleWeek, error) {
	ret := ScheduleWeek{}
	query := ScheduleWeekFindQuery{userID: uID, year: y, week: w}
	err := mongoDBClient.Session.DB(constants.DataBaseName).C(constants.ScheduleCollection).Find(query).One(&ret)
	if err != nil {
		logrus.WithFields(utils.Locate()).Error(err.Error())
		return ret, err
	}
	return ret, nil
}
