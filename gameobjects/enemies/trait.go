package enemies

type PersonalityTrait struct {
	Name string
}

func createProfligate() PersonalityTrait {
	return PersonalityTrait{
		Name: "Profligate",
	}
}
