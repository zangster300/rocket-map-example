package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"rocket-map-example/config"
	"rocket-map-example/web/resources"

	"github.com/evanw/esbuild/pkg/api"
	"golang.org/x/sync/errgroup"

	"github.com/fsnotify/fsnotify"
)

var watch = false

func init() {
	flag.BoolVar(&watch, "watch", watch, "Enable watcher mode")
	flag.Parse()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)
	defer cancel()

	if err := run(ctx); err != nil {
		slog.Error("failure", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return build(egctx)
	})

	eg.Go(func() error {
		return notify(egctx)
	})

	return eg.Wait()
}

func build(ctx context.Context) error {
	opts := api.BuildOptions{
		EntryPointsAdvanced: []api.EntryPoint{
			{
				InputPath:  resources.StylesDirectoryPath + "/styles.css",
				OutputPath: "css/index",
			},
		},
		Bundle:            true,
		Format:            api.FormatESModule,
		LogLevel:          api.LogLevelInfo,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		MinifyWhitespace:  true,
		Outdir:            resources.StaticDirectoryPath,

		Sourcemap: api.SourceMapLinked,
		Target:    api.ESNext,
		Write:     true,
	}

	if watch {
		slog.Info("watching...")

		opts.Plugins = append(opts.Plugins, api.Plugin{
			Name: "hotreload",
			Setup: func(build api.PluginBuild) {
				build.OnEnd(func(result *api.BuildResult) (api.OnEndResult, error) {
					slog.Info("build complete", "errors", len(result.Errors), "warnings", len(result.Warnings))
					if len(result.Errors) == 0 {
						http.Get(fmt.Sprintf("http://%s:%s/hotreload", config.Global.Host, config.Global.Port))
					}
					return api.OnEndResult{}, nil
				})
			},
		})

		buildCtx, err := api.Context(opts)
		if err != nil {
			return err
		}
		defer buildCtx.Dispose()

		if err := buildCtx.Watch(api.WatchOptions{}); err != nil {
			return err
		}

		<-ctx.Done()
		return nil
	}

	slog.Info("building...")

	result := api.Build(opts)

	if len(result.Errors) > 0 {
		errs := make([]error, len(result.Errors))
		for i, err := range result.Errors {
			errs[i] = errors.New(err.Text)
		}
		return errors.Join(errs...)
	}

	return nil
}

func notify(ctx context.Context) error {
	if !watch {
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					slog.Info("modified file", slog.String("name", event.Name))
					http.Get(fmt.Sprintf("http://%s:%s/hotreload", config.Global.Host, config.Global.Port))
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				slog.Info(fmt.Sprintf("error: %v", err))
			}
		}
	}()

	if err := watcher.Add(resources.StaticDirectoryPath + "/rocket"); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
