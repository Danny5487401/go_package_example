package handlers

import "github.com/Danny5487401/go_package_example/42_go-openapi/nbi/gen/server/models"

// This file holds the code that stores the PhoneBook entries
// This is a very simplistic implementations
// Two entries are hard coded so that all GET return something
// Then at run-time you can add more entries, they will be lost when the process exits

// PhoneBookDb is a map.
// Keys are strings build like this <firstName>-<LastName>.
// Values are pointers to PhoneBookEntry.
// It plays the role of our run-time database
type PbDb map[string]*models.PhoneBookEntry

var PhoneBookDb = PbDb{}

// Package initializer
func init() {
	PbDbInit()
}

// PbDbInit populate the Phone Book DB with its first 2 hard-coded entries
func PbDbInit() {
	// Hard coded fake entry 1
	epAddr := models.AddressEntry{CivicNumber: 3764, Street: "Highway 51 South", City: "Memphis", State: "Tennessee", Zip: 38116}
	ep := models.PhoneBookEntry{FirstName: "Elvis", LastName: "Presley", PhoneNumber: "901-555-5823", Address: &epAddr}

	// Hard coded fake entry 2
	brAddr := models.AddressEntry{CivicNumber: 100, Street: "Ocean Drive", City: "Los Angeles", State: "California", Zip: 90201}
	bw := models.PhoneBookEntry{FirstName: "Barry", LastName: "White", PhoneNumber: "310-555-9274", Address: &brAddr}

	// add hard-coded entries into the phone book
	PhoneBookDb[buildDbKey(ep.FirstName, ep.LastName)] = &ep
	PhoneBookDb[buildDbKey(bw.FirstName, bw.LastName)] = &bw
}

// buildDbKey return a key suitable to access the phone book entry for first and last name
func buildDbKey(first, last string) string {
	return first + "-" + last // strings.Builder{} would be more efficient, but for this demo that is OK
}

// Entries a utility functions that returns all entries in the PhoneBookDb map as a slice of pointers to PhoneBookEntry
func (pb PbDb) Entries() []*models.PhoneBookEntry {
	entries := make([]*models.PhoneBookEntry, 0, len(pb))
	for _, e := range pb {
		entries = append(entries, e)
	}
	return entries
}

func (pb PbDb) AddEntry(e *models.PhoneBookEntry) {
	key := buildDbKey(e.FirstName, e.LastName)
	pb[key] = e
}

func (pb PbDb) GetEntry(first, last string) *models.PhoneBookEntry {
	key := buildDbKey(first, last)

	entry, found := pb[key]
	if !found {
		return nil
	}

	return entry
}
