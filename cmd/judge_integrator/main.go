package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"lms_judge_integrator/internal/db"
	"lms_judge_integrator/internal/handler"
	"lms_judge_integrator/internal/repository"
	"lms_judge_integrator/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

const judge0UrlEnvVarName = "judge.service.url"
const serverPortEnvKey = "server.port"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file")
	}

	dbConnection, err := db.InitDB()
	if err != nil {
		log.Fatalf("Unable to connect to the db: \n%v\n", dbConnection.Config().ConnConfig)
	}
	defer dbConnection.Close()

	repo := repository.NewPostgresRepository(dbConnection)

	judge0URL := os.Getenv(judge0UrlEnvVarName)
	judgeService := service.NewJudgeService(repo, judge0URL)

	submissionHandler := handler.NewSubmissionHandler(judgeService)
	mux := http.NewServeMux()
	mux.Handle("POST /api/submissions", submissionHandler)

	srv := &http.Server{
		Addr:    ":" + os.Getenv(serverPortEnvKey),
		Handler: mux,
	}

	fmt.Printf("Judge-Integrator server started at %s", srv.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
}
