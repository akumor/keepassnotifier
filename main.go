package main

import (
	"fmt"
	"os"
	"time"

	gokeepasslib "github.com/tobischo/gokeepasslib/v3"
)

func main() {
	now := time.Now()
	proximity, err := time.ParseDuration("6h")
	if err != nil {
		fmt.Printf("failed to parse provided duration: %v", err)
		os.Exit(1)
	}
	file, _ := os.Open("./example.kdbx")

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials("example")
	err = gokeepasslib.NewDecoder(file).Decode(db)
	if err != nil {
		fmt.Printf("failed to decode specified keepass database file: %v", err)
		os.Exit(1)
	}

	db.UnlockProtectedEntries()

	nearlyExpiredEntries := make([]gokeepasslib.Entry, 0)
	for index := range db.Content.Root.Groups {
		newEntries, err := getNearlyExpiredEntries(db.Content.Root.Groups[index], now, proximity)
		if err != nil {
			fmt.Printf("failed to get nearly expired entries: %v", err)
			os.Exit(1)
		}
		nearlyExpiredEntries = append(nearlyExpiredEntries, newEntries...)
	}

	for index := range nearlyExpiredEntries {
		fmt.Println(nearlyExpiredEntries[index].GetTitle())
	}
	//fmt.Printf("%v\n", nearlyExpiredEntries)
}

func getNearlyExpiredEntries(group gokeepasslib.Group, now time.Time, proximity time.Duration) ([]gokeepasslib.Entry, error) {
	entries := make([]gokeepasslib.Entry, 0)
	if len(group.Groups) > 0 {
		// loop over groups and recursively call the function for each one
		for index := range group.Groups {
			newEntries, err := getNearlyExpiredEntries(group.Groups[index], now, proximity)
			if err != nil {
				return []gokeepasslib.Entry{}, err
			}
			entries = append(entries, newEntries...)
		}
	}
	if len(group.Entries) > 0 {
		// loop over entries and append them to slice that gets returned if the entry is expired
		for index := range group.Entries {
			if group.Entries[index].Times.Expires.Bool {
				// check if entry is expired
				if group.Entries[index].Times.ExpiryTime.Time.Sub(now) <= proximity {
					entries = append(entries, group.Entries[index])
				}
			}
		}
	}
	return entries, nil
}
