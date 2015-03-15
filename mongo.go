package blackmagic

import (
	"gopkg.in/mgo.v2"
)

type Database interface {
	Connect() error
	Collection(string) (Collection, error)
}

type database struct {
	connectionString string
	session          *mgo.Session
	store            *mgo.Database
	databaseName     string
}

func NewDatabase(mongoURL string, dbName string) *database {
	return &database{connectionString: mongoURL, databaseName: dbName}
}

func (db *database) Connect() error {
	s, err := mgo.Dial(db.connectionString)
	if err != nil {
		return err
	}

	db.session = s

	db.store = s.DB(db.databaseName)
	return nil
}

func (db *database) Collection(cName string) (Collection, error) {
	c := db.store.C(cName)
	return &collection{store: c}, nil
}
