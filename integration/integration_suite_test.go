package integration

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core"

	"testing"
	"time"

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

var localPersist *fakes.Persist

var _ = BeforeSuite(func() {
	var err error

	// Choose a WebDriver:

	//agoutiDriver, err = PhantomJS()
	// agoutiDriver, err = Selenium()
	agoutiDriver, err = Chrome()

	Expect(err).NotTo(HaveOccurred())
	Expect(agoutiDriver.Start()).To(Succeed())

	localPersist = &fakes.Persist{}
	s := server.BuildServer(localPersist)
	go s.Run(fmt.Sprintf(":%d", port))

	baseUrl = fmt.Sprintf("http://localhost:%d", port)
	waitToBoot(baseUrl)

})

var _ = AfterSuite(func() {
	agoutiDriver.Stop()
})

func waitToBoot(route string) {
	fmt.Printf("Waiting for server to boot on %s\n", route)
	timer := time.After(0 * time.Second)
	timeout := time.After(10 * time.Second)
	for {
		select {
		case <-timeout:
			panic("Failed to boot!")
		case <-timer:
			resp, err := http.Get(route)
			defer resp.Body.Close()
			if err == nil {
				fmt.Printf("Test server booted\n")
				return
			}
			timer = time.After(1 * time.Second)
		}
	}
}
