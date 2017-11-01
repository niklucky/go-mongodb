package mongodb

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/*
Mapper - MongoDB mapper
*/
type Mapper struct {
	DBConfig   DBConfig
	Collection string
	Conn       *mgo.Session
}

/*
DBConfig - config for connection
*/
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

/*
New - mapper constructor
*/
func New(config DBConfig, c string) Mapper {
	return Mapper{
		DBConfig:   config,
		Collection: c,
	}
}

/*
Connect - connecting to DB
*/
func (m *Mapper) Connect() error {
	m.log("Connecting to: ", m.getDBInfo())
	session, err := mgo.Dial(m.prepareConnectionString())
	m.Conn = session
	if err != nil {
		return err
	}
	return nil
}

/*
Create - inserting new enity
*/
func (m *Mapper) Create(data interface{}) error {
	if m.Conn == nil {
		err := m.Connect()
		if err != nil {
			fmt.Println("Error connecting: ", err)
			return nil
		}
	}
	c := m.Conn.DB(m.DBConfig.Database).C(m.Collection)
	return c.Insert(data)
}

/*
Search - Searching data in db between fromDate and toDate and substr in "message"
*/
func (m *Mapper) Search(dateFrom, dateTo time.Time, substr string) (interface{}, error) {
	if m.Conn == nil {
		err := m.Connect()
		if err != nil {
			fmt.Println("Error connecting: ", err)
			return nil, err
		}
	}
	c := m.Conn.DB(m.DBConfig.Database).C(m.Collection)
	var logs []interface{}

	if dateTo == time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC) {
		dateTo = time.Now()
	}
	fmt.Println(dateFrom, dateTo, substr)
	var err error
	if substr != "" {
		err = c.Find(
			bson.M{
				"timestamp": bson.M{
					"$gte": dateFrom.Format(time.RFC3339Nano),
					"$lte": dateTo.Format(time.RFC3339Nano),
				},
				"message": &bson.RegEx{Pattern: substr, Options: "i"},
			}).All(&logs)
	} else {
		err = c.Find(
			bson.M{
				"timestamp": bson.M{
					"$gte": dateFrom.Format(time.RFC3339Nano),
					"$lte": dateTo.Format(time.RFC3339Nano),
				},
			}).All(&logs)
	}
	if err != nil {
		fmt.Println("Error request: ", err)
	}
	fmt.Println("Data: ", logs)
	return logs, nil
}

/*
Close - closing connection
*/
func (m *Mapper) Close() error {
	if m.Conn == nil {
		return nil
	}
	m.Conn.Close()
	return nil
}

/*
Converts db config into connection string
[mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
*/
func (m *Mapper) prepareConnectionString() string {
	c := m.DBConfig
	str := "mongodb://"
	if c.User != "" {
		str += c.User + c.Password + "@"
	}
	str += c.Host
	if c.Port != "" {
		str += ":" + c.Port
	}
	if c.Database != "" {
		str += "/" + c.Database
	}
	return str
}

func (m *Mapper) getDBInfo() string {
	c := m.DBConfig
	return "mongodb://" + c.User + "***@" + c.Host + ":" + c.Port + "/" + c.Database
}

func (m *Mapper) log(data ...interface{}) error {
	fmt.Println(data...)
	return nil
}
