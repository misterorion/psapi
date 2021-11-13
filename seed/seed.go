package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/BurntSushi/toml"
)

var ctx = context.Background()
var fs = newFireStoreClient(ctx)
var projectID = "mechapower"

func newFireStoreClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

type tomlConfig struct {
	Characters map[string]Character
	Games      map[string]Game
}

type Character struct {
	Name   string   `firestore:"name,omitempty"`
	Game   string   `firestore:"games,omitempty"`
	Born   string   `firestore:"born,omitempty"`
	Gender string   `firestore:"gender,omitempty"`
	Age    int      `firestore:"age,omitempty"`
	Race   string   `firestore:"race,omitempty"`
	Spells []string `firestore:"spells,omitempty"`
}

type Game struct {
	Title          string   `firestore:"title,omitempty"`
	Title_Japanese string   `firestore:"title_japanese,omitempty"`
	Released       string   `firestore:"released,omitempty"`
	Characters     []string `firestore:"characters,omitempty"`
}

func SeedData(ctx context.Context, client *firestore.Client) error {

	var config tomlConfig

	if _, err := toml.DecodeFile("seed.toml", &config); err != nil {
		fmt.Println(err)
		return err
	}

	for id, c := range config.Characters {
		_, err := client.Collection("PSDB").Doc("api").Collection("characters").Doc(id).Set(ctx, c)
		if err != nil {
			return err
		}
	}

	for id, g := range config.Games {
		_, err := client.Collection("PSDB").Doc("api").Collection("games").Doc(id).Set(ctx, g)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {

	defer fs.Close()

	fmt.Println("About to SEED")
	err := SeedData(ctx, fs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Done SEEDING")
}
