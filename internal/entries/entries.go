package entries

import (
	"time"

	gokeepasslib "github.com/tobischo/gokeepasslib/v3"
)

func GetNearlyExpiredEntries(group gokeepasslib.Group, now time.Time, proximity time.Duration) ([]gokeepasslib.Entry, error) {
	entries := make([]gokeepasslib.Entry, 0)
	if len(group.Groups) > 0 {
		// loop over groups and recursively call the function for each one
		for index := range group.Groups {
			newEntries, err := GetNearlyExpiredEntries(group.Groups[index], now, proximity)
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
