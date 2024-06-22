package types

type Game struct {
	Cards          []Card
	Complete       bool
	NumberOfRounds int
	Rounds         []Round
	Score          int
	Round          int
}

type Round struct {
	Character string
	Cards     []Card
	Correct   bool
}

type Card struct {
	Chinese string
	Pinyin  string
	English string
}
