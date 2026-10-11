package main

import (
	"log"
	"net/http"

	"vldGorepo/internal/handlers"
	"vldGorepo/internal/models"
)

func main() {
	store := models.NewExpenseStore()
	h := handlers.New(store)

	mux := http.NewServeMux()

	// статика: /static/style.css -> web/static/style.css
	mux.Handle("GET /static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("web/static"))))

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/expenses", http.StatusSeeOther)
	})
	mux.HandleFunc("GET /expenses", h.List)
	mux.HandleFunc("GET /expenses/new", h.NewForm)
	mux.HandleFunc("POST /expenses", h.Create)

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
