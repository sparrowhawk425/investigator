package times

type DayOfTheWeek int

const (
	Monday DayOfTheWeek = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func (day DayOfTheWeek) String() string {
	switch day {
	case Monday:
		return "Monday"
	case Tuesday:
		return "Tuesday"
	case Wednesday:
		return "Wednesday"
	case Thursday:
		return "Thursday"
	case Friday:
		return "Friday"
	case Saturday:
		return "Saturday"
	case Sunday:
		return "Sunday"
	}
	return "Someday"
}

func (day DayOfTheWeek) NextDay() DayOfTheWeek {
	return (day + 1) % 7
}

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

func (tod TimeOfDay) GetName() string {
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
