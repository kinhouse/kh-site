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
)

func CreateGoogleCloudContext() (context.Context, error) {
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

func DatastoreTest(cloudContext context.Context) string {
	//cloudContext = datastore.WithNamespace(cloudContext, "")
	query := datastore.NewQuery("")
	iter := query.Run(cloudContext)
	for {
		k, e := iter.Next(nil)
		if e == datastore.Done {
			break
		}
		if e != nil {
			panic(e)
		}
		fmt.Printf("\nkey is: %+v\n", k)
	}
	type Rsvp struct {
		name   string
		email  string
		guests int `datastore:",noindex"`
	}
	return "foo"

	key := datastore.NewIncompleteKey(cloudContext, "rsvp", nil)
	key, err := datastore.Put(cloudContext, key, &Rsvp{
		name:   "Some Person",
		email:  "example@example.com",
		guests: 2,
	})

	if err != nil {
		panic(err)
	}

	return key.Name()
}
