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
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
)

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

	db, err := retryDbUntilReady(&logger)
	if err != nil {
		logger.Fatal(newLog("Timed out waiting for database connection"))
		os.Exit(1)
	}

	router := mux.NewRouter()

	healthHandler := handlers.NewHealth(&logger, db)
	router.Handle("/health", healthHandler).Methods("GET")

	resourceHandler := handlers.NewResource(&logger, db)
	router.Handle("/resources", resourceHandler).Methods("GET")

	userHandler := handlers.NewUser(&logger, db)
	router.HandleFunc("/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/signin", userHandler.SignIn).Methods("POST")

	lambdaHandler := handlers.NewLambda(&logger, db)
	router.Handle("/lambdas", isAuthorizedMiddleware(lambdaHandler.CreateLambda)).Methods("POST")
	router.Handle("/lambdas", isAuthorizedMiddleware(lambdaHandler.GetLambdas)).Methods("GET")
	router.Handle("/lambdas/{id}", isAuthorizedMiddleware(lambdaHandler.GetLambda)).Methods("GET")
	router.Handle("/lambdas/{id}", isAuthorizedMiddleware(lambdaHandler.UpdateLambda)).Methods("PUT")
	router.Handle("/lambdas/{id}", isAuthorizedMiddleware(lambdaHandler.DeleteLambda)).Methods("DELETE")

	vmHandler := handlers.NewVirtualMachine(&logger, db)
	router.Handle("/virtual-machines", isAuthorizedMiddleware(vmHandler.CreateVirtualMachine)).Methods("POST")
	router.Handle("/virtual-machines", isAuthorizedMiddleware(vmHandler.GetVirtualMachines)).Methods("GET")
	router.Handle("/virtual-machines/{id}", isAuthorizedMiddleware(vmHandler.GetVirtualMachine)).Methods("GET")
	router.Handle("/virtual-machines/{id}", isAuthorizedMiddleware(vmHandler.UpdateVirtualMachine)).Methods("PUT")
	router.Handle("/virtual-machines/{id}", isAuthorizedMiddleware(vmHandler.DeleteVirtualMachine)).Methods("DELETE")

	sqldbHandler := handlers.NewSQLDatabase(&logger, db)
	router.Handle("/sql-databases", isAuthorizedMiddleware(sqldbHandler.CreateSQLDatabase)).Methods("POST")
	router.Handle("/sql-databases", isAuthorizedMiddleware(sqldbHandler.GetSQLDatabases)).Methods("GET")
	router.Handle("/sql-databases/{id}", isAuthorizedMiddleware(sqldbHandler.GetSQLDatabase)).Methods("GET")
	router.Handle("/sql-databases/{id}", isAuthorizedMiddleware(sqldbHandler.UpdateSQLDatabase)).Methods("PUT")
	router.Handle("/sql-databases/{id}", isAuthorizedMiddleware(sqldbHandler.DeleteSQLDatabase)).Methods("DELETE")

	nosqldbHandler := handlers.NewNoSQLDatabase(&logger, db)
	router.Handle("/nosql-databases", isAuthorizedMiddleware(nosqldbHandler.CreateNoSQLDatabase)).Methods("POST")
	router.Handle("/nosql-databases", isAuthorizedMiddleware(nosqldbHandler.GetNoSQLDatabases)).Methods("GET")
	router.Handle("/nosql-databases/{id}", isAuthorizedMiddleware(nosqldbHandler.GetNoSQLDatabase)).Methods("GET")
	router.Handle("/nosql-databases/{id}", isAuthorizedMiddleware(nosqldbHandler.UpdateNoSQLDatabase)).Methods("PUT")
	router.Handle("/nosql-databases/{id}", isAuthorizedMiddleware(nosqldbHandler.DeleteNoSQLDatabase)).Methods("DELETE")

	logger.Info(newLog("Routing configured, starting server on %s", conf.BindAddress))
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
