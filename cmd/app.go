package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/av-belyakov/placeholder_misp/cmd/coremodule"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/appversion"
	"github.com/av-belyakov/placeholder_misp/internal/dicontainer"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

type App struct {
	diContainer *dicontainer.DiContainer
	coremodule  *coremodule.CoreHandler
}

func NewApp(ctx context.Context) *App {
	rootPath, err := supportingfunctions.GetRootPath(constants.Root_Dir)
	if err != nil {
		log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
	}

	// сервер для отладки
	if os.Getenv("GO_PHMISP_MAIN") == "test" || os.Getenv("GO_PHMISP_MAIN") == "development" {
		httpServer := &http.Server{
			Addr: fmt.Sprintf("%s:%d", "localhost", 6161),
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
		}

		g, gCtx := errgroup.WithContext(ctx)
		g.Go(func() error {
			return httpServer.ListenAndServe()
		})
		g.Go(func() error {
			<-gCtx.Done()

			return httpServer.Shutdown(context.Background())
		})

		if err := g.Wait(); err != nil {
			log.Fatal("error debugging server:", err)
		}
	}

	app := &App{
		diContainer: dicontainer.NewDIContainer(rootPath, make(chan commoninterfaces.Messager)),
	}
	app.coremodule = coremodule.NewCoreHandler(
		app.diContainer.Logger(ctx),
		app.diContainer.Counter(ctx),
		app.diContainer.Rules(ctx),
	)

	version, err := appversion.GetAppVersion()
	if err != nil {
		log.Println(err)
	}

	// вывод информационного сообщения при старте приложения
	msg := getInformationMessage(version)
	app.diContainer.SimpleLogger(ctx).Write("info", strings.ToLower(msg))

	// старт приложения
	app.coremodule.Start(ctx, app.diContainer.NatsConnecter(ctx), app.diContainer.MispConnecter(ctx), app.diContainer.DB(ctx))

	return app
}
