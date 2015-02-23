package db_test

import (
	"fmt"
	"time"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/event"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type dbSharedBehaviorInput struct {
	db.DB
}

func dbSharedBehavior(database *dbSharedBehaviorInput) func() {
	return func() {
		It("initially reports zero builds for a job", func() {
			builds, err := database.GetAllJobBuilds("some-job")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(builds).Should(BeEmpty())
		})

		It("initially has no current build for a job", func() {
			_, err := database.GetCurrentBuild("some-job")
			Ω(err).Should(Equal(db.ErrNoBuild))
		})

		It("initially has no pending build for a job", func() {
			_, _, err := database.GetNextPendingBuild("some-job")
			Ω(err).Should(Equal(db.ErrNoBuild))
		})

		Context("when a build is created for a job", func() {
			var build1 db.Build

			BeforeEach(func() {
				var err error

				build1, err = database.CreateJobBuild("some-job")
				Ω(err).ShouldNot(HaveOccurred())

				Ω(build1.ID).ShouldNot(BeZero())
				Ω(build1.JobName).Should(Equal("some-job"))
				Ω(build1.Name).Should(Equal("1"))
				Ω(build1.Status).Should(Equal(db.StatusPending))
			})

			It("can be read back as the same object", func() {
				gotBuild, err := database.GetBuild(build1.ID)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(gotBuild).Should(Equal(build1))
			})

			It("becomes the current build", func() {
				currentBuild, err := database.GetCurrentBuild("some-job")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(currentBuild).Should(Equal(build1))
			})

			It("becomes the next pending build", func() {
				nextPending, _, err := database.GetNextPendingBuild("some-job")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(nextPending).Should(Equal(build1))
			})

			It("is not reported as a started build", func() {
				Ω(database.GetAllStartedBuilds()).Should(BeEmpty())
			})

			It("is returned in the job's builds", func() {
				Ω(database.GetAllJobBuilds("some-job")).Should(ConsistOf([]db.Build{build1}))
			})

			It("is returned in the set of all builds", func() {
				Ω(database.GetAllBuilds()).Should(Equal([]db.Build{build1}))
			})

			Describe("aborting", func() {
				It("notifies listeners", func() {
					notifier, err := database.AbortNotifier(build1.ID)
					Ω(err).ShouldNot(HaveOccurred())

					defer notifier.Close()

					Consistently(notifier.Notify()).ShouldNot(Receive())

					err = database.AbortBuild(build1.ID)
					Ω(err).ShouldNot(HaveOccurred())

					Eventually(notifier.Notify(), 5).Should(Receive())
				})

				It("updates the build's status", func() {
					err := database.AbortBuild(build1.ID)
					Ω(err).ShouldNot(HaveOccurred())

					build, err := database.GetBuild(build1.ID)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(build.Status).Should(Equal(db.StatusAborted))
				})

				It("immediately notifies new listeners", func() {
					err := database.AbortBuild(build1.ID)
					Ω(err).ShouldNot(HaveOccurred())

					notifier, err := database.AbortNotifier(build1.ID)
					Ω(err).ShouldNot(HaveOccurred())

					Eventually(notifier.Notify()).Should(Receive())
				})
			})

			Context("when scheduled", func() {
				BeforeEach(func() {
					scheduled, err := database.ScheduleBuild(build1.ID, false)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(scheduled).Should(BeTrue())
				})

				It("remains the current build", func() {
					currentBuild, err := database.GetCurrentBuild("some-job")
					Ω(err).ShouldNot(HaveOccurred())
					Ω(currentBuild).Should(Equal(build1))
				})

				It("remains the next pending build", func() {
					nextPending, _, err := database.GetNextPendingBuild("some-job")
					Ω(err).ShouldNot(HaveOccurred())
					Ω(nextPending).Should(Equal(build1))
				})
			})

			Context("when started", func() {
				BeforeEach(func() {
					started, err := database.StartBuild(build1.ID, "some-engine", "some-metadata")
					Ω(err).ShouldNot(HaveOccurred())
					Ω(started).Should(BeTrue())
				})

				It("saves the updated status, and the engine and engine metadata", func() {
					currentBuild, err := database.GetCurrentBuild("some-job")
					Ω(err).ShouldNot(HaveOccurred())
					Ω(currentBuild.Status).Should(Equal(db.StatusStarted))
					Ω(currentBuild.Engine).Should(Equal("some-engine"))
					Ω(currentBuild.EngineMetadata).Should(Equal("some-metadata"))
				})

				It("is not reported as a started build", func() {
					startedBuilds, err := database.GetAllStartedBuilds()
					Ω(err).ShouldNot(HaveOccurred())

					ids := make([]int, len(startedBuilds))
					for i, b := range startedBuilds {
						ids[i] = b.ID
					}

					Ω(ids).Should(ConsistOf(build1.ID))
				})

				It("can have its engine metadata saved", func() {
					err := database.SaveBuildEngineMetadata(build1.ID, "some-updated-metadata")
					Ω(err).ShouldNot(HaveOccurred())

					currentBuild, err := database.GetCurrentBuild("some-job")
					Ω(err).ShouldNot(HaveOccurred())
					Ω(currentBuild.EngineMetadata).Should(Equal("some-updated-metadata"))
				})

				It("saves the build's start time", func() {
					currentBuild, err := database.GetCurrentBuild("some-job")
					Ω(err).ShouldNot(HaveOccurred())
					Ω(currentBuild.StartTime.Unix()).Should(BeNumerically("~", time.Now().Unix(), 3))
				})
			})

			Context("when the build finishes", func() {
				BeforeEach(func() {
					err := database.FinishBuild(build1.ID, db.StatusSucceeded)
					Ω(err).ShouldNot(HaveOccurred())
				})

				It("sets the build's status and end time", func() {
					currentBuild, err := database.GetCurrentBuild("some-job")
					Ω(err).ShouldNot(HaveOccurred())
					Ω(currentBuild.Status).Should(Equal(db.StatusSucceeded))
					Ω(currentBuild.EndTime.Unix()).Should(BeNumerically("~", time.Now().Unix(), 3))
				})
			})

			Context("and another is created for the same job", func() {
				var build2 db.Build

				BeforeEach(func() {
					var err error
					build2, err = database.CreateJobBuild("some-job")
					Ω(err).ShouldNot(HaveOccurred())

					Ω(build2.ID).ShouldNot(BeZero())
					Ω(build2.ID).ShouldNot(Equal(build1.ID))
					Ω(build2.Name).Should(Equal("2"))
					Ω(build2.Status).Should(Equal(db.StatusPending))
				})

				It("can also be read back as the same object", func() {
					gotBuild, err := database.GetBuild(build2.ID)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(gotBuild).Should(Equal(build2))
				})

				It("is returned in the job's builds, before the rest", func() {
					Ω(database.GetAllJobBuilds("some-job")).Should(Equal([]db.Build{
						build2,
						build1,
					}))
				})

				It("is returned in the set of all builds, before the rest", func() {
					Ω(database.GetAllBuilds()).Should(Equal([]db.Build{build2, build1}))
				})

				Describe("the first build", func() {
					It("remains the next pending build", func() {
						nextPending, _, err := database.GetNextPendingBuild("some-job")
						Ω(err).ShouldNot(HaveOccurred())
						Ω(nextPending).Should(Equal(build1))
					})

					It("remains the current build", func() {
						currentBuild, err := database.GetCurrentBuild("some-job")
						Ω(err).ShouldNot(HaveOccurred())
						Ω(currentBuild).Should(Equal(build1))
					})
				})
			})

			Context("and another is created for a different job", func() {
				var otherJobBuild db.Build

				BeforeEach(func() {
					var err error

					otherJobBuild, err = database.CreateJobBuild("some-other-job")
					Ω(err).ShouldNot(HaveOccurred())

					Ω(otherJobBuild.ID).ShouldNot(BeZero())
					Ω(otherJobBuild.Name).Should(Equal("1"))
					Ω(otherJobBuild.Status).Should(Equal(db.StatusPending))
				})

				It("shows up in its job's builds", func() {
					Ω(database.GetAllJobBuilds("some-other-job")).Should(Equal([]db.Build{otherJobBuild}))
				})

				It("does not show up in the first build's job's builds", func() {
					Ω(database.GetAllJobBuilds("some-job")).Should(Equal([]db.Build{build1}))
				})

				It("is returned in the set of all builds, before the rest", func() {
					Ω(database.GetAllBuilds()).Should(Equal([]db.Build{otherJobBuild, build1}))
				})
			})
		})

		It("saves and propagates events correctly", func() {
			build, err := database.CreateJobBuild("some-job")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(build.Name).Should(Equal("1"))

			By("allowing you to subscribe when no events have yet occurred")
			events, err := database.GetBuildEvents(build.ID, 0)
			Ω(err).ShouldNot(HaveOccurred())

			defer events.Close()

			By("saving them in order")
			err = database.SaveBuildEvent(build.ID, event.Log{
				Payload: "some ",
			})
			Ω(err).ShouldNot(HaveOccurred())

			Ω(events.Next()).Should(Equal(event.Log{
				Payload: "some ",
			}))

			err = database.SaveBuildEvent(build.ID, event.Log{
				Payload: "log",
			})
			Ω(err).ShouldNot(HaveOccurred())

			Ω(events.Next()).Should(Equal(event.Log{
				Payload: "log",
			}))

			By("allowing you to subscribe from an offset")
			eventsFrom1, err := database.GetBuildEvents(build.ID, 1)
			Ω(err).ShouldNot(HaveOccurred())

			defer eventsFrom1.Close()

			Ω(eventsFrom1.Next()).Should(Equal(event.Log{
				Payload: "log",
			}))

			By("notifying those waiting on events as soon as they're saved")
			nextEvent := make(chan atc.Event)
			nextErr := make(chan error)

			go func() {
				event, err := events.Next()
				if err != nil {
					nextErr <- err
				} else {
					nextEvent <- event
				}
			}()

			Consistently(nextEvent).ShouldNot(Receive())
			Consistently(nextErr).ShouldNot(Receive())

			err = database.SaveBuildEvent(build.ID, event.Log{
				Payload: "log 2",
			})
			Ω(err).ShouldNot(HaveOccurred())

			Eventually(nextEvent).Should(Receive(Equal(event.Log{
				Payload: "log 2",
			})))

			By("returning ErrBuildEventStreamClosed for Next calls after Close")
			events3, err := database.GetBuildEvents(build.ID, 0)
			Ω(err).ShouldNot(HaveOccurred())

			events3.Close()

			_, err = events3.Next()
			Ω(err).Should(Equal(db.ErrBuildEventStreamClosed))
		})

		It("saves and emits status events", func() {
			build, err := database.CreateJobBuild("some-job")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(build.Name).Should(Equal("1"))

			By("allowing you to subscribe when no events have yet occurred")
			events, err := database.GetBuildEvents(build.ID, 0)
			Ω(err).ShouldNot(HaveOccurred())

			defer events.Close()

			By("emitting a status event when started")
			started, err := database.StartBuild(build.ID, "engine", "metadata")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(started).Should(BeTrue())

			startedBuild, err := database.GetBuild(build.ID)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(events.Next()).Should(Equal(event.Status{
				Status: atc.StatusStarted,
				Time:   startedBuild.StartTime.Unix(),
			}))

			By("emitting a status event when finished")
			err = database.FinishBuild(build.ID, db.StatusSucceeded)
			Ω(err).ShouldNot(HaveOccurred())

			finishedBuild, err := database.GetBuild(build.ID)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(events.Next()).Should(Equal(event.Status{
				Status: atc.StatusSucceeded,
				Time:   finishedBuild.EndTime.Unix(),
			}))

			By("ending the stream when finished")
			_, err = events.Next()
			Ω(err).Should(Equal(db.ErrEndOfBuildEventStream))
		})

		It("can keep track of workers", func() {
			Ω(database.Workers()).Should(BeEmpty())

			infoA := db.WorkerInfo{
				Addr:             "1.2.3.4:7777",
				ActiveContainers: 42,
			}

			infoB := db.WorkerInfo{
				Addr:             "1.2.3.4:8888",
				ActiveContainers: 42,
			}

			By("persisting workers with no TTLs")
			err := database.SaveWorker(infoA, 0)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(database.Workers()).Should(ConsistOf(infoA))

			By("being idempotent")
			err = database.SaveWorker(infoA, 0)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(database.Workers()).Should(ConsistOf(infoA))

			By("expiring TTLs")
			ttl := 1 * time.Second

			err = database.SaveWorker(infoB, ttl)
			Ω(err).ShouldNot(HaveOccurred())

			Consistently(database.Workers, ttl/2).Should(ConsistOf(infoA, infoB))
			Eventually(database.Workers, 2*ttl).Should(ConsistOf(infoA))

			By("overwriting TTLs")
			err = database.SaveWorker(infoA, ttl)
			Ω(err).ShouldNot(HaveOccurred())

			Consistently(database.Workers, ttl/2).Should(ConsistOf(infoA))
			Eventually(database.Workers, 2*ttl).Should(BeEmpty())
		})

		It("can create one-off builds with increasing names", func() {
			oneOff, err := database.CreateOneOffBuild()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(oneOff.ID).ShouldNot(BeZero())
			Ω(oneOff.JobName).Should(BeZero())
			Ω(oneOff.Name).Should(Equal("1"))
			Ω(oneOff.Status).Should(Equal(db.StatusPending))

			oneOffGot, err := database.GetBuild(oneOff.ID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(oneOffGot).Should(Equal(oneOff))

			jobBuild, err := database.CreateJobBuild("some-other-job")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(jobBuild.Name).Should(Equal("1"))

			nextOneOff, err := database.CreateOneOffBuild()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(nextOneOff.ID).ShouldNot(BeZero())
			Ω(nextOneOff.ID).ShouldNot(Equal(oneOff.ID))
			Ω(nextOneOff.JobName).Should(BeZero())
			Ω(nextOneOff.Name).Should(Equal("2"))
			Ω(nextOneOff.Status).Should(Equal(db.StatusPending))

			allBuilds, err := database.GetAllBuilds()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(allBuilds).Should(Equal([]db.Build{nextOneOff, jobBuild, oneOff}))
		})

		It("can report a job's latest running and finished builds", func() {
			finished, next, err := database.GetJobFinishedAndNextBuild("some-job")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(next).Should(BeNil())
			Ω(finished).Should(BeNil())

			finishedBuild, err := database.CreateJobBuild("some-job")
			Ω(err).ShouldNot(HaveOccurred())

			err = database.FinishBuild(finishedBuild.ID, db.StatusSucceeded)
			Ω(err).ShouldNot(HaveOccurred())

			finished, next, err = database.GetJobFinishedAndNextBuild("some-job")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(next).Should(BeNil())
			Ω(finished.ID).Should(Equal(finishedBuild.ID))

			nextBuild, err := database.CreateJobBuild("some-job")
			Ω(err).ShouldNot(HaveOccurred())

			started, err := database.StartBuild(nextBuild.ID, "some-engine", "meta")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(started).Should(BeTrue())

			finished, next, err = database.GetJobFinishedAndNextBuild("some-job")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(next.ID).Should(Equal(nextBuild.ID))
			Ω(finished.ID).Should(Equal(finishedBuild.ID))

			anotherRunningBuild, err := database.CreateJobBuild("some-job")
			Ω(err).ShouldNot(HaveOccurred())

			finished, next, err = database.GetJobFinishedAndNextBuild("some-job")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(next.ID).Should(Equal(nextBuild.ID)) // not anotherRunningBuild
			Ω(finished.ID).Should(Equal(finishedBuild.ID))

			started, err = database.StartBuild(anotherRunningBuild.ID, "some-engine", "meta")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(started).Should(BeTrue())

			finished, next, err = database.GetJobFinishedAndNextBuild("some-job")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(next.ID).Should(Equal(nextBuild.ID)) // not anotherRunningBuild
			Ω(finished.ID).Should(Equal(finishedBuild.ID))

			err = database.FinishBuild(nextBuild.ID, db.StatusSucceeded)
			Ω(err).ShouldNot(HaveOccurred())

			finished, next, err = database.GetJobFinishedAndNextBuild("some-job")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(next.ID).Should(Equal(anotherRunningBuild.ID))
			Ω(finished.ID).Should(Equal(nextBuild.ID))
		})

		Describe("locking", func() {
			It("can be done generically with a unique name", func() {
				lock, err := database.AcquireWriteLock([]db.NamedLock{db.ResourceCheckingLock("a-name")})
				Ω(err).ShouldNot(HaveOccurred())

				secondLockCh := make(chan db.Lock, 1)

				go func() {
					defer GinkgoRecover()

					secondLock, err := database.AcquireWriteLock([]db.NamedLock{db.ResourceCheckingLock("a-name")})
					Ω(err).ShouldNot(HaveOccurred())

					secondLockCh <- secondLock
				}()

				Consistently(secondLockCh).ShouldNot(Receive())

				err = lock.Release()
				Ω(err).ShouldNot(HaveOccurred())

				var secondLock db.Lock
				Eventually(secondLockCh).Should(Receive(&secondLock))

				err = secondLock.Release()
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("can be done without waiting", func() {
				lock, err := database.AcquireWriteLockImmediately([]db.NamedLock{db.ResourceCheckingLock("a-name")})
				Ω(err).ShouldNot(HaveOccurred())

				secondLock, err := database.AcquireWriteLockImmediately([]db.NamedLock{db.ResourceCheckingLock("a-name")})
				Ω(err).Should(HaveOccurred())
				Ω(secondLock).Should(BeNil())

				err = lock.Release()
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("does not let anyone write if someone is reading", func() {
				lock, err := database.AcquireReadLock([]db.NamedLock{db.ResourceCheckingLock("a-name")})
				Ω(err).ShouldNot(HaveOccurred())

				secondLockCh := make(chan db.Lock, 1)

				go func() {
					defer GinkgoRecover()

					secondLock, err := database.AcquireWriteLock([]db.NamedLock{db.ResourceCheckingLock("a-name")})
					Ω(err).ShouldNot(HaveOccurred())

					secondLockCh <- secondLock
				}()

				Consistently(secondLockCh).ShouldNot(Receive())

				err = lock.Release()
				Ω(err).ShouldNot(HaveOccurred())

				var secondLock db.Lock
				Eventually(secondLockCh).Should(Receive(&secondLock))

				err = secondLock.Release()
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("does not let anyone read if someone is writing", func() {
				lock, err := database.AcquireWriteLock([]db.NamedLock{db.ResourceCheckingLock("a-name")})
				Ω(err).ShouldNot(HaveOccurred())

				secondLockCh := make(chan db.Lock, 1)

				go func() {
					defer GinkgoRecover()

					secondLock, err := database.AcquireReadLock([]db.NamedLock{db.ResourceCheckingLock("a-name")})
					Ω(err).ShouldNot(HaveOccurred())

					secondLockCh <- secondLock
				}()

				Consistently(secondLockCh).ShouldNot(Receive())

				err = lock.Release()
				Ω(err).ShouldNot(HaveOccurred())

				var secondLock db.Lock
				Eventually(secondLockCh).Should(Receive(&secondLock))

				err = secondLock.Release()
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("lets many reads simultaneously", func() {
				lock, err := database.AcquireReadLock([]db.NamedLock{db.ResourceCheckingLock("a-name")})
				Ω(err).ShouldNot(HaveOccurred())

				secondLockCh := make(chan db.Lock, 1)

				go func() {
					defer GinkgoRecover()

					secondLock, err := database.AcquireReadLock([]db.NamedLock{db.ResourceCheckingLock("a-name")})
					Ω(err).ShouldNot(HaveOccurred())

					secondLockCh <- secondLock
				}()

				var secondLock db.Lock
				Eventually(secondLockCh).Should(Receive(&secondLock))

				err = secondLock.Release()
				Ω(err).ShouldNot(HaveOccurred())

				err = lock.Release()
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("can be done multiple times if using different locks", func() {
				lock, err := database.AcquireWriteLock([]db.NamedLock{db.ResourceCheckingLock("name-1")})
				Ω(err).ShouldNot(HaveOccurred())

				var secondLock db.Lock
				secondLockCh := make(chan db.Lock, 1)

				go func() {
					defer GinkgoRecover()

					secondLock, err := database.AcquireWriteLock([]db.NamedLock{db.ResourceCheckingLock("name-2")})
					Ω(err).ShouldNot(HaveOccurred())

					secondLockCh <- secondLock
				}()

				Eventually(secondLockCh).Should(Receive(&secondLock))

				err = lock.Release()
				Ω(err).ShouldNot(HaveOccurred())

				err = secondLock.Release()
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("can be done for multiple locks at a time", func() {
				lock, err := database.AcquireWriteLock([]db.NamedLock{db.ResourceCheckingLock("name-1"), db.ResourceCheckingLock("name-2")})
				Ω(err).ShouldNot(HaveOccurred())

				secondLockCh := make(chan db.Lock, 1)

				go func() {
					defer GinkgoRecover()

					secondLock, err := database.AcquireWriteLock([]db.NamedLock{db.ResourceCheckingLock("name-1")})
					Ω(err).ShouldNot(HaveOccurred())

					secondLockCh <- secondLock
				}()

				Consistently(secondLockCh).ShouldNot(Receive())

				thirdLockCh := make(chan db.Lock, 1)

				go func() {
					defer GinkgoRecover()

					thirdLock, err := database.AcquireWriteLock([]db.NamedLock{db.ResourceCheckingLock("name-2")})
					Ω(err).ShouldNot(HaveOccurred())

					thirdLockCh <- thirdLock
				}()

				Consistently(thirdLockCh).ShouldNot(Receive())

				err = lock.Release()
				Ω(err).ShouldNot(HaveOccurred())

				var secondLock db.Lock
				Eventually(secondLockCh).Should(Receive(&secondLock))

				var thirdLock db.Lock
				Eventually(thirdLockCh).Should(Receive(&thirdLock))

				err = secondLock.Release()
				Ω(err).ShouldNot(HaveOccurred())

				err = thirdLock.Release()
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("cleans up after releasing", func() {
				lock, err := database.AcquireWriteLock([]db.NamedLock{db.ResourceCheckingLock("a-name")})
				Ω(err).ShouldNot(HaveOccurred())
				Ω(database.ListLocks()).Should(ContainElement(db.ResourceCheckingLock("a-name").Name()))

				secondLockCh := make(chan db.Lock, 1)

				go func() {
					defer GinkgoRecover()

					secondLock, err := database.AcquireWriteLock([]db.NamedLock{db.ResourceCheckingLock("a-name")})
					Ω(err).ShouldNot(HaveOccurred())

					secondLockCh <- secondLock
				}()

				Consistently(secondLockCh).ShouldNot(Receive())

				err = lock.Release()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(database.ListLocks()).Should(ContainElement(db.ResourceCheckingLock("a-name").Name()))

				var secondLock db.Lock
				Eventually(secondLockCh).Should(Receive(&secondLock))

				err = secondLock.Release()
				Ω(err).ShouldNot(HaveOccurred())

				Ω(database.ListLocks()).Should(BeEmpty())
			})
		})

		Describe("saving build inputs", func() {
			buildMetadata := []db.MetadataField{
				{
					Name:  "meta1",
					Value: "value1",
				},
				{
					Name:  "meta2",
					Value: "value2",
				},
			}

			vr1 := db.VersionedResource{
				Resource: "some-resource",
				Type:     "some-type",
				Source:   db.Source{"some": "source"},
				Version:  db.Version{"ver": "1"},
				Metadata: buildMetadata,
			}

			vr2 := db.VersionedResource{
				Resource: "some-other-resource",
				Type:     "some-type",
				Source:   db.Source{"some": "other-source"},
				Version:  db.Version{"ver": "2"},
			}

			It("saves build's inputs and outputs as versioned resources", func() {
				build, err := database.CreateJobBuild("some-job")
				Ω(err).ShouldNot(HaveOccurred())

				input1 := db.BuildInput{
					Name:              "some-input",
					VersionedResource: vr1,
				}

				input2 := db.BuildInput{
					Name:              "some-other-input",
					VersionedResource: vr2,
				}

				otherInput := db.BuildInput{
					Name:              "some-random-input",
					VersionedResource: vr2,
				}

				_, err = database.SaveBuildInput(build.ID, input1)
				Ω(err).ShouldNot(HaveOccurred())

				_, err = database.GetJobBuildForInputs("some-job", []db.BuildInput{
					input1,
					input2,
				})
				Ω(err).Should(HaveOccurred())

				_, err = database.SaveBuildInput(build.ID, otherInput)
				Ω(err).ShouldNot(HaveOccurred())

				_, err = database.GetJobBuildForInputs("some-job", []db.BuildInput{
					input1,
					input2,
				})
				Ω(err).Should(HaveOccurred())

				_, err = database.SaveBuildInput(build.ID, input2)
				Ω(err).ShouldNot(HaveOccurred())

				foundBuild, err := database.GetJobBuildForInputs("some-job", []db.BuildInput{
					input1,
					input2,
				})
				Ω(err).ShouldNot(HaveOccurred())
				Ω(foundBuild).Should(Equal(build))

				_, err = database.SaveBuildOutput(build.ID, vr1)
				Ω(err).ShouldNot(HaveOccurred())

				modifiedVR2 := vr2
				modifiedVR2.Version = db.Version{"ver": "3"}

				_, err = database.SaveBuildOutput(build.ID, modifiedVR2)
				Ω(err).ShouldNot(HaveOccurred())

				_, err = database.SaveBuildOutput(build.ID, vr2)
				Ω(err).ShouldNot(HaveOccurred())

				inputs, outputs, err := database.GetBuildResources(build.ID)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(inputs).Should(ConsistOf([]db.BuildInput{
					{Name: "some-input", VersionedResource: vr1, FirstOccurrence: true},
					{Name: "some-other-input", VersionedResource: vr2, FirstOccurrence: true},
					{Name: "some-random-input", VersionedResource: vr2, FirstOccurrence: true},
				}))
				Ω(outputs).Should(ConsistOf([]db.BuildOutput{
					{VersionedResource: modifiedVR2},
				}))

				duplicateBuild, err := database.CreateJobBuild("some-job")
				Ω(err).ShouldNot(HaveOccurred())

				_, err = database.SaveBuildInput(duplicateBuild.ID, db.BuildInput{
					Name:              "other-build-input",
					VersionedResource: vr1,
				})
				Ω(err).ShouldNot(HaveOccurred())

				_, err = database.SaveBuildInput(duplicateBuild.ID, db.BuildInput{
					Name:              "other-build-other-input",
					VersionedResource: vr2,
				})
				Ω(err).ShouldNot(HaveOccurred())

				inputs, _, err = database.GetBuildResources(duplicateBuild.ID)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(inputs).Should(ConsistOf([]db.BuildInput{
					{Name: "other-build-input", VersionedResource: vr1, FirstOccurrence: false},
					{Name: "other-build-other-input", VersionedResource: vr2, FirstOccurrence: false},
				}))

				newBuildInOtherJob, err := database.CreateJobBuild("some-other-job")
				Ω(err).ShouldNot(HaveOccurred())

				_, err = database.SaveBuildInput(newBuildInOtherJob.ID, db.BuildInput{
					Name:              "other-job-input",
					VersionedResource: vr1,
				})
				Ω(err).ShouldNot(HaveOccurred())

				_, err = database.SaveBuildInput(newBuildInOtherJob.ID, db.BuildInput{
					Name:              "other-job-other-input",
					VersionedResource: vr2,
				})
				Ω(err).ShouldNot(HaveOccurred())

				inputs, _, err = database.GetBuildResources(newBuildInOtherJob.ID)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(inputs).Should(ConsistOf([]db.BuildInput{
					{Name: "other-job-input", VersionedResource: vr1, FirstOccurrence: true},
					{Name: "other-job-other-input", VersionedResource: vr2, FirstOccurrence: true},
				}))
			})

			It("updates metadata of existing inputs resources", func() {
				build, err := database.CreateJobBuild("some-job")
				Ω(err).ShouldNot(HaveOccurred())

				_, err = database.SaveBuildInput(build.ID, db.BuildInput{
					Name:              "some-input",
					VersionedResource: vr2,
				})
				Ω(err).ShouldNot(HaveOccurred())

				inputs, _, err := database.GetBuildResources(build.ID)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(inputs).Should(ConsistOf([]db.BuildInput{
					{Name: "some-input", VersionedResource: vr2, FirstOccurrence: true},
				}))

				withMetadata := vr2
				withMetadata.Metadata = buildMetadata

				_, err = database.SaveBuildInput(build.ID, db.BuildInput{
					Name:              "some-other-input",
					VersionedResource: withMetadata,
				})
				Ω(err).ShouldNot(HaveOccurred())

				inputs, _, err = database.GetBuildResources(build.ID)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(inputs).Should(ConsistOf([]db.BuildInput{
					{Name: "some-input", VersionedResource: withMetadata, FirstOccurrence: true},
					{Name: "some-other-input", VersionedResource: withMetadata, FirstOccurrence: true},
				}))

				_, err = database.SaveBuildInput(build.ID, db.BuildInput{
					Name:              "some-input",
					VersionedResource: withMetadata,
				})
				Ω(err).ShouldNot(HaveOccurred())

				inputs, _, err = database.GetBuildResources(build.ID)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(inputs).Should(ConsistOf([]db.BuildInput{
					{Name: "some-input", VersionedResource: withMetadata, FirstOccurrence: true},
					{Name: "some-other-input", VersionedResource: withMetadata, FirstOccurrence: true},
				}))
			})

			It("can be done on build creation", func() {
				inputs := []db.BuildInput{
					{Name: "first-input", VersionedResource: vr1},
					{Name: "second-input", VersionedResource: vr2},
				}

				pending, err := database.CreateJobBuildWithInputs("some-job", inputs)
				Ω(err).ShouldNot(HaveOccurred())

				foundBuild, err := database.GetJobBuildForInputs("some-job", inputs)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(foundBuild).Should(Equal(pending))

				nextPending, pendingInputs, err := database.GetNextPendingBuild("some-job")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(nextPending).Should(Equal(pending))
				Ω(pendingInputs).Should(ConsistOf([]db.BuildInput{
					{Name: "first-input", VersionedResource: vr1, FirstOccurrence: true},
					{Name: "second-input", VersionedResource: vr2, FirstOccurrence: true},
				}))
			})
		})

		Describe("saving versioned resources", func() {
			It("updates the latest versioned resource", func() {
				err := database.SaveResourceVersions(atc.ResourceConfig{
					Name:   "some-resource",
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}, []atc.Version{{"version": "1"}})
				Ω(err).ShouldNot(HaveOccurred())

				savedVR, err := database.GetLatestVersionedResource("some-resource")
				Ω(err).ShouldNot(HaveOccurred())

				Ω(savedVR.VersionedResource).Should(Equal(db.VersionedResource{
					Resource: "some-resource",
					Type:     "some-type",
					Source:   db.Source{"some": "source"},
					Version:  db.Version{"version": "1"},
				}))

				err = database.SaveResourceVersions(atc.ResourceConfig{
					Name:   "some-resource",
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}, []atc.Version{{"version": "2"}, {"version": "3"}})
				Ω(err).ShouldNot(HaveOccurred())

				savedVR, err = database.GetLatestVersionedResource("some-resource")
				Ω(err).ShouldNot(HaveOccurred())

				Ω(savedVR.VersionedResource).Should(Equal(db.VersionedResource{
					Resource: "some-resource",
					Type:     "some-type",
					Source:   db.Source{"some": "source"},
					Version:  db.Version{"version": "3"},
				}))
			})
		})

		Describe("enabling and disabling versioned resources", func() {
			resource := "some-resource"

			It("returns an error if the resource or version is bogus", func() {
				err := database.EnableVersionedResource(42)
				Ω(err).Should(HaveOccurred())

				err = database.DisableVersionedResource(42)
				Ω(err).Should(HaveOccurred())
			})

			It("does not affect explicitly fetching the latest version", func() {
				err := database.SaveResourceVersions(atc.ResourceConfig{
					Name:   "some-resource",
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}, []atc.Version{{"version": "1"}})
				Ω(err).ShouldNot(HaveOccurred())

				savedVR, err := database.GetLatestVersionedResource(resource)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(savedVR.VersionedResource).Should(Equal(db.VersionedResource{
					Resource: "some-resource",
					Type:     "some-type",
					Source:   db.Source{"some": "source"},
					Version:  db.Version{"version": "1"},
				}))

				err = database.DisableVersionedResource(savedVR.ID)
				Ω(err).ShouldNot(HaveOccurred())

				disabledVR := savedVR
				disabledVR.Enabled = false

				Ω(database.GetLatestVersionedResource(resource)).Should(Equal(disabledVR))

				err = database.EnableVersionedResource(savedVR.ID)
				Ω(err).ShouldNot(HaveOccurred())

				enabledVR := savedVR
				enabledVR.Enabled = true

				Ω(database.GetLatestVersionedResource(resource)).Should(Equal(enabledVR))
			})

			It("prevents the resource version from being a candidate for build inputs", func() {
				err := database.SaveResourceVersions(atc.ResourceConfig{
					Name:   resource,
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}, []atc.Version{{"version": "1"}})
				Ω(err).ShouldNot(HaveOccurred())

				savedVR1, err := database.GetLatestVersionedResource(resource)
				Ω(err).ShouldNot(HaveOccurred())

				err = database.SaveResourceVersions(atc.ResourceConfig{
					Name:   resource,
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}, []atc.Version{{"version": "2"}})
				Ω(err).ShouldNot(HaveOccurred())

				savedVR2, err := database.GetLatestVersionedResource(resource)
				Ω(err).ShouldNot(HaveOccurred())

				inputConfigs := []atc.JobInputConfig{
					{
						Resource: resource,
					},
				}

				Ω(database.GetLatestInputVersions(inputConfigs)).Should(Equal(db.SavedVersionedResources{
					savedVR2,
				}))

				err = database.DisableVersionedResource(savedVR2.ID)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(database.GetLatestInputVersions(inputConfigs)).Should(Equal(db.SavedVersionedResources{
					savedVR1,
				}))

				err = database.DisableVersionedResource(savedVR1.ID)
				Ω(err).ShouldNot(HaveOccurred())

				_, err = database.GetLatestInputVersions(inputConfigs)
				Ω(err).Should(HaveOccurred())

				err = database.EnableVersionedResource(savedVR1.ID)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(database.GetLatestInputVersions(inputConfigs)).Should(Equal(db.SavedVersionedResources{
					savedVR1,
				}))

				err = database.EnableVersionedResource(savedVR2.ID)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(database.GetLatestInputVersions(inputConfigs)).Should(Equal(db.SavedVersionedResources{
					savedVR2,
				}))
			})
		})

		Describe("determining the inputs for a job", func() {
			It("ensures that versions from jobs mentioned in two input's 'passed' sections came from the same builds", func() {
				j1b1, err := database.CreateJobBuild("job-1")
				Ω(err).ShouldNot(HaveOccurred())

				j2b1, err := database.CreateJobBuild("job-2")
				Ω(err).ShouldNot(HaveOccurred())

				sb1, err := database.CreateJobBuild("shared-job")
				Ω(err).ShouldNot(HaveOccurred())

				_, err = database.GetLatestInputVersions([]atc.JobInputConfig{
					{
						Resource: "resource-1",
						Passed:   []string{"shared-job", "job-1"},
					},
					{
						Resource: "resource-2",
						Passed:   []string{"shared-job", "job-2"},
					},
				})
				Ω(err).Should(Equal(db.ErrNoVersions))

				_, err = database.SaveBuildOutput(sb1.ID, db.VersionedResource{
					Resource: "resource-1",
					Type:     "some-type",
					Version:  db.Version{"v": "r1-common-to-shared-and-j1"},
				})
				Ω(err).ShouldNot(HaveOccurred())

				_, err = database.SaveBuildOutput(sb1.ID, db.VersionedResource{
					Resource: "resource-2",
					Type:     "some-type",
					Version:  db.Version{"v": "r2-common-to-shared-and-j2"},
				})
				Ω(err).ShouldNot(HaveOccurred())

				savedVR1, err := database.SaveBuildOutput(j1b1.ID, db.VersionedResource{
					Resource: "resource-1",
					Type:     "some-type",
					Version:  db.Version{"v": "r1-common-to-shared-and-j1"},
				})
				Ω(err).ShouldNot(HaveOccurred())

				savedVR2, err := database.SaveBuildOutput(j2b1.ID, db.VersionedResource{
					Resource: "resource-2",
					Type:     "some-type",
					Version:  db.Version{"v": "r2-common-to-shared-and-j2"},
				})
				Ω(err).ShouldNot(HaveOccurred())

				Ω(database.GetLatestInputVersions([]atc.JobInputConfig{
					{
						Resource: "resource-1",
						Passed:   []string{"shared-job", "job-1"},
					},
					{
						Resource: "resource-2",
						Passed:   []string{"shared-job", "job-2"},
					},
				})).Should(Equal(db.SavedVersionedResources{savedVR1, savedVR2}))

				sb2, err := database.CreateJobBuild("shared-job")
				Ω(err).ShouldNot(HaveOccurred())

				j1b2, err := database.CreateJobBuild("job-1")
				Ω(err).ShouldNot(HaveOccurred())

				j2b2, err := database.CreateJobBuild("job-2")
				Ω(err).ShouldNot(HaveOccurred())

				savedCommonVR1, err := database.SaveBuildOutput(sb2.ID, db.VersionedResource{
					Resource: "resource-1",
					Type:     "some-type",
					Version:  db.Version{"v": "new-r1-common-to-shared-and-j1"},
				})
				Ω(err).ShouldNot(HaveOccurred())

				_, err = database.SaveBuildOutput(sb2.ID, db.VersionedResource{
					Resource: "resource-2",
					Type:     "some-type",
					Version:  db.Version{"v": "new-r2-common-to-shared-and-j2"},
				})
				Ω(err).ShouldNot(HaveOccurred())

				savedCommonVR1, err = database.SaveBuildOutput(j1b2.ID, db.VersionedResource{
					Resource: "resource-1",
					Type:     "some-type",
					Version:  db.Version{"v": "new-r1-common-to-shared-and-j1"},
				})
				Ω(err).ShouldNot(HaveOccurred())

				// do NOT save resource-2 as an output of job-2

				Ω(database.GetLatestInputVersions([]atc.JobInputConfig{
					{
						Resource: "resource-1",
						Passed:   []string{"shared-job", "job-1"},
					},
					{
						Resource: "resource-2",
						Passed:   []string{"shared-job", "job-2"},
					},
				})).Should(Equal(db.SavedVersionedResources{savedVR1, savedVR2}))

				// now save the output of resource-2 job-2
				savedCommonVR2, err := database.SaveBuildOutput(j2b2.ID, db.VersionedResource{
					Resource: "resource-2",
					Type:     "some-type",
					Version:  db.Version{"v": "new-r2-common-to-shared-and-j2"},
				})
				Ω(err).ShouldNot(HaveOccurred())

				Ω(database.GetLatestInputVersions([]atc.JobInputConfig{
					{
						Resource: "resource-1",
						Passed:   []string{"shared-job", "job-1"},
					},
					{
						Resource: "resource-2",
						Passed:   []string{"shared-job", "job-2"},
					},
				})).Should(Equal(db.SavedVersionedResources{savedCommonVR1, savedCommonVR2}))

				// save newer versions; should be new latest
				for i := 0; i < 10; i++ {
					version := fmt.Sprintf("version-%d", i+1)

					savedCommonVR1, err := database.SaveBuildOutput(sb1.ID, db.VersionedResource{
						Resource: "resource-1",
						Type:     "some-type",
						Version:  db.Version{"v": version + "-r1-common-to-shared-and-j1"},
					})
					Ω(err).ShouldNot(HaveOccurred())

					savedCommonVR2, err := database.SaveBuildOutput(sb1.ID, db.VersionedResource{
						Resource: "resource-2",
						Type:     "some-type",
						Version:  db.Version{"v": version + "-r2-common-to-shared-and-j2"},
					})
					Ω(err).ShouldNot(HaveOccurred())

					savedCommonVR1, err = database.SaveBuildOutput(j1b1.ID, db.VersionedResource{
						Resource: "resource-1",
						Type:     "some-type",
						Version:  db.Version{"v": version + "-r1-common-to-shared-and-j1"},
					})
					Ω(err).ShouldNot(HaveOccurred())

					savedCommonVR2, err = database.SaveBuildOutput(j2b1.ID, db.VersionedResource{
						Resource: "resource-2",
						Type:     "some-type",
						Version:  db.Version{"v": version + "-r2-common-to-shared-and-j2"},
					})
					Ω(err).ShouldNot(HaveOccurred())

					Ω(database.GetLatestInputVersions([]atc.JobInputConfig{
						{
							Resource: "resource-1",
							Passed:   []string{"shared-job", "job-1"},
						},
						{
							Resource: "resource-2",
							Passed:   []string{"shared-job", "job-2"},
						},
					})).Should(Equal(db.SavedVersionedResources{savedCommonVR1, savedCommonVR2}))
				}
			})
		})

		Context("when the first build is created", func() {
			var firstBuild db.Build

			var job string

			BeforeEach(func() {
				var err error

				job = "some-job"

				firstBuild, err = database.CreateJobBuild(job)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(firstBuild.Name).Should(Equal("1"))
				Ω(firstBuild.Status).Should(Equal(db.StatusPending))
			})

			Context("and then aborted", func() {
				BeforeEach(func() {
					err := database.FinishBuild(firstBuild.ID, db.StatusAborted)
					Ω(err).ShouldNot(HaveOccurred())
				})

				It("changes the state to aborted", func() {
					build, err := database.GetJobBuild(job, firstBuild.Name)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(build.Status).Should(Equal(db.StatusAborted))
				})

				Describe("scheduling the build", func() {
					It("fails", func() {
						scheduled, err := database.ScheduleBuild(firstBuild.ID, false)
						Ω(err).ShouldNot(HaveOccurred())
						Ω(scheduled).Should(BeFalse())
					})
				})
			})

			Context("and then scheduled", func() {
				BeforeEach(func() {
					scheduled, err := database.ScheduleBuild(firstBuild.ID, false)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(scheduled).Should(BeTrue())
				})

				Context("and then aborted", func() {
					BeforeEach(func() {
						err := database.FinishBuild(firstBuild.ID, db.StatusAborted)
						Ω(err).ShouldNot(HaveOccurred())
					})

					It("changes the state to aborted", func() {
						build, err := database.GetJobBuild(job, firstBuild.Name)
						Ω(err).ShouldNot(HaveOccurred())
						Ω(build.Status).Should(Equal(db.StatusAborted))
					})

					Describe("starting the build", func() {
						It("fails", func() {
							started, err := database.StartBuild(firstBuild.ID, "some-engine", "some-meta")
							Ω(err).ShouldNot(HaveOccurred())
							Ω(started).Should(BeFalse())
						})
					})
				})
			})

			Describe("scheduling the build", func() {
				It("succeeds", func() {
					scheduled, err := database.ScheduleBuild(firstBuild.ID, false)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(scheduled).Should(BeTrue())
				})

				Describe("twice", func() {
					It("succeeds idempotently", func() {
						scheduled, err := database.ScheduleBuild(firstBuild.ID, false)
						Ω(err).ShouldNot(HaveOccurred())
						Ω(scheduled).Should(BeTrue())

						scheduled, err = database.ScheduleBuild(firstBuild.ID, false)
						Ω(err).ShouldNot(HaveOccurred())
						Ω(scheduled).Should(BeTrue())
					})
				})

				Context("serially", func() {
					It("succeeds", func() {
						scheduled, err := database.ScheduleBuild(firstBuild.ID, true)
						Ω(err).ShouldNot(HaveOccurred())
						Ω(scheduled).Should(BeTrue())
					})

					Describe("twice", func() {
						It("succeeds idempotently", func() {
							scheduled, err := database.ScheduleBuild(firstBuild.ID, true)
							Ω(err).ShouldNot(HaveOccurred())
							Ω(scheduled).Should(BeTrue())

							scheduled, err = database.ScheduleBuild(firstBuild.ID, true)
							Ω(err).ShouldNot(HaveOccurred())
							Ω(scheduled).Should(BeTrue())
						})
					})
				})
			})

			Context("and a second build is created", func() {
				var secondBuild db.Build

				Context("for a different job", func() {
					BeforeEach(func() {
						var err error

						secondBuild, err = database.CreateJobBuild("some-other-job")
						Ω(err).ShouldNot(HaveOccurred())
						Ω(secondBuild.Name).Should(Equal("1"))
						Ω(secondBuild.Status).Should(Equal(db.StatusPending))
					})

					Describe("scheduling the second build", func() {
						It("succeeds", func() {
							scheduled, err := database.ScheduleBuild(secondBuild.ID, false)
							Ω(err).ShouldNot(HaveOccurred())
							Ω(scheduled).Should(BeTrue())
						})

						Describe("serially", func() {
							It("succeeds", func() {
								scheduled, err := database.ScheduleBuild(secondBuild.ID, true)
								Ω(err).ShouldNot(HaveOccurred())
								Ω(scheduled).Should(BeTrue())
							})
						})
					})
				})

				Context("for the same job", func() {
					BeforeEach(func() {
						var err error

						secondBuild, err = database.CreateJobBuild(job)
						Ω(err).ShouldNot(HaveOccurred())
						Ω(secondBuild.Name).Should(Equal("2"))
						Ω(secondBuild.Status).Should(Equal(db.StatusPending))
					})

					Describe("scheduling the second build", func() {
						It("succeeds", func() {
							scheduled, err := database.ScheduleBuild(secondBuild.ID, false)
							Ω(err).ShouldNot(HaveOccurred())
							Ω(scheduled).Should(BeTrue())
						})

						Describe("serially", func() {
							It("fails", func() {
								scheduled, err := database.ScheduleBuild(secondBuild.ID, true)
								Ω(err).ShouldNot(HaveOccurred())
								Ω(scheduled).Should(BeFalse())
							})
						})
					})

					Describe("after the first build schedules", func() {
						BeforeEach(func() {
							scheduled, err := database.ScheduleBuild(firstBuild.ID, false)
							Ω(err).ShouldNot(HaveOccurred())
							Ω(scheduled).Should(BeTrue())
						})

						Context("when the second build is scheduled serially", func() {
							It("fails", func() {
								scheduled, err := database.ScheduleBuild(secondBuild.ID, true)
								Ω(err).ShouldNot(HaveOccurred())
								Ω(scheduled).Should(BeFalse())
							})
						})

						for _, s := range []db.Status{db.StatusSucceeded, db.StatusFailed, db.StatusErrored} {
							status := s

							Context("and the first build's status changes to "+string(status), func() {
								BeforeEach(func() {
									err := database.FinishBuild(firstBuild.ID, status)
									Ω(err).ShouldNot(HaveOccurred())
								})

								Context("and the second build is scheduled serially", func() {
									It("succeeds", func() {
										scheduled, err := database.ScheduleBuild(secondBuild.ID, true)
										Ω(err).ShouldNot(HaveOccurred())
										Ω(scheduled).Should(BeTrue())
									})
								})
							})
						}
					})

					Describe("after the first build is aborted", func() {
						BeforeEach(func() {
							err := database.FinishBuild(firstBuild.ID, db.StatusAborted)
							Ω(err).ShouldNot(HaveOccurred())
						})

						Context("when the second build is scheduled serially", func() {
							It("succeeds", func() {
								scheduled, err := database.ScheduleBuild(secondBuild.ID, true)
								Ω(err).ShouldNot(HaveOccurred())
								Ω(scheduled).Should(BeTrue())
							})
						})
					})

					Context("and a third build is created", func() {
						var thirdBuild db.Build

						BeforeEach(func() {
							var err error

							thirdBuild, err = database.CreateJobBuild(job)
							Ω(err).ShouldNot(HaveOccurred())
							Ω(thirdBuild.Name).Should(Equal("3"))
							Ω(thirdBuild.Status).Should(Equal(db.StatusPending))
						})

						Context("and the first build finishes", func() {
							BeforeEach(func() {
								err := database.FinishBuild(firstBuild.ID, db.StatusSucceeded)
								Ω(err).ShouldNot(HaveOccurred())
							})

							Context("and the third build is scheduled serially", func() {
								It("fails, as it would have jumped the queue", func() {
									scheduled, err := database.ScheduleBuild(thirdBuild.ID, true)
									Ω(err).ShouldNot(HaveOccurred())
									Ω(scheduled).Should(BeFalse())
								})
							})
						})

						Context("and then scheduled", func() {
							It("succeeds", func() {
								scheduled, err := database.ScheduleBuild(thirdBuild.ID, false)
								Ω(err).ShouldNot(HaveOccurred())
								Ω(scheduled).Should(BeTrue())
							})

							Describe("serially", func() {
								It("fails", func() {
									scheduled, err := database.ScheduleBuild(thirdBuild.ID, true)
									Ω(err).ShouldNot(HaveOccurred())
									Ω(scheduled).Should(BeFalse())
								})
							})
						})
					})
				})
			})
		})
	}
}
