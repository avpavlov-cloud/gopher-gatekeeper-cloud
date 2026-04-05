package main

import (
	"log"
	"net/http"
	"notes-api/internal/auth"
	"notes-api/internal/handlers"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем системные переменные")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	r := chi.NewRouter()

	// Базовые middleware для всех
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Публичный маршрут
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Service is up"))
	})

	// Группа защищенных маршрутов
	r.Group(func(r chi.Router) {
		// Наш кастомный Middleware, который проверяет токен от Keycloak
		r.Use(auth.KeycloakMiddleware)

		r.Route("/api/v1/notes", func(r chi.Router) {
			r.Get("/", handlers.GetAllNotes) // GET /api/v1/notes
			r.Post("/", handlers.CreateNote) // POST /api/v1/notes

			// Вложенная группа только для админов
			r.Group(func(r chi.Router) {
				r.Use(auth.AdminOnlyMiddleware)
				r.Delete("/{id}", handlers.DeleteNote) // DELETE /api/v1/notes/123
			})
		})
	})

	http.ListenAndServe(":8000", r)
}
