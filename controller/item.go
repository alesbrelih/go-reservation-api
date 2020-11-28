package controller

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/alesbrelih/go-reservation-api/db"
	"github.com/alesbrelih/go-reservation-api/interfaces"
	"github.com/alesbrelih/go-reservation-api/middleware"
	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/alesbrelih/go-reservation-api/pkg/myutil"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var validate = validator.New()

type ItemController interface {
	interfaces.Controller
}

type DefaultItemController struct{}

func (h *DefaultItemController) GetAll(w http.ResponseWriter, req *http.Request) {
	myDb := db.Connect()
	defer myDb.Close()

	items, err := h.getItemsFromDb(myDb)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	items.ToJSON(w)
}

func (h *DefaultItemController) getItemsFromDb(myDb *sql.DB) (models.Items, error) {
	items := []models.Item{}

	query := "SELECT * FROM item"

	rows, err := myDb.Query(query)

	if err != nil {
		log.Printf("Query error: %v, err: %v", query, err)
		return nil, err
	}

	for rows.Next() {
		var item models.Item

		err = rows.Scan(&item.Id, &item.Title, &item.ShowFrom, &item.ShowTo)
		if err != nil {
			log.Printf("Query scan error: %v", err)
			return nil, err
		}

		items = append(items, item)
	}
	return items, nil

}

func (h *DefaultItemController) GetOne(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Printf("Can't parse id parameter to int: %v. Error: %v", params["id"], err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	myDb := db.Connect()
	defer myDb.Close()

	var item models.Item
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

	// func (u User) MarshalJSON() ([]byte, error) {
	// 	type user User // prevent recursion < transform to type so it doesnt use recursion!
	// 	x := user(u) -> because uses MarshalJson is called by json.Marshal
	// 	x.Password = ""
	// 	return json.Marshal(x)
	// }
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

	// .( ) <- type assertion
	item := req.Context().Value(&middleware.ItemBodyKeyType{}).(*models.Item)

	err := validate.Struct(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	w.WriteHeader(http.StatusCreated)
}

func (h *DefaultItemController) Update(w http.ResponseWriter, r *http.Request) {

	item := r.Context().Value(&middleware.ItemBodyKeyType{}).(*models.Item)

	err := validate.Struct(item)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	myDb := db.Connect()
	defer myDb.Close()

	stmt := "UPDATE item SET title=$1, show_from=$2, show_to=$3 WHERE id = $4"
	res, err := myDb.Exec(stmt, item.Title, item.ShowFrom, item.ShowTo, item.Id)
	if err != nil {
		log.Printf("Error saving to db: %v", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = myutil.ValidateRowsAffected(res, w, &log.Logger{})
	if err != nil {
		return
	}
}

func (h *DefaultItemController) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Printf("Can't parse id parameter to int: %v. Error: %v", params["id"], err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	myDb := db.Connect()
	defer myDb.Close()

	stmt := "DELETE FROM item WHERE id = $1"
	res, err := myDb.Exec(stmt, int64(id))
	if err != nil {
		log.Printf("Error deleting item id %v from db. Error: %v", id, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	err = myutil.ValidateRowsAffected(res, w, &log.Logger{})
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func NewItemRouter(controller ItemController) *mux.Router {
	r := mux.NewRouter()

	middleware := middleware.NewItemMiddleware(&log.Logger{})

	getSubrouter := r.Methods(http.MethodGet).Subrouter()
	getSubrouter.HandleFunc("/item", controller.GetAll)
	getSubrouter.HandleFunc("/item/{id:[\\d]+}", controller.GetOne)

	postSubrouter := r.Methods(http.MethodPost).Subrouter()
	postSubrouter.HandleFunc("/item", controller.Create)
	postSubrouter.Use(middleware.GetBody)

	putSubrouter := r.Methods(http.MethodPut).Subrouter()
	putSubrouter.HandleFunc("/item", controller.Update)
	putSubrouter.Use(middleware.GetBody)

	deleteSubgrouter := r.Methods(http.MethodDelete).Subrouter()
	deleteSubgrouter.HandleFunc("/item/{id:[\\d]+}", controller.Delete)

	return r
}
