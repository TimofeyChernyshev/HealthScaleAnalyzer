package domain

type BMICategory string

var (
	SeriousUnder BMICategory = "Серьезный недобор"
	LightUnder   BMICategory = "Легкий недобор"
	Normal       BMICategory = "Норма"
	OverWeight   BMICategory = "Лишний вес"
	Obese        BMICategory = "Ожирение"
)
