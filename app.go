package main

import (
	"encoding/json"
	"log"
	"net/http"

	. "darkseid/config"
	. "darkseid/dao"
	. "darkseid/models"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

var config = Config{}
var dao = LocalsDAO{}

func AllLocalsEndPoint(w http.ResponseWriter, r *http.Request) {
	locals, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, locals)
}

func FindLocalEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	local, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Local ID")
		return
	}
	respondWithJson(w, http.StatusOK, local)
}

func CreateLocalEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var local Local
	if err := json.NewDecoder(r.Body).Decode(&local); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	local.ID = bson.NewObjectId()
	if err := dao.Insert(local); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, local)
}

func UpdateLocalEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var local Local
	if err := json.NewDecoder(r.Body).Decode(&local); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(local); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func DeleteLocalEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var local Local
	if err := json.NewDecoder(r.Body).Decode(&local); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(local); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
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

func init() {
	log.Print("Lendo configurações...")
	config.Read()
	dao.Server = config.Server
	dao.Database = config.Database
	log.Print("Conectando ao banco de dados...")
	dao.Connect()
}

func main() {
	log.Print("Conectando ao banco de dados...")
	r := mux.NewRouter()
	r.HandleFunc("/locals", AllLocalsEndPoint).Methods("GET")
	r.HandleFunc("/locals", CreateLocalEndPoint).Methods("POST")
	r.HandleFunc("/locals", UpdateLocalEndPoint).Methods("PUT")
	r.HandleFunc("/locals", DeleteLocalEndPoint).Methods("DELETE")
	r.HandleFunc("/locals/{id}", FindLocalEndPoint).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
