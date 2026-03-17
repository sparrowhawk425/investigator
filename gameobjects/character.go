package gameobjects

import (
	"math/rand/v2"
	"strings"
	"time"

	"github.com/sparrowhawk425/investigators/internal/nameapi"
)

type EyeColor string

const (
	Blue   EyeColor = "Blue"
	Green  EyeColor = "Green"
	Brown  EyeColor = "Brown"
	Amber  EyeColor = "Amber"
	Hazel  EyeColor = "Hazel"
	Gray_E EyeColor = "Gray"
	// Red/Violet (due to albinism)?
	// Heterochromia?
)

var EyeColors = []EyeColor{
	Blue, Green, Brown, Amber, Hazel, Gray_E,
}

type HairColor string

const (
	Blond       HairColor = "Blond"
	DarkBlond   HairColor = "Dark Blond"
	MediumBrown HairColor = "Medium Brown"
	DarkBrown   HairColor = "Dark Brown"
	Black       HairColor = "Black"
	Auburn      HairColor = "Auburn"
	Red         HairColor = "Red"
	Gray_H      HairColor = "Gray"
	White       HairColor = "White"
)

var HairColors = []HairColor{
	Blond, DarkBlond, MediumBrown, DarkBrown, Black, Auburn, Red, Gray_H, White,
}

type DateOfBirth struct {
	Date time.Time `json:"date"`
	Age  int
}

type Characteristics struct {
	Dob         DateOfBirth
	Nationality string
	Gender      string
	EyeColor    EyeColor
	HairColor   HairColor
	Height      int //cm
	Weight      int //lb
}

type Name struct {
	First string
	Last  string
}

type Character struct {
	Name   Name
	Traits Characteristics
}

func CreateRandomCharacter(apiChar nameapi.Character) Character {

	eyeColor := EyeColors[rand.IntN(len(EyeColors))]
	hairColor := HairColors[rand.IntN(len(HairColors))]
	return Character{
		Name: Name{
			First: apiChar.Name.First,
			Last:  apiChar.Name.Last,
		},
		Traits: Characteristics{
			Dob: DateOfBirth{
				Date: apiChar.DateOfBirth.Date,
				Age:  apiChar.DateOfBirth.Age,
			},
			Nationality: apiChar.Nationality,
			EyeColor:    eyeColor,
			HairColor:   hairColor,
			Gender:      apiChar.Gender,
			Height:      getHeight(apiChar.Gender),
			Weight:      getWeight(apiChar.Gender),
		},
	}
}

// Avg Male weight: 200lb
// Avg Female weight: 170lb
func getWeight(gender string) int {
	avg := 170
	if strings.ToLower(gender) == "male" {
		avg = 200
	}
	return calcDiff(avg)
}

// Avg Male height: 171cm (5'7")
// Avg Female height: 159cm (5'3")
func getHeight(gender string) int {
	avg := 159
	if strings.ToLower(gender) == "male" {
		avg = 171
	}
	return calcDiff(avg)
}

func calcDiff(avg int) int {
	pctDiff := rand.IntN(11) // 0-10%
	diff := avg * pctDiff / 100
	if rand.IntN(2)%2 == 0 { // + or -
		return avg + diff
	}
	return avg - diff
}
