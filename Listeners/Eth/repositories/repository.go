package repositories

import (
	"github.com/globalsign/mgo"
	"log"
	"time"
)

type IRepository interface {
}

type repository struct {
	dbConnect *mgo.Database
}

func New(url, dbName string) (IRepository, error) {
	var err error
	var conn *mgo.Session
	for i := 1; i < 6; i++ {
		log.Printf("Attempt %d to connect to MongoDB at %s", i, url)
		conn, err = mgo.Dial(url)
		if err != nil {
			log.Printf("Failed to connect to MongoDB at %s: %v", url, err)
			time.Sleep(1 * time.Second)
		} else {
			log.Println("Connected successfully to MongoDB at", url)
			return &repository{dbConnect: conn.DB(dbName)}, nil
		}
	}
	return nil, err
}
