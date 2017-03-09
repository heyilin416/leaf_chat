package accountC

import (
	"gopkg.in/mgo.v2/bson"
	"loginServer/db/mongodb"
)

type Account struct {
	Id       bson.ObjectId `bson:"_id"`
	UserName string
	Password string
}

func HasAccount(userName string) (*Account, error) {
	session := mongodb.Context.Ref()
	defer mongodb.Context.UnRef(session)

	result := &Account{}
	collection := session.DB("login").C("account")
	err := collection.Find(bson.M{"username": userName}).One(result)
	return result, err
}

func CreateAccount(account *Account) error {
	session := mongodb.Context.Ref()
	defer mongodb.Context.UnRef(session)

	collection := session.DB("login").C("account")
	return collection.Insert(account)
}