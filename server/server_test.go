package server_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/kinhouse/kh-site/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	It("responds to GET /event with the event info", func() {
		serverConfig := BuildServerConfig(nil)
		router := serverConfig.BuildRouter()

		writer := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/event", nil)
		if err != nil {
			panic(err)
		}

		router.ServeHTTP(writer, request)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(ContainSubstring("Coopers Hall"))
	})
})
