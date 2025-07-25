package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Db     *Database
		Server *Server
		Rabbit *RabbitMQ
		Minio  *Minio
	}

	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DbName   string
		SslMode  string
		TimeZone string
	}

	Server struct {
		BaseDomain            string
		Port                  int
		InteractionServiceURL string
	}

	RabbitMQ struct {
		Host     string
		Port     int
		User     string
		Password string
	}

	Minio struct {
		MinioEndpoint   string
		MinioAccessKey  string
		MinioSecretKey  string
		MinioUseSSL     bool
		MinioBucketName string
	}
)

var (
	once   sync.Once
	config *Config
)

func GetConfig() *Config {
	once.Do(func() {
		config = &Config{
			Db:     &Database{},
			Server: &Server{},
			Rabbit: &RabbitMQ{},
			Minio:  &Minio{},
		}

		err := godotenv.Load()

		if err != nil {
			log.Println("Warning: Could not load .env file")
		}

		config.Db.Host = os.Getenv("DB_HOST")
		config.Db.Port, err = strconv.Atoi(os.Getenv("DB_PORT"))

		if err != nil {
			panic(err)
		}

		config.Db.User = os.Getenv("DB_USER")
		config.Db.Password = os.Getenv("DB_PASSWORD")
		config.Db.DbName = os.Getenv("DB_NAME")
		config.Db.SslMode = os.Getenv("DB_SSLMODE")
		config.Db.TimeZone = os.Getenv("DB_TIMEZONE")

		config.Server.BaseDomain = os.Getenv("BASE_DOMAIN")
		config.Server.Port, err = strconv.Atoi(os.Getenv("PORT"))
		config.Server.InteractionServiceURL = os.Getenv("INTERACTION_SERVICE_URL")

		if err != nil {
			panic(err)
		}

		config.Rabbit.Host = os.Getenv("RABBITMQ_HOST")
		config.Rabbit.Port, _ = strconv.Atoi(os.Getenv("RABBITMQ_PORT"))
		config.Rabbit.User = os.Getenv("RABBITMQ_USER")
		config.Rabbit.Password = os.Getenv("RABBITMQ_PASSWORD")

		config.Minio.MinioEndpoint = os.Getenv("MINIO_ENDPOINT")
		config.Minio.MinioAccessKey = os.Getenv("MINIO_ACCESS_KEY")
		config.Minio.MinioSecretKey = os.Getenv("MINIO_SECRET_KEY")
		config.Minio.MinioUseSSL, _ = strconv.ParseBool(os.Getenv("MINIO_USE_SSL"))
		config.Minio.MinioBucketName = os.Getenv("MINIO_BUCKET_NAME")
	})

	return config
}
