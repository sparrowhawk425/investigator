package nameapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// v1.4
// Countries: AU, BR, CA, CH, DE, DK, ES, FI, FR, GB, IE, IN, IR, MX, NL, NO, NZ, RS, TR, UA, US
type Country struct {
	Name string
	Code string
}

var Countries = []Country{
	{Name: "Australia", Code: "AU"}, {Name: "Brazil", Code: "BR"}, {Name: "Canada", Code: "CA"}, {Name: "Switzerland", Code: "CH"},
	{Name: "Germany", Code: "DE"}, {Name: "Denmark", Code: "DK"}, {Name: "Spain", Code: "ES"}, {Name: "Finland", Code: "FI"},
	{Name: "France", Code: "FR"}, {Name: "Great Britain", Code: "GB"}, {Name: "Ireland", Code: "IE"}, {Name: "India", Code: "IN"},
	{Name: "Iran", Code: "IR"}, {Name: "Mexico", Code: "MX"}, {Name: "Netherlands", Code: "NL"}, {Name: "Norway", Code: "NO"},
	{Name: "New Zealand", Code: "NZ"}, {Name: "Serbia", Code: "RS"}, {Name: "Türkiye", Code: "TR"}, {Name: "Ukraine", Code: "UA"},
	{Name: "United States", Code: "US"},
}

const nameApi = "https://randomuser.me/api"

func MakeHTTPGetRequest(country Country, numResults int) ([]Character, error) {

	url := fmt.Sprintf("%s?nat=%s&results=%d", nameApi, country.Code, numResults)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var name GeneratedName
	if err = json.Unmarshal(data, &name); err != nil {
		return nil, err
	}
	names := make([]Character, len(name.Results))
	for i, result := range name.Results {
		names[i] = Character{
			Gender:      result.Gender,
			Name:        result.Name,
			Location:    result.Location,
			DateOfBirth: result.Dob,
			Email:       result.Email,
			Nationality: result.Nat,
		}
	}
	return names, nil
}

type Character struct {
	Gender      string   `json:"gender"`
	Name        Name     `json:"name"`
	Location    Location `json:"location"`
	Email       string   `json:"email"`
	DateOfBirth struct {
		Date time.Time `json:"date"`
		Age  int       `json:"age"`
	} `json:"dob"`
	Nationality string `json:"nat"`
}

type Name struct {
	First string `json:"first"`
	Last  string `json:"last"`
}
type Street struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
}
type Location struct {
	Street   Street          `json:"street"`
	City     string          `json:"city"`
	State    string          `json:"state"`
	Country  string          `json:"country"`
	Postcode json.RawMessage `json:"postcode"`
}

type GeneratedName struct {
	Results []struct {
		Gender   string   `json:"gender"`
		Name     Name     `json:"name"`
		Location Location `json:"location"`
		Email    string   `json:"email"`
		Dob      struct {
			Date time.Time `json:"date"`
			Age  int       `json:"age"`
		} `json:"dob"`
		Nat string `json:"nat"`
	} `json:"results"`
	Info struct {
		Seed    string `json:"seed"`
		Results int    `json:"results"`
		Page    int    `json:"page"`
		Version string `json:"version"`
	} `json:"info"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
