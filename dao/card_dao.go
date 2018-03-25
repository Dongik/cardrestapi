package dao

import (
	"log"

	. "github.com/mlabouardy/cards-restapi/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CardsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "cards"
)

// Establish a connection to database
func (m *CardsDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of cards
func (m *CardsDAO) FindAll() ([]Movie, error) {
	var cards []Card
	err := db.C(COLLECTION).Find(bson.M{}).All(&cards)
	return cards, err
}

// Find a card by its id
func (m *CardsDAO) FindById(id string) (Movie, error) {
	var card Card
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&card)
	return card, err
}

// Insert a card into database
func (m *CardsDAO) Insert(card Movie) error {
	err := db.C(COLLECTION).Insert(&card)
	return err
}

// Delete an existing card
func (m *CardsDAO) Delete(card Movie) error {
	err := db.C(COLLECTION).Remove(&card)
	return err
}

// Update an existing card
func (m *CardsDAO) Update(card Movie) error {
	err := db.C(COLLECTION).UpdateId(card.ID, &movie)
	return err
}
