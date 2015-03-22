package blackmagic

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Collection interface {
	All(map[string]interface{}, []interface{}) error
	Find(map[string]interface{}) *mgo.Query
	FindId(bson.ObjectId, interface{}) error
	Insert(interface{}) error
	UpdateId(bson.ObjectId, map[string]interface{}) error
	UpsertId(interface{}, interface{}) (*mgo.ChangeInfo, error)
}

type collection struct {
	store *mgo.Collection
}

func (c *collection) All(query map[string]interface{}, result []interface{}) error {
	err := c.store.Find(query).All(&result)
	return err
}

func (c *collection) Find(query map[string]interface{}) *mgo.Query {
	return c.store.Find(query)
}

func (c *collection) FindId(id bson.ObjectId, result interface{}) error {
	err := c.store.FindId(id).One(&result)
	return err
}

func (c *collection) Insert(newObject interface{}) error {
	err := c.store.Insert(newObject)
	return err
}

func (c *collection) UpdateId(id bson.ObjectId, update map[string]interface{}) error {
	return c.store.UpdateId(id, update)
}

func (c *collection) UpsertId(id interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	return c.store.UpsertId(id, update)
}
