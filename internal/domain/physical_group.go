package domain

type PhysicalGroup string

var (
	Default PhysicalGroup = "основная"
	Prepare PhysicalGroup = "подготовительная"
	Special PhysicalGroup = "специальная"
)

var AvailablePhysicalGroups = []PhysicalGroup{Default, Prepare, Special}
