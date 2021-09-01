package main

import (
	"fmt"
	"github.com/link33/sidecar/internal"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/link33/sidecar/internal/app"
	"github.com/link33/sidecar/internal/loggers"
	"github.com/link33/sidecar/internal/repo"
	"github.com/meshplus/bitxhub-kit/log"
	"github.com/urfave/cli"
)

var (
	startCMD = cli.Command{
		Name:   "start",
		Usage:  "Start a long-running daemon process",
		Action: start,
	}
)

func start(ctx *cli.Context) error {
	fmt.Println(getVersion(true))

	repoRoot, err := repo.PathRootWithDefault(ctx.GlobalString("repo"))
	if err != nil {
		return err
	}

	repo.SetPath(repoRoot)

	config, err := repo.UnmarshalConfig(repoRoot)
	if err != nil {
		return fmt.Errorf("init config error: %s", err)
	}

	err = log.Initialize(
		log.WithReportCaller(config.Log.ReportCaller),
		log.WithPersist(true),
		log.WithFilePath(filepath.Join(repoRoot, config.Log.Dir)),
		log.WithFileName(config.Log.Filename),
		log.WithMaxSize(2*1024*1024),
		log.WithMaxAge(24*time.Hour),
		log.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return fmt.Errorf("log initialize: %w", err)
	}
	// init loggers map for sidecar
	loggers.InitializeLogger(config)

	var sidecar internal.Launcher

	sidecar, err = app.NewSidecar(repoRoot, config)
	if err != nil {
		return err
	}

	runPProf(config.Port.PProf)

	var wg sync.WaitGroup
	wg.Add(1)
	handleShutdown(sidecar, &wg)

	if err := sidecar.Start(); err != nil {
		return err
	}

	wg.Wait()

	logger.Info("Sidecar exits")
	return nil
}

func handleShutdown(sidecar internal.Launcher, wg *sync.WaitGroup) {
	var stop = make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	go func() {
		<-stop
		fmt.Println("received interrupt signal, shutting down...")
		if err := sidecar.Stop(); err != nil {
			logger.Error("sidecar stop: ", err)
		}

		wg.Done()
		os.Exit(0)
	}()
}

func runPProf(port int64) {
	go func() {
		addr := fmt.Sprintf("localhost:%d", port)
		fmt.Printf("Pprof on localhost:%d\n\n", port)
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}()
}

func checkPlugin(pluginName string) error {
	// check if plugin exists
	pluginRoot, err := repo.PluginPath()
	if err != nil {
		return err
	}

	pluginPath := filepath.Join(pluginRoot, pluginName)
	_, err = os.Stat(pluginPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("plugin file `%s` is required", pluginPath)
		}

		return fmt.Errorf("get plugin file state error: %w", err)
	}

	return nil
}
