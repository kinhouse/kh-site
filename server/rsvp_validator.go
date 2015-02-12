package server

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/kinhouse/kh-site/types"
)

func ValidateRsvp(rsvp types.Rsvp) error {
	if strings.TrimSpace(rsvp.FullName) == "" {
		return errors.New("Please provide your name!")
	}
	_, err := mail.ParseAddress(rsvp.Email)
	if err != nil {
		return errors.New("Please provide a valid email address.")
	}
	if rsvp.Count > 9 {
		return errors.New("That's quite a group!  Please get in touch with us.")
	}
	if rsvp.Count < 0 {
		return errors.New("Your guest count is...negative?")
	}
	if !rsvp.Decline && rsvp.Count == 0 {
		return errors.New("Please indicate if you're a Yes or a No.")
	}
	if rsvp.Decline && rsvp.Count != 0 {
		return errors.New("If you're declining, please set your group size to 0.")
	}
	return nil
}
