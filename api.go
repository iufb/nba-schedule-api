package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/accounts", makeHttpHandleFunc(s.handleAccountWithoutParams))
	router.HandleFunc("/accounts/{id}", makeHttpHandleFunc(s.handleAccountWithParams))
	log.Println("Running on port : ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccountWithParams(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)

	default:
		return fmt.Errorf("Invalid method %s", r.Method)
	}
}

func (s *APIServer) handleAccountWithoutParams(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return s.handleCreateAccount(w, r)
	default:
		return fmt.Errorf("Invalid method %s", r.Method)

	}
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getIdFromParams(r)
	if err != nil {
		return err
	}
	acc, err := s.store.GetAccountById(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, acc)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccRq := &CreateAccountRequest{}
	err := json.NewDecoder(r.Body).Decode(createAccRq)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	account := NewAccount(createAccRq.Username, createAccRq.Timezone)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusCreated, AccountRouteResponse{Status: "Created"})
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getIdFromParams(r)
	if err != nil {
		return err
	}
	_, err = s.store.GetAccountById(id)
	if err != nil {
		return err
	}
	err = s.store.DeleteAccount(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, AccountRouteResponse{Status: "Deleted"})
}

type ApiError struct {
	Error string `json:"error"`
}
type apiFunc func(w http.ResponseWriter, r *http.Request) error

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle errors
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getIdFromParams(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}
