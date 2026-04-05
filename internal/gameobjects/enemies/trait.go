package enemies

type PersonalityTrait struct {
	Name string
	Desc string
}

func createProfligate() PersonalityTrait {
	return PersonalityTrait{
		Name: "Profligate",
		Desc: "Prone to spend money and stay in expensive locations",
	}
}

func createSquatter() PersonalityTrait {
	return PersonalityTrait{
		Name: "Squatter",
		Desc: "Lives in unoccupied buildings",
	}
}
