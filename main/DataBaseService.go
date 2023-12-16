package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type DataBase struct {
	db *sql.DB
}

type DataBaseService interface {
	Insert(game GameScreenDTO) (bool, error)
	GetAll() ([]GameDatabaseDTO, error)
	AddInfoToGame(game GameDatabaseDTO) (bool, error)
	GetById(id string) (GameDatabaseDTO, error)
	ClearTable()
}

func NewDataBaseService(db *sql.DB, err error) DataBaseService {
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS game (game_id INTEGER PRIMARY KEY, game_tag varchar(50) NOT NULL , game_year varchar(10), game_desc varchar(10000))")
	if err != nil {
		panic(err)
	}
	log.Println("Подключение к БД выполнено!")
	return &DataBase{
		db: db,
	}
}

func (d *DataBase) Insert(game GameScreenDTO) (bool, error) {
	log.Println("Начало записи в БД игры: ", game.Game)
	_, err := d.db.Exec("INSERT INTO game VALUES ($1, $2)", game.Id, game.Game)
	log.Println("Успешный конец записи в БД игры: ", game.Game)
	return true, err
}

func (d *DataBase) AddInfoToGame(game GameDatabaseDTO) (bool, error) {
	log.Println("Начало обновления в БД игры: ", game.Game)
	_, err := d.db.Exec("UPDATE game SET game_year = $1, game_desc = $2 WHERE game_id = $3; ", game.Year, game.Description, game.Id)
	return true, err
}

func (d *DataBase) ClearTable() {
	_, _ = d.db.Exec("DELETE FROM game")
}

func (d *DataBase) GetAll() ([]GameDatabaseDTO, error) {
	log.Println("Начало поиска всех объектов типа Game")
	facts, err := d.db.Query("SELECT * FROM game")
	if err != nil {
		log.Println("В процессе поиска поймали исключение", err.Error())
		panic(err)
	}
	var foundGames []GameDatabaseDTO
	for facts.Next() {
		var h GameDatabaseDTO
		err = facts.Scan(&h.Id, &h.Game, &h.Description, &h.Year)
		if err != nil {
			log.Println("В процессе десериализации объекта поймали исключение", err.Error())
			panic(err)
		}
		foundGames = append(foundGames, GameDatabaseDTO{
			Id:          h.Id,
			Game:        h.Game,
			Description: h.Description,
			Year:        h.Year})
	}
	log.Println("Успешный конец поиска всех объектов типа Game")
	return foundGames, err
}

func (d *DataBase) GetById(id string) (GameDatabaseDTO, error) {
	log.Println("Начало поиска объекта типа Game")
	game, err := d.db.Query("SELECT * FROM game WHERE game_id = $1", id)
	if err != nil {
		log.Println("В процессе поиска поймали исключение", err.Error())
		panic(err)
	}
	var foundGame GameDatabaseDTO
	var h GameDatabaseDTO
	for game.Next() {
		err = game.Scan(&h.Id, &h.Game, &h.Description, &h.Year)
		if err != nil {
			err = game.Scan(&h.Id, &h.Game)
		}
		foundGame = GameDatabaseDTO{
			Id:          h.Id,
			Game:        h.Game,
			Description: h.Description,
			Year:        h.Year,
		}
	}
	log.Println("Успешный конец поиска объекта типа Game")
	return foundGame, err
}
