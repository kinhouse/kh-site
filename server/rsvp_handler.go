package server

import (
	"fmt"

	"github.com/kinhouse/kh-site/types"
)

func words(count int) string {
	switch count {
	case 1:
		return "you'll be"
	case 2:
		return "the two of you will be"
	case 3:
		return "the three of you will be"
	case 4:
		return "the four of you will be"
	case 5:
		return "all five (!) of you will be"
	}
	return "you will be"
}

func RsvpHandler(rsvp types.Rsvp) string {
	if rsvp.Decline {
		return "We're sorry you won't be making it, but thanks for letting us know!"
	}

	response := "Wonderful!  We're thrilled that %s joining us.  " +
		"Poke around this site for more information about the wedding and our story."
	return fmt.Sprintf(response, words(rsvp.Count))
}
