package main

import (
	"context"
	"database/sql"
	"fmt"
	nr "github.com/newrelic/go-agent/v3/newrelic"
	billscan "github.com/okcredit/billscan/service"
	"github.com/okcredit/billscan/service/database/postgres"
	billscan_http "github.com/okcredit/billscan/service/http"
	"github.com/okcredit/go-common"
	"github.com/okcredit/go-common/config"
	"github.com/okcredit/go-common/log"
	"github.com/okcredit/go-common/newrelic"
	"github.com/okcredit/go-common/shutdown"
	"github.com/okcredit/nap"
	"net/http"
	"time"
)

func main() {
	// config
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	// database
	var dbConf config.DatabaseConfig
	err = conf.UnmarshalKey("database", &dbConf)
	if err != nil {
		panic(err)
	}
	db := initDatabase(&dbConf)

	// service
	srv := &billscan.Service{
		Database: postgres.New(db),
	}

	httpSrv := (*http.Server)(nil)
	go func() {
		httpSrv = &http.Server{
			Addr:         ":8080",
			Handler:      billscan_http.New(srv),
			IdleTimeout:  30 * time.Second,
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
		}

		log.Println("starting http server...")
		if err := httpSrv.ListenAndServe(); err != nil {
			log.Fatalf("http server error: %v", err)
		}
	}()

	<-shutdown.Handle(func() {
		log.Println("shutting down...")

		// shutdown http server
		if httpSrv != nil {
			ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
			if err := httpSrv.Shutdown(ctx); err != nil {
				log.Printf("failed to shutdown http server: %v", err)
			} else {
				log.Printf("http server shutdown")
			}
		}
	})
}

func initDatabase(databaseConfig *config.DatabaseConfig) *nap.DB {
	host := databaseConfig.Host
	port := databaseConfig.Port
	username := databaseConfig.Username
	password := databaseConfig.Password
	dbname := databaseConfig.Dbname

	dataSrc := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, dbname)
	db, err := sql.Open("postgres", dataSrc)
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(databaseConfig.IdleConn)
	db.SetMaxOpenConns(databaseConfig.MaxConn)

	extensions := []common.DBExtension{
		newrelic.GetDBMiddleware(nr.DatastorePostgres, dbname),
		common.GetQueryLoggerExtension(),
	}

	return nap.Of([]*sql.DB{db}, "postgres", nap.WithExtensions(extensions))
}
