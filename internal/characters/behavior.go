package characters

type Behavior struct {
	Name string
	Desc string
}

func CreateFrugal() Behavior {
	return Behavior{
		Name: "Frugal",
		Desc: "Prone to conserving money and prefer cheap locations",
	}
}

func CreateProfligate() Behavior {
	return Behavior{
		Name: "Profligate",
		Desc: "Prone to spend money and prefer expensive locations",
	}
}

func CreateGambler() Behavior {
	return Behavior{
		Name: "Gambler",
		Desc: "Tends to spend free time in Casinos. More willing to take risks",
	}
}

var RegularBehaviors = []Behavior{
	CreateFrugal(), CreateProfligate(), CreateGambler(),
}

func CreateSquatter() Behavior {
	return Behavior{
		Name: "Squatter",
		Desc: "Lives in unoccupied buildings",
	}
}
