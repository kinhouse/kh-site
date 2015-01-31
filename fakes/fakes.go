package fakes

import "github.com/kinhouse/kh-site/types"

type Persist struct{}

func (p Persist) GetAllRSVPs() ([]types.Rsvp, error) {
	return []types.Rsvp{}, nil
}

func (p Persist) InsertNewRSVP(types.Rsvp) (int64, error) {
	return 0, nil
}
