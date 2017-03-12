package tokenDB

import (
	"gopkg.in/mgo.v2/bson"
	"loginServer/db/mongodb"
	"gopkg.in/mgo.v2"
	lmongodb "github.com/name5566/leaf/db/mongodb"
)

type Data struct {
	Token     bson.ObjectId `bson:"_id"`
	AccountId bson.ObjectId
	FrontName string
}

func getCollection(session *lmongodb.Session) *mgo.Collection {
	return session.DB("login").C("token")
}

func Create(accountId bson.ObjectId, frontName string) (token bson.ObjectId, err error) {
	session := mongodb.Context.Ref()
	defer mongodb.Context.UnRef(session)

	token = bson.NewObjectId()
	data := &Data{Token: token, AccountId: accountId, FrontName: frontName}
	return token, getCollection(session).Insert(data)
}

func Check(token bson.ObjectId, frontName string) (accountId bson.ObjectId, err error) {
	session := mongodb.Context.Ref()
	defer mongodb.Context.UnRef(session)

	result := &Data{}
	collection := getCollection(session)
	err = collection.FindId(token).One(result)
	if err == nil {
		collection.RemoveId(token)
		if result.FrontName == frontName {
			accountId = result.AccountId
		}
	}
	return
}
