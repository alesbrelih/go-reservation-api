package item

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/alesbrelih/go-reservation-api/db"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type ItemController interface {
	GetAll(w http.ResponseWriter, req *http.Request)
	GetOne(w http.ResponseWriter, req *http.Request)
	Create(w http.ResponseWriter, req *http.Request)
}

type DefaultItemController struct{}

func (h *DefaultItemController) GetAll(w http.ResponseWriter, req *http.Request) {
	myDb := db.Connect()
	defer myDb.Close()

	items := []Item{}

	query := "SELECT * FROM item"

	rows, err := myDb.Query(query)

	if err != nil {
		log.Printf("Query error: %v, err: %v", query, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	for rows.Next() {
		var item Item

		err = rows.Scan(&item.Id, &item.Title, &item.ShowFrom, &item.ShowTo)
		if err != nil {
			log.Printf("Query scan error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}

		items = append(items, item)
	}

	itemsJson, err := json.Marshal(items)

	if err != nil {
		log.Printf("Json marshal %v. Error: %v", items, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write(itemsJson)
}

func (h *DefaultItemController) GetOne(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Printf("Can't parse id parameter to int: %v. Error: ", params["id"], err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	myDb := db.Connect()
	defer myDb.Close()

	var item Item
	stmt := "SELECT * FROM item WHERE id = $1"
	row := myDb.QueryRow(stmt, int64(id))
	err = row.Scan(&item.Id, &item.Title, &item.ShowFrom, &item.ShowTo)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
			w.Write(nil)
			return
		default:
			log.Printf("Internal server error: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
	}

	itemJson, err := json.Marshal(item)
	if err != nil {
		log.Printf("Json marshal %v. Error: %v", item, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write([]byte(itemJson))
}

func (h *DefaultItemController) Create(w http.ResponseWriter, req *http.Request) {
	itemBytes, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Printf("Cant read Create item body: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	var item Item
	err = json.Unmarshal(itemBytes, &item)
	if err != nil {
		log.Printf("Cannot unmarshal item: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	if item.Title == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing property: Title"))
		return
	}

	if item.ShowFrom == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing property: ShowFrom"))
		return
	}

	if item.ShowTo == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing property: ShowTo"))
		return
	}

	myDb := db.Connect()
	defer myDb.Close()

	stmt := "INSERT INTO item (title, show_from, show_to) VALUES ($1, $2, $3)"
	_, err = myDb.Exec(stmt, item.Title, item.ShowFrom, item.ShowTo)

	if err != nil {
		log.Printf("Error saving to db: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func Router(controller ItemController) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/item", controller.GetAll).Methods("GET")
	r.HandleFunc("/item/{id:[\\d]+}", controller.GetOne).Methods("GET")
	r.HandleFunc("/item", controller.Create).Methods("POST")
	return r
}
