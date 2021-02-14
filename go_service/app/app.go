package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/andrewawni/chatsystem/internal/rabbit"
	"github.com/andrewawni/chatsystem/internal/redis"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// App type
type App struct {
	router       *mux.Router
	rabbitClient *rabbit.Client
	redisClient  *redis.Client
}

type applicationType struct {
	ApplicationToken string `json:"application_token"`
	ApplicationName  string `json:"application_name"`
}

type chatType struct {
	ApplicationToken string `json:"application_token"`
	ChatNumber       int    `json:"chat_number"`
	ChatName         string `json:"chat_name"`
}

type messageType struct {
	ApplicationToken string `json:"application_token"`
	ChatNumber       int    `json:"chat_number"`
	MessageNumber    int    `json:"message_number"`
	MessageContent   string `json:"message_content"`
}

// Init - initialize application
func (app *App) Init(amqpURL string, redisURL string) {
	app.router = mux.NewRouter()
	rabbitClient, err := rabbit.CreateClient(amqpURL)
	if err != nil {
		log.Fatal(err)
	}
	app.rabbitClient = rabbitClient
	redisClient, err := redis.CreateClient(redisURL)
	if err != nil {
		log.Fatal(err)
	}
	app.redisClient = redisClient
	app.initRoutes()
}

// Run - run application
func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(":8000", app.router))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *App) initRoutes() {
	// app.router.HandleFunc("/api/", indexHandler).Methods("POST")
	app.router.HandleFunc("/api/applications", app.createApplication).Methods("POST")
	app.router.HandleFunc("/api/applications/{application_token}/chats", app.createChat).Methods("POST")
	app.router.HandleFunc("/api/applications/{application_token}/chats/{chat_number}/messages", app.createMessage).Methods("POST")
}

func (app *App) createApplication(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	token := uuid.New().String()
	err := app.redisClient.SetInt(token, 0)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	// body := parseBody(r)
	var body applicationType
	err = json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		panic(err)
	}

	application := applicationType{
		ApplicationToken: token,
		ApplicationName:  body.ApplicationName,
	}

	app.rabbitClient.Publish("applications", application)
	respondWithJSON(w, http.StatusCreated, application)
}

func (app *App) createChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["application_token"]
	if token == "" {
		respondWithError(w, http.StatusBadRequest, "no token provided")
	}
	n, err := app.redisClient.IncrementIntAndSetNewKey(token, 0)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid token")
	}
	var body chatType
	err = json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		panic(err)
	}

	chat := chatType{
		ApplicationToken: token,
		ChatNumber:       n,
		ChatName:         body.ChatName,
	}

	app.rabbitClient.Publish("applications", chat)
	respondWithJSON(w, http.StatusCreated, chat)
}

func (app *App) createMessage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	token := vars["application_token"]
	if token == "" {
		respondWithError(w, http.StatusBadRequest, "no token provided")
	}
	chatNumber, err := strconv.Atoi(vars["chat_number"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chat number")
	}

	key := fmt.Sprint(token, ":", chatNumber)
	n, err := app.redisClient.IncrementInt(key)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid token")
	}

	var body messageType
	err = json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		panic(err)
	}

	message := messageType{
		ApplicationToken: token,
		ChatNumber:       chatNumber,
		MessageNumber:    n,
		MessageContent:   body.MessageContent,
	}
	app.rabbitClient.Publish("applications", message)
	respondWithJSON(w, http.StatusCreated, message)
}
