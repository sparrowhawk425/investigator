package gameobjects

type LootType string

const (
	Jewelry     LootType = "Jewelry"
	Money       LootType = "Money"
	Art         LootType = "Art"
	Electronics LootType = "Electronics"
	Cars        LootType = "Cars"
)

var LootTypeList = []LootType{
	Jewelry, Money, Art, Electronics, Cars,
}

func (lt LootType) IsType() bool {
	return true
}

func (lt LootType) String() string {
	return string(lt)
}

func (lt LootType) GetValue() int {
	switch lt {
	case Jewelry:
		return 10
	case Money:
		return 5
	case Art:
		return 20
	case Electronics:
		return 30
	case Cars:
		return 40
	default:
		return 1
	}
}

type Loot struct {
	Type     LootType
	Quantity int
	Value    int
	IsStolen bool
}
