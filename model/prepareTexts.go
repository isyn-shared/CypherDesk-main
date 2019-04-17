package model

import (
	"CypherDesk-main/stemmer"
	"regexp"
	"strings"
)

type Unigram struct {
	text     string
	unigrams []string
}

func (u *Unigram) init(text string) {
	u.text = text
	words := strings.Split(text, " ")

	reg, _ := regexp.Compile("[^a-zA-Z]+")
	for i, w := range words {
		words[i] = string(stemmer.Stem([]byte(reg.ReplaceAllString(w, ""))))
	}

	u.unigrams = words
}
