// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/concourse/atc/worker"
)

type FakeClient struct {
	CreateContainerStub        func(worker.Identifier, worker.ContainerSpec) (worker.Container, error)
	createContainerMutex       sync.RWMutex
	createContainerArgsForCall []struct {
		arg1 worker.Identifier
		arg2 worker.ContainerSpec
	}
	createContainerReturns struct {
		result1 worker.Container
		result2 error
	}
	LookupContainerStub        func(worker.Identifier) (worker.Container, error)
	lookupContainerMutex       sync.RWMutex
	lookupContainerArgsForCall []struct {
		arg1 worker.Identifier
	}
	lookupContainerReturns struct {
		result1 worker.Container
		result2 error
	}
}

func (fake *FakeClient) CreateContainer(arg1 worker.Identifier, arg2 worker.ContainerSpec) (worker.Container, error) {
	fake.createContainerMutex.Lock()
	fake.createContainerArgsForCall = append(fake.createContainerArgsForCall, struct {
		arg1 worker.Identifier
		arg2 worker.ContainerSpec
	}{arg1, arg2})
	fake.createContainerMutex.Unlock()
	if fake.CreateContainerStub != nil {
		return fake.CreateContainerStub(arg1, arg2)
	} else {
		return fake.createContainerReturns.result1, fake.createContainerReturns.result2
	}
}

func (fake *FakeClient) CreateContainerCallCount() int {
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	return len(fake.createContainerArgsForCall)
}

func (fake *FakeClient) CreateContainerArgsForCall(i int) (worker.Identifier, worker.ContainerSpec) {
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	return fake.createContainerArgsForCall[i].arg1, fake.createContainerArgsForCall[i].arg2
}

func (fake *FakeClient) CreateContainerReturns(result1 worker.Container, result2 error) {
	fake.CreateContainerStub = nil
	fake.createContainerReturns = struct {
		result1 worker.Container
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) LookupContainer(arg1 worker.Identifier) (worker.Container, error) {
	fake.lookupContainerMutex.Lock()
	fake.lookupContainerArgsForCall = append(fake.lookupContainerArgsForCall, struct {
		arg1 worker.Identifier
	}{arg1})
	fake.lookupContainerMutex.Unlock()
	if fake.LookupContainerStub != nil {
		return fake.LookupContainerStub(arg1)
	} else {
		return fake.lookupContainerReturns.result1, fake.lookupContainerReturns.result2
	}
}

func (fake *FakeClient) LookupContainerCallCount() int {
	fake.lookupContainerMutex.RLock()
	defer fake.lookupContainerMutex.RUnlock()
	return len(fake.lookupContainerArgsForCall)
}

func (fake *FakeClient) LookupContainerArgsForCall(i int) worker.Identifier {
	fake.lookupContainerMutex.RLock()
	defer fake.lookupContainerMutex.RUnlock()
	return fake.lookupContainerArgsForCall[i].arg1
}

func (fake *FakeClient) LookupContainerReturns(result1 worker.Container, result2 error) {
	fake.LookupContainerStub = nil
	fake.lookupContainerReturns = struct {
		result1 worker.Container
		result2 error
	}{result1, result2}
}

var _ worker.Client = new(FakeClient)
