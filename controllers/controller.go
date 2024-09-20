package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"sorteio-daily/models"

	"github.com/gin-gonic/gin"
)

var pessoas []models.Pessoa

func CriaNovaPessoa(c *gin.Context) {
	CarregarPessoas()
	var model models.Pessoa
	err := c.ShouldBindJSON(&model)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	file, err := os.Create("data/data.json")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	defer file.Close()
	encoder := json.NewEncoder(file)
	id := len(pessoas)
	model.ID = id + 1
	pessoas = append(pessoas, model)
	if err := encoder.Encode(pessoas); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"pessoa": model})
}

func CarregarPessoas() {
	file, err := os.Open("data/data.json")
	if err != nil {
		fmt.Println("error to open json")
		return
	}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&pessoas); err != nil {
		fmt.Println("error decoding JSON")
		return
	}
}

func ListPessoas(c *gin.Context) {
	CarregarPessoas()
	if pessoas != nil {
		c.JSON(200, pessoas)
		return
	}
	c.JSON(204, nil)
}

func SorteiaPessoa(c *gin.Context) {
	idSorteado := SorteiaId()
	file, err := os.Create("data/data.json")
	encoder := json.NewEncoder(file)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	if idSorteado == 0 {
		c.JSON(400, "don't have pessoas")
		return
	}
	for i, p := range pessoas {
		if p.ID == idSorteado {
			p.Sorteado = true
			pessoas[i] = p
			if err := encoder.Encode(pessoas); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, p)
			return
		}
	}
	c.JSON(400, gin.H{"not found pessoa with id: ": idSorteado})
}

func SorteiaId() int {
	CarregarPessoas()
	if pessoas != nil {
		idSorteado := rand.IntN(len(pessoas)) + 1
		return idSorteado
	}
	fmt.Println("not found pessoas")
	return 0
}
