package server_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/kinhouse/kh-site/fakes"
	. "github.com/kinhouse/kh-site/server"

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
})
