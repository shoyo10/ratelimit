package cmd

import (
	"context"
	"ratelimit/internal/httpserver/delivery/http"
	"ratelimit/pkg/config"
	"ratelimit/pkg/echorouter"
	ratelimitRedis "ratelimit/pkg/ratelimit/redis"
	"ratelimit/pkg/redis"
	"ratelimit/pkg/zerolog"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// ServerCmd ...
var ServerCmd = &cobra.Command{
	Run: run,
	Use: "server",
}

func run(cmd *cobra.Command, args []string) {
	defer recoverFn()

	rand.Seed(time.Now().UnixNano())

	cfg, err := config.New()
	if err != nil {
		os.Exit(1)
		return
	}
	zerolog.Init(cfg.Log)

	app := fx.New(
		fx.Supply(*cfg),
		fx.Provide(
			redis.New,
			echorouter.Start,
			ratelimitRedis.NewRateLimit,
		),
		fx.Invoke(http.NewHandler),
	)
	startApp(cmd.Name(), app)
}

func recoverFn() {
	if r := recover(); r != nil {
		var msg string
		for i := 2; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			msg += fmt.Sprintf("%s:%d\n", file, line)
		}
		log.Error().Msgf("========== PANIC Start ==========\n%s\n%s========== PANIC End ==========", r, msg)
	}
}

func startApp(name string, app *fx.App) {
	if err := app.Start(context.Background()); err != nil {
		log.Error().Msg(err.Error())
		return
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
	log.Info().Msgf("main: shutting down %s...", name)

	stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Error().Msg(err.Error())
	}
}
