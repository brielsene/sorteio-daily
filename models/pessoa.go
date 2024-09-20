package models

type Pessoa struct {
	ID       int    `json:"id"`
	Nome     string `json:"nome"`
	Sorteado bool   `json:"sorteado"`
}
