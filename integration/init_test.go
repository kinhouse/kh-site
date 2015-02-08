package integration

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core"

	"testing"

	"github.com/kinhouse/kh-site/fakes"
	"github.com/kinhouse/kh-site/server"
)

func TestKhSite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KhSite Suite")
}

var agoutiDriver WebDriver

var baseUrl string

const port = 5555

var _ = BeforeSuite(func() {
	var err error

	// Choose a WebDriver:

	agoutiDriver, err = PhantomJS()
	// agoutiDriver, err = Selenium()
	// agoutiDriver, err = Chrome()

	Expect(err).NotTo(HaveOccurred())
	Expect(agoutiDriver.Start()).To(Succeed())

	persist := fakes.Persist{}
	serverConfig := server.BuildServerConfig(persist)
	router := serverConfig.BuildRouter()
	go router.Run(fmt.Sprintf(":%d", port))

	baseUrl = fmt.Sprintf("http://localhost:%d", port)
	WaitToBoot(baseUrl)

})

var _ = AfterSuite(func() {
	agoutiDriver.Stop()
})

func WaitToBoot(route string) {
	fmt.Printf("Waiting for test server to boot on %s\n", route)
	timer := time.After(0 * time.Second)
	timeout := time.After(10 * time.Second)
	for {
		select {
		case <-timeout:
			panic("Failed to boot!")
		case <-timer:
			_, err := http.Get(route)
			if err == nil {
				fmt.Printf("Test server booted\n")
				return
			}
			timer = time.After(1 * time.Second)
		}
	}
}
