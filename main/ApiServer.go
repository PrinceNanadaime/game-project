package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type ApiServer struct {
	game          GameInfoService
	determination DetermineInfoService
	database      DataBaseService
}

func NewApiServer(human GameInfoService, determination DetermineInfoService, database DataBaseService) ApiServer {
	return ApiServer{
		game:          human,
		determination: determination,
		database:      database,
	}
}

func (s *ApiServer) Start() error {
	http.HandleFunc("/info", s.HandleGetStartInfo)
	http.HandleFunc("/", s.HandleGetById)
	http.HandleFunc("/all", s.HandleGetAll)
	http.HandleFunc("/db/clear", s.HandleClearDatabase)

	http.HandleFunc("/determine-info", s.HandleDetermination)
	http.HandleFunc("/db/insert-base", s.HandleInsertionToDatabase)
	http.HandleFunc("/db/get/", s.HandleGetByIdFromDatabase)
	http.HandleFunc("/db/insert-info", s.HandleUpdateToDatabase)
	http.HandleFunc("/db/get-all", s.HandleGetAllFromDatabase)
	return http.ListenAndServe("localhost:8080", nil)
}

func (s *ApiServer) HandleGetStartInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на вывод общей информации об играх 2000-х!")
	games := s.game.GetNewStartGameInfo()
	WriteJson(w, http.StatusOK, games)
	log.Println("Запрос на вывод общей информации об играх 2000-х выполнен!")
}

func (s *ApiServer) HandleGetById(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	game := s.game.GetById(id)
	WriteJson(w, http.StatusOK, game)
}

func (s *ApiServer) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	games := s.game.GetAll()
	WriteJson(w, http.StatusOK, games)
}

func (s *ApiServer) HandleGetByIdFromDatabase(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/db/get/")
	game, _ := s.database.GetById(id)
	WriteJson(w, http.StatusOK, game)
}

func (s *ApiServer) HandleDetermination(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на определение доп. данных об игре!")
	body := GameDatabaseDTO{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("В процессе определения доп. данных об игре поймали исключение:", err.Error())
		WriteJson(w, http.StatusInternalServerError, err.Error())
	}
	fact := s.determination.DetermineInfo(body)
	WriteJson(w, http.StatusOK, fact)
	log.Println("Запрос на определение доп. данных об игре успешно выполнен!")
}

func (s *ApiServer) HandleInsertionToDatabase(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на запись новой игры в БД!")
	body := GameScreenDTO{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("В процессе записи данных в БД поймали исключение:", err.Error())
		WriteJson(w, http.StatusInternalServerError, err.Error())
	}
	fact, _ := s.database.Insert(body)
	WriteJson(w, http.StatusOK, fact)
	log.Println("Запрос на запись новой игры в БД успешно выполнен!")
}

func (s *ApiServer) HandleUpdateToDatabase(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на обновление новой игры в БД!")
	body := GameDatabaseDTO{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("В процессе записи данных в БД поймали исключение:", err.Error())
		WriteJson(w, http.StatusInternalServerError, err.Error())
	}
	game, err := s.database.AddInfoToGame(body)
	if err != nil {
		log.Println("В процессе записи новой игры поймали исключение:", err.Error())
		WriteJson(w, http.StatusInternalServerError, err.Error())
	}
	WriteJson(w, http.StatusOK, game)
	log.Println("Запрос на запись новой игры в БД успешно выполнен!")
}

func (s *ApiServer) HandleGetAllFromDatabase(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на вывод общей информации об играх 2000-х из БД!")
	games, err := s.database.GetAll()
	if err != nil {
		log.Println("В процессе вывода общей информации об играх 2000-х из БД поймали исключение:", err.Error())
		WriteJson(w, http.StatusInternalServerError, err.Error())
	}
	WriteJson(w, http.StatusOK, games)
	log.Println("Запрос на вывод общей информации об играх 2000-х из БД успешно выполнен!")
}

func (s *ApiServer) HandleClearDatabase(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на очистку БД!")
	s.database.ClearTable()
	WriteJson(w, http.StatusOK, "База данных очищена")
	log.Println("Запрос на очистку БД выполнен!")
}

func WriteJson(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	return encoder.Encode(v)
}
