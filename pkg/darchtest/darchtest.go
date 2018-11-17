package darchtest

import (
	"context"
	"github.com/containerd/containerd"
	controlapi "github.com/moby/buildkit/api/services/control"
	"github.com/moby/buildkit/control"
	"github.com/moby/buildkit/frontend"
	"github.com/moby/buildkit/frontend/dockerfile/builder"
	"github.com/moby/buildkit/frontend/gateway"
	"github.com/moby/buildkit/frontend/gateway/forwarder"
	"github.com/moby/buildkit/identity"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/session/filesync"
	"github.com/moby/buildkit/session/testutil"
	"github.com/moby/buildkit/solver/bboltcachestorage"
	"github.com/moby/buildkit/worker"
	"github.com/moby/buildkit/worker/base"
	containerdworker "github.com/moby/buildkit/worker/containerd"
	"golang.org/x/sync/errgroup"
	"log"
	"strings"
)

// Test .
func Test() error {
	// id, err := base.ID("/home/pknopf/scratch")
	// if err != nil {
	// 	return err
	// }

	opt, err := containerdworker.NewWorkerOpt("/home/pknopf/scratch", "/var/run/containerd/containerd.sock", containerd.DefaultSnapshotter, "darch", make(map[string]string))
	if err != nil {
		return err
	}

	sm, err := session.NewManager()
	if err != nil {
		return err
	}

	opt.SessionManager = sm

	w, err := base.NewWorker(opt)
	if err != nil {
		return err
	}

	wc := &worker.Controller{}
	if err = wc.Add(w); err != nil {
		return err
	}

	frontends := map[string]frontend.Frontend{}
	frontends["dockerfile.v0"] = forwarder.NewGatewayForwarder(wc, builder.Build)
	frontends["gateway.v0"] = gateway.NewGatewayFrontend(wc)

	// Create the cache storage
	cacheStorage, err := bboltcachestorage.NewStore("/home/pknopf/scratch/cache.db")
	if err != nil {
		return err
	}

	controller, err := control.NewController(control.Opt{
		SessionManager:   sm,
		WorkerController: wc,
		Frontends:        frontends,
		CacheKeyStorage:  cacheStorage,
		// No cache importer/exporter
	})

	id := identity.NewID()

	// sessionManager, err := session.NewManager()
	// if err != nil {
	// 	return err
	// }

	s, err := session.NewSession(context.Background(), "darchtest", "")
	if err != nil {
		return err
	}

	syncedDirs := []filesync.SyncedDir{
		{
			Name: "context",
			Dir:  "/home/pknopf/scratch",
		},
		{
			Name: "dockerfile",
			Dir:  "/home/pknopf/scratch",
		},
	}
	// }

	// 	syncedDirs = append(syncedDirs, filesync.SyncedDir{Name: name, Dir: d})
	// }
	s.Allow(filesync.NewFSSyncProvider(syncedDirs))
	s.Allow(authprovider.NewDockerAuthProvider())

	dialer := session.Dialer(testutil.TestStream(testutil.Handler(sm.HandleConn)))

	frontendAttrs := map[string]string{
		// We use the base for filename here because we already set up the local dirs which sets the path in createController.
		"filename": "Dockerfile",
		"target":   "",
	}
	eg, ctx := errgroup.WithContext(context.Background())

	//ch := make(chan *controlapi.StatusResponse)
	eg.Go(func() error {
		return s.Run(ctx, dialer)
	})

	eg.Go(func() error {
		defer s.Close()
		resp, err := controller.Solve(context.Background(),
			&controlapi.SolveRequest{
				Ref:      id,
				Session:  s.ID(),
				Exporter: "image",
				ExporterAttrs: map[string]string{
					"name": strings.Join([]string{"docker.io/library/test-image:latest"}, ","),
				},
				Frontend:      "dockerfile.v0",
				FrontendAttrs: frontendAttrs,
			})
		log.Println(resp.ExporterResponse["container.image.digest"])
		log.Println(resp.String())
		//sha256:689ee9840b989ef0937060239eb341578bf58fb1ea25c9232c02d144801980dd
		return err
	})

	//log.Println(resp)

	return eg.Wait()
}

// func solve(ctx context.Context, req *controlapi.SolveRequest, ch chan *controlapi.StatusResponse) error {
// 	defer close(ch)

// 	statusCtx, cancelStatus := context.WithCancel(context.Background())
// 	eg, ctx := errgroup.WithContext(ctx)
// 	eg.Go(func() error {
// 		defer func() { // make sure the Status ends cleanly on build errors
// 			go func() {
// 				<-time.After(3 * time.Second)
// 				cancelStatus()
// 			}()
// 		}()
// 		_, err := c.controller.Solve(ctx, req)
// 		if err != nil {
// 			return errors.Wrap(err, "failed to solve")
// 		}
// 		return nil
// 	})

// 	eg.Go(func() error {
// 		srv := &controlStatusServer{
// 			ctx: statusCtx,
// 			ch:  ch,
// 		}
// 		return c.controller.Status(&controlapi.StatusRequest{
// 			Ref: req.Ref,
// 		}, srv)
// 	})
// 	return eg.Wait()
// }
