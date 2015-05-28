package server_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"

	"github.com/kinhouse/kh-site/fakes"
	. "github.com/kinhouse/kh-site/server"
	"github.com/kinhouse/kh-site/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	var writer *httptest.ResponseRecorder

	BeforeEach(func() {
		writer = httptest.NewRecorder()
	})

	It("delegates to the page factory for static pages", func() {
		fakePageFactory := &fakes.PageFactory{
			FD_StaticPages: map[string][]byte{"whatever": []byte("Some <form></form> goes here.")},
		}
		serverConfig := ServerConfig{
			Data:        nil,
			AssetNames:  []string{},
			PageFactory: fakePageFactory,
		}
		router := serverConfig.BuildRouter()

		request, err := http.NewRequest("GET", "/whatever", nil)
		if err != nil {
			panic(err)
		}

		router.ServeHTTP(writer, request)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(ContainSubstring("<form"))
	})

	Describe("GET to /rsvp/all", func() {
		It("Returns all RSVP data when given valid credentials", func() {
			fakePersist := &fakes.Persist{
				Rsvps: []types.Rsvp{
					types.Rsvp{
						FullName: "Some one",
						Email:    "someone@example.com",
						Decline:  false,
						Count:    3,
					},
					types.Rsvp{
						FullName: "Another person",
						Email:    "",
						Decline:  true,
						Count:    0,
					},
				},
			}
			serverConfig := ServerConfig{
				Data:                fakePersist,
				AssetNames:          []string{},
				PageFactory:         &fakes.PageFactory{},
				RsvpListCredentials: map[string]string{"username": "foobar!"},
			}
			router := serverConfig.BuildRouter()

			request, err := http.NewRequest("GET", "/rsvp/all", nil)
			if err != nil {
				panic(err)
			}
			request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("username:foobar!")))

			router.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusOK))
			Expect(writer.HeaderMap["Content-Type"]).To(Equal([]string{"text/csv"}))
			expectedData := `FullName,Email,Count
Some one,someone@example.com,3
Another person,,0
`
			body := writer.Body.String()
			Expect(body).To(Equal(expectedData))
		})

		It("Returns a 401 when credentials are missing", func() {
			serverConfig := ServerConfig{
				Data:                &fakes.Persist{},
				AssetNames:          []string{},
				PageFactory:         &fakes.PageFactory{},
				RsvpListCredentials: map[string]string{"username": "foobar!"},
			}
			router := serverConfig.BuildRouter()

			request, err := http.NewRequest("GET", "/rsvp/all", nil)
			if err != nil {
				panic(err)
			}

			router.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusUnauthorized))

		})

	})
})
