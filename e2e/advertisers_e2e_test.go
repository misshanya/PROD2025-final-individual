package e2e_test

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/config"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/server"
)

func TestE2ECreateGetAdvertisers(t *testing.T) {
	cfg := config.NewConfig()
	cfg.ServerAddress = "127.0.0.1:0"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server, err := server.NewServer(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	ts := httptest.NewServer(server.HttpServer.Handler)
	defer ts.Close()

	client := &http.Client{Timeout: 5 * time.Second}

	reqBody := `[
  		{
    		"advertiser_id": "ac240a28-36d7-448f-bbc5-62ed05bb433d",
    		"name": "Т-Банк"
  		}
	]`

	resp, err := client.Post(ts.URL+"/advertisers/bulk", "application/json", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса на создание рекламодателей: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Ошибка чтения тела ответа при создании: %v", err)
	}
	log.Printf("Ответ сервера на создание рекламодателей: %s", string(body))

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Ожидался статус 201, а получил %d", resp.StatusCode)
	}

	var createdAdvertisers []domain.Advertiser
	err = json.Unmarshal(body, &createdAdvertisers)
	if err != nil {
		t.Fatalf("Ошибка при парсинге тела ответа при создании рекламодателей: %v", err)
	}

	respGet, err := client.Get(ts.URL + "/advertisers/" + createdAdvertisers[0].ID.String())
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса на получение рекламодателя: %v", err)
	}
	defer respGet.Body.Close()

	bodyGet, err := io.ReadAll(respGet.Body)
	if err != nil {
		t.Fatalf("Ошибка чтения тела ответа при получении: %v", err)
	}
	log.Printf("Ответ сервера на получение: %s", string(bodyGet))

	if respGet.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200, а получил %d", respGet.StatusCode)
	}

	var response domain.Advertiser
	err = json.Unmarshal(bodyGet, &response)
	if err != nil {
		t.Fatalf("Ошибка при парсинге тела ответа при получении рекламодателя: %v", err)
	}

	if response.ID != createdAdvertisers[0].ID {
		t.Fatalf("Ожидался ID рекламодателя '%v', а получили %v", createdAdvertisers[0].ID, response.ID)
	}

	t.Log("Тест создания и получения рекламодателей прошел успешно!")
}

func TestE2ECreateUpdateMLScore(t *testing.T) {
	cfg := config.NewConfig()
	cfg.ServerAddress = "127.0.0.1:0"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server, err := server.NewServer(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	ts := httptest.NewServer(server.HttpServer.Handler)
	defer ts.Close()

	client := &http.Client{Timeout: 5 * time.Second}

	reqBody := `{
  		"advertiser_id": "ac240a28-36d7-448f-bbc5-62ed05bb433d",
  		"client_id": "8f696638-f877-433e-9562-2d2910f8ea9b",
  		"score": 3
	}`

	resp, err := client.Post(ts.URL+"/ml-scores", "application/json", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Ошибка при запросе на добавление ML скора")
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Ожидалось 200, получили: %v", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Не удалось прочитать тело ответа")
	}

	var response domain.MLScore
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		t.Fatalf("Ошибка при парсинге тела ответа: %v", err)
	}

	if response.Score != 3 {
		t.Fatalf("Несовпадение ML скора")
	}

	t.Log("Тест создания ML скора прошел успешно")
}
