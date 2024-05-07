package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"

	"strconv"
	"sync"
	"syscall"

	"google.golang.org/grpc"

	config "github.com/SPA2006/Distributed-Calculator-of-Arithmetic-Expressions/internal/config"
	database "github.com/SPA2006/Distributed-Calculator-of-Arithmetic-Expressions/internal/database"
	prt_agent "github.com/SPA2006/Distributed-Calculator-of-Arithmetic-Expressions/proto/agent/agent"
)

const (
	envLocal = "local"

	httpPort = ":8081"
	gRPCPort = ":8079"
)

type Storage struct {
	DB *sql.DB
}

func main() {
	var wg sync.WaitGroup
	mux := http.NewServeMux()
	ctx := context.Background()
	database.InitDB()
	logger := setupLogger(envLocal)
	cfg := config.MustLoad()

	logger.Log(ctx, slog.LevelDebug, cfg.Env)
	logger.Log(ctx, slog.LevelDebug, cfg.StoragePath)
	logger.Log(ctx, slog.LevelDebug, strconv.Itoa(cfg.GRPC.Port))
	logger.Log(ctx, slog.LevelDebug, cfg.GRPC.Timeout.String())
	logger.Log(ctx, slog.LevelDebug, cfg.HTTP.Addr)
	logger.Log(ctx, slog.LevelDebug, strconv.Itoa(cfg.HTTP.MaxHeaderBytes))
	logger.Log(ctx, slog.LevelDebug, cfg.HTTP.ReadTimeout.String())
	logger.Log(ctx, slog.LevelDebug, cfg.HTTP.WriteTimeout.String())

	mux.HandleFunc("/register", database.RegisterHandler(ctx))
	mux.HandleFunc("/login", database.LoginHandler(ctx))
	mux.HandleFunc("/expression", database.CreateExpressionHandler(ctx))
	mux.HandleFunc("/get_exp", database.GetExpressionByIDHandler(ctx))
	mux.HandleFunc("/get_user_exp", database.GetExpressionByUserIDHandler(ctx))
	mux.HandleFunc("/get_all_exp", database.GetAllExpressionsHandler(ctx))
	mux.HandleFunc("/delete_exp", database.DeleteExpressionHandler(ctx))

	grpcServer := grpc.NewServer()
	wg.Add(2)

	go func() {
		defer wg.Done()

		lis, err := net.Listen("tcp", gRPCPort)
		if err != nil {
			log.Fatalf("failed to listen %v", err)
		}
		defer lis.Close()
		server := prt_agent.UnimplementedAgentServer{}
		prt_agent.RegisterAgentServer(grpcServer, server)

		log.Println("Starting gRPC Server")
		log.Printf("gRPC server starts litening on port %s", lis.Addr().String())

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to connect to gRPC: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		log.Println("Starting REST server")
		log.Println("HTTP server starts listening on port ", httpPort)
		log.Println(http.ListenAndServe(httpPort, mux))
	}()

	wg.Wait()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logger.Info("Shutting down server")
	var db database.Storage
	db.Closee()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}
