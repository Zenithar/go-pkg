package diagnostic

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/google/gops/agent"
	"go.uber.org/zap"

	"go.zenithar.org/pkg/log"
)

// Register adds diagnostic tools to main process
func Register(ctx context.Context, conf Config, r *http.ServeMux) error {

	if conf.GOPS.Enabled {
		// Start diagnostic handler
		if conf.GOPS.RemoteURL != "" {
			log.For(ctx).Info("Starting gops agent", zap.String("url", conf.GOPS.RemoteURL))
			if err := agent.Listen(agent.Options{Addr: conf.GOPS.RemoteURL}); err != nil {
				log.For(ctx).Error("Error on starting gops agent", zap.Error(err))
			}
		} else {
			log.For(ctx).Info("Starting gops agent locally")
			if err := agent.Listen(agent.Options{}); err != nil {
				log.For(ctx).Error("Error on starting gops agent locally", zap.Error(err))
			}
		}
	}

	if conf.PProf.Enabled {
		r.HandleFunc("/diag/pprof", pprof.Index)
		r.HandleFunc("/diag/cmdline", pprof.Cmdline)
		r.HandleFunc("/diag/profile", pprof.Profile)
		r.HandleFunc("/diag/symbol", pprof.Symbol)
		r.HandleFunc("/diag/trace", pprof.Trace)
		r.Handle("/diag/goroutine", pprof.Handler("goroutine"))
		r.Handle("/diag/heap", pprof.Handler("heap"))
		r.Handle("/diag/threadcreate", pprof.Handler("threadcreate"))
		r.Handle("/diag/block", pprof.Handler("block"))
	}

	// No error
	return nil
}
