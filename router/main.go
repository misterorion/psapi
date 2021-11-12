package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/api/iterator"
)

var ctx = context.Background()
var fs = newFireStoreClient(ctx)
var projectID = "mechapower"

type jsonCharacter struct {
	Name   string   `json:"name,omitempty"`
	Race   string   `json:"race,omitempty"`
	Gender string   `json:"gender,omitempty"`
	Age    int      `json:"age,omitempty"`
	Born   string   `json:"born,omitempty"`
	Spells []string `json:"spells,omitempty"`
}

type jsonGame struct {
	Title          string   `json:"title,omitempty"`
	Title_Japanese string   `json:"title_japanese,omitempty"`
	Released       string   `json:"released,omitempty"`
	Characters     []string `json:"characters,omitempty"`
}

type myChar string
type myGame string

func newFireStoreClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

func dbGetCharacter(id string) (*jsonCharacter, error) {
	dsnap, err := fs.Collection("PSDB").Doc("api").Collection("characters").Doc(id).Get(ctx)
	// j, err := json.MarshalIndent(dsnap.Data(), "", "    ")
	if err != nil {
		return nil, err
	}
	var c jsonCharacter
	dsnap.DataTo(&c)
	return &c, nil
}

func dbGetGame(id string) (*jsonGame, error) {
	dsnap, err := fs.Collection("PSDB").Doc("api").Collection("games").Doc(id).Get(ctx)
	// j, err := json.MarshalIndent(dsnap.Data(), "", "    ")
	if err != nil {
		return nil, err
	}
	var g jsonGame
	dsnap.DataTo(&g)
	return &g, nil
}

func dbGetCollection(collection string) []byte {
	iter := fs.Collection("PSDB").Doc("api").Collection(collection).Documents(ctx)

	type M map[string]interface{}

	var collectionMap []M

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		collectionMap = append(collectionMap, doc.Data())
	}
	j, _ := json.Marshal(collectionMap)
	return j
}

func main() {
	defer fs.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Welcome to the Phantasy Star API!"))
	})

	r.Route("/characters", func(r chi.Router) {
		r.Get("/", getCharacters)
		r.Route("/{characterID}", func(r chi.Router) {
			r.Use(CharCtx)
			r.Get("/", getCharacter)
		})
	})

	r.Route("/games", func(r chi.Router) {
		r.Get("/", getGames)
		r.Route("/{gameID}", func(r chi.Router) {
			r.Use(GameCtx)
			r.Get("/", getGame)
		})
	})

	fmt.Println("Listening for eonnections on port 8080")
	http.ListenAndServe(":8080", r)
}

func getCharacters(w http.ResponseWriter, r *http.Request) {
	j := dbGetCollection("characters")
	w.Write(j)
}

func getGames(w http.ResponseWriter, r *http.Request) {
	j := dbGetCollection("games")
	w.Write(j)
}

func CharCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		charID := chi.URLParam(r, "characterID")
		char, err := dbGetCharacter(charID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		var c myChar
		ctx := context.WithValue(r.Context(), c, char)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GameCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gameID := chi.URLParam(r, "gameID")
		char, err := dbGetGame(gameID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		var g myGame
		ctx := context.WithValue(r.Context(), g, char)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getCharacter(w http.ResponseWriter, r *http.Request) {
	var c myChar
	character, ok := r.Context().Value(c).(*jsonCharacter)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	j, _ := json.Marshal(character)
	w.Write(j)
}

func getGame(w http.ResponseWriter, r *http.Request) {
	var g myGame
	game, ok := r.Context().Value(g).(*jsonGame)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	j, _ := json.Marshal(game)
	w.Write(j)
}
