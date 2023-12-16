package main

type Game struct {
	Game        string `json:"game"`
	Description string `json:"description"`
	Year        string `json:"year"`
}

type GameScreenDTO struct {
	Id   int
	Game string `json:"game"`
}

type GameDatabaseDTO struct {
	Id          int
	Game        string `json:"game"`
	Description string `json:"description"`
	Year        string `json:"year"`
}
