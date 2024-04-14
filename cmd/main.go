package main

import (
	avitotest "avito"
	"avito/internal/cache"
	"avito/internal/handler"
	"avito/internal/repository"
	"avito/internal/service"
	"log"

	_ "github.com/lib/pq"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Написать в README что мы кешируем отдельный запрос

func main() {
  if err:= initConfig(); err!= nil{
		logrus.Fatalf("error initializing config: %s", err.Error())
	}
  // Не забыть настроить конфиг
  db,err :=repository.NewPostgresDB(repository.Config{
   Host: viper.GetString("db.host"),
	 Port: viper.GetString("db.post"),
	 Username: viper.GetString("db.username"),
	 Password: viper.GetString("db.password"),
	 DBname: viper.GetString("db.dbname"),
	 SSLmode: viper.GetString("db.sslmode"),
	})
if err!= nil{
	log.Fatalf("failed to initialized db %s", err.Error())
}
  simple:=cache.NewSimple()
  cache:= cache.NewCache(simple)
  repos:= repository.NewRepository(db)
  services:=service.NewService(repos, cache)
  handlers:= handler.NewHandler(services)
	server:= new(avitotest.Server)
	if err:= server.Run(viper.GetString("port"), handlers.InitRoutes()); err!= nil{
		logrus.Fatalf("error server init %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
  viper.SetConfigName("config")
	return viper.ReadInConfig()
}