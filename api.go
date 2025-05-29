package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)



func WriteJSON(w http.ResponseWriter , status int, v any ) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type apiError struct {
	Error  string
}	


func makeHttpHandlerFunc( f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, apiError{Error: err.Error(),} )
		}
	}
}

type API struct {
	listenAddr string
}

func NewAPI(listenAddr string) *API {
	return &API{
		listenAddr: listenAddr,
	}
}

func (s *API) Start()  {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHttpHandlerFunc(s.handleAccount))
	
	router.HandleFunc("/account/{id}", makeHttpHandlerFunc(s.handleGetAccount))
	
	log.Println("Bank API running on port:", s.listenAddr)
	
	http.ListenAndServe(s.listenAddr, router)
}


func (s *API) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("methode not allowed: %s", r.Method)
}

func (s *API) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r) ["id"]
	fmt.Println(id)
	// Here you would typically fetch the account from a database !
	return WriteJSON(w , http.StatusOK , &Account{})
}

func (s *API) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *API) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *API) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}