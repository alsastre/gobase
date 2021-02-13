package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alsastre/gobase/internal/config"
	"github.com/alsastre/gobase/internal/data"
	s "github.com/alsastre/gobase/internal/server"
	"github.com/upper/db/v4"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

func main() {
	var server s.Server
	var err error
	// Setup config
	viperCfg, err := config.New()
	if err != nil {
		fmt.Println("A problem occurred initalizing the configuration: ", err.Error())
		os.Exit(1)
	}
	server.ViperCfg = viperCfg

	// Setup Logger
	logger, err := SetupLogger(server)
	if err != nil {
		fmt.Println("A problem occurred initalizing the logger: ", err.Error())
		os.Exit(1)
	}
	server.Logger = logger
	logger.Info("Logger Configured")

	// Setup DB
	// TODO get DB config and define type from config
	var typeDB = "postgresql"
	var dbSession db.Session

	if typeDB == "postgresql" {
		dbSession, err = data.NewPostgresqlSession()
	}

	// From this point forward we should not care the DB Type selected

	if err != nil {
		fmt.Println("A problem occurred initalizing the database: ", err.Error())
		os.Exit(1)
	}

	defer dbSession.Close()

	// Validate the connection to the DB
	if err := dbSession.Ping(); err != nil {
		logger.Error("Ping: ", zap.Error(err))
	}

	logger.Info("Successfully connected to database", zap.String("connectionName", dbSession.Name()))
	server.Data = data.NewData(dbSession, logger)

	// Fill DB with a Dummy Something collection
	server.Data.FillSession()

	// Setup Routes
	r := SetupRouter(&server)
	logger.Info("Router Configured")

	// Start server
	http.ListenAndServe(":3333", r)
}

// SetupLogger setups the logger configuration
func SetupLogger(server s.Server) (*zap.Logger, error) {
	zapCfg := zap.NewDevelopmentConfig()
	err := zapCfg.Level.UnmarshalText([]byte(server.ViperCfg.GetString(config.LogLevelKey)))
	if err != nil {
		return nil, err
	}

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}

	defer logger.Sync() // flushes buffer, if any
	return logger, nil
}

// SetupRouter setups the router configuration
func SetupRouter(server *s.Server) *chi.Mux {
	r := chi.NewRouter()

	// Add RequestID middleware to add a request ID to every request
	r.Use(middleware.RequestID)
	// Add a RequestLogger middleware to log all the requestc
	r.Use(RequestLogger(server.Logger))
	// Add Recoverer middleware to recover from panics and print the stacktrace on the log
	r.Use(middleware.Recoverer)
	// Add SetContentType middleware to always set the ContentType header
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Dummy endpoint
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Something objects endpoint, to retrieve/update/list/...
	r.Mount("/something", server.Router())
	return r
}

// RequestLogger is a Request Middleware that will log every request processed using the zap logger provided
func RequestLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				if ww.Status() >= 200 && ww.Status() < 300 {
					// Log with Debug
					logger.Debug("Served",
						zap.String("protocol", r.Proto),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.Duration("latency", time.Since(t1)),
						zap.Int("status", ww.Status()),
						zap.Int("size", ww.BytesWritten()),
						zap.String("reqId", middleware.GetReqID(r.Context())))
				} else if ww.Status() >= 500 {
					// Log with Error
					logger.Error("Served",
						zap.String("protocol", r.Proto),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.Duration("latency", time.Since(t1)),
						zap.Int("status", ww.Status()),
						zap.Int("size", ww.BytesWritten()),
						zap.String("reqId", middleware.GetReqID(r.Context())))

				} else {
					// Log with Warning
					logger.Warn("Served",
						zap.String("protocol", r.Proto),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.Duration("latency", time.Since(t1)),
						zap.Int("status", ww.Status()),
						zap.Int("size", ww.BytesWritten()),
						zap.String("reqId", middleware.GetReqID(r.Context())))
				}
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
