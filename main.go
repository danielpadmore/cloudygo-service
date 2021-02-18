package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/danielpadmore/cloudygo-service/config"
	"github.com/danielpadmore/cloudygo-service/data"
	"github.com/danielpadmore/cloudygo-service/handlers"
	"github.com/danielpadmore/cloudygo-service/logs"
	"github.com/danielpadmore/cloudygo-service/validation"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
)

// TODO - consistent logs
// TODO - tests

// Config contains configuration data for the application
type Config struct {
	DbConnection string `json:"db_connection"`
	BindAddress  string `json:"bind_address"`
}

var conf *Config

var configFile = env.String("CONFIG_FILE", false, "./conf.json", "Path to JSON encoded config file")

func newLog(message string, a ...interface{}) logs.LogStruct {
	return logs.NewLog("SETUP", fmt.Sprintf(message, a...))
}

func main() {

	logger := logs.NewStdLogger(logs.LogLevelInfo)

	logger.Info(newLog("Initiating CloudyGo service..."))

	err := env.Parse()
	if err != nil {
		logger.Fatal(newLog("Error parsing flags: %s", err.Error()))
		os.Exit(1)
	}

	conf = &Config{}

	c, err := config.New(*configFile, conf, func() {
		logger.Info(newLog("Config file updated"))
	})

	if err != nil {
		logger.Fatal(newLog("Unable to load config file: %s", err.Error()))
		os.Exit(1)
	}
	defer c.Close()

	validator := validation.New(logger)

	db, err := retryDbUntilReady(logger)
	if err != nil {
		logger.Fatal(newLog("Timed out waiting for database connection"))
		os.Exit(1)
	}

	router := mux.NewRouter()
	registerRoutes(router, logger, validator, db)

	logger.Info(newLog("Starting server on %s", conf.BindAddress))
	err = http.ListenAndServe(conf.BindAddress, router)
	if err != nil {
		logger.Fatal(newLog("Unable to start server on %s. Error: %s", conf.BindAddress, err.Error()))
	}

}

func isAuthorizedMiddleware(next func(userID string, w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")

		if authToken == "" {
			fmt.Println("No Authorization provided")
			http.Error(w, "No Authorization provided", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Println("Unable to parse JWT token", "path", r.URL.Path)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return nil, nil
			}
			return []byte(handlers.JWTSecret), nil
		})

		if err != nil {
			fmt.Println("Unauthorized", "error", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["user_id"].(string)
			next(userID, w, r)
			return
		}

		fmt.Println("Invalid token", "error", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	})
}

func retryDbUntilReady(logger logs.Logger) (data.Connection, error) {
	st := time.Now()
	dt := 1 * time.Second
	mt := 60 * time.Second

	for {
		db, err := data.New(logger, conf.DbConnection)
		if err == nil {
			logger.Info(newLog(fmt.Sprintf("Successfully connected to database in %fs", time.Now().Sub(st).Seconds())))
			return db, nil
		}

		logger.Warning(newLog(fmt.Sprintf("Unable to connect to the database. Error: %s", err.Error())))
		fmt.Println("Unable to connect to the database", "error", err)

		if time.Now().Sub(st) > mt {
			return nil, err
		}

		time.Sleep(dt)
	}
}

func registerRoutes(router *mux.Router, logger logs.Logger, validator validation.Validator, db data.Connection) {

	logger.Debug(newLog("Registering routes"))

	healthHandler := handlers.NewHealth(logger, db)
	router.Handle("/health", healthHandler).Methods("GET")

	resourceHandler := handlers.NewResource(logger, db)
	router.Handle("/resources", resourceHandler).Methods("GET")

	userHandler := handlers.NewUser(logger, db)
	router.HandleFunc("/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/signin", userHandler.SignIn).Methods("POST")

	lambdaHandler := handlers.NewLambda(logger, validator, db)
	lambdaRouter := router.PathPrefix("/lambdas").Subrouter()
	lambdaRouter.Handle("", isAuthorizedMiddleware(lambdaHandler.CreateLambda)).Methods("POST")
	lambdaRouter.Handle("", isAuthorizedMiddleware(lambdaHandler.GetLambdas)).Methods("GET")
	lambdaRouter.Handle("/{id}", isAuthorizedMiddleware(lambdaHandler.GetLambda)).Methods("GET")
	lambdaRouter.Handle("/{id}", isAuthorizedMiddleware(lambdaHandler.UpdateLambda)).Methods("PUT")
	lambdaRouter.Handle("/{id}", isAuthorizedMiddleware(lambdaHandler.DeleteLambda)).Methods("DELETE")

	vmHandler := handlers.NewVirtualMachine(logger, db)
	vmRouter := router.PathPrefix("/virtual-machines").Subrouter()
	vmRouter.Handle("", isAuthorizedMiddleware(vmHandler.CreateVirtualMachine)).Methods("POST")
	vmRouter.Handle("", isAuthorizedMiddleware(vmHandler.GetVirtualMachines)).Methods("GET")
	vmRouter.Handle("/{id}", isAuthorizedMiddleware(vmHandler.GetVirtualMachine)).Methods("GET")
	vmRouter.Handle("/{id}", isAuthorizedMiddleware(vmHandler.UpdateVirtualMachine)).Methods("PUT")
	vmRouter.Handle("/{id}", isAuthorizedMiddleware(vmHandler.DeleteVirtualMachine)).Methods("DELETE")

	sqldbHandler := handlers.NewSQLDatabase(logger, db)
	sqldbRouter := router.PathPrefix("/sql-databases").Subrouter()
	sqldbRouter.Handle("", isAuthorizedMiddleware(sqldbHandler.CreateSQLDatabase)).Methods("POST")
	sqldbRouter.Handle("", isAuthorizedMiddleware(sqldbHandler.GetSQLDatabases)).Methods("GET")
	sqldbRouter.Handle("/{id}", isAuthorizedMiddleware(sqldbHandler.GetSQLDatabase)).Methods("GET")
	sqldbRouter.Handle("/{id}", isAuthorizedMiddleware(sqldbHandler.UpdateSQLDatabase)).Methods("PUT")
	sqldbRouter.Handle("/{id}", isAuthorizedMiddleware(sqldbHandler.DeleteSQLDatabase)).Methods("DELETE")

	nosqldbHandler := handlers.NewNoSQLDatabase(logger, db)
	nosqldbRouter := router.PathPrefix("/nosql-databases").Subrouter()
	nosqldbRouter.Handle("", isAuthorizedMiddleware(nosqldbHandler.CreateNoSQLDatabase)).Methods("POST")
	nosqldbRouter.Handle("", isAuthorizedMiddleware(nosqldbHandler.GetNoSQLDatabases)).Methods("GET")
	nosqldbRouter.Handle("/{id}", isAuthorizedMiddleware(nosqldbHandler.GetNoSQLDatabase)).Methods("GET")
	nosqldbRouter.Handle("/{id}", isAuthorizedMiddleware(nosqldbHandler.UpdateNoSQLDatabase)).Methods("PUT")
	nosqldbRouter.Handle("/{id}", isAuthorizedMiddleware(nosqldbHandler.DeleteNoSQLDatabase)).Methods("DELETE")

	logger.Debug(newLog("Routes registered"))

}
