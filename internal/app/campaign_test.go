package app

import (
	"testing"

	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
)

func TestValidateTargeting(t *testing.T) {
	validGender := "MALE"
	validAgeFrom := int32(3)
	validAgeTo := int32(23)
	validTargeting := domain.Targeting{
		Gender:  &validGender,
		AgeFrom: &validAgeFrom,
		AgeTo:   &validAgeTo,
	}

	invalidGender := "smth"
	invalidGenderTargeting := domain.Targeting{
		Gender:  &invalidGender,
		AgeFrom: &validAgeFrom,
		AgeTo:   &validAgeTo,
	}

	invalidAgeFrom := int32(-5)
	invalidAgeFromTargeting := domain.Targeting{
		Gender:  &validGender,
		AgeFrom: &invalidAgeFrom,
		AgeTo:   &validAgeTo,
	}

	invalidAgeTo := int32(-5)
	invalidAgeToTargeting := domain.Targeting{
		Gender:  &validGender,
		AgeFrom: &validAgeFrom,
		AgeTo:   &invalidAgeTo,
	}

	secondAgeTo := int32(2)
	ageFromGreaterAgeToTargeting := domain.Targeting{
		Gender:  &validGender,
		AgeFrom: &validAgeFrom,
		AgeTo:   &secondAgeTo,
	}

	if !validateTargeting(validTargeting) {
		t.Fatal("Валидный таргетинг не прошел валидацию")
	}
	if validateTargeting(invalidGenderTargeting) {
		t.Fatalf("Таргетинг с невалидным гендером прошел валидацию")
	}
	if validateTargeting(invalidAgeFromTargeting) {
		t.Fatalf("Таргетинг с невалидным ageFrom прошел валидацию")
	}
	if validateTargeting(invalidAgeToTargeting) {
		t.Fatalf("Таргетинг с невалидным ageTo прошел валидацию")
	}
	if validateTargeting(ageFromGreaterAgeToTargeting) {
		t.Fatalf("Таргетинг, где ageFrom > ageTo, прошел валидацию")
	}

	t.Log("Тест валидации таргетинга успешно пройден!")
}
