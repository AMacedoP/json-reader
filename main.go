package main

import (
	_ "embed"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	DatabasePassword string `json:"database_password"`
}

type Response struct {
	DatabasePassword string `json:"password"`
}

var UsersJson Users

//go:embed users.json
var jsonString []byte

func readJson() (users Users) {
	err := json.Unmarshal(jsonString, &users)
	if err != nil {
		log.Fatalf("Failed to unmarshal json")
	}

	return
}

func getCredentials(w http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var userSent User
	err = json.Unmarshal(reqBody, &userSent)
	if err != nil {
		log.Println("Failed to unmarshal json, user error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i := range UsersJson.Users {
		if UsersJson.Users[i].Username == userSent.Username {
			var userFound = UsersJson.Users[i]
			if userFound.Password == userSent.Password {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(&Response{userFound.DatabasePassword})
				return
			}
		}
	}

	log.Println("User or password combination not found")
	w.WriteHeader(http.StatusBadRequest)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/get_credentials", getCredentials)

	log.Fatal(http.ListenAndServe(":80", myRouter))
}

func main() {
	UsersJson = readJson()
	handleRequests()
}
