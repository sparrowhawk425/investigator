package characters

import (
	"time"

	"github.com/sparrowhawk425/investigators/internal/config"
)

type Gender string

const (
	MaleGender      Gender = "Male"
	FemaleGender    Gender = "Female"
	NonbinaryGender Gender = "Nonbinary"
)

var Genders = []Gender{
	MaleGender, FemaleGender, NonbinaryGender,
}

func (g Gender) String() string {
	return string(g)
}

type EyeColor string

const (
	BlueEyes  EyeColor = "Blue"
	GreenEyes EyeColor = "Green"
	BrownEyes EyeColor = "Brown"
	AmberEyes EyeColor = "Amber"
	HazelEyes EyeColor = "Hazel"
	GrayEyes  EyeColor = "Gray"
	// Red/Violet (due to albinism)?
	// Heterochromia?
)

var EyeColors = []EyeColor{
	BlueEyes, GreenEyes, BrownEyes, AmberEyes, HazelEyes, GrayEyes,
}

func (ec EyeColor) String() string {
	return string(ec)
}

type HairColor string

const (
	BlondHair       HairColor = "Blond"
	DarkBlondHair   HairColor = "Dark Blond"
	MediumBrownHair HairColor = "Medium Brown"
	DarkBrownHair   HairColor = "Dark Brown"
	BlackHair       HairColor = "Black"
	AuburnHair      HairColor = "Auburn"
	RedHair         HairColor = "Red"
	GrayHair        HairColor = "Gray"
	WhiteHair       HairColor = "White"
)

var HairColors = []HairColor{
	BlondHair, DarkBlondHair, MediumBrownHair, DarkBrownHair, BlackHair, AuburnHair, RedHair, GrayHair, WhiteHair,
}

func (hc HairColor) String() string {
	return string(hc)
}

type HairLength string

const (
	BaldHair   HairLength = "Bald"
	ShortHair  HairLength = "Short"
	MediumHair HairLength = "Medium"
	LongHair   HairLength = "Long"
)

var HairLengths = []HairLength{
	BaldHair, ShortHair, MediumHair, LongHair,
}

func (hl HairLength) String() string {
	return string(hl)
}

type ShoeSize string

const (
	SmallShoe  ShoeSize = "Small"
	MediumShoe ShoeSize = "Medium"
	LargeShoe  ShoeSize = "Large"
)

var ShoeSizes = []ShoeSize{
	SmallShoe, MediumShoe, LargeShoe,
}

func (ss ShoeSize) String() string {
	return string(ss)
}

type Height string

const (
	ShortHeight   Height = "Short"
	AverageHeight Height = "Average"
	TallHeight    Height = "Tall"
)

var Heights = []Height{
	ShortHeight, AverageHeight, TallHeight,
}

func (h Height) String() string {
	return string(h)
}

type Weight string

const (
	ThinWeight    Weight = "Thin"
	AverageWeight Weight = "Average"
	OverWeight    Weight = "Heavy"
)

var Weights = []Weight{
	ThinWeight, AverageWeight, OverWeight,
}

func (w Weight) String() string {
	return string(w)
}

type DateOfBirth struct {
	Date time.Time `json:"date"`
	Age  int
}

type Characteristics struct {
	Dob         DateOfBirth
	Nationality config.CountryCode
	Gender      Gender
	EyeColor    EyeColor
	HairColor   HairColor
	HairLength  HairLength
	Height      Height
	Weight      Weight
	ShoeSize    ShoeSize
}
