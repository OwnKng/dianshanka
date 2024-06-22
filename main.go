package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
)

var(
	name string 
	beverage string
)

func main() {
	form := huh.NewForm(
		huh.NewGroup(
		huh.NewInput().Title("What's your name?").Value(&name),
		),
		huh.NewGroup(
			huh.NewSelect[string]().Title("Chose beverage").Options(
				huh.NewOption("Beer, please!", "beer"), 
				huh.NewOption("Wine would be fine", "wine"),
				huh.NewOption("Water, just water", "water"),
			).Value(&beverage),
		),
	)

	err := form.Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hey, %s!\n", name)
	fmt.Printf("One %s, comming right up", beverage)
}