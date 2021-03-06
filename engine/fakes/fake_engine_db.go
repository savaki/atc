// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/engine"
)

type FakeEngineDB struct {
	SaveBuildEventStub        func(buildID int, event atc.Event) error
	saveBuildEventMutex       sync.RWMutex
	saveBuildEventArgsForCall []struct {
		buildID int
		event   atc.Event
	}
	saveBuildEventReturns struct {
		result1 error
	}
	FinishBuildStub        func(buildID int, status db.Status) error
	finishBuildMutex       sync.RWMutex
	finishBuildArgsForCall []struct {
		buildID int
		status  db.Status
	}
	finishBuildReturns struct {
		result1 error
	}
	SaveBuildEngineMetadataStub        func(buildID int, metadata string) error
	saveBuildEngineMetadataMutex       sync.RWMutex
	saveBuildEngineMetadataArgsForCall []struct {
		buildID  int
		metadata string
	}
	saveBuildEngineMetadataReturns struct {
		result1 error
	}
	SaveBuildInputStub        func(buildID int, input db.BuildInput) (db.SavedVersionedResource, error)
	saveBuildInputMutex       sync.RWMutex
	saveBuildInputArgsForCall []struct {
		buildID int
		input   db.BuildInput
	}
	saveBuildInputReturns struct {
		result1 db.SavedVersionedResource
		result2 error
	}
	SaveBuildOutputStub        func(buildID int, vr db.VersionedResource) (db.SavedVersionedResource, error)
	saveBuildOutputMutex       sync.RWMutex
	saveBuildOutputArgsForCall []struct {
		buildID int
		vr      db.VersionedResource
	}
	saveBuildOutputReturns struct {
		result1 db.SavedVersionedResource
		result2 error
	}
}

func (fake *FakeEngineDB) SaveBuildEvent(buildID int, event atc.Event) error {
	fake.saveBuildEventMutex.Lock()
	fake.saveBuildEventArgsForCall = append(fake.saveBuildEventArgsForCall, struct {
		buildID int
		event   atc.Event
	}{buildID, event})
	fake.saveBuildEventMutex.Unlock()
	if fake.SaveBuildEventStub != nil {
		return fake.SaveBuildEventStub(buildID, event)
	} else {
		return fake.saveBuildEventReturns.result1
	}
}

func (fake *FakeEngineDB) SaveBuildEventCallCount() int {
	fake.saveBuildEventMutex.RLock()
	defer fake.saveBuildEventMutex.RUnlock()
	return len(fake.saveBuildEventArgsForCall)
}

func (fake *FakeEngineDB) SaveBuildEventArgsForCall(i int) (int, atc.Event) {
	fake.saveBuildEventMutex.RLock()
	defer fake.saveBuildEventMutex.RUnlock()
	return fake.saveBuildEventArgsForCall[i].buildID, fake.saveBuildEventArgsForCall[i].event
}

func (fake *FakeEngineDB) SaveBuildEventReturns(result1 error) {
	fake.SaveBuildEventStub = nil
	fake.saveBuildEventReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeEngineDB) FinishBuild(buildID int, status db.Status) error {
	fake.finishBuildMutex.Lock()
	fake.finishBuildArgsForCall = append(fake.finishBuildArgsForCall, struct {
		buildID int
		status  db.Status
	}{buildID, status})
	fake.finishBuildMutex.Unlock()
	if fake.FinishBuildStub != nil {
		return fake.FinishBuildStub(buildID, status)
	} else {
		return fake.finishBuildReturns.result1
	}
}

func (fake *FakeEngineDB) FinishBuildCallCount() int {
	fake.finishBuildMutex.RLock()
	defer fake.finishBuildMutex.RUnlock()
	return len(fake.finishBuildArgsForCall)
}

func (fake *FakeEngineDB) FinishBuildArgsForCall(i int) (int, db.Status) {
	fake.finishBuildMutex.RLock()
	defer fake.finishBuildMutex.RUnlock()
	return fake.finishBuildArgsForCall[i].buildID, fake.finishBuildArgsForCall[i].status
}

