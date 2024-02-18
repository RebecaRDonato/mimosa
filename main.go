package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	ws := gin.Default()
	ws.GET("/clientes/:id/extrato", getExtrato)
	ws.POST("/clientes/:id/transacoes", postTransacao)
	ws.Run() // listen and serve on 0.0.0.0:8080
}

func postTransacao(ctx *gin.Context) {
	type postTransacaoResposta struct {
		Valor     int    `json:"valor"`
		Tipo      string `json:"tipo"`
		Descricao string `json:"descricao"`
	}
	resposta := postTransacaoResposta{
		Valor:     1000,
		Tipo:      "c",
		Descricao: "descricao",
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
