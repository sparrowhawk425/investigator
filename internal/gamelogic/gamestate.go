package gamelogic

import (
	"bufio"
	"fmt"
	"log"
	"math/rand/v2"
	"slices"

	"github.com/fatih/color"
	"github.com/samber/lo"

	"github.com/sparrowhawk425/investigators/internal/characters"
	"github.com/sparrowhawk425/investigators/internal/config"
	"github.com/sparrowhawk425/investigators/internal/gameobjects"
	"github.com/sparrowhawk425/investigators/internal/nameapi"
	"github.com/sparrowhawk425/investigators/internal/times"
)

type GameState struct {
	Scanner     *bufio.Scanner
	DayNumber   int
	WeekDay     times.DayOfTheWeek
	TimeOfDay   times.TimeOfDay
	Player      characters.Player
	Bosses      []characters.Character
	Places      []gameobjects.Location
	CrimePlaces []gameobjects.Location
	People      []characters.Character
	Criminals   []*characters.Character
	Escaping    []characters.Character
	Escaped     []characters.Character
	Caught      []characters.Character
	Crimes      []Crime

	Bolos []characters.Character

	newCrimes bool
}

func (gs GameState) PrintDay() {
	fmt.Printf("\n%s, Day: %d Time: %s\n", gs.WeekDay, gs.DayNumber, gs.TimeOfDay.GetName())
}

func (gs GameState) GetTimeOfDay() times.TimeOfDay {
	return gs.TimeOfDay
}

func (gs GameState) GetDayOfTheWeek() times.DayOfTheWeek {
	return gs.WeekDay
}

func (gs *GameState) NextDay() {
	gs.DayNumber++
	gs.WeekDay = gs.WeekDay.NextDay()
}

func (gs *GameState) BuildGame(cfg *config.Config) {

	// Select country
	countryNames := lo.Map(cfg.Countries, func(c config.Country, i int) string { return c.GetName().String() })
	idx := MenuSelect(gs.Scanner, "Select a Country to begin your investigation:", countryNames)
	country := cfg.Countries[idx]
	color.Green("Travelling to %s...", country.GetName())

	gs.Bosses = createBosses(cfg.GetNumBosses(), slices.DeleteFunc(cfg.Countries, func(c config.Country) bool { return c.Name == country.Name }))
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
}

