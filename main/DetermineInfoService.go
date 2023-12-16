package main

import (
	"os"
	"strings"
)

type DetermineInfoService interface {
	DetermineInfo(game GameDatabaseDTO) GameDatabaseDTO
}

func NewDetermineInfoService() DetermineInfoService {
	return GameDatabaseDTO{}
}

func (h GameDatabaseDTO) DetermineInfo(game GameDatabaseDTO) GameDatabaseDTO {
	games, _ := os.ReadFile("main/Games.txt")
	newGame := strings.Split(string(games), "\n")[game.Id]
	preparedGame := strings.Replace(newGame, "\r", "", -1)
	year := strings.Split(preparedGame, "*")[1][0:4]
	description := strings.Split(preparedGame, "*")[1][7:]
	return GameDatabaseDTO{
		Id:          game.Id,
		Game:        game.Game,
		Year:        year,
		Description: description,
	}
}
