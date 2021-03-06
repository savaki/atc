// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/web/getresource"
)

type FakeResourcesDB struct {
	GetPipelineNameStub        func() string
	getPipelineNameMutex       sync.RWMutex
	getPipelineNameArgsForCall []struct{}
	getPipelineNameReturns struct {
		result1 string
	}
	GetConfigStub        func() (atc.Config, db.ConfigVersion, error)
	getConfigMutex       sync.RWMutex
	getConfigArgsForCall []struct{}
	getConfigReturns struct {
		result1 atc.Config
		result2 db.ConfigVersion
		result3 error
	}
	GetResourceStub        func(string) (db.SavedResource, error)
	getResourceMutex       sync.RWMutex
	getResourceArgsForCall []struct {
		arg1 string
	}
	getResourceReturns struct {
		result1 db.SavedResource
		result2 error
	}
	GetResourceHistoryStub        func(string) ([]*db.VersionHistory, error)
	getResourceHistoryMutex       sync.RWMutex
	getResourceHistoryArgsForCall []struct {
		arg1 string
	}
	getResourceHistoryReturns struct {
		result1 []*db.VersionHistory
		result2 error
	}
}

func (fake *FakeResourcesDB) GetPipelineName() string {
	fake.getPipelineNameMutex.Lock()
	fake.getPipelineNameArgsForCall = append(fake.getPipelineNameArgsForCall, struct{}{})
	fake.getPipelineNameMutex.Unlock()
	if fake.GetPipelineNameStub != nil {
		return fake.GetPipelineNameStub()
	} else {
		return fake.getPipelineNameReturns.result1
	}
}

func (fake *FakeResourcesDB) GetPipelineNameCallCount() int {
	fake.getPipelineNameMutex.RLock()
	defer fake.getPipelineNameMutex.RUnlock()
	return len(fake.getPipelineNameArgsForCall)
}

func (fake *FakeResourcesDB) GetPipelineNameReturns(result1 string) {
	fake.GetPipelineNameStub = nil
	fake.getPipelineNameReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeResourcesDB) GetConfig() (atc.Config, db.ConfigVersion, error) {
	fake.getConfigMutex.Lock()
	fake.getConfigArgsForCall = append(fake.getConfigArgsForCall, struct{}{})
	fake.getConfigMutex.Unlock()
	if fake.GetConfigStub != nil {
		return fake.GetConfigStub()
	} else {
		return fake.getConfigReturns.result1, fake.getConfigReturns.result2, fake.getConfigReturns.result3
	}
}

func (fake *FakeResourcesDB) GetConfigCallCount() int {
	fake.getConfigMutex.RLock()
	defer fake.getConfigMutex.RUnlock()
	return len(fake.getConfigArgsForCall)
}

func (fake *FakeResourcesDB) GetConfigReturns(result1 atc.Config, result2 db.ConfigVersion, result3 error) {
	fake.GetConfigStub = nil
	fake.getConfigReturns = struct {
		result1 atc.Config
		result2 db.ConfigVersion
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeResourcesDB) GetResource(arg1 string) (db.SavedResource, error) {
	fake.getResourceMutex.Lock()
	fake.getResourceArgsForCall = append(fake.getResourceArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.getResourceMutex.Unlock()
	if fake.GetResourceStub != nil {
		return fake.GetResourceStub(arg1)
	} else {
		return fake.getResourceReturns.result1, fake.getResourceReturns.result2
	}
}

func (fake *FakeResourcesDB) GetResourceCallCount() int {
	fake.getResourceMutex.RLock()
	defer fake.getResourceMutex.RUnlock()
	return len(fake.getResourceArgsForCall)
}

func (fake *FakeResourcesDB) GetResourceArgsForCall(i int) string {
	fake.getResourceMutex.RLock()
	defer fake.getResourceMutex.RUnlock()
	return fake.getResourceArgsForCall[i].arg1
}

func (fake *FakeResourcesDB) GetResourceReturns(result1 db.SavedResource, result2 error) {
	fake.GetResourceStub = nil
	fake.getResourceReturns = struct {
		result1 db.SavedResource
		result2 error
	}{result1, result2}
}

func (fake *FakeResourcesDB) GetResourceHistory(arg1 string) ([]*db.VersionHistory, error) {
	fake.getResourceHistoryMutex.Lock()
	fake.getResourceHistoryArgsForCall = append(fake.getResourceHistoryArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.getResourceHistoryMutex.Unlock()
	if fake.GetResourceHistoryStub != nil {
		return fake.GetResourceHistoryStub(arg1)
	} else {
		return fake.getResourceHistoryReturns.result1, fake.getResourceHistoryReturns.result2
	}
}

func (fake *FakeResourcesDB) GetResourceHistoryCallCount() int {
	fake.getResourceHistoryMutex.RLock()
	defer fake.getResourceHistoryMutex.RUnlock()
	return len(fake.getResourceHistoryArgsForCall)
}

func (fake *FakeResourcesDB) GetResourceHistoryArgsForCall(i int) string {
	fake.getResourceHistoryMutex.RLock()
	defer fake.getResourceHistoryMutex.RUnlock()
	return fake.getResourceHistoryArgsForCall[i].arg1
}

func (fake *FakeResourcesDB) GetResourceHistoryReturns(result1 []*db.VersionHistory, result2 error) {
	fake.GetResourceHistoryStub = nil
	fake.getResourceHistoryReturns = struct {
		result1 []*db.VersionHistory
		result2 error
	}{result1, result2}
}

var _ getresource.ResourcesDB = new(FakeResourcesDB)
