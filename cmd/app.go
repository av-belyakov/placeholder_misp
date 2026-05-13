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
	"github.com/av-belyakov/placeholder_misp/internal/wrappers"
)

type App struct {
	diContainer *dicontainer.DiContainer
	coremodule  *coremodule.CoreHandler
	ctx         context.Context
}

func NewApp(ctx context.Context) *App {
	rootPath, err := supportingfunctions.GetRootPath(constants.Root_Dir)
	if err != nil {
		log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
	}

	ch := make(chan commoninterfaces.Messager)
	app := &App{
		ctx:         ctx,
		diContainer: dicontainer.NewDIContainer(rootPath, ch),
	}
	app.coremodule = coremodule.NewCoreHandler(
		app.diContainer.Logger(ctx),
		app.diContainer.Counter(ctx),
		app.diContainer.Rules(ctx),
	)

	// настройка обёртки для взаимодействия с Zabbix
	zabbixSettings := wrappers.WrappersZabbixInteractionSettings{
		NetworkPort: app.diContainer.Configer().GetCommonApp().Zabbix.NetworkPort,
		NetworkHost: app.diContainer.Configer().GetCommonApp().Zabbix.NetworkHost,
		ZabbixHost:  app.diContainer.Configer().GetCommonApp().Zabbix.ZabbixHost,
		EventTypes:  make([]wrappers.EventType, len(app.diContainer.Configer().GetCommonApp().Zabbix.EventTypes)),
	}
	for _, v := range app.diContainer.Configer().GetCommonApp().Zabbix.EventTypes {
		zabbixSettings.EventTypes = append(zabbixSettings.EventTypes, wrappers.EventType{
			IsTransmit: v.IsTransmit,
			EventType:  v.EventType,
			ZabbixKey:  v.ZabbixKey,
			Handshake: wrappers.Handshake{
				TimeInterval: v.Handshake.TimeInterval,
				Message:      v.Handshake.Message,
			},
		})
	}
	// обертка для взаимодействия с Zabbix
	wrappers.WrappersZabbixInteraction(ctx, zabbixSettings, app.diContainer.SimpleLogger(ctx), ch)

	return app
}

func (a *App) Start() {
	// сервер для отладки
	if os.Getenv("GO_PHMISP_MAIN") == "test" || os.Getenv("GO_PHMISP_MAIN") == "development" {
		go func() {
			debugServerHost := "localhost"
			debugServerPort := 6161

			httpServer := &http.Server{
				Addr: fmt.Sprintf("%s:%d", debugServerHost, debugServerPort),
				BaseContext: func(_ net.Listener) context.Context {
					return a.ctx
				},
			}

			g, gCtx := errgroup.WithContext(a.ctx)
			g.Go(func() error {
				return httpServer.ListenAndServe()
			})
			g.Go(func() error {
				<-gCtx.Done()

				return httpServer.Shutdown(context.Background())
			})

			log.Printf("%vdebug server %v%s:%d%v\n", constants.Ansi_Bright_Green, constants.Ansi_Dark_Gray, debugServerHost, debugServerPort, constants.Ansi_Reset)

			if err := g.Wait(); err != nil {
				log.Fatal("error debugging server:", err)
			}
		}()
	}

	version, err := appversion.GetAppVersion()
	if err != nil {
		log.Println(err)
	}

	// вывод информационного сообщения при старте приложения
	msg := getInformationMessage(version)
	a.diContainer.SimpleLogger(a.ctx).Write("info", strings.ToLower(msg))

	// старт приложения
	a.coremodule.Start(a.ctx, a.diContainer.NatsConnecter(a.ctx), a.diContainer.MispConnecter(a.ctx), a.diContainer.DB(a.ctx))
}