func (fake *FakeEngineDB) FinishBuildReturns(result1 error) {
	fake.FinishBuildStub = nil
	fake.finishBuildReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeEngineDB) SaveBuildEngineMetadata(buildID int, metadata string) error {
	fake.saveBuildEngineMetadataMutex.Lock()
	fake.saveBuildEngineMetadataArgsForCall = append(fake.saveBuildEngineMetadataArgsForCall, struct {
		buildID  int
		metadata string
	}{buildID, metadata})
	fake.saveBuildEngineMetadataMutex.Unlock()
	if fake.SaveBuildEngineMetadataStub != nil {
		return fake.SaveBuildEngineMetadataStub(buildID, metadata)
	} else {
		return fake.saveBuildEngineMetadataReturns.result1
	}
}

func (fake *FakeEngineDB) SaveBuildEngineMetadataCallCount() int {
	fake.saveBuildEngineMetadataMutex.RLock()
	defer fake.saveBuildEngineMetadataMutex.RUnlock()
	return len(fake.saveBuildEngineMetadataArgsForCall)
}

func (fake *FakeEngineDB) SaveBuildEngineMetadataArgsForCall(i int) (int, string) {
	fake.saveBuildEngineMetadataMutex.RLock()
	defer fake.saveBuildEngineMetadataMutex.RUnlock()
	return fake.saveBuildEngineMetadataArgsForCall[i].buildID, fake.saveBuildEngineMetadataArgsForCall[i].metadata
}

func (fake *FakeEngineDB) SaveBuildEngineMetadataReturns(result1 error) {
	fake.SaveBuildEngineMetadataStub = nil
	fake.saveBuildEngineMetadataReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeEngineDB) SaveBuildInput(buildID int, input db.BuildInput) (db.SavedVersionedResource, error) {
	fake.saveBuildInputMutex.Lock()
	fake.saveBuildInputArgsForCall = append(fake.saveBuildInputArgsForCall, struct {
		buildID int
		input   db.BuildInput
	}{buildID, input})
	fake.saveBuildInputMutex.Unlock()
	if fake.SaveBuildInputStub != nil {
		return fake.SaveBuildInputStub(buildID, input)
	} else {
		return fake.saveBuildInputReturns.result1, fake.saveBuildInputReturns.result2
	}
}

func (fake *FakeEngineDB) SaveBuildInputCallCount() int {
	fake.saveBuildInputMutex.RLock()
	defer fake.saveBuildInputMutex.RUnlock()
	return len(fake.saveBuildInputArgsForCall)
}

func (fake *FakeEngineDB) SaveBuildInputArgsForCall(i int) (int, db.BuildInput) {
	fake.saveBuildInputMutex.RLock()
	defer fake.saveBuildInputMutex.RUnlock()
	return fake.saveBuildInputArgsForCall[i].buildID, fake.saveBuildInputArgsForCall[i].input
}

func (fake *FakeEngineDB) SaveBuildInputReturns(result1 db.SavedVersionedResource, result2 error) {
	fake.SaveBuildInputStub = nil
	fake.saveBuildInputReturns = struct {
		result1 db.SavedVersionedResource
		result2 error
	}{result1, result2}
}

func (fake *FakeEngineDB) SaveBuildOutput(buildID int, vr db.VersionedResource) (db.SavedVersionedResource, error) {
	fake.saveBuildOutputMutex.Lock()
	fake.saveBuildOutputArgsForCall = append(fake.saveBuildOutputArgsForCall, struct {
		buildID int
		vr      db.VersionedResource
	}{buildID, vr})
	fake.saveBuildOutputMutex.Unlock()
	if fake.SaveBuildOutputStub != nil {
		return fake.SaveBuildOutputStub(buildID, vr)
	} else {
		return fake.saveBuildOutputReturns.result1, fake.saveBuildOutputReturns.result2
	}
}

func (fake *FakeEngineDB) SaveBuildOutputCallCount() int {
	fake.saveBuildOutputMutex.RLock()
	defer fake.saveBuildOutputMutex.RUnlock()
	return len(fake.saveBuildOutputArgsForCall)
}

func (fake *FakeEngineDB) SaveBuildOutputArgsForCall(i int) (int, db.VersionedResource) {
	fake.saveBuildOutputMutex.RLock()
	defer fake.saveBuildOutputMutex.RUnlock()
	return fake.saveBuildOutputArgsForCall[i].buildID, fake.saveBuildOutputArgsForCall[i].vr
}

func (fake *FakeEngineDB) SaveBuildOutputReturns(result1 db.SavedVersionedResource, result2 error) {
	fake.SaveBuildOutputStub = nil
	fake.saveBuildOutputReturns = struct {
		result1 db.SavedVersionedResource
		result2 error
	}{result1, result2}
}

var _ engine.EngineDB = new(FakeEngineDB)
