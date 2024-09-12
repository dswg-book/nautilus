package server

import (
	"fmt"
	"math/rand"
)

var adjectives = []string{
	"gentle",
	"robust",
	"violent",
	"harmless",
	"wobbly",
	"woke",
	"old",
	"young",
	"shiny",
	"dingy",
	"slow",
	"rapid",
}

var pronouns = []string{
	"puppy",
	"kitty",
	"dog",
	"cat",
	"turtle",
	"rabbit",
	"lion",
	"tiger",
	"hawk",
	"falcon",
	"poodle",
	"hound",
	"goat",
	"lamb",
}

func generateName() string {
	a := adjectives[rand.Intn(len(adjectives)-1)]
	b := pronouns[rand.Intn(len(pronouns)-1)]
	c := rand.Intn(1000)
	return fmt.Sprintf("%s-%s%d", a, b, c)
}
