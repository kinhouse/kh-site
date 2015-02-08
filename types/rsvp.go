package types

type Rsvp struct {
	FullName string `form:"FullName" json:"fullname"`
	Email    string `form:"Email" json:"email"`
	Accept   bool   `form:"AcceptInvite" json:"accept"`
	Guests   int    `form:"Guests" json:"guests" datastore:",noindex"`
}
