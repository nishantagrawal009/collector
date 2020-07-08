package main

import (
	"collector/config"
	"collector/log"
	"collector/middleware"
	profefe "collector/profile-collector"
	"collector/storage"
	"context"
	"expvar"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/pprof"
	"os"

	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)
func main() {
	var conf config.Config
	conf.RegisterFlags(flag.CommandLine)

	flag.Parse()

	logger, err := conf.Logger.Build()
	if err != nil {
		panic(err)
	}

	if err := run(context.Background(), logger, conf, os.Stdout); err != nil {
		logger.Error(err)
	}
}

func run(ctx context.Context, logger *log.Logger, conf config.Config, stdout io.Writer) error {
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	ctx, cancelTopMostCtx := context.WithCancel(ctx)
	defer cancelTopMostCtx()

	collector, querier, closer, err := initProfefe(logger, conf)
	if err != nil {
		return err
	}
	defer closer()

	mux := http.NewServeMux()

	profefe.SetupRoutes(mux, logger, prometheus.DefaultRegisterer, collector, querier)

	setupDebugRoutes(mux)

	// TODO(narqo) hardcoded stdout when setup logging middleware
	h := middleware.LoggingHandler(stdout, mux)
	h = middleware.RecoveryHandler(h)

	server := http.Server{
		Addr:    conf.Addr,
		Handler: h,
	}

	errc := make(chan error, 1)
	go func() {
		logger.Infow("server is running", "addr", server.Addr)
		errc <- server.ListenAndServe()
	}()


	select {
	case <-sigs:
		logger.Info("exiting")
	case <-ctx.Done():
		logger.Info("exiting", zap.Error(ctx.Err()))
	case err := <-errc:
		if err != http.ErrServerClosed {
			return fmt.Errorf("terminated: %w", err)
		}
	}

	cancelTopMostCtx()

	// create new context because top-most is already canceled
	ctx, cancel := context.WithTimeout(context.Background(), conf.ExitTimeout)
	defer cancel()

	return server.Shutdown(ctx)
}
func initProfefe(
	logger *log.Logger,
	conf config.Config,
) (
	collector *profefe.Collector,
	querier *profefe.Querier,
	closer func(),
	err error,
) {
	stypes, err := conf.StorageType()
	if err != nil {
		return nil, nil, nil, err
	}

	if len(stypes) > 1 {
		logger.Infof("WARNING: several storage types specified: %s. Only first one %q is used for querying", stypes, stypes[0])
	}

	var (
		writers []storage.Writer
		reader  storage.Reader
		closers []io.Closer
	)

	assembleStorage := func(sw storage.Writer, sr storage.Reader, closer io.Closer) {
		writers = append(writers, sw)
		// only the first reader is used
		if reader == nil {
			reader = sr
		}
		if closer != nil {
			closers = append(closers, closer)
		}
	}

	initStorage := func(stype string) error {
		logger := logger.With(zap.String("storage", stype))

		switch stype {
		case config.StorageTypeBadger:
			st, closer, err := conf.Badger.CreateStorage(logger)
			if err == nil {
				assembleStorage(st, st, closer)
			}
			return err
		default:
			return fmt.Errorf("unknown storage type %q, config %v", stype, conf)
		}
	}

	for _, stype := range stypes {
		if err := initStorage(stype); err != nil {
			return nil, nil, nil, fmt.Errorf("could not init storage %q: %w", stype, err)
		}
	}

	closer = func() {
		for _, closer := range closers {
			if err := closer.Close(); err != nil {
				logger.Error(err)
			}
		}
	}

	var writer storage.Writer
	if len(writers) == 1 {
		writer = writers[0]
	} else {
		writer = storage.NewMultiWriter(writers...)
	}
	return profefe.NewCollector(logger, writer), profefe.NewQuerier(logger, reader), closer, nil
}
func setupDebugRoutes(mux *http.ServeMux) {
	// pprof handlers, see https://github.com/golang/go/blob/release-branch.go1.13/src/net/http/pprof/pprof.go
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))

	// expvar handlers, see https://github.com/golang/go/blob/release-branch.go1.13/src/expvar/expvar.go
	mux.Handle("/debug/vars", expvar.Handler())

	// prometheus handlers
	mux.Handle("/debug/metrics", promhttp.Handler())
}
