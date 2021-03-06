// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/concourse/atc/resource"
)

type FakeArtifactSource struct {
	StreamToStub        func(resource.ArtifactDestination) error
	streamToMutex       sync.RWMutex
	streamToArgsForCall []struct {
		arg1 resource.ArtifactDestination
	}
	streamToReturns struct {
		result1 error
	}
}

func (fake *FakeArtifactSource) StreamTo(arg1 resource.ArtifactDestination) error {
	fake.streamToMutex.Lock()
	fake.streamToArgsForCall = append(fake.streamToArgsForCall, struct {
		arg1 resource.ArtifactDestination
	}{arg1})
	fake.streamToMutex.Unlock()
	if fake.StreamToStub != nil {
		return fake.StreamToStub(arg1)
	} else {
		return fake.streamToReturns.result1
	}
}

func (fake *FakeArtifactSource) StreamToCallCount() int {
	fake.streamToMutex.RLock()
	defer fake.streamToMutex.RUnlock()
	return len(fake.streamToArgsForCall)
}

func (fake *FakeArtifactSource) StreamToArgsForCall(i int) resource.ArtifactDestination {
	fake.streamToMutex.RLock()
	defer fake.streamToMutex.RUnlock()
	return fake.streamToArgsForCall[i].arg1
}

func (fake *FakeArtifactSource) StreamToReturns(result1 error) {
	fake.StreamToStub = nil
	fake.streamToReturns = struct {
		result1 error
	}{result1}
}

var _ resource.ArtifactSource = new(FakeArtifactSource)
