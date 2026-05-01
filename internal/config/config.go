package config

import (
	"bufio"
	"log"
	"log/slog"
	"os"
	"slices"
)

type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)

// Available countries derive from nameapi package, but don't want to pass that package around
type CountryName string

const (
	AustraliaName    CountryName = "Australia"
	BrazilName       CountryName = "Brazil"
	CanadaName       CountryName = "Canada"
	SwitzerlandName  CountryName = "Switzerland"
	GermanyName      CountryName = "Germany"
	DenmarkName      CountryName = "Denmark"
	SpainName        CountryName = "Spain"
	FinlandName      CountryName = "Finland"
	FranceName       CountryName = "France"
	GreatBritainName CountryName = "Great Britain"
	IrelandName      CountryName = "Ireland"
	IndiaName        CountryName = "India"
	MexicoName       CountryName = "Mexico"
	NetherlandsName  CountryName = "Netherlands"
	NorwayName       CountryName = "Norway"
	NewZealandName   CountryName = "New Zealand"
	SerbiaName       CountryName = "Serbia"
	TürkiyeName      CountryName = "Türkiye"
	UkraineName      CountryName = "Ukraine"
	UnitedStatesName CountryName = "United States"
)

var CountryNames = []CountryName{
	AustraliaName, BrazilName, CanadaName, SwitzerlandName, GermanyName, DenmarkName, SpainName, FinlandName, FranceName, GreatBritainName,
	IrelandName, IndiaName, MexicoName, NetherlandsName, NorwayName, NewZealandName, SerbiaName, UkraineName, UnitedStatesName,
}

func (cn CountryName) String() string {
	return string(cn)
}

type CountryCode string

const (
	AustraliaCode    CountryCode = "AU"
	BrazilCode       CountryCode = "BR"
	CanadaCode       CountryCode = "CA"
	SwitzerlandCode  CountryCode = "CH"
	GermanyCode      CountryCode = "DE"
	DenmarkCode      CountryCode = "DK"
	SpainCode        CountryCode = "ES"
	FinlandCode      CountryCode = "FI"
	FranceCode       CountryCode = "FR"
	GreatBritainCode CountryCode = "GB"
	IrelandCode      CountryCode = "IE"
	IndiaCode        CountryCode = "IN"
	MexicoCode       CountryCode = "MX"
	NetherlandsCode  CountryCode = "NL"
	NorwayCode       CountryCode = "NO"
	NewZealandCode   CountryCode = "NZ"
	SerbiaCode       CountryCode = "RS"
	TürkiyeCode      CountryCode = "TR"
	UkraineCode      CountryCode = "UA"
	UnitedStatesCode CountryCode = "US"
)

var CountryCodes = []CountryCode{
	AustraliaCode, BrazilCode, CanadaCode, SwitzerlandCode, GermanyCode, DenmarkCode, SpainCode, FinlandCode, GreatBritainCode,
	IrelandCode, IndiaCode, MexicoCode, NetherlandsCode, NorwayCode, NewZealandCode, SerbiaCode, TürkiyeCode, UkraineCode, UnitedStatesCode,
}

func (cc CountryCode) String() string {
	return string(cc)
}

type Country struct {
	Name CountryName
	Code CountryCode
}

var Countries = []Country{
	{Name: AustraliaName, Code: AustraliaCode}, {Name: BrazilName, Code: BrazilCode}, {Name: CanadaName, Code: CanadaCode}, {Name: SwitzerlandName, Code: SwitzerlandCode},
	{Name: GermanyName, Code: GermanyCode}, {Name: DenmarkName, Code: DenmarkCode}, {Name: SpainName, Code: SpainCode}, {Name: FinlandName, Code: FinlandCode},
	{Name: FranceName, Code: FranceCode}, {Name: GreatBritainName, Code: GreatBritainCode}, {Name: IrelandName, Code: IrelandCode}, {Name: IndiaName, Code: IndiaCode},
	{Name: MexicoName, Code: MexicoCode}, {Name: NetherlandsName, Code: NetherlandsCode}, {Name: NorwayName, Code: NorwayCode}, {Name: NewZealandName, Code: NewZealandCode},
	{Name: SerbiaName, Code: SerbiaCode}, {Name: TürkiyeName, Code: TürkiyeCode}, {Name: UkraineName, Code: UkraineCode}, {Name: UnitedStatesName, Code: UnitedStatesCode},
}

func (c Country) GetName() CountryName {
	return c.Name
}

func (c Country) GetCode() CountryCode {
	return c.Code
}

