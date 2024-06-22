package vocab

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"ownkng.dev/cli/types"
)

var Cards []types.Card

func readVocab() ([]types.Card, error) {
	basepath, _ := os.Getwd()
	content, err := os.ReadFile(basepath + "/data/vocab.json")

	if err != nil {
		return nil, fmt.Errorf("error reading file")
	}

	var vocab []types.Card
	err = json.Unmarshal(content, &vocab)

	if err != nil {
		return nil, fmt.Errorf("error during Unmarshall")
	}

	return vocab, nil
}

func NewGame(rounds int) types.Game {
	cards, _ := readVocab()

	game := types.Game{NumberOfRounds: rounds, Cards: cards, Complete: false, Score: 0}

	for i := 0; i < rounds; i++ {
		newRound(&game)
	}

	return game
}

func randomCharacter(g *types.Game) string {
	index := rand.Intn(len(g.Cards))

	return g.Cards[index].Chinese
}

func randomCards(g *types.Game) []types.Card {
	var cards []types.Card

	for i := 0; i < 4; i++ {
		index := rand.Intn(len(g.Cards))
		cards = append(cards, g.Cards[index])
	}

	return cards
}

func newRound(g *types.Game) {
	cards := randomCards(g)

	r := types.Round{Character: cards[0].Chinese, Cards: cards}

	g.Rounds = append(g.Rounds, r)
}

func CheckAnswer(r *types.Round, c string) {
	if r.Character == c {
		r.Correct = true
	} else {
		r.Correct = false
	}
}

func NextRound(g *types.Game) {
	if g.Round < g.NumberOfRounds {
		newRound(g)
		g.Round++
	} else {
		g.Complete = true
	}
}

func UpdateScore(g *types.Game) {
	for _, round := range g.Rounds {
		if round.Correct {
			g.Score++
		}
	}
}
