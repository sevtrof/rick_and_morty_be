package model

type CharactersWithInfo struct {
	Info       Info        `json:"info"`
	Characters []Character `json:"results"`
}
