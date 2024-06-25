package vocab

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

type Game struct {
	Cards          []Card
	Complete       bool
	NumberOfRounds int
	Rounds         []Round
	Score          int
	Round          int
}

type Round struct {
	Card    Card
	Cards   []Card
	Correct bool
}

type Card struct {
	Chinese string
	Pinyin  string
	English string
}

var Cards []Card

func readVocab() ([]Card, error) {
	basepath, _ := os.Getwd()
	content, err := os.ReadFile(basepath + "/data/vocab.json")

	if err != nil {
		return nil, fmt.Errorf("error reading file")
	}

	var vocab []Card
	err = json.Unmarshal(content, &vocab)

	if err != nil {
		return nil, fmt.Errorf("error during Unmarshall")
	}

	return vocab, nil
}

func NewGame(rounds int) Game {
	cards, _ := readVocab()

	game := Game{NumberOfRounds: rounds, Cards: cards, Complete: false, Score: 0, Round: 0}

	for i := 0; i < rounds; i++ {
		newRound(&game)
	}

	return game
}

func randomCards(g *Game) []Card {
	var cards []Card

	for i := 0; i < 4; i++ {
		index := rand.Intn(len(g.Cards))
		cards = append(cards, g.Cards[index])
	}

	return cards
}

func newRound(g *Game) {
	cards := randomCards(g)

	index := rand.Intn(len(cards))

	r := Round{Card: cards[index], Cards: cards}

	g.Rounds = append(g.Rounds, r)
}

func (g *Game) MarkAnswer(c string) {
	if g.Rounds[g.Round].Card.Chinese == c {
		g.Rounds[g.Round].Correct = true
	} else {
		g.Rounds[g.Round].Correct = false
	}
}

func (g *Game) NextRound() {
	if g.Round < g.NumberOfRounds-1 {
		g.Round++
	} else {
		g.Complete = true
	}
}

func (g *Game) GetScore() float32 {
	for _, round := range g.Rounds {
		if round.Correct {
			g.Score++
		}
	}

	return float32(g.Score) / float32(g.NumberOfRounds)
}
