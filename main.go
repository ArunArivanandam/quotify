package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	router := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.handlerReset)
	apiRouter.Post("/validate_quote", handlerQuotesValidate)
	router.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)
	router.Mount("/admin", adminRouter)

	CorsMux := middlewareCors(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: CorsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}

func handlerQuotesValidate(w http.ResponseWriter, r *http.Request) {
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

	const maxQuoteLength = 300
	if len(params.Quote) > maxQuoteLength {
		respondWithError(w, http.StatusBadRequest, "Quote is too long")
		return
	}

	if params.Author == "" {
		params.Author = "unknown"
	}

	badWords := map[string]struct{}{
		"fuck":   {},
		"shit":   {},
		"badass": {},
		"dick":   {},
		"sex":    {},
	}

	cleaned := getCleanedQuote(params.Quote, badWords)

	respondWithJSON(w, http.StatusOK, parameters{
		Quote:  cleaned,
		Author: params.Author,
	})

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

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5xx error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