func GetCountryByName(name CountryName) Country {
	idx := slices.IndexFunc(Countries, func(c Country) bool { return c.Name == name })
	return Countries[idx]
}

func GetCountryByCode(code CountryCode) Country {
	idx := slices.IndexFunc(Countries, func(c Country) bool { return c.Code == code })
	return Countries[idx]
}

type PostCard struct {
	CountryName
	Image []string
}

type PostCards map[CountryCode]PostCard

func (pc PostCards) Get(code CountryCode) PostCard {
	card, ok := pc[code]
	if !ok {
		slog.Warn("Post Card not found", slog.String("countryName", code.String()))
	}
	return card
}

type Config struct {
	Difficulty
	PostCards

	Level      int
	PlayerName string
	Countries  []Country
}

func (c Config) GetNumCriminals() int {
	switch c.Difficulty {
	case Medium:
		return 6
	case Hard:
		return 9
	default:
		return 3
	}
}

func (c Config) GetNumPlaces() int {
	switch c.Difficulty {
	case Medium:
		return 80
	case Hard:
		return 120
	default:
		return 40
	}
}

func (c Config) GetNumBosses() int {
	switch c.Difficulty {
	case Medium:
		return 2
	case Hard:
		return 3
	default:
		return 1
	}
}

func (c Config) GetFinalLevel() int {
	switch c.Difficulty {
	case Medium:
		return 4
	case Hard:
		return 8
	default:
		return 2
	}
}

func LoadPostCards() PostCards {
	cards := make(PostCards)

	cards[AustraliaCode] = PostCard{CountryName: AustraliaName, Image: loadPostCard("SydneyOperaHouse.txt")}
	cards[BrazilCode] = PostCard{CountryName: BrazilName, Image: loadPostCard("CristoRedentor.txt")}
	cards[CanadaCode] = PostCard{CountryName: CanadaName, Image: loadPostCard("CNTower.txt")}
	cards[SwitzerlandCode] = PostCard{CountryName: SwitzerlandName, Image: loadPostCard("Matterhorn.txt")}
	cards[GermanyCode] = PostCard{CountryName: GermanyName, Image: loadPostCard("BrandenburgGate.txt")}
	cards[DenmarkCode] = PostCard{CountryName: DenmarkName, Image: loadPostCard("LittleMermaid.txt")}
	cards[SpainCode] = PostCard{CountryName: SpainName, Image: loadPostCard("LaSagradaFamilia.txt")}
	cards[FinlandCode] = PostCard{CountryName: FinlandName, Image: loadPostCard("HelsinkiCathedral.txt")}
	cards[FranceCode] = PostCard{CountryName: FranceName, Image: loadPostCard("EiffelTower.txt")}
	cards[GreatBritainCode] = PostCard{CountryName: GreatBritainName, Image: loadPostCard("BigBen.txt")}
	cards[IrelandCode] = PostCard{CountryName: IrelandName, Image: loadPostCard("BlarneyCastle.txt")}
	cards[IndiaCode] = PostCard{CountryName: IndiaName, Image: loadPostCard("TajMahal.txt")}
	cards[MexicoCode] = PostCard{CountryName: MexicoName, Image: loadPostCard("ChichenItza.txt")}
	cards[NetherlandsCode] = PostCard{CountryName: NetherlandsName, Image: loadPostCard("Kinderdijk.txt")}
	cards[NorwayCode] = PostCard{CountryName: NorwayName, Image: loadPostCard("BorgundStaveChurch.txt")}
	cards[NewZealandCode] = PostCard{CountryName: NewZealandName, Image: loadPostCard("Beehive.txt")}
	cards[SerbiaCode] = PostCard{CountryName: SerbiaName, Image: loadPostCard("SaintSavia.txt")}
	cards[TürkiyeCode] = PostCard{CountryName: TürkiyeName, Image: loadPostCard("HagiaSophia.txt")}
	cards[UkraineCode] = PostCard{CountryName: UkraineName, Image: loadPostCard("MotherUkraine.txt")}
	cards[UnitedStatesCode] = PostCard{CountryName: UnitedStatesName, Image: loadPostCard("StatueOfLiberty.txt")}
	return cards
}

func loadPostCard(fileName string) []string {
	filePath := "assets/ascii/" + fileName
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error reading postcard file %s: %v", filePath, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	image := []string{}
	for scanner.Scan() {
		image = append(image, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		log.Fatalf("error reading file contents: %v", err)
	}
	return image
}
