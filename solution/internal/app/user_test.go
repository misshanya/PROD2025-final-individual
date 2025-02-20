package app

import (
	"testing"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
)

func TestValidateUser(t *testing.T) {
	validUser := domain.User{
		ID:     uuid.New(),
		Login:  "lotty",
		Age:    3,
		Gender: "MALE",
	}
	invalidGenderUser := domain.User{
		ID:     uuid.New(),
		Login:  "lotty",
		Age:    3,
		Gender: "smth",
	}
	invalidAgeUser := domain.User{
		ID:     uuid.New(),
		Login:  "lotty",
		Age:    -1,
		Gender: "MALE",
	}

	if err := validateUser(&validUser); err != nil {
		t.Fatalf("Валидный юзер не прошел валидацию: %v", err)
	}
	if err := validateUser(&invalidGenderUser); err == nil {
		t.Fatalf("Юзер с невалидным гендером прошел валидацию")
	}
	if err := validateUser(&invalidAgeUser); err == nil {
		t.Fatalf("Юзер с невалидным возрастом прошел валидацию")
	}

	t.Log("Тест валидации юзера пройден успешно!")
}
