package main

import (
	"event/config"
	"event/db"
	"event/genproto"
	"event/logger"
	"event/service"
	"event/storage"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"
)

func main() {

	c := config.Load()

	log := logger.New(c.LogLevel, "event service")
	defer func(l logger.Logger) {
		err := logger.Cleanup(l)
		if err != nil {
			log.Fatal("failed cleanup logger", logger.Error(err))
		}
	}(log)

	log.Info("main: sqlxConfig",
		logger.String("host", c.PostgresHost),
		logger.Int("port", c.PostgresPort),
		logger.String("database", c.PostgresDatabase))

	connDb, err := db.ConnectionToDB(c)
	if err != nil {
		log.Fatal("cant conect to db")
	}
	pgStorage := storage.NewEventConnectSql(connDb)

	service := service.NewEventService(pgStorage, log)
	lis, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal("cant listen the port")
	}

	s := grpc.NewServer()

	genproto.RegisterEventServiceServer(s, service)
	reflection.Register(s)
	log.Info("main: server running",
		logger.String("port", c.RPCPort))
	if err := s.Serve(lis); err != nil {
		log.Fatal("cant serve")
	}
}
