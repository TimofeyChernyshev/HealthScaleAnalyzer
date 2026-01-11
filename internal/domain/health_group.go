package domain

type HealthGroup int

var (
	First  HealthGroup = 1
	Second HealthGroup = 2
	Third  HealthGroup = 3
	Fourth HealthGroup = 4
	Fifth  HealthGroup = 5
)

var AvailableHealthGroups = []HealthGroup{First, Second, Third, Fourth, Fifth}
