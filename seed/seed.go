package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/BurntSushi/toml"
	"google.golang.org/api/iterator"
)

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
	Title         string   `firestore:"title,omitempty"`
	TitleJapanese string   `firestore:"titleJapanese,omitempty"`
	Released      string   `firestore:"released,omitempty"`
	Characters    []string `firestore:"characters,omitempty"`
}

func main() {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "mechapower")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	var config tomlConfig
	_, err = toml.DecodeFile("seed.toml", &config)
	if err != nil {
		log.Fatalf("Failed to decode toml: %v", err)
	}

	log.Println("Cleaning data")
	err = deleteCollection(ctx, client, client.Collection("PSDB/api/characters"), 10)
	if err != nil {
		log.Fatalf("Failed to clean characters: %v", err)
	}

	err = deleteCollection(ctx, client, client.Collection("PSDB/api/games"), 10)
	if err != nil {
		log.Fatalf("Failed to clean games: %v", err)
	}

	log.Println("Seeding data")
	for id, c := range config.Characters {
		_, err := client.Collection("PSDB/api/characters").Doc(id).Set(ctx, c)
		if err != nil {
			log.Fatalf("Failed to seed characters: %v", err)
		}
	}
	for id, g := range config.Games {
		_, err := client.Collection("PSDB/api/games").Doc(id).Set(ctx, g)
		if err != nil {
			log.Fatalf("Failed to seed games: %v", err)
		}
	}
	log.Println("Done!")
}

func deleteCollection(ctx context.Context, client *firestore.Client,
	ref *firestore.CollectionRef, batchSize int) error {

	for {
		iter := ref.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		batch := client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}
}
