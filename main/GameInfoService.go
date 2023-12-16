package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type GameInfoService interface {
	GetNewStartGameInfo() []GameScreenDTO
	GetById(id string) *Game
	GetAll() []*Game
}

func NewGameInfoService() GameInfoService {
	return &Game{}
}

func (g Game) GetNewStartGameInfo() []GameScreenDTO {
	gamesString, _ := os.ReadFile("main/Games.txt")
	games := strings.Split(string(gamesString), "\n")
	var startInfos []GameScreenDTO
	for index, game := range games {
		preparedGame := strings.Replace(game, "\r", "", -1)
		gameScreenDTO := GameScreenDTO{
			Id:   index,
			Game: strings.Split(preparedGame, "*")[0],
		}
		data, _ := json.Marshal(gameScreenDTO)
		InsertGameToDatabase(data)
		startInfos = append(startInfos, gameScreenDTO)
	}
	return startInfos
}

func (g Game) GetAll() []*Game {
	var games []*Game
	for i := 0; i < 11; i++ {
		games = append(games, GetGameById(strconv.Itoa(i)))
	}
	return games
}
func (g *Game) GetById(id string) *Game {
	return GetGameById(id)
}

func GetGameById(id string) *Game {
	game, _ := GetByIdFromDatabase(id)
	data, _ := json.Marshal(game)
	updatedGame, _ := DetermineGameInfo(data)
	data, _ = json.Marshal(updatedGame)
	UpdateGameToDatabase(data)
	return &Game{
		Game:        updatedGame.Game,
		Description: updatedGame.Description,
		Year:        updatedGame.Year,
	}
}

func InsertGameToDatabase(game []byte) {
	dbRs, err := http.Post("http://localhost:8080/db/insert-base", "application/json", bytes.NewBuffer(game))
	if err != nil {
		panic(err)
	}
	defer dbRs.Body.Close()
}

func GetByIdFromDatabase(id string) (*GameDatabaseDTO, error) {
	url := "http://localhost:8080/db/get/" + id
	dbRs, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer dbRs.Body.Close()
	updatedGame := &GameDatabaseDTO{}
	if err := json.NewDecoder(dbRs.Body).Decode(updatedGame); err != nil {
		return nil, err
	}
	return updatedGame, err
}

func GetAllFromDatabase() ([]GameDatabaseDTO, error) {
	dbRs, err := http.Get("http://localhost:8080/db/get-all")
	if err != nil {
		panic(err)
	}
	defer dbRs.Body.Close()
	var updatedGames []GameDatabaseDTO
	if err := json.NewDecoder(dbRs.Body).Decode(&updatedGames); err != nil {
		return nil, err
	}
	return updatedGames, err
}

func DetermineGameInfo(game []byte) (*GameDatabaseDTO, error) {
	dbRs, err := http.Post("http://localhost:8080/determine-info", "application/json", bytes.NewBuffer(game))
	if err != nil {
		panic(err)
	}
	defer dbRs.Body.Close()
	updatedGame := &GameDatabaseDTO{}
	if err := json.NewDecoder(dbRs.Body).Decode(updatedGame); err != nil {
		return nil, err
	}
	return updatedGame, err
}

func UpdateGameToDatabase(game []byte) {
	dbRs, err := http.Post("http://localhost:8080/db/insert-info", "application/json", bytes.NewBuffer(game))
	if err != nil {
		panic(err)
	}
	defer dbRs.Body.Close()
}
