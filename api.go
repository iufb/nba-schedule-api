package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	router.HandleFunc("/login", makeHttpHandleFunc(s.handleLogin))
	router.HandleFunc("/register", makeHttpHandleFunc(s.handleRegister))
	router.HandleFunc("/accounts", makeHttpHandleFunc(s.handleAccountWithoutParams))
	router.HandleFunc("/accounts/{id}", AuthGuard(makeHttpHandleFunc(s.handleAccountWithParams)))
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
	// case "POST":
	// 	return s.handleCreateAccount(w, r)
	default:
		return fmt.Errorf("Invalid method %s", r.Method)
	}
}

func (s *APIServer) handleRegister(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("Unallowed method %s : ", r.Method)
	}

	registerRq := &RegisterRequest{}
	err := BodyDecoder(registerRq, r.Body)
	if err != nil {
		return err
	}
	acc, err := NewAccount(registerRq.Username, registerRq.Password)
	if err != nil {
		return err
	}
	if err := s.store.CreateAccount(acc); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusCreated, WithStatusResponse{Status: "Registered successfully."})
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("Unallowed method %s : ", r.Method)
	}
	loginRq := &LoginRequest{}
	err := BodyDecoder(loginRq, r.Body)
	if err != nil {
		return err
	}
	acc, err := s.store.GetAccountByUsername(loginRq.Username)
	if err != nil {
		return err
	}
	if isValid := acc.ValidateAccount(loginRq.Password); isValid == false {
		return WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Invalid password"})
	}
	token, err := CreateJWT(acc)
	if err != nil {
		return err
	}
	authCookie := http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
	}
	http.SetCookie(w, &authCookie)
	return WriteJSON(w, http.StatusOK, WithStatusResponse{Status: "Logged"})
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

// func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
// 	createAccRq := &CreateAccountRequest{}
// 	err := BodyDecoder(createAccRq, r.Body)
// 	if err != nil {
// 		return err
// 	}
// 	account := NewAccount(createAccRq.Username, createAccRq.Timezone)
// 	if err := s.store.CreateAccount(account); err != nil {
// 		return err
// 	}
// 	return WriteJSON(w, http.StatusCreated, WithStatusResponse{Status: "Created"})
// }

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
	return WriteJSON(w, http.StatusOK, WithStatusResponse{Status: "Deleted"})
}

type ApiError struct {
	Error string `json:"error"`
}
type apiFunc func(w http.ResponseWriter, r *http.Request) error

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ipAddr := r.RemoteAddr

		fmt.Println("IP Address:", ipAddr)
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

func BodyDecoder(v interface{}, body io.ReadCloser) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(v)
}
