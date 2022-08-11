package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"manager_service/internal/config"
	"manager_service/internal/manager"
	"manager_service/internal/manager/db"
	"manager_service/pkg/client/mongodb"
	"manager_service/pkg/logging"
	"manager_service/pkg/metrics"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

type ManagerFunc func(ctx context.Context, fq *manager.FuncArray)

type App struct {
	managerFunc ManagerFunc
	funcQueue   manager.FuncArray
	cfg         *config.Config
	logger      *logging.Logger
	router      *httprouter.Router
	httpServer  *http.Server
}

func managerFunc(ctx context.Context, fa *manager.FuncArray) {
	bf := backoff.NewExponentialBackOff()
	bf.InitialInterval = 2 * time.Millisecond
	bf.MaxInterval = 10 * time.Second
	for {
		next := bf.NextBackOff()
		time.Sleep(next)
		lrs, err := fa.ManagerService.GetAll(ctx)
		if err != nil {
			continue
		}
		if len(lrs) == 0 {
			continue
		}
		for _, lr := range lrs {
			if !lr.Expired() {
				continue
			}
			updateTime, err := fa.ManagerService.UpdateTime(ctx, lr.LobbyID)
			if err != nil {
				log.Printf("failed to update time due to: %v", err)
				continue
				//err = fa.ManagerService.Delete(ctx, lr.ID)
				//if err != nil {
				//	panic("failed to delete lr")
				//}
			}
			lr.Expiration = updateTime
			err = fa.Update(lr)
			if err != nil {
				continue
			}
		}
	}
}

func NewApp(cfg *config.Config, logger *logging.Logger) (App, error) {
	logger.Println("router initializing")
	router := httprouter.New()
	logger.Println("swagger docs initialization")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	logger.Println("heartbeat metric initializing")
	metricHandler := metrics.Handler{}
	metricHandler.Register(router)

	mongodbClient, err := mongodb.NewClient(context.Background(), cfg.MongoDB.Host, cfg.MongoDB.Port,
		cfg.MongoDB.Username, cfg.MongoDB.Password, cfg.MongoDB.Database, cfg.MongoDB.AuthDB)
	if err != nil {
		panic(err)
	}

	storage := db.NewStorage(mongodbClient, "managers", logger)
	service, err := manager.NewService(storage, *logger)
	if err != nil {
		panic(err)
	}

	managersHandler := manager.Handler{
		Logger:         logging.GetLogger(cfg.AppConfig.LogLevel),
		ManagerService: service,
	}
	managersHandler.Register(router)

	return App{
		managerFunc,
		manager.GetFuncQueue(service),
		cfg,
		logger,
		router,
		nil,
	}, nil
}

func (a *App) Run() {
	go a.managerFunc(context.Background(), &a.funcQueue)
	a.startHTTP()
}

func (a *App) startHTTP() {
	a.logger.Info("start HTTP")

	var listener net.Listener

	if a.cfg.Listen.Type == config.ListenTypeSock {
		appDir, err := filepath.Abs(os.Args[0])
		if err != nil {
			a.logger.Fatal(err)
		}
		socketPath := path.Join(appDir, a.cfg.Listen.SocketFile)
		a.logger.Infof("socket path: %s", socketPath)

		a.logger.Info("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			a.logger.Fatal(err)
		}
	} else {
		a.logger.Infof("bind application to host: %s and port: %s", a.cfg.Listen.BindIP, a.cfg.Listen.Port)
		var err error
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
		if err != nil {
			a.logger.Fatal(err)
		}
	}

	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodOptions, http.MethodDelete},
		AllowedOrigins:     []string{"https://localhost:3000", "https://localhost:8080"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Authorization", "Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Access-Token", "Refresh-Token", "Location", "Authorization", "Content-Disposition"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	a.logger.Println("application completely initialized and started")
	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warn("server shutdown")
		default:
			a.logger.Fatal(err)
		}
	}
	err := a.httpServer.Shutdown(context.Background())
	if err != nil {
		a.logger.Fatal(err)
	}

}
