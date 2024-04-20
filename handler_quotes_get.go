package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi"
)

func (cfg *apiConfig) handlerQuotesRetrieve(w http.ResponseWriter, r *http.Request) {

	dbQuotes, err := cfg.DB.GetQuotes()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve quotes")
		return
	}
	quotes := []Quote{}
	for _, dbQuote := range dbQuotes {
		quote := Quote{
			Id:     dbQuote.Id,
			Quote:  dbQuote.Quote,
			Author: dbQuote.Author,
		}
		quotes = append(quotes, quote)
	}
	sort.Slice(quotes, func(i, j int) bool {
		return quotes[i].Id < quotes[j].Id
	})
	respondWithJSON(w, http.StatusOK, quotes)
}

func (cfg *apiConfig) handlerQuotesGet(w http.ResponseWriter, r *http.Request) {
	quoteIdString := chi.URLParam(r, "quoteId")
	quoteId, err := strconv.Atoi(quoteIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid quote ID")
		return
	}

	dbQuote, err := cfg.DB.GetQuote(quoteId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve quote")
		return
	}
	respondWithJSON(w, http.StatusOK, Quote{
		Id:     dbQuote.Id,
		Quote:  dbQuote.Quote,
		Author: dbQuote.Author,
	})
}
