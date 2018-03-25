package main

import (
	"encoding/json"
	"log"
	"net/http"
	"html/template"
	"fmt"
	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	. "github.com/dongik/cards-restapi/config"
	. "github.com/dongik/cards-restapi/dao"
	. "github.com/dongik/cards-restapi/models"
)

var config = Config{}
var dao = CardsDAO{}

type CardPageData struct {
	PageTitle string
	Cards	[]Card
}

//Index EndPoint SHOW list of cards
func IndexEndPoint(w http.ResponseWriter, r *http.Request) {
	cards, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Printf("cards")
	for _, card := range cards{
		fmt.Printf("name:%s, owner:%s, holder:%s\n",
			card.Name,card.Holder,card.Owner)
	}

	data := CardPageData{
		PageTitle: "CardList",
		Cards: cards,
	}
	tmpl := template.Must(template.ParseFiles("static/index.html"))
	tmpl.Execute(w, data)
}


// GET list of cards
func AllCardsEndPoint(w http.ResponseWriter, r *http.Request) {
	cards, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, cards)
}

// GET a card by its ID
func FindCardEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	card, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Card ID")
		return
	}
	respondWithJson(w, http.StatusOK, card)
}



// POST a new card
func CreateCardEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var card Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	card.ID = bson.NewObjectId()
	if err := dao.Insert(card); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//fmt.Fprint("Posted")
	respondWithJson(w, http.StatusCreated, card)
}

// PUT update an existing card
func UpdateCardEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var card Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(card); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing card
func DeleteCardEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var card Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(card); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func ShowTablesEndPoint(w http.ResponseWriter, r *http.Request) {

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

// Define HTTP request routes
func main() {
	r := mux.NewRouter()

	r.HandleFunc("/",IndexEndPoint)

	//restapi
	r.HandleFunc("/cards", AllCardsEndPoint).Methods("GET")
	r.HandleFunc("/cards", CreateCardEndPoint).Methods("POST")
	r.HandleFunc("/cards", UpdateCardEndPoint).Methods("PUT")
	r.HandleFunc("/cards", DeleteCardEndPoint).Methods("DELETE")
	r.HandleFunc("/cards/{id}", FindCardEndPoint).Methods("GET")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
