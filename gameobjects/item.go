package gameobjects

type Loot string

const (
	Jewelry     Loot = "Jewelry"
	Money       Loot = "Money"
	Art         Loot = "Art"
	Electronics Loot = "Electronics"
	Cars        Loot = "Cars"
)

var LootList = []Loot{
	Jewelry, Money, Art, Electronics, Cars,
}
