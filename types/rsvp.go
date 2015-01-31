package types

type Rsvp struct {
	FullName string
	Email    string
	Guests   int `datastore:",noindex"`
}
