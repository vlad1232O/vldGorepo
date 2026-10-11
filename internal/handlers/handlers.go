package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"vldGorepo/internal/models"
)

type PageData struct {
	Title    string
	Expenses []models.Expense
}

type Handler struct {
	store *models.ExpenseStore
	pages map[string]*template.Template
}

func New(store *models.ExpenseStore) *Handler {
	layout := "web/templates/layout.html"
	pages := make(map[string]*template.Template)
	for _, name := range []string{"expenses", "add"} {
		pages[name] = template.Must(
			template.ParseFiles(layout, "web/templates/"+name+".html"),
		)
	}
	return &Handler{store: store, pages: pages}
}

func (h *Handler) render(w http.ResponseWriter, page string, data PageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.pages[page].ExecuteTemplate(w, "layout", data); err != nil {
		log.Println("ошибка шаблона:", err)
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	h.render(w, "expenses", PageData{
		Title:    "Список трат",
		Expenses: h.store.All(),
	})
}

func (h *Handler) NewForm(w http.ResponseWriter, r *http.Request) {
	h.render(w, "add", PageData{Title: "Добавить трату"})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Некорректная сумма", http.StatusBadRequest)
		return
	}

	description := strings.TrimSpace(r.FormValue("description"))
	if description == "" {
		http.Error(w, "Описание не может быть пустым", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		http.Error(w, "Некорректная дата", http.StatusBadRequest)
		return
	}

	h.store.Add(models.Expense{
		Amount:      amount,
		Description: description,
		Date:        date,
	})

	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}
