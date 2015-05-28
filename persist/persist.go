package persist

import (
	"fmt"
	"os"

	"github.com/kinhouse/kh-site/types"

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

func (p Persist) GetAllRSVPs() ([]types.Rsvp, error) {
	var results []types.Rsvp

	query := datastore.NewQuery(KindRsvp)
	_, err := query.GetAll(p.context, &results)
	return results, err
}
