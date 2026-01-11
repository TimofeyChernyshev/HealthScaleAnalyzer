package bmiconfig

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/TimofeyChernyshev/HealthScaleAnalyzer/internal/domain"
)

//go:embed *.json
var configFS embed.FS

type BMIConfigLoader struct{}

func NewBMIConfigLoader() BMIConfigLoader {
	return BMIConfigLoader{}
}

// LoadBMIConfig загружает конфиг из embed файлов
func (l BMIConfigLoader) LoadBMIConfig() (*domain.BMIConfig, *domain.BMIConfig, error) {
	maleFile := "bmi_male.json"
	femaleFile := "bmi_female.json"

	data, err := configFS.ReadFile(maleFile)
	if err != nil {
		return nil, nil, fmt.Errorf("не удалось загрузить конфиг %s: %w", maleFile, err)
	}

	var configMale domain.BMIConfig
	if err := json.Unmarshal(data, &configMale); err != nil {
		return nil, nil, fmt.Errorf("ошибка парсинга конфига %s: %w", maleFile, err)
	}

	data, err = configFS.ReadFile(femaleFile)
	if err != nil {
		return nil, nil, fmt.Errorf("не удалось загрузить конфиг %s: %w", femaleFile, err)
	}

	var configFemale domain.BMIConfig
	if err := json.Unmarshal(data, &configFemale); err != nil {
		return nil, nil, fmt.Errorf("ошибка парсинга конфига %s: %w", femaleFile, err)
	}

	return &configMale, &configFemale, nil
}
