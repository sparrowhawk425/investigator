package gamelogic

import (
	"github.com/sparrowhawk425/investigators/gameobjects"
)

type Name struct {
	First string
	Last  string
}

type Dossier struct {
	Name    string
	Target  gameobjects.Character
	Profile string
	Notes   []string
}

type Player struct {
	Name     Name
	Dossiers []Dossier
}
