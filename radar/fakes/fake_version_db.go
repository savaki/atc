// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/radar"
)

type FakeVersionDB struct {
	SaveResourceVersionsStub        func(atc.ResourceConfig, []atc.Version) error
	saveResourceVersionsMutex       sync.RWMutex
	saveResourceVersionsArgsForCall []struct {
		arg1 atc.ResourceConfig
		arg2 []atc.Version
	}
	saveResourceVersionsReturns struct {
		result1 error
	}
	GetLatestVersionedResourceStub        func(string) (db.SavedVersionedResource, error)
	getLatestVersionedResourceMutex       sync.RWMutex
	getLatestVersionedResourceArgsForCall []struct {
		arg1 string
	}
	getLatestVersionedResourceReturns struct {
		result1 db.SavedVersionedResource
		result2 error
	}
}

func (fake *FakeVersionDB) SaveResourceVersions(arg1 atc.ResourceConfig, arg2 []atc.Version) error {
	fake.saveResourceVersionsMutex.Lock()
	fake.saveResourceVersionsArgsForCall = append(fake.saveResourceVersionsArgsForCall, struct {
		arg1 atc.ResourceConfig
		arg2 []atc.Version
	}{arg1, arg2})
	fake.saveResourceVersionsMutex.Unlock()
	if fake.SaveResourceVersionsStub != nil {
		return fake.SaveResourceVersionsStub(arg1, arg2)
	} else {
		return fake.saveResourceVersionsReturns.result1
	}
}

func (fake *FakeVersionDB) SaveResourceVersionsCallCount() int {
	fake.saveResourceVersionsMutex.RLock()
	defer fake.saveResourceVersionsMutex.RUnlock()
	return len(fake.saveResourceVersionsArgsForCall)
}

func (fake *FakeVersionDB) SaveResourceVersionsArgsForCall(i int) (atc.ResourceConfig, []atc.Version) {
	fake.saveResourceVersionsMutex.RLock()
	defer fake.saveResourceVersionsMutex.RUnlock()
	return fake.saveResourceVersionsArgsForCall[i].arg1, fake.saveResourceVersionsArgsForCall[i].arg2
}

func (fake *FakeVersionDB) SaveResourceVersionsReturns(result1 error) {
	fake.SaveResourceVersionsStub = nil
	fake.saveResourceVersionsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeVersionDB) GetLatestVersionedResource(arg1 string) (db.SavedVersionedResource, error) {
	fake.getLatestVersionedResourceMutex.Lock()
	fake.getLatestVersionedResourceArgsForCall = append(fake.getLatestVersionedResourceArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.getLatestVersionedResourceMutex.Unlock()
	if fake.GetLatestVersionedResourceStub != nil {
		return fake.GetLatestVersionedResourceStub(arg1)
	} else {
		return fake.getLatestVersionedResourceReturns.result1, fake.getLatestVersionedResourceReturns.result2
	}
}

func (fake *FakeVersionDB) GetLatestVersionedResourceCallCount() int {
	fake.getLatestVersionedResourceMutex.RLock()
	defer fake.getLatestVersionedResourceMutex.RUnlock()
	return len(fake.getLatestVersionedResourceArgsForCall)
}

func (fake *FakeVersionDB) GetLatestVersionedResourceArgsForCall(i int) string {
	fake.getLatestVersionedResourceMutex.RLock()
	defer fake.getLatestVersionedResourceMutex.RUnlock()
	return fake.getLatestVersionedResourceArgsForCall[i].arg1
}

func (fake *FakeVersionDB) GetLatestVersionedResourceReturns(result1 db.SavedVersionedResource, result2 error) {
	fake.GetLatestVersionedResourceStub = nil
	fake.getLatestVersionedResourceReturns = struct {
		result1 db.SavedVersionedResource
		result2 error
	}{result1, result2}
}

var _ radar.VersionDB = new(FakeVersionDB)
