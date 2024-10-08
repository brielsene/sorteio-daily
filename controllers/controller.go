package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"sorteio-daily/models"
	"strconv"

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
	CarregarPessoas()
	var listPessoaNSorteadas []models.Pessoa
	for _, ps := range pessoas {
		if ps.Sorteado == false {
			listPessoaNSorteadas = append(listPessoaNSorteadas, ps)

		}
	}
	if listPessoaNSorteadas == nil {
		for i, p := range pessoas {
			p.Sorteado = false
			pessoas[i] = p
		}
		for _, ps := range pessoas {
			if ps.Sorteado == false {
				listPessoaNSorteadas = append(listPessoaNSorteadas, ps)

			}
		}
	}
	idSorteado := rand.IntN(len(listPessoaNSorteadas))

	fmt.Println(listPessoaNSorteadas)
	fmt.Println(idSorteado)
	fmt.Println(listPessoaNSorteadas[idSorteado])

	pessoaSorteada := listPessoaNSorteadas[idSorteado]
	for i, p := range pessoas {
		if p.ID == pessoaSorteada.ID {

			pessoaSorteada.Sorteado = true
			pessoas[i] = pessoaSorteada
			file, err := os.Create("data/data.json")
			if err != nil {
				c.JSON(400, gin.H{"error": "error to open json" + err.Error()})
				return
			}
			defer file.Close()
			encoder := json.NewEncoder(file)
			if err := encoder.Encode(&pessoas); err != nil {
				c.JSON(400, gin.H{"error": "error to encode: " + err.Error()})
				return
			}
			c.JSON(200, &pessoaSorteada)
		}
	}

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

func DeletePessoaById(c *gin.Context) {
	CarregarPessoas()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "error to format to int: " + err.Error()})
		return
	}
	for i, pessoa := range pessoas {
		if pessoa.ID == id {
			file, err := os.Create("data/data.json")
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			defer file.Close()
			pessoas = append(pessoas[:i], pessoas[i+1:]...)
			encoder := json.NewEncoder(file)
			if err := encoder.Encode(pessoas); err != nil {
				c.JSON(400, gin.H{"error to encode file": err.Error()})
				return
			}
			c.JSON(200, gin.H{"Pessoa deletada com sucesso": pessoas})
			return
		}
	}

}

func UpdatedPessoaById(c *gin.Context) {
	idPessoa, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "error to conver idPessoa" + err.Error()})
		return
	}

	var pessoaModel models.Pessoa
	err = c.ShouldBindJSON(&pessoaModel)
	if err != nil {
		c.JSON(400, gin.H{"error": "should bind json pessoa: " + err.Error()})
		return
	}

	for i, p := range pessoas {
		if p.ID == idPessoa {
			p.ID = idPessoa
			pessoas[i] = p
			c.JSON(200, pessoas[i])
			return
		}
	}
	c.JSON(404, gin.H{"message": "pessoa not found"})

}
