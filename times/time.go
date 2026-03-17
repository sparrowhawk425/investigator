package times

type TimeOfDay int

const (
	Morning TimeOfDay = iota
	Afternoon
	Evening
	Night
)

func TransitionTimeOfDay(tod TimeOfDay) TimeOfDay {
	switch tod {
	case Morning:
		return Afternoon
	case Afternoon:
		return Evening
	case Evening:
		return Night
	case Night:
		return Morning
	}
	// Can't get here
	return tod
}

func GetTimeOfDayName(tod TimeOfDay) string {
	switch tod {
	case Morning:
		return "Morning"
	case Afternoon:
		return "Afternoon"
	case Evening:
		return "Evening"
	case Night:
		return "Night"
	}
	return "Witching Hour"
}
