package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Quote struct {
	Id     int    `json:"id"`
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

func (cfg *apiConfig) handlerQuotesCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Quote  string `json:"quote"`
		Author string `json:"author"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleaned, err := validateQuote(params.Quote)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if params.Author == "" {
		params.Author = "unknown"
	}

	quote, err := cfg.DB.CreateQuote(cleaned, params.Author)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create quote")
		return
	}

	respondWithJSON(w, http.StatusCreated, Quote{
		Id:     quote.Id,
		Quote:  quote.Quote,
		Author: quote.Author,
	})
}

func validateQuote(body string) (string, error) {

	const maxQuoteLength = 300
	if len(body) > maxQuoteLength {
		return "", errors.New("Quote is too long")
	}

	badWords := map[string]struct{}{
		"fuck":   {},
		"shit":   {},
		"badass": {},
		"sex":    {},
	}

	cleaned := getCleanedQuote(body, badWords)

	return cleaned, nil

}

func getCleanedQuote(quote string, badWords map[string]struct{}) string {
	words := strings.Split(quote, " ")
	fmt.Println(words)
	for i, word := range words {
		if word[len(word)-1] == ',' {
			newWord := word[0 : len(word)-1]
			lowWord := strings.ToLower(newWord)
			if _, ok := badWords[lowWord]; ok {
				re := regexp.MustCompile("[a-zA-Z]")
				replacedStr := re.ReplaceAllString(string(lowWord[2:]), "*")
				words[i] = string(lowWord[0]) + replacedStr + string(lowWord[len(lowWord)-1:]) + ","
			}

		}
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			re := regexp.MustCompile("[a-zA-Z]")
			replacedStr := re.ReplaceAllString(string(loweredWord[2:]), "*")
			words[i] = string(loweredWord[0]) + replacedStr + string(loweredWord[len(loweredWord)-1:])
		}
	}

	cleaned := strings.Join(words, " ")
	return cleaned
}
