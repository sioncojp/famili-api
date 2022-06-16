package application

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	v1todos "github.com/sioncojp/famili-api/application/v1/todos"
	"github.com/sioncojp/famili-api/utils/config"
	"github.com/sioncojp/famili-api/utils/log"
)

// HttpHandler...http_response serverを立ち上げるため必要なstruct
type HttpHandler struct {
	AppConfig *config.AppConfig
	Router
	// ServeMux...HTTP request multiplexer. リクエストを登録済みのURLパターンリストと照合して、マッチしたHandlerを呼び出す
	ServeMux *chi.Mux
}

// Router...ルーティング情報
type Router struct {
	V1
}

// V1Handler.../v1 で利用するstructを格納
type V1 struct {
	TodosHandler v1todos.Handler
}

// RunServer...サーバ起動
func (s *HttpHandler) RunServer() {
	log.Log.Debugf("start run server port :%s", s.AppConfig.Server.Port)

	r := chi.NewRouter()
	r.Mount("/", s.ServeMux)

	server := &http.Server{
		Addr:         fmt.Sprintf(":" + s.AppConfig.Server.Port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Log.Fatalf("could not start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Log.Info("server shutdown signal has been received, the service will exit in 30 seconds.")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// graceful shutdown http.Server
	if err := server.Shutdown(ctx); err != nil {
		log.Log.Fatalf("Could not gracefully shutdown the server:%v", err)
	}
	log.Log.Info("server is graceful shutdown now, new request will be rejected.")

	// waiting for ctx.Done(). timeout of 30 seconds.
	<-ctx.Done()

	log.Log.Info("server shutdown")
	// TODO: swagger分岐欲しい
}
