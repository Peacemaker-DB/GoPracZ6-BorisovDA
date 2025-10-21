package httpapi

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func BuildRouter(d *gorm.DB) *chi.Mux {
	r := chi.NewRouter()
	h := NewHandlers(d)

	r.Get("/health", h.Health)

	// Пользователи (упрощённо)
	r.Post("/users", h.CreateUser)
	r.Get("/users/{id}", h.GetUserByID) 
	r.Delete("/users/{id}", h.DeleteUser)
	r.Put("/users/{id}", h.UpdateUser)
	// Заметки
	r.Post("/notes", h.CreateNote)      // создаём заметку с тегами
	r.Get("/notes/{id}", h.GetNoteByID) // получаем заметку с автором и тегами
	r.Delete("/notes/{id}", h.DeleteNote)
	r.Put("/notes/{id}", h.UpdateNote)
	return r
}
