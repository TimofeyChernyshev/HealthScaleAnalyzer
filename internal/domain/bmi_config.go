package domain

// BMIRangeConfig конфигурация одного диапазона
type BMIRangeConfig struct {
	MinAge       int     `json:"min_age"`
	MaxAge       int     `json:"max_age"`
	SeriousUnder float64 `json:"serious_under"`
	LightUnder   float64 `json:"light_under"`
	Normal       float64 `json:"normal"`
	OverWeight   float64 `json:"overweight"`
	Obese        float64 `json:"obese"`
}

// BMIConfig полный конфиг для одного пола
type BMIConfig struct {
	Ranges []BMIRangeConfig `json:"ranges"`
}
