package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Primeiro = limite - Segundo = saldo
var clientes = map[string][]int{
	"1": {100000, 0},
	"2": {80000, 0},
	"3": {1000000, 0},
	"4": {10000000, 0},
	"5": {500000, 0},
	"7": {10, 0},
}

func main() {
	ws := gin.Default()
	ws.GET("/clientes/:id/extrato", getExtrato)
	ws.POST("/clientes/:id/transacoes", postTransacao)
	ws.Run() // listen and serve on 0.0.0.0:8080
}

func postTransacao(ctx *gin.Context) {
	type postTransacaoResposta struct {
		Limite int `json:"limite"`
		Saldo  int `json:"saldo"`
	}

	type postTransacaoRequestBody struct {
		Valor     int    `json:"valor"`
		Tipo      string `json:"tipo"`
		Descricao string `json:"descricao"`
	}

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

	ctx.JSON(http.StatusOK, resposta)

}

func getExtrato(ctx *gin.Context) {
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
	resposta := getExtratoResposta{
		Saldo: saldo{
			Total:       -9098,
			DataExtrato: "2024-01-17T02:34:41.217753Z",
			Limite:      100000,
		},
		UltimasTransacoes: []ultimaTransacao{
			{
				Valor:       10,
				Tipo:        "c",
				Descricao:   "descricao",
				RealizadaEm: "2024-01-17T02:34:38.543030Z",
			},
			{
				Valor:       90000,
				Tipo:        "d",
				Descricao:   "descricao",
				RealizadaEm: "2024-01-17T02:34:38.543030Z",
			},
		},
	}
	ctx.JSON(http.StatusOK, resposta)
}
