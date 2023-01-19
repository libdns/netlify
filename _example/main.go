package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/libdns/libdns"
	netlify "github.com/libdns/netlify"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("NETLIFY_TOKEN")
	if token == "" {
		fmt.Printf("NETLIFY_TOKEN not set\n")
		return
	}
	zone := os.Getenv("ZONE")
	if zone == "" {
		fmt.Printf("ZONE not set\n")
		return
	}
	provider := netlify.Provider{
		PersonalAccessToken: token,
	}

	records, err := provider.GetRecords(context.TODO(), zone)
	if err != nil {
		log.Fatalln("ERROR: %s\n", err.Error())
	}

	testName := "_acme-challenge.home"
	hasTestName := false

	for _, record := range records {
		fmt.Printf("%s (.%s): %s, %s\n", record.Name, zone, record.Value, record.Type)
		if record.Name == testName {
			hasTestName = true
		}
	}

	if !hasTestName {
		appendedRecords, err := provider.AppendRecords(context.TODO(), zone, []libdns.Record{
			libdns.Record{
				Type:  "TXT",
				Name:  testName + "." + zone,
				TTL:   0,
				Value: "20HnRk5p6rZd7TXhiMoVEYSjt5OpetC6mdovlTfJ4As",
			},
		})

		if err != nil {
			log.Fatalln("ERROR: %s\n", err.Error())
		}

		fmt.Println("appendedRecords")
		fmt.Println(appendedRecords)

		editedRecords, err := provider.SetRecords(context.TODO(), zone, []libdns.Record{
			libdns.Record{
				Type:  "TXT",
				Name:  testName,
				TTL:   0,
				Value: "6MK3aXrkLz29tOsD5RQ0n1yx9mCTS8SAtjJptIPw",
			},
		})

		if err != nil {
			log.Fatalln("ERROR: %s\n", err.Error())
		}

		records, err := provider.GetRecords(context.TODO(), zone)

		if err != nil {
			log.Fatalln("ERROR: %s\n", err.Error())
		}

		for _, record := range records {
			if record.Name == testName {
				if record.Value == "6MK3aXrkLz29tOsD5RQ0n1yx9mCTS8SAtjJptIPw" {
					fmt.Println("editedRecord")
					fmt.Println(editedRecords)
				}
			}
		}

		editedNoExistRecords, err := provider.SetRecords(context.TODO(), zone, []libdns.Record{
			libdns.Record{
				Type:  "TXT",
				Name:  testName + ".no_exist",
				TTL:   0,
				Value: "wZF6K7OHXcJvhFCwKnEx8zq3PzObeIYa3YSYdL9t",
			},
		})

		if err != nil {
			log.Fatalln("ERROR: %s\n", err.Error())
		}

		records, err = provider.GetRecords(context.TODO(), zone)

		if err != nil {
			log.Fatalln("ERROR: %s\n", err.Error())
		}

		for _, record := range records {
			if record.Name == testName + ".no_exist" {
				if record.Value == "wZF6K7OHXcJvhFCwKnEx8zq3PzObeIYa3YSYdL9t" {
					fmt.Println("editedNoExistRecord")
					fmt.Println(editedNoExistRecords)
				}
			}
		}

	} else {
		deleteRecords, err := provider.DeleteRecords(context.TODO(), zone, []libdns.Record{
			libdns.Record{
				Type: "TXT",
				Name: testName,
			},
		})

		if err != nil {
			log.Fatalln("ERROR: %s\n", err.Error())
		}

		fmt.Println("deleteRecords")
		fmt.Println(deleteRecords)

		deleteRecords, err = provider.DeleteRecords(context.TODO(), zone, []libdns.Record{
			libdns.Record{
				Type: "TXT",
				Name: testName + ".no_exist",
			},
		})

		if err != nil {
			log.Fatalln("ERROR: %s\n", err.Error())
		}

		fmt.Println("deleteRecords")
		fmt.Println(deleteRecords)
	}
}
