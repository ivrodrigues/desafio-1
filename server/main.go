package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Cotacao struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type Preco struct {
	ID    int `gorm:"primaryKey"`
	Preco string
	gorm.Model
}

func main() {
	db := initDatabase()

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		clientContext := r.Context()
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
		defer cancel()
		cotacao, err := getCotacao(ctx)
		if err != nil {
			http.Error(w, "Erro ao obter cotação", http.StatusInternalServerError)
			return
		}

		dbCtx, dbCancel := context.WithTimeout(context.Background(), time.Millisecond*10)
		defer dbCancel()
		db.WithContext(dbCtx).Create(&Preco{
			Preco: *cotacao,
		})

		json.NewEncoder(w).Encode(&cotacao)

		select {
		case <-dbCtx.Done():
			log.Println("Timeout: Conexão fechada com o banco de dados devido a timeout")
		case <-clientContext.Done():
			log.Println("Timeout: Conexão fechada com o cliente devido a timeout")
		case <-ctx.Done():
			log.Println("Timeout: Contexto cancelado pela API de cotação devido a timeout")
		}

	},
	)

	http.ListenAndServe(":8080", nil)

}

func initDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Preco{})

	return db
}

func getCotacao(ctx context.Context) (*string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)

	if err != nil {
		return nil, err
	}

	return &cotacao.Usdbrl.Bid, nil
}
