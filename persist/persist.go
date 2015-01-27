package persist

import (
	"fmt"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/datastore"
)

const (
	CloudKeyEnvVar = "GOOGLE_CLOUD_KEY"
	ProjectId      = "kh-site"
	KindRsvp       = "rsvp"
)

type Persist struct {
	context context.Context
}

func NewPersist() (Persist, error) {
	ctx, err := createGoogleCloudContext()
	if err != nil {
		return Persist{}, err
	}
	return Persist{context: ctx}, nil
}

func createGoogleCloudContext() (context.Context, error) {
	jsonString := os.Getenv(CloudKeyEnvVar)
	if jsonString == "" {
		return nil, fmt.Errorf("missing env var %q", CloudKeyEnvVar)
	}
	conf, err := google.JWTConfigFromJSON([]byte(jsonString), datastore.ScopeDatastore, datastore.ScopeUserEmail)
	if err != nil {
		return nil, err
	}
	return cloud.NewContext(ProjectId, conf.Client(oauth2.NoContext)), nil
}

type Rsvp struct {
	FullName string
	Email    string
	Guests   int `datastore:",noindex"`
}

func (p Persist) GetAllRSVPs() ([]Rsvp, error) {
	var results []Rsvp

	query := datastore.NewQuery(KindRsvp)
	_, err := query.GetAll(p.context, &results)
	return results, err
}

func (p Persist) InsertNewRSVP(rsvp Rsvp) (int64, error) {
	key := datastore.NewIncompleteKey(p.context, KindRsvp, nil)
	key, err := datastore.Put(p.context, key, &rsvp)
	fmt.Printf("just put: %+v\n", key)
	return key.ID(), err
}
