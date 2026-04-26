package characters

import "time"

type Gender string

const (
	MaleGender      Gender = "Male"
	FemaleGender    Gender = "Female"
	NonbinaryGender Gender = "Nonbinary"
	UnknownGender   Gender = "Unknown"
)

var Genders = []Gender{
	MaleGender, FemaleGender,
}

func (g Gender) String() string {
	return string(g)
}

type Nationality string

const (
	Australia          Nationality = "AU"
	Brazil             Nationality = "BR"
	Canada             Nationality = "CA"
	Switzerland        Nationality = "CH"
	Germany            Nationality = "DE"
	Denmark            Nationality = "DK"
	Spain              Nationality = "ES"
	Finland            Nationality = "FI"
	France             Nationality = "FR"
	GreatBritain       Nationality = "GB"
	Ireland            Nationality = "IE"
	India              Nationality = "IN"
	Mexico             Nationality = "MX"
	Netherlands        Nationality = "NL"
	Norway             Nationality = "NO"
	NewZealand         Nationality = "NZ"
	Serbia             Nationality = "RS"
	Türkiye            Nationality = "TR"
	Ukraine            Nationality = "UA"
	UnitedStates       Nationality = "US"
	UnknownNationality Nationality = "XX"
)

var Nationalities = []Nationality{
	Australia, Brazil, Canada, Switzerland, Germany, Denmark, Spain, Finland, GreatBritain,
	Ireland, India, Mexico, Netherlands, Norway, NewZealand, Serbia, Türkiye, Ukraine, UnitedStates,
}

func (n Nationality) String() string {
	return string(n)
}

type EyeColor string

const (
	BlueEyes    EyeColor = "Blue"
	GreenEyes   EyeColor = "Green"
	BrownEyes   EyeColor = "Brown"
	AmberEyes   EyeColor = "Amber"
	HazelEyes   EyeColor = "Hazel"
	GrayEyes    EyeColor = "Gray"
	UnknownEyes EyeColor = "Unknown"
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
	BlondHair        HairColor = "Blond"
	DarkBlondHair    HairColor = "Dark Blond"
	MediumBrownHair  HairColor = "Medium Brown"
	DarkBrownHair    HairColor = "Dark Brown"
	BlackHair        HairColor = "Black"
	AuburnHair       HairColor = "Auburn"
	RedHair          HairColor = "Red"
	GrayHair         HairColor = "Gray"
	WhiteHair        HairColor = "White"
	UnknownHairColor HairColor = "Unknown"
)

var HairColors = []HairColor{
	BlondHair, DarkBlondHair, MediumBrownHair, DarkBrownHair, BlackHair, AuburnHair, RedHair, GrayHair, WhiteHair,
}

func (hc HairColor) String() string {
	return string(hc)
}

type HairLength string

const (
	BaldHair          HairLength = "Bald"
	ShortHair         HairLength = "Short"
	MediumHair        HairLength = "Medium"
	LongHair          HairLength = "Long"
	UnknownHairLength HairLength = "Unknown"
)

var HairLengths = []HairLength{
	BaldHair, ShortHair, MediumHair, LongHair,
}

func (hl HairLength) String() string {
	return string(hl)
}

type ShoeSize string

const (
	SmallShoe   ShoeSize = "Small"
	MediumShoe  ShoeSize = "Medium"
	LargeShoe   ShoeSize = "Large"
	UnknownShoe ShoeSize = "Unknown"
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
	UnknownHeight Height = "Unknown"
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
	UnknownWeight Weight = "Unknown"
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
	Nationality Nationality
	Gender      Gender
	EyeColor    EyeColor
	HairColor   HairColor
	HairLength  HairLength
	Height      Height
	Weight      Weight
	ShoeSize    ShoeSize
}
