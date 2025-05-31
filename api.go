package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type API struct {
	listenAddr string
	store  	   Storage
}

func NewAPI(listenAddr string, store Storage) *API {
	return &API{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *API) Start()  {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHttpHandlerFunc(s.handleLogin))
	router.HandleFunc("/account", makeHttpHandlerFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", JWTAuth(makeHttpHandlerFunc(s.handleGetAccountByID), s.store))
	router.HandleFunc("/transfer" , makeHttpHandlerFunc(s.handleTransfer))
	
	log.Println("Bank API running on port:", s.listenAddr)
	
	http.ListenAndServe(s.listenAddr, router)
}



func (s *API) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	
	case "GET":
		return s.handleGetAccount(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	}

	return fmt.Errorf("methode not allowed: %s", r.Method)
}


func (s *API) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed: %s", r.Method)
	}
	var loginReq LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		return err 
	}
	account, err := s.store.GetAccountByNumber(int(loginReq.Number))
	
	if err != nil {
		return nil 
	}
	
	if !account.validPassword(loginReq.Password) {
		return fmt.Errorf("invalid password")
	}
	
	token,err := createJWT(account)
	
	if err != nil {
		return err 
	}
	
	loginResp := LoginResponse{
		Token: token,
		Number: account.Number,
	}
	return WriteJSON(w, http.StatusOK, loginResp)
}




func (s *API) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAcccountReq:= new (createAcccountRequest)
	
	if err := json.NewDecoder(r.Body).Decode(createAcccountReq); err != nil {
		return err
	}

	account,err := NewAccount(createAcccountReq.FirstName, createAcccountReq.LastName , createAcccountReq.Password)

	if err != nil {
		return err
	}
	
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, account)
}



func (s *API) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts , err := s.store.GetAllAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}



func (s *API) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id , err := getID(r)
		
		if err != nil {
			return err
		}
		
		account, err := s.store.GetAccountByID(id)
		
		if err != nil {
			return err 
		}
		return WriteJSON(w, http.StatusOK, account)
	}
	
	if r.Method == "PUT" {
		return s.handleUpdateAccount(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	
	return fmt.Errorf("method not allowed: %s", r.Method)
}



func (s *API) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error{
	id ,err := getID(r)
	
	if err != nil {
		return err
	}

	var account Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return err
	}
	
	account.ID = id

	if err := s.store.UpdateAccount(id, &account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"account updated": id})
}


func (s *API) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id , err := getID(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK , map[string]int{"account deleted": id})
}


func (s *API) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	TransferRequest := new(TransferRequest)

	if err := json.NewDecoder(r.Body).Decode(TransferRequest); err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJSON(w , http.StatusOK, TransferRequest)
}


func WriteJSON(w http.ResponseWriter , status int, v any ) error {
	w.Header().Add("Content-Type", "application/json")
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


func createJWT (account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt" : 15000 ,
		"accountNumber" : account.Number,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}


func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, apiError{Error: "permission denied"})
}


func JWTAuth(handlerFunc http.HandlerFunc , s Storage ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("JWTAuth middleware called")
		TokenString := r.Header.Get("x-jwt-token")
		
		token,err := verifyJWT(TokenString)
		
		if err != nil {
			permissionDenied(w)
			return
		}
		
		if !token.Valid {
			permissionDenied(w)
			return
		}
		
		userID , err := getID(r)
		
		if err != nil {
			permissionDenied(w)
			return
		}
		
		account, err := s.GetAccountByID(userID)
		
		if err != nil {
			permissionDenied(w)
			return
		}
		
		claims := token.Claims.(jwt.MapClaims)

		if account.Number != int64(claims["accountNumber"].(float64)) {
			permissionDenied(w)
			return
		}

		handlerFunc(w, r) 
	}
}		

func verifyJWT(TokenString string) (*jwt.Token , error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(TokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}


func getID (r *http.Request) (int , error) {
	idString :=mux.Vars(r)["id"]
	
	id, err := strconv.Atoi(idString) 

	if err != nil {
		return id , fmt.Errorf("invalid id: %s", idString)
	}
	return id, nil
}

