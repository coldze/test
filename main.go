package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/coldze/test/logic"
	"github.com/coldze/test/logic/handles"
	"github.com/coldze/test/logic/sources"
	"github.com/coldze/test/logs"
	"github.com/coldze/test/utils"
)

const (
	HEALTH_CHECK_PATH   = "/ping"
	CONTACT_ID_VARIABLE = "contactid"
	CONTACT_ROUTE       = "/contact"
	API_VERSION         = "v1"
)

func newDataSource(cfg *appCfg) (sources.DataSource, error) {
	httpDataSource := sources.NewDefaultHttpDataSource(cfg.Api)
	rWrap, err := sources.NewRedisWrap(cfg.GetRedisOptions())
	if err != nil {
		return nil, err
	}
	cacheSource := sources.NewRedisCacheSource(rWrap, cfg.GetCacheTtl())
	return sources.NewCachedDataSource(httpDataSource, cacheSource), nil
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(fmt.Sprintf("Health check at: %v\n", time.Now().UTC())))
	if r.Body == nil {
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()
	_, _ = ioutil.ReadAll(r.Body)
}

//I've put it here, because it is dependent on gorilla/mux, I didn't want to spoil business-logic code with such dependencies
//This function can be replaced by our own implementation with regex or other manipulations with strings
//Didn't want to re-implement that logic, as this code is already dependent on gorilla mux, decided to use it's feature
func NewGetVariableFromRequest(varName string) logic.RequestDataExtractor {
	return func(r *http.Request) ([]byte, error) {
		vars := mux.Vars(r)
		value, ok := vars[varName]
		if !ok {
			return nil, fmt.Errorf("Failed to find  '%v'", varName)
		}
		return []byte(value), nil
	}
}

func buildRoutes(dataSource sources.DataSource, logger logs.Logger) http.Handler {
	getData := NewGetVariableFromRequest(CONTACT_ID_VARIABLE)
	getHandler := handles.NewGetHandler(handles.NewDefaultLoggerFactory(logs.NewPrefixedLogger(logger, "[GET]")), dataSource, getData)
	createHandler := handles.NewPostHandler(handles.NewDefaultLoggerFactory(logs.NewPrefixedLogger(logger, "[POST]")), dataSource)
	updateHandler := handles.NewPutHandler(handles.NewDefaultLoggerFactory(logs.NewPrefixedLogger(logger, "[PUT]")), dataSource)

	router := mux.NewRouter()
	router.Path(HEALTH_CHECK_PATH).HandlerFunc(healthCheck)
	sr := router.PathPrefix(fmt.Sprintf("/%s", API_VERSION)).Subrouter()

	sr.HandleFunc(CONTACT_ROUTE, createHandler).Methods(http.MethodPost)
	sr.HandleFunc(CONTACT_ROUTE, updateHandler).Methods(http.MethodPut)
	sr.HandleFunc(fmt.Sprintf("%s/{%s}", CONTACT_ROUTE, CONTACT_ID_VARIABLE), getHandler).Methods(http.MethodGet)
	return router
}

func newMainFunc(cfg *appCfg) utils.MainFunc {
	return func(logger logs.Logger, stop <-chan struct{}) int {
		dataSource, err := newDataSource(cfg)
		if err != nil {
			logger.Errorf("Failed to create data-source. Error: %v", err)
			return 1
		}

		router := buildRoutes(dataSource, logger)

		bind := cfg.GetBind()
		srv, err := utils.NewService(bind, router)
		if err != nil {
			logger.Errorf("Failed to start service. Error: %v", err)
			return 1
		}
		defer func() {
			cErr := srv.Stop()
			if cErr != nil {
				logger.Errorf("Failed to stop service: %+v", cErr)
			}
		}()
		logger.Infof("Ready. Listening at '%s'", bind)
		<-stop
		return 0
	}
}

func main() {
	configPath := flag.String("config", "./config.json", "service's configuration in JSON format")
	redisPwd := flag.String("redispwd", "", "Redis password")
	flag.Parse()
	logger := logs.NewStdLogger()
	logger.Infof("Starting...")
	cfg, err := getConfig(*configPath, *redisPwd)
	if err != nil {
		logger.Errorf("Failed to load config. Error: %v", err)
		return
	}
	utils.Run(cfg.GetAppTimeout(), newMainFunc(cfg), logger)
	logger.Infof("Done")
}
