package gameloop

import (
	"fmt"
	"log"
	"math/rand/v2"
	"slices"

	"github.com/fatih/color"
	"github.com/samber/lo"
	"github.com/sparrowhawk425/investigators/internal/characters"
	"github.com/sparrowhawk425/investigators/internal/commands"
	"github.com/sparrowhawk425/investigators/internal/config"
	"github.com/sparrowhawk425/investigators/internal/functions"
	"github.com/sparrowhawk425/investigators/internal/gamelogic"
	"github.com/sparrowhawk425/investigators/internal/gameobjects"
	"github.com/sparrowhawk425/investigators/internal/nameapi"
	"github.com/sparrowhawk425/investigators/internal/times"
)

type Level struct {
	Country   config.Country
	Days      int
	Caught    []characters.Character
	Escaped   []characters.Character
	Bosses    []characters.Character
	PostCards []characters.CriminalPostCard
}

func SelectLevel(cfg *config.Config) config.Country {
	// Select country
	countryNames := lo.Map(cfg.Countries, func(c config.Country, i int) string { return c.GetName().String() })
	idx := gamelogic.MenuSelect(cfg.Scanner, "Select a Country to begin your investigation:", countryNames)
	return cfg.Countries[idx]
}

func PlayLevel(country config.Country, bosses []characters.Character, cfg *config.Config) Level {

	// Build the game
	color.Green("Travelling to %s...", country.GetName())
	gs := BuildLevel(country, bosses, cfg)
	fmt.Printf("Our sources indicate %d Syndicate members are currently in the area.\n", len(gs.Criminals))

	// REPL game loop
	commandMap := commands.GetCommandMap()
	var err error
	for len(gs.Criminals) > 0 {
		gs.PrintDay()

		// Get player input
		fmt.Print("What do you wish to do? > ")
		gs.Scanner.Scan()
		cleanText := functions.CleanInput(gs.Scanner.Text())

		// Run command
		update := false
		cmd, exists := commandMap[cleanText[0]]
		if exists {
			update, err = cmd.Callback(&gs, cleanText[1:])
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
		if update {
			gs.Update()
		}
	}
	commands.ReviewLevel(&gs)

	return Level{
		Days:      gs.DayNumber,
		Country:   gs.Country,
		Caught:    gs.Caught,
		Escaped:   gs.Escaped,
		Bosses:    gs.Bosses,
		PostCards: gs.Player.GetPostCards(),
	}
}

func BuildLevel(country config.Country, bosses []characters.Character, cfg *config.Config) gamelogic.GameState {

	gs := gamelogic.GameState{
		Scanner:   cfg.Scanner,
		DayNumber: 1,
		TimeOfDay: times.Morning,
		Country:   country,
	}
	gs.Player.Name = cfg.PlayerName

	// Handle bosses (including bosses from a previous level)
	gs.Bosses = createBosses(cfg.GetNumBosses(), slices.DeleteFunc(cfg.Countries, func(c config.Country) bool { return c.Name == country.Name }))
	localBossIdx := slices.IndexFunc(bosses, func(c characters.Character) bool { return c.Traits.Nationality == country.GetCode() })
	if localBossIdx != -1 {
		color.Green("You have successfully tracked a Syndicate boss to their home base.")
		gs.People = append(gs.People, bosses[localBossIdx])
		gs.Criminals = append(gs.Criminals, &gs.People[len(gs.People)-1])
		// TODO: Add a Mega boss postcard track?
	}
	// Add locations and people to game
	results, err := nameapi.MakeHTTPGetRequest(country, cfg.GetNumPlaces())
	if err != nil {
		log.Fatalf("Error getting locations from API: %v", err)
	}
	half := len(results) / 2
	locationResults := results[:half]
	peopleResults := results[half:]

	// Create locations
	numCrimeLocs := len(gameobjects.CrimeLocations) * 2
	crimeLocations := locationResults[:numCrimeLocs]
	for i, apiLoc := range crimeLocations {
		gs.CrimePlaces = append(gs.CrimePlaces, gameobjects.CreateLocation(apiLoc.Location, gameobjects.CrimeLocations[i%len(gameobjects.CrimeLocations)], true))
	}
	apiLocations := lo.Map(locationResults[numCrimeLocs:], func(character nameapi.Character, i int) nameapi.Location { return character.Location })
	gs.Places = gameobjects.CreateRandomLocations(apiLocations)

	// Create people
	bossIdx := 0
	for i := range peopleResults {
		if i < cfg.GetNumCriminals() {
			c := characters.CreateRandomCharacter(peopleResults[i], characters.CriminalRoles[rand.IntN(len(characters.CriminalRoles))])
			c.Boss = &gs.Bosses[bossIdx%len(gs.Bosses)]
			fmt.Printf("Boss: %s, Nationality: %s\n", c.Boss.GetName(), c.Boss.Traits.Nationality)
			countryCode := c.Boss.Traits.Nationality
			addPostCards(cfg, countryCode, &c)
			bossIdx++
			gs.People = append(gs.People, c)
			gs.Criminals = append(gs.Criminals, &c)
		} else {
			c := characters.CreateRandomCharacter(peopleResults[i], characters.RegularRoles[rand.IntN(len(characters.RegularRoles))])
			c.SetJobLocation(c.FindTarget(gs.Places))
			gs.People = append(gs.People, c)
		}
	}

	// Add residences
	for i := range gs.People {
		gs.Places = append(gs.Places, gs.People[i].Address)
	}
	// Sort lists
	slices.SortFunc(gs.Places, func(l1, l2 gameobjects.Location) int {
		if l1.Address.Number < l2.Address.Number {
			return -1
		} else if l1.Address.Number > l2.Address.Number {
			return 1
		}
		if l1.Address.Name < l2.Address.Name {
			return -1
		} else if l1.Address.Name > l2.Address.Name {
			return 1
		}
		return 0
	})
	slices.SortFunc(gs.People, func(p1, p2 characters.Character) int {
		if p1.GetFirstName() < p2.GetFirstName() {
			return -1
		} else if p1.GetFirstName() > p2.GetFirstName() {
			return 1
		}
		if p1.GetLastName() < p2.GetLastName() {
			return -1
		} else if p1.GetLastName() > p2.GetLastName() {
			return 1
		}
		return 0
	})
	return gs
}

// Create bosses from a different country that the local members will report to
func createBosses(numBosses int, countries []config.Country) []characters.Character {
	bosses := make([]characters.Character, numBosses)
	for i := range numBosses {
		results, err := nameapi.MakeHTTPGetRequest(countries[rand.IntN(len(countries))], 1)
		if err != nil {
			log.Fatalf("Error getting locations from API: %v", err)
		}
		bosses[i] = characters.CreateRandomCharacter(results[0], characters.CreateCapo())
	}
	return bosses
}

func addPostCards(cfg *config.Config, countryCode config.CountryCode, person *characters.Character) {
	criminalCards := []characters.CriminalPostCard{}

	card := cfg.PostCards.Get(countryCode)
	criminalCards = append(criminalCards, card.Image)
	otherCountries := slices.DeleteFunc(cfg.Countries, func(c config.Country) bool { return c.GetCode() == countryCode })
	for range 2 {
		postCard := cfg.PostCards.Get(otherCountries[rand.IntN(len(otherCountries))].GetCode())
		criminalCards = append(criminalCards, postCard.Image)
	}
	slices.SortFunc(criminalCards, func(c1, c2 characters.CriminalPostCard) int {
		if c1[0] < c2[0] {
			return -1
		} else if c1[0] > c2[0] {
			return 1
		}
		return 0
	})
	person.SetPostCards(criminalCards)
}
