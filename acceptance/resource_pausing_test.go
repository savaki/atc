package acceptance_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/sclevine/agouti"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	. "github.com/sclevine/agouti/matchers"

	"github.com/cloudfoundry/gunk/urljoiner"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
)

var _ = Describe("Resource Pausing", func() {
	var atcProcess ifrit.Process
	var dbListener *pq.Listener
	var pipelineDBFactory db.PipelineDBFactory
	var atcPort uint16

	BeforeEach(func() {
		atcBin, err := gexec.Build("github.com/concourse/atc/cmd/atc")
		Ω(err).ShouldNot(HaveOccurred())

		dbLogger := lagertest.NewTestLogger("test")
		postgresRunner.CreateTestDB()
		dbConn = postgresRunner.Open()
		dbListener = pq.NewListener(postgresRunner.DataSourceName(), time.Second, time.Minute, nil)
		bus := db.NewNotificationsBus(dbListener)

		sqlDB = db.NewSQL(dbLogger, dbConn, bus)
		Ω(err).ShouldNot(HaveOccurred())

		pipelineDBFactory = db.NewPipelineDBFactory(dbLogger, dbConn, bus, sqlDB)

		atcProcess, atcPort = startATC(atcBin, 1)
	})

	AfterEach(func() {
		ginkgomon.Interrupt(atcProcess)

		Ω(dbConn.Close()).Should(Succeed())
		Ω(dbListener.Close()).Should(Succeed())

		postgresRunner.DropTestDB()
	})

	Describe("pausing a resource", func() {
		var page *agouti.Page

		BeforeEach(func() {
			var err error
			page, err = agoutiDriver.NewPage()
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(page.Destroy()).To(Succeed())
		})

		homepage := func() string {
			return fmt.Sprintf("http://127.0.0.1:%d/pipelines/%s", atcPort, atc.DefaultPipelineName)
		}

		withPath := func(path string) string {
			return urljoiner.Join(homepage(), path)
		}

		Context("with a resource in the configuration", func() {
			var pipelineDB db.PipelineDB

			BeforeEach(func() {
				// job build data
				_, err := sqlDB.SaveConfig(atc.DefaultPipelineName, atc.Config{
					Jobs: atc.JobConfigs{
						{
							Name: "job-name",
							Plan: atc.PlanSequence{
								{
									Get: "resource-name",
								},
							},
						},
					},
					Resources: atc.ResourceConfigs{
						{Name: "resource-name"},
					},
				}, db.ConfigVersion(1), db.PipelineUnpaused)
				Ω(err).ShouldNot(HaveOccurred())

				pipelineDB, err = pipelineDBFactory.BuildDefault()
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("can view the resource", func() {
				// homepage -> resource detail
				Expect(page.Navigate(homepage())).To(Succeed())
				Eventually(page.FindByLink("resource-name")).Should(BeFound())
				Expect(page.FindByLink("resource-name").Click()).To(Succeed())

				// resource detail -> paused resource detail
				Expect(page).Should(HaveURL(withPath("/resources/resource-name")))
				Expect(page.Find("h1")).To(HaveText("resource-name"))

				Authenticate(page, "admin", "password")

				Expect(page.Find(".js-resource .js-pauseUnpause").Click()).To(Succeed())
				Eventually(page.Find(".header i.fa-play")).Should(BeFound())

				page.Refresh()

				Eventually(page.Find(".header i.fa-play")).Should(BeFound())

				resource, err := pipelineDB.GetResource("resource-name")
				Ω(err).ShouldNot(HaveOccurred())

				err = pipelineDB.SetResourceCheckError(resource, errors.New("failed to foo the bar"))
				Ω(err).ShouldNot(HaveOccurred())

				page.Refresh()

				Eventually(page.Find(".header h3")).Should(HaveText("checking failed"))
				Eventually(page.Find(".build-step .step-body")).Should(HaveText("failed to foo the bar"))
			})
		})
	})
})
