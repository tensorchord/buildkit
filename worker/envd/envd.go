package envd

import (
	mobyworker "github.com/docker/docker/builder/builder-next/worker"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/source"
	"github.com/moby/buildkit/source/git"
	"github.com/moby/buildkit/source/http"
	"github.com/moby/buildkit/source/local"
	"github.com/sirupsen/logrus"
)

// Worker is a local worker instance with dedicated snapshotter, cache, and so on.
// TODO: s/Worker/OpWorker/g ?
type Worker struct {
	mobyworker.Worker
	version client.BuildkitVersion
}

func (w *Worker) BuildkitVersion() client.BuildkitVersion {
	return w.version
}

// NewWorker instantiates a local worker
func NewWorker(opt mobyworker.Opt, v client.BuildkitVersion) (*Worker, error) {
	sm, err := source.NewManager()
	if err != nil {
		return nil, err
	}

	cm := opt.CacheManager
	sm.Register(opt.ImageSource)

	gs, err := git.NewSource(git.Opt{
		CacheAccessor: cm,
	})
	if err == nil {
		sm.Register(gs)
	} else {
		logrus.Warnf("Could not register builder git source: %s", err)
	}

	hs, err := http.NewSource(http.Opt{
		CacheAccessor: cm,
		Transport:     opt.Transport,
	})
	if err == nil {
		sm.Register(hs)
	} else {
		logrus.Warnf("Could not register builder http source: %s", err)
	}

	ss, err := local.NewSource(local.Opt{
		CacheAccessor: cm,
	})
	if err == nil {
		sm.Register(ss)
	} else {
		logrus.Warnf("Could not register builder local source: %s", err)
	}

	return &Worker{
		version: v,
		Worker: mobyworker.Worker{
			Opt:           opt,
			SourceManager: sm,
		},
	}, nil
}
