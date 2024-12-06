package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		log.Println(err)
		return
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Erro ao fazer requisição:", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("Erro na resposta: %s\n", res.Status)
		return
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Println(err)
	}

	var bid *string
	err = json.Unmarshal(body, &bid)

	if err != nil {
		log.Println(err)
	}

	generateFile(*bid)
}

func getCotacao() {

}

func generateFile(value string) {

	f, err := os.Create("file.txt")

	if err != nil {
		log.Println(err)
	}

	_, err = f.WriteString(fmt.Sprintf("Dólar: %s", value))

	if err != nil {
		log.Println(err)
	}
}
