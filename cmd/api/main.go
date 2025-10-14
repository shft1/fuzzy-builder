package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexm/fuzzy-builder/internal/database/postgresql"
	"github.com/alexm/fuzzy-builder/internal/repositories"
	"github.com/alexm/fuzzy-builder/internal/services"
	rest "github.com/alexm/fuzzy-builder/internal/transport/rest"
)

func main() {
	addr := getEnv("HTTP_ADDR", ":8080")
	dsn := getEnv("DATABASE_URL", "postgresql://localhost:5432/fuzzy_builder")
	jwtSecret := getEnv("JWT_SECRET", "dev-secret-change-me")
	uploadDir := getEnv("UPLOAD_DIR", "uploads")

	ctx := context.Background()
	pool, err := postgresql.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	defer pool.Close()

	usersRepo := repositories.NewUserRepository(pool)
	projectsRepo := repositories.NewProjectRepository(pool)
	defectsRepo := repositories.NewDefectRepository(pool)
	attachRepo := repositories.NewAttachmentRepository(pool)
	hasher := services.NewPasswordHasher()
	jwt := services.NewJWTIssuer(jwtSecret, "fuzzy-builder", 24*time.Hour)

	defectSvc := services.NewDefectService(defectsRepo)
	srv := rest.NewServer(usersRepo, projectsRepo, defectsRepo, attachRepo, defectSvc, uploadDir, hasher, jwt)
	handler := srv.Router()

	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("HTTP server listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	log.Printf("server stopped")
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
