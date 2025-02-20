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

func TestE2ECreateGetCampaign(t *testing.T) {
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
		"ad_title": "PROD",
  		"ad_text": "Для участия хватит школьных знаний!",
  		"impressions_limit": 1000000,
  		"clicks_limit": 10000,
  		"cost_per_impression": 0.05,
  		"cost_per_click": 0.2,
  		"start_date": 8,
  		"end_date": 15,
  		"targeting": {
    		"age_from": 13,
    		"age_to": 18,
    		"gender": "ALL"
  		}
	}`

	resp, err := client.Post(ts.URL+"/advertisers/ac240a28-36d7-448f-bbc5-62ed05bb433d/campaigns", "application/json", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса на создание кампании: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Ошибка чтения тела ответа при создании: %v", err)
	}
	log.Printf("Ответ сервера на создание кампании: %s", string(body))

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Ожидался статус 201, а получили %d", resp.StatusCode)
	}

	var createdCampaign domain.Campaign
	err = json.Unmarshal(body, &createdCampaign)
	if err != nil {
		t.Fatalf("Ошибка при парсинге тела ответа при создании кампании: %v", err)
	}

	respGet, err := client.Get(ts.URL + "/advertisers/ac240a28-36d7-448f-bbc5-62ed05bb433d/campaigns/" + createdCampaign.ID.String())
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса на получение кампании: %v", err)
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

	var response domain.Campaign
	err = json.Unmarshal(bodyGet, &response)
	if err != nil {
		t.Fatalf("Ошибка при парсинге тела ответа при получении кампании: %v", err)
	}

	if response.ID != createdCampaign.ID {
		t.Fatalf("Ожидался ID кампании '%v', а получили '%v'", createdCampaign.ID, response.ID)
	}

	t.Log("Тест создания и получения кампании прошел успешно!")
}
