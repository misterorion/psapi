package main

// type Character struct {
// 	Name   string   `firestore:"name,omitempty"`
// 	Race   string   `firestore:"race,omitempty"`
// 	Gender string   `firestore:"gender,omitempty"`
// 	Age    int      `firestore:"age,omitempty"`
// 	Born   string   `firestore:"born,omitempty"`
// 	Spells []string `firestore:"spells,omitempty"`
// }

// func SeedCharacters(ctx context.Context, client *firestore.Client) error {
// 	characters := []struct {
// 		id string
// 		c  Character
// 	}{
// 		{id: "1", c: Character{Name: "Alis Landale", Race: "Human", Age: 15, Gender: "Female", Born: "AW 327, 5.25", Spells: []string{"Heal", "Bye", "Chat", "Fire", "Rope", "Fly"}}},
// 		{id: "2", c: Character{Name: "Myau", Race: "Musk Cat", Gender: "Male"}},
// 		{id: "3", c: Character{Name: "Odin", Race: "Human", Gender: "Male", Born: "AW 314, 2.26", Age: 28}},
// 		{id: "4", c: Character{Name: "Noah", Race: "Human", Gender: "Male", Born: " AW 315, 3.24", Age: 27}},
// 	}

// 	for _, c := range characters {
// 		_, err := client.Collection("PSDB").Doc("ps1").Collection("characters").Doc(c.id).Set(ctx, c.c)
// 		// _, err := client.Collection("characters").Doc(c.id).Set(ctx, c.c)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
