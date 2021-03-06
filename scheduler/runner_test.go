package scheduler_test

import (
	"errors"
	"sync"
	"time"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	dbfakes "github.com/concourse/atc/db/fakes"
	. "github.com/concourse/atc/scheduler"
	"github.com/concourse/atc/scheduler/fakes"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Runner", func() {
	var (
		locker     *fakes.FakeLocker
		pipelineDB *dbfakes.FakePipelineDB
		scheduler  *fakes.FakeBuildScheduler
		noop       bool

		lock *dbfakes.FakeLock

		initialConfig atc.Config

		process ifrit.Process
	)

	BeforeEach(func() {
		locker = new(fakes.FakeLocker)
		pipelineDB = new(dbfakes.FakePipelineDB)
		scheduler = new(fakes.FakeBuildScheduler)
		noop = false

		scheduler.TryNextPendingBuildStub = func(lager.Logger, atc.JobConfig, atc.ResourceConfigs) Waiter {
			return new(sync.WaitGroup)
		}

		initialConfig = atc.Config{
			Jobs: atc.JobConfigs{
				{
					Name: "some-job",
				},
				{
					Name: "some-other-job",
				},
			},

			Resources: atc.ResourceConfigs{
				{
					Name:   "some-resource",
					Type:   "git",
					Source: atc.Source{"uri": "git://some-resource"},
				},
				{
					Name:   "some-dependant-resource",
					Type:   "git",
					Source: atc.Source{"uri": "git://some-dependant-resource"},
				},
			},
		}

		pipelineDB.ScopedNameStub = func(thing string) string {
			return "pipeline:" + thing
		}
		pipelineDB.GetConfigReturns(initialConfig, 1, nil)

		lock = new(dbfakes.FakeLock)
		locker.AcquireWriteLockImmediatelyReturns(lock, nil)
	})

	JustBeforeEach(func() {
		process = ginkgomon.Invoke(&Runner{
			Logger:    lagertest.NewTestLogger("test"),
			Locker:    locker,
			DB:        pipelineDB,
			Scheduler: scheduler,
			Noop:      noop,
			Interval:  100 * time.Millisecond,
		})
	})

	AfterEach(func() {
		ginkgomon.Interrupt(process)
	})

	It("acquires the build scheduling lock for each job", func() {
		Eventually(locker.AcquireWriteLockImmediatelyCallCount).Should(Equal(2))

		job := locker.AcquireWriteLockImmediatelyArgsForCall(0)
		Ω(job).Should(Equal([]db.NamedLock{db.JobSchedulingLock("pipeline:some-job")}))

		job = locker.AcquireWriteLockImmediatelyArgsForCall(1)
		Ω(job).Should(Equal([]db.NamedLock{db.JobSchedulingLock("pipeline:some-other-job")}))
	})

	Context("whe it can't get the lock for the first job", func() {
		BeforeEach(func() {
			locker.AcquireWriteLockImmediatelyStub = func(locks []db.NamedLock) (db.Lock, error) {
				if locker.AcquireWriteLockImmediatelyCallCount() == 1 {
					return nil, errors.New("can't aqcuire lock")
				}
				return lock, nil
			}
		})

		It("follows on to the next job", func() {
			Eventually(locker.AcquireWriteLockImmediatelyCallCount).Should(Equal(2))

			_, job, resources := scheduler.TryNextPendingBuildArgsForCall(0)
			Ω(job).Should(Equal(atc.JobConfig{Name: "some-other-job"}))
			Ω(resources).Should(Equal(initialConfig.Resources))
		})
	})

	It("schedules pending builds", func() {
		Eventually(scheduler.TryNextPendingBuildCallCount).Should(Equal(2))

		_, job, resources := scheduler.TryNextPendingBuildArgsForCall(0)
		Ω(job).Should(Equal(atc.JobConfig{Name: "some-job"}))
		Ω(resources).Should(Equal(initialConfig.Resources))

		_, job, resources = scheduler.TryNextPendingBuildArgsForCall(1)
		Ω(job).Should(Equal(atc.JobConfig{Name: "some-other-job"}))
		Ω(resources).Should(Equal(initialConfig.Resources))
	})

	It("schedules builds for new inputs", func() {
		Eventually(scheduler.BuildLatestInputsCallCount).Should(Equal(2))

		_, job, resources := scheduler.BuildLatestInputsArgsForCall(0)
		Ω(job).Should(Equal(atc.JobConfig{Name: "some-job"}))
		Ω(resources).Should(Equal(initialConfig.Resources))

		_, job, resources = scheduler.BuildLatestInputsArgsForCall(1)
		Ω(job).Should(Equal(atc.JobConfig{Name: "some-other-job"}))
		Ω(resources).Should(Equal(initialConfig.Resources))
	})

	Context("when in noop mode", func() {
		BeforeEach(func() {
			noop = true
		})

		It("does not start scheduling builds", func() {
			Consistently(scheduler.TryNextPendingBuildCallCount).Should(Equal(0))
			Consistently(scheduler.BuildLatestInputsCallCount).Should(Equal(0))
		})
	})
})
