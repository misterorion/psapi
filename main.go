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

type myChar string

type Character struct {
	Name   string   `firestore:"name,omitempty"`
	Race   string   `firestore:"race,omitempty"`
	Gender string   `firestore:"gender,omitempty"`
	Age    int      `firestore:"age,omitempty"`
	Born   string   `firestore:"born,omitempty"`
	Spells []string `firestore:"spells,omitempty"`
}

func SeedCharacters(ctx context.Context, client *firestore.Client) error {
	characters := []struct {
		id string
		c  Character
	}{
		{id: "1", c: Character{Name: "Alis Landale", Race: "Human", Age: 15, Gender: "Female", Born: "AW 327, 5.25", Spells: []string{"Heal", "Bye", "Chat", "Fire", "Rope", "Fly"}}},
		{id: "2", c: Character{Name: "Myau", Race: "Musk Cat", Gender: "Male"}},
		{id: "3", c: Character{Name: "Odin", Race: "Human", Gender: "Male", Born: "AW 314, 2.26", Age: 28}},
		{id: "4", c: Character{Name: "Noah", Race: "Human", Gender: "Male", Born: " AW 315, 3.24", Age: 27}},
	}

	for _, c := range characters {
		_, err := client.Collection("PSDB").Doc("ps1").Collection("characters").Doc(c.id).Set(ctx, c.c)
		// _, err := client.Collection("characters").Doc(c.id).Set(ctx, c.c)
		if err != nil {
			return err
		}
	}

	return nil
}

func newFireStoreClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

func dbGetCharacter(id string) (*jsonCharacter, error) {
	dsnap, err := fs.Collection("PSDB").Doc("ps1").Collection("characters").Doc(id).Get(ctx)
	// j, err := json.MarshalIndent(dsnap.Data(), "", "    ")
	if err != nil {
		return nil, err
	}
	var c jsonCharacter
	dsnap.DataTo(&c)
	return &c, nil
}

func dbGetCollection(collection string) []byte {
	iter := fs.Collection("PSDB").Doc("ps1").Collection(collection).Documents(ctx)

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

	// fmt.Println("About to SEED")
	// err := SeedCharacters(ctx, fs)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("Done SEEDING")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Welcome to the Phantasy Star Database!"))
	})

	r.Route("/characters", func(r chi.Router) {
		r.Get("/", getChars)
		r.Route("/{characterID}", func(r chi.Router) {
			r.Use(CharCtx)
			r.Get("/", getCharacter)
		})
	})

	fmt.Println("Listening for eonnections on port 8080")
	http.ListenAndServe(":8080", r)
}

func getChars(w http.ResponseWriter, r *http.Request) {
	j := dbGetCollection("characters")
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
