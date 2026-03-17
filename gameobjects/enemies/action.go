package enemies

type Action struct {
	Name string
	Risk int //percent
}

func CreateSleepAction() Action {
	return Action{
		Name: "Sleep",
		Risk: 1,
	}
}
