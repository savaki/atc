package engine

import (
	"errors"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/pivotal-golang/lager"
)

var ErrBuildNotFound = errors.New("build not found")

//go:generate counterfeiter . Engine

type Engine interface {
	Name() string

	CreateBuild(db.Build, atc.Plan) (Build, error)
	LookupBuild(db.Build) (Build, error)
}

//go:generate counterfeiter . EngineDB

type EngineDB interface {
	SaveBuildEvent(buildID int, event atc.Event) error

	FinishBuild(buildID int, status db.Status) error

	SaveBuildEngineMetadata(buildID int, metadata string) error

	SaveBuildInput(buildID int, input db.BuildInput) (db.SavedVersionedResource, error)
	SaveBuildOutput(buildID int, vr db.VersionedResource) (db.SavedVersionedResource, error)
}

//go:generate counterfeiter . Build

type Build interface {
	Metadata() string

	Abort() error
	Resume(lager.Logger)
}

type Engines []Engine

func (engines Engines) Lookup(name string) (Engine, bool) {
	for _, e := range engines {
		if e.Name() == name {
			return e, true
		}
	}

	return nil, false
}
