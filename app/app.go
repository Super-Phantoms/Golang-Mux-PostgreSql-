package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangdevm/fullstack/domain"
	"github.com/golangdevm/fullstack/service"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"

	"github.com/golangdevm/fullstack/logger"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql database driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // sqlite database driver
)

func Start() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	sanityCheck()
	router := mux.NewRouter().StrictSlash(true)
	// r := mux.NewRouter().StrictSlash(true)
	// api :=router.PathPrefix("/api").Subrouter()
	auth_router := router.PathPrefix("/auth").Subrouter()
	api_router := router.PathPrefix("/api").Subrouter()

	sqlDbClient := getSqlDbClient()
	dbClient := getDbClient()
	authRepository := domain.NewAuthRepository(dbClient, sqlDbClient)
	ah := AuthHandler{service.NewLoginService(authRepository, domain.GetRolePermissions())}

	auth_router.HandleFunc("/login", ah.Login).Methods(http.MethodPost)       //done
	auth_router.HandleFunc("/register", ah.Register).Methods(http.MethodPost) //done
	auth_router.HandleFunc("/refresh", ah.Refresh).Methods(http.MethodPost)   //done
	auth_router.HandleFunc("/verify", ah.Verify).Methods(http.MethodGet)

	// define routes

	customerRepositoryDb := domain.NewCustomerRepositoryDb(dbClient)
	accountRepositoryDb := domain.NewAccountRepositoryDb(dbClient)
	ch := CustomerHandlers{service.NewCustomerService(customerRepositoryDb)}
	ac := AccountHandler{service.NewAccountService(accountRepositoryDb)}

	api_router.
		HandleFunc("/customers", ch.getAllCustomers).
		Methods(http.MethodGet).
		Name("GetAllCustomers")
	api_router.
		HandleFunc("/customers/{customer_id:[0-9]+}", ch.getCustomer).
		Methods(http.MethodGet).
		Name("GetCustomer")
	api_router.
		HandleFunc("/customers/{customer_id:[0-9]+}/account", ac.NewAccount).
		Methods(http.MethodPost).
		Name("NewAccount")
	api_router.
		HandleFunc("/customers/{customer_id:[0-9]+}/account/{account_id:[0-9]+}", ac.MakeTransaction).
		Methods(http.MethodPost).
		Name("NewTransaction")

	am := AuthMiddleware{domain.NewAuthRepository(dbClient, sqlDbClient)}
	api_router.Use(am.authorizationHandler())

	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	logger.Info(fmt.Sprintf("Starting server on %s:%s ...", address, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router))
}

func getDbClient() *gorm.DB {
	DBDriver := os.Getenv("DB_DRIVER")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbHost := os.Getenv("DB_HOST")
	DbPort := os.Getenv("DB_PORT")
	DbName := os.Getenv("DB_NAME")

	fmt.Println("db_driver::::", os.Getenv("DB_DRIVER"))

	var err error
	var client *gorm.DB
	if DBDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		client, err = gorm.Open(DBDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", DBDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("Successfully connected to the %s database", DBDriver)
		}
	}
	//database migration
	client.Debug().AutoMigrate(
		&domain.User{},
		&domain.Post{},
		&domain.RefreshToken{},
		&domain.Account{},
		&domain.Customer{},
		&domain.Transaction{},
	)
	return client
}

func getSqlDbClient() *sqlx.DB {
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWORD")
	dbAddr := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPasswd, dbAddr, dbPort, dbName)
	client, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}
func sanityCheck() {
	envProps := []string{
		"SERVER_ADDRESS",
		"SERVER_PORT",
		"DB_HOST",
		"DB_DRIVER",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_PORT",
	}
	for _, k := range envProps {
		if os.Getenv(k) == "" {
			logger.Error(fmt.Sprintf("Environment variable %s not defined. Terminating application...", k))
		}
	}
}
