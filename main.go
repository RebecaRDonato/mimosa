package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type postTransacaoResposta struct {
	Limite int `json:"limite"`
	Saldo  int `json:"saldo"`
}

type postTransacaoRequestBody struct {
	Valor     int    `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

type saldo struct {
	Total       int    `json:"total"`
	DataExtrato string `json:"data_extrato"`
	Limite      int    `json:"limite"`
}

type ultimaTransacao struct {
	Valor       int    `json:"valor"`
	Tipo        string `json:"tipo"`
	Descricao   string `json:"descricao"`
	RealizadaEm string `json:"realizada_em"`
}

type getExtratoResposta struct {
	Saldo             saldo             `json:"saldo"`
	UltimasTransacoes []ultimaTransacao `json:"ultimas_transacoes"`
}

// Primeiro = limite - Segundo = saldo
var clientes = map[string][]int{
	"1": {100000, 0},
	"2": {80000, 0},
	"3": {1000000, 0},
	"4": {10000000, 0},
	"5": {500000, 0},
	"7": {10, 0},
}

var transacoes = map[string][]ultimaTransacao{
	"1": {},
	"2": {},
	"3": {},
	"4": {},
	"5": {},
	"7": {},
}

func main() {
	ws := gin.Default()
	ws.GET("/clientes/:id/extrato", getExtrato)
	ws.POST("/clientes/:id/transacoes", postTransacao)
	ws.Run() // listen and serve on 0.0.0.0:8080
}

func postTransacao(ctx *gin.Context) {
	id := ctx.Param("id")
	cliente, ok := clientes[id]

	if ok == false {
		ctx.Status(404)
		return
	}

	requestBody := postTransacaoRequestBody{}

	err := ctx.BindJSON(&requestBody)
	if err != nil {
		ctx.Status(422)
		return
	}

	if requestBody.Tipo != "c" && requestBody.Tipo != "d" {
		ctx.Status(422)
		return
	}

	if len(requestBody.Descricao) < 1 || len(requestBody.Descricao) > 10 {
		ctx.Status(422)
		return
	}

	if requestBody.Tipo == "c" {
		cliente[1] = cliente[1] + requestBody.Valor
	}

	if requestBody.Tipo == "d" {
		if cliente[0]+cliente[1] == 0 {
			ctx.Status(422)
			return
		}
		cliente[1] = cliente[1] - requestBody.Valor
	}

	resposta := postTransacaoResposta{
		Limite: cliente[0],
		Saldo:  cliente[1],
	}

	transacoes[id] = append(transacoes[id], ultimaTransacao{
		Valor:       requestBody.Valor,
		Tipo:        requestBody.Tipo,
		Descricao:   requestBody.Descricao,
		RealizadaEm: time.Now().UTC().Format(time.RFC3339Nano),
	})

	ctx.JSON(http.StatusOK, resposta)
}

func getExtrato(ctx *gin.Context) {
	id := ctx.Param("id")
	cliente, ok := clientes[id]

	if ok == false {
		ctx.Status(404)
		return
	}

	resposta := getExtratoResposta{
		Saldo: saldo{
			Total:       cliente[1],
			DataExtrato: time.Now().UTC().Format(time.RFC3339Nano),
			Limite:      cliente[0],
		},
		UltimasTransacoes: transacoes[id],
	}

	ctx.JSON(http.StatusOK, resposta)
}