// Create bosses from a different country that the local members will report to
func createBosses(numBosses int, countries []config.Country) []characters.Character {
	bosses := make([]characters.Character, numBosses)
	for i := range numBosses {
		char, err := nameapi.MakeHTTPGetRequest(countries[rand.IntN(len(countries))], 1)
		if err != nil {
			log.Fatalf("Error getting locations from API: %v", err)
		}
		bosses[i] = characters.CreateRandomCharacter(char[0], characters.CreateCapo())
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

func (gs GameState) GetLocations() []gameobjects.Location {
	return gs.Places
}

func (gs GameState) GetLocationsByType(locTypes []gameobjects.LocationType) []gameobjects.Location {
	filters := []func(gameobjects.Location, int) bool{}
	filters = append(filters, gameobjects.FilterLocationsByType(locTypes))
	return gs.GetFilteredLocations(filters)
}

func (gs GameState) GetLocationsByLootType(lootTypes []gameobjects.LootType) []gameobjects.Location {
	filters := []func(gameobjects.Location, int) bool{}
	filters = append(filters, gameobjects.FilterLocationsByLootType(lootTypes))
	return gs.GetFilteredLocations(filters)
}

func (gs GameState) GetFilteredLocations(filters []func(gameobjects.Location, int) bool) []gameobjects.Location {
	locations := gs.Places
	for _, filter := range filters {
		locations = lo.Filter(locations, filter)
	}
	return locations
}

func (gs *GameState) AddCharacterToLocation(location gameobjects.Location, character characters.Character) {
	for i := range gs.Places {
		if gs.Places[i].Equals(location) {
			gs.Places[i].Visitors = append(gs.Places[i].Visitors, character)
			return
		}
	}
}

func (gs *GameState) CreateClue(location gameobjects.Location, clue string) {
	for i := range gs.Places {
		if gs.Places[i].Equals(location) {
			gs.Places[i].AddClue(fmt.Sprintf("Day %d: %s", gs.DayNumber, clue))
			return
		}
	}
}

func (gs *GameState) CreateCrime(location gameobjects.Location, name string, loot []gameobjects.Loot) {

	gs.newCrimes = true
	gs.Crimes = append(gs.Crimes, Crime{
		Day:        gs.DayNumber,
		TimeOfDay:  gs.TimeOfDay,
		Location:   location,
		Type:       name,
		StolenLoot: loot,
	})
}

func (gs *GameState) TransferItems(lootType gameobjects.LootType, amount int, src gameobjects.ItemHolder, dest gameobjects.ItemHolder) {
	src.RemoveItems(lootType, amount)
	dest.AddItems(lootType, amount)
}

func (gs GameState) GetPlayer() characters.Player {
	return gs.Player
}

func (gs GameState) HasBolo(person characters.Character) bool {
	return slices.ContainsFunc(gs.Bolos, person.Equals)
}

func (gs GameState) CreateBoloAlert(location *gameobjects.Location, name string) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("%s %s spotted near %s\n", red("BOLO Alert:"), name, location.Address)
	if !slices.ContainsFunc(gs.Places, location.Equals) {
		color.Green("Location not previously monitored. Adding %s to tracked locations", location.Address)
		gs.Places = append(gs.Places, *location)
	}
}

func matchCriminalFunc(person *characters.Character) func(*characters.Character) bool {
	return func(c *characters.Character) bool {
		return c.Equals(*person)
	}
}

func (gs *GameState) ArrestCriminal(target characters.Character) {

	fmt.Printf("Arresting %s...\n", target.GetName())
	if slices.ContainsFunc(gs.Criminals, matchCriminalFunc(&target)) {
		gs.RemoveCriminal("You have successfully identified and arrested a member of the Syndicate. Well done.", target, true)
		// Remove criminal from the list of active BOLOs
		if slices.ContainsFunc(gs.Bolos, target.Equals) {
			gs.Bolos = slices.DeleteFunc(gs.Bolos, target.Equals)
		}
		fmt.Print("Continue...")
		gs.Scanner.Scan()
		color.Green("Searching their hideout, you find 3 postcards:")
		for _, card := range target.GetPostCards() {
			card.Print()
			gs.Player.AddPostCard(card)
		}
	} else {
		fmt.Printf("Unfortunately, %s is not a member of the Syndicate\n", target.GetName())
	}
}

func (gs *GameState) SetCriminalEscaping(person characters.Character) {
	gs.Escaping = append(gs.Escaping, person)
}

func (gs *GameState) RemoveCriminal(msg string, person characters.Character, wasCaught bool) {
	fmt.Println(msg)
	gs.People = slices.DeleteFunc(gs.People, person.Equals)
	gs.Criminals = slices.DeleteFunc(gs.Criminals, matchCriminalFunc(&person))
	if wasCaught {
		gs.Caught = append(gs.Caught, person)
	} else {
		gs.Escaped = append(gs.Escaped, person)
	}
	if len(gs.Criminals) == 0 {
		fmt.Println("There are no more Sydicate members nearby")
	} else if len(gs.Criminals) == 1 {
		fmt.Println("There is still 1 more Syndicate member in the area")
	} else {
		fmt.Printf("There are %d more Syndicate members in the area\n", len(gs.Criminals))
	}
}

func (gs GameState) GetPostCards() {
	// TODO: Pick three post cards to show
}

func (gs *GameState) Update() {

	// Reset location visitors
	for i := range gs.Places {
		gs.Places[i].Visitors = nil
	}
	for i := range gs.People {
		gs.People[i].PerformAction(gs)
	}
	// Announce criminal milestones (50% to reaching goal)
	for i := range gs.Criminals {
		for _, alert := range gs.Criminals[i].Goal.Alerts {
			color.Yellow(alert)
		}
		gs.Criminals[i].Goal.Alerts = nil
	}

	// TODO: Need better way to incorporate the player in the update loop
	if gs.Player.Action != nil {
		gs.Player.Action.Act(gs, &gs.Player.Character)
	}
	if gs.Player.CurrentLocation != nil {
		fmt.Printf("%s is currently at:\n", gs.Player.Name)
		gs.Player.CurrentLocation.Print()
	}
	// Criminals escape at the end of the sequence
	for _, escapee := range gs.Escaping {
		gs.RemoveCriminal("A member of the Syndicate has left the area...", escapee, false)
	}
	gs.Escaping = nil

	gs.TimeOfDay = times.TransitionTimeOfDay(gs.TimeOfDay)
	if gs.TimeOfDay == times.Morning {
		gs.NextDay()
		if gs.newCrimes {
			gs.newCrimes = false
			color.Red("A new crime has been reported!")
		}
		if gs.GetDayOfTheWeek() == times.Monday {
			// TODO: Add building loot
			gs.updateLoot()
		}
	}
}

func (gs *GameState) updateLoot() {
	for i := range gs.Places {
		switch gs.Places[i].Type {
		case gameobjects.PawnShop:
			gs.Places[i].AddItems(gameobjects.Money, 50)
		case gameobjects.Hotel:
			gs.Places[i].AddItems(gameobjects.Money, 25)
		case gameobjects.Casino:
			gs.Places[i].AddItems(gameobjects.Money, 100)
		}
	}
}
