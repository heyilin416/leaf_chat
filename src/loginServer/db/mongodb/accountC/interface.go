package accountC

import (
	"gopkg.in/mgo.v2/bson"
	"loginServer/db/mongodb"
	"gopkg.in/mgo.v2"
	lmongodb "github.com/name5566/leaf/db/mongodb"
)

type AccountData struct {
	Id       bson.ObjectId `bson:"_id"`
	Name     string
	Password string
}

func init() {
	session := mongodb.Context.Ref()
	defer mongodb.Context.UnRef(session)

	GetAccountCollection(session).EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
		Sparse: true,
	})
}

func GetAccountCollection(session *lmongodb.Session) *mgo.Collection {
	return session.DB("login").C("account")
}

func GetAccount(name string) (*AccountData, error) {
	session := mongodb.Context.Ref()
	defer mongodb.Context.UnRef(session)

	result := &AccountData{}
	err := GetAccountCollection(session).Find(bson.M{"name": name}).One(result)
	return result, err
}

func CreateAccount(account *AccountData) error {
	session := mongodb.Context.Ref()
	defer mongodb.Context.UnRef(session)

	return GetAccountCollection(session).Insert(account)
}
