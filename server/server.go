package server

import (
	"encoding/csv"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kinhouse/kh-site/types"

	"fmt"
	"net/http"
)

type PageFactoryInterface interface {
	GenerateDynamicPage(string, string) []byte
	StaticPages() map[string][]byte
}

type PersistInterface interface {
	GetAllRSVPs() ([]types.Rsvp, error)
	InsertNewRSVP(types.Rsvp) (int64, error)
}

type ServerConfig struct {
	Data                PersistInterface
	AssetNames          []string
	RsvpHandler         func(types.Rsvp) string
	RsvpValidator       func(types.Rsvp) error
	PageFactory         PageFactoryInterface
	AssetProvider       AssetProviderInterface
	RsvpListCredentials map[string]string
}

func (s ServerConfig) AddPageRoutes(e *gin.Engine) {
	for route, pageContent := range s.PageFactory.StaticPages() {
		d := pageContent // copy variable, for closure
		e.GET("/"+route, func(c *gin.Context) {
			c.Data(http.StatusOK, gin.MIMEHTML, d)
		})
	}
}

func (s ServerConfig) AddStaticAssetRoutes(e *gin.Engine) {
	for _, name := range s.AssetNames {
		assetPath := s.AssetProvider.GetAssetPath(name)
		e.GET("/"+name, func(c *gin.Context) { c.File(assetPath) })
	}
}

func (s ServerConfig) AddRsvpPostHandler(e *gin.Engine) {
	e.POST("/rsvp", func(c *gin.Context) {
		var rsvp types.Rsvp
		if !c.Bind(&rsvp) {
			return
		}

		err := s.RsvpValidator(rsvp)
		if err != nil {
			responseHTML := s.PageFactory.GenerateDynamicPage("Incomplete RSVP",
				fmt.Sprintf("<h1>There was a problem with your RSVP</h1><p>%s</p>", err.Error()))
			c.Data(http.StatusBadRequest, gin.MIMEHTML, responseHTML)
			return
		}

		id, err := s.Data.InsertNewRSVP(rsvp)
		if err != nil {
			panic("persisting rsvp " + err.Error())
		}

		fmt.Printf("Inserted new RSVP (%d) : %+v\n", id, rsvp)

		responseText := s.RsvpHandler(rsvp)
		responseTitle := "♡"
		if rsvp.Decline {
			responseTitle = "☹"
		}
		responseHtml := s.PageFactory.GenerateDynamicPage(responseTitle, responseText)
		c.Data(http.StatusCreated, gin.MIMEHTML, responseHtml)
	})
}

func writeRSVPsAsCSV(rsvps []types.Rsvp, resp http.ResponseWriter) error {
	w := csv.NewWriter(resp)
	w.Write([]string{"FullName", "Email", "Count"})
	for _, r := range rsvps {
		w.Write([]string{r.FullName, r.Email, strconv.FormatInt(int64(r.Count), 10)})
	}
	w.Flush()
	return w.Error()
}

func (s ServerConfig) AddRsvpListingHandler(e *gin.Engine) {
	if s.RsvpListCredentials == nil {
		fmt.Println("No credentials given for RSVP listing.  This feature won't be available")
		return
	}
	e.GET("/rsvp/all",
		gin.BasicAuth(s.RsvpListCredentials),
		func(c *gin.Context) {
			rsvps, err := s.Data.GetAllRSVPs()
			if err != nil {
				panic("reading rsvps " + err.Error())
			}
			c.Writer.Header().Set("Content-Type", "text/csv")
			c.Writer.WriteHeader(http.StatusOK)
			err = writeRSVPsAsCSV(rsvps, c.Writer)
			if err != nil {
				panic("error writing RSVPs as CSV file: " + err.Error())
			}
		})
}

func (s ServerConfig) BuildRouter() *gin.Engine {
	r := gin.Default()

	s.AddStaticAssetRoutes(r)
	s.AddPageRoutes(r)
	s.AddRsvpPostHandler(r)
	s.AddRsvpListingHandler(r)

	return r
}
