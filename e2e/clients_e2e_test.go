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

	"github.com/joho/godotenv"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/config"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/server"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}
}

func TestE2ECreateClients(t *testing.T) {
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
    		"client_id": "d28b03b3-c25d-49de-afb8-3c6a30a4d122",
    		"login": "lotty",
    		"age": 3,
    		"location": "Moscow",
    		"gender": "MALE"
  		},
		{
 			"client_id": "8f696638-f877-433e-9562-2d2910f8ea9b",
    		"login": "shakhov_ad",
    		"age": 28,
    		"location": "Saint-Petersburg",
    		"gender": "MALE"
  		}
	]`

	resp, err := client.Post(ts.URL+"/clients/bulk", "application/json", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Ошибка чтения тела ответа: %v", err)
	}
	log.Printf("Ответ сервера: %s", string(body))

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Ожидался статус 201, а получил %d", resp.StatusCode)
	}

	var response []domain.User
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Ошибка при парсинге тела ответа: %v", err)
	}

	if response[0].ID.String() != "d28b03b3-c25d-49de-afb8-3c6a30a4d122" {
		t.Fatalf("Ожидалось поле 'id' (первого клиента) со значением 'd28b03b3-c25d-49de-afb8-3c6a30a4d122', а получили %v", response[0].ID)
	}

	t.Log("Тест создания клиентов прошел успешно!")
}

func TestE2EGetClients(t *testing.T) {
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

	// Create client to get
	reqBody := `[
		{
			"client_id": "d28b03b3-c25d-49de-afb8-3c6a30a4d122",
			"login": "lotty",
			"age": 3,
			"location": "Moscow",
			"gender": "MALE"
		}
	]`

	resp, err := client.Post(ts.URL+"/clients/bulk", "application/json", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса на создание: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Ошибка чтения тела ответа при создании: %v", err)
	}
	log.Printf("Ответ сервера на создание: %s", string(body))

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Ожидался статус 201, а получил %d", resp.StatusCode)
	}

	var createdClients []domain.User
	err = json.Unmarshal(body, &createdClients)
	if err != nil {
		t.Fatalf("Ошибка при парсинге тела ответа при создании клиента: %v", err)
	}

	// Get client
	respGet, err := client.Get(ts.URL + "/clients/" + createdClients[0].ID.String())
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса на получение клиента: %v", err)
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

	var response domain.User
	err = json.Unmarshal(bodyGet, &response)
	if err != nil {
		t.Fatalf("Ошибка при парсинге тела ответа при получении клиента: %v", err)
	}

	if response.ID.String() != createdClients[0].ID.String() {
		t.Fatalf("Ожидался ID клиента '%v', а получили %v", createdClients[0].ID, response.ID)
	}

	if response.Login != "lotty" {
		t.Fatalf("Ожидался логин 'lotty', а получили %v", response.Login)
	}

	t.Log("Тест получения клиента прошел успешно!")
}
