package types

type Rsvp struct {
	FullName string `form:"FullName" json:"fullname"`
	Email    string `form:"Email" json:"email"`
	Decline  bool   `form:"DeclineInvite" json:"decline"`
	Count    int    `form:"Count" json:"count" datastore:",noindex"`
}
