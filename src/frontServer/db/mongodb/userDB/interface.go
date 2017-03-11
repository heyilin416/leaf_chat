package userDB

import (
	"gopkg.in/mgo.v2/bson"
	"frontServer/db/mongodb"
	lmongodb "github.com/name5566/leaf/db/mongodb"
	"gopkg.in/mgo.v2"
)

type UserData struct {
	Id        bson.ObjectId `bson:"_id"`
	AccountId bson.ObjectId
	Name      string
}

func init() {
	session := mongodb.Context.Ref()
	defer mongodb.Context.UnRef(session)

	getCollection(session).EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
		Sparse: true,
	})
}

func getCollection(session *lmongodb.Session) *mgo.Collection {
	return session.DB("game").C("user")
}

func Get(AccountId bson.ObjectId) (*UserData, error) {
	session := mongodb.Context.Ref()
	defer mongodb.Context.UnRef(session)

	result := &UserData{}
	err := getCollection(session).Find(bson.M{"accountid": AccountId}).One(result)
	return result, err
}

func Create(user *UserData) error {
	session := mongodb.Context.Ref()
	defer mongodb.Context.UnRef(session)

	return getCollection(session).Insert(user)
}
