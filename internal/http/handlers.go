package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/pz6-gorm/internal/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Handlers struct{ db *gorm.DB }

func NewHandlers(db *gorm.DB) *Handlers { return &Handlers{db: db} }

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type createUserReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var in createUserReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Name == "" || in.Email == "" {
		writeErr(w, http.StatusBadRequest, "name and email are required")
		return
	}
	u := models.User{Name: in.Name, Email: in.Email}
	if err := h.db.Create(&u).Error; err != nil {
		writeErr(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, u)
}

// GetUserByID получает пользователя по ID
func (h *Handlers) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		writeErr(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// UpdateUser обновляет пользователя
func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		writeErr(w, http.StatusNotFound, "user not found")
		return
	}

	var in struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if in.Name != nil {
		user.Name = *in.Name
	}
	if in.Email != nil {
		user.Email = *in.Email
	}

	if err := h.db.Save(&user).Error; err != nil {
		writeErr(w, http.StatusConflict, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// DeleteUser удаляет пользователя и все его заметки
func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		writeErr(w, http.StatusNotFound, "user not found")
		return
	}

	// Используем транзакцию для безопасного удаления
	err = h.db.Transaction(func(tx *gorm.DB) error {
		// Находим все заметки пользователя
		var notes []models.Note
		if err := tx.Where("user_id = ?", user.ID).Find(&notes).Error; err != nil {
			return err
		}

		// Для каждой заметки удаляем связи с тегами
		for _, note := range notes {
			if err := tx.Model(&note).Association("Tags").Clear(); err != nil {
				return err
			}
		}

		// Удаляем все заметки пользователя
		if err := tx.Where("user_id = ?", user.ID).Delete(&models.Note{}).Error; err != nil {
			return err
		}

		// Удаляем самого пользователя
		if err := tx.Delete(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "user and all associated notes deleted successfully"})
}

type createNoteReq struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	UserID  uint     `json:"userId"`
	Tags    []string `json:"tags"`
}

func (h *Handlers) CreateNote(w http.ResponseWriter, r *http.Request) {
	var in createNoteReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Title == "" || in.UserID == 0 {
		writeErr(w, http.StatusBadRequest, "title and userId are required")
		return
	}

	// Проверка существования пользователя
	var user models.User
	if err := h.db.First(&user, in.UserID).Error; err != nil {
		writeErr(w, http.StatusBadRequest, "user not found")
		return
	}

	var tags []models.Tag
	for _, name := range in.Tags {
		if name == "" {
			continue
		}
		t := models.Tag{Name: name}
		if err := h.db.FirstOrCreate(&t, models.Tag{Name: name}).Error; err == nil {
			tags = append(tags, t)
		}
	}

	note := models.Note{
		Title:   in.Title,
		Content: in.Content,
		UserID:  in.UserID,
		Tags:    tags,
	}
	if err := h.db.Create(&note).Error; err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.db.Preload("User").Preload("Tags").First(&note, note.ID).Error; err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, note)
}

func (h *Handlers) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	var note models.Note
	if err := h.db.Preload("User").Preload("Tags").First(&note, id).Error; err != nil {
		writeErr(w, http.StatusNotFound, "note not found")
		return
	}
	writeJSON(w, http.StatusOK, note)
}

// UpdateNote обрабатывает обновление заметки
func (h *Handlers) UpdateNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}

	var note models.Note
	if err := h.db.Preload("Tags").First(&note, id).Error; err != nil {
		writeErr(w, http.StatusNotFound, "note not found")
		return
	}

	var in struct {
		Title   *string  `json:"title"`
		Content *string  `json:"content"`
		UserID  *uint    `json:"userId"`
		Tags    []string `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// Обновляем поля если они предоставлены
	if in.Title != nil {
		note.Title = *in.Title
	}
	if in.Content != nil {
		note.Content = *in.Content
	}
	if in.UserID != nil {
		// Проверяем существование пользователя
		var user models.User
		if err := h.db.First(&user, *in.UserID).Error; err != nil {
			writeErr(w, http.StatusBadRequest, "user not found")
			return
		}
		note.UserID = *in.UserID
	}

	// Обновляем теги если они предоставлены
	if in.Tags != nil {
		var tags []models.Tag
		for _, name := range in.Tags {
			if name == "" {
				continue
			}
			t := models.Tag{Name: name}
			if err := h.db.FirstOrCreate(&t, models.Tag{Name: name}).Error; err == nil {
				tags = append(tags, t)
			}
		}
		// Заменяем теги через ассоциации
		h.db.Model(&note).Association("Tags").Replace(tags)
	}

	if err := h.db.Save(&note).Error; err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Загружаем обновленные связи
	h.db.Preload("User").Preload("Tags").First(&note, note.ID)
	writeJSON(w, http.StatusOK, note)
}

// DeleteNote обрабатывает удаление заметки
func (h *Handlers) DeleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}

	var note models.Note
	if err := h.db.First(&note, id).Error; err != nil {
		writeErr(w, http.StatusNotFound, "note not found")
		return
	}

	// Удаляем связи многие-ко-многим перед удалением заметки
	if err := h.db.Model(&note).Association("Tags").Clear(); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.db.Delete(&note).Error; err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "note deleted successfully"})
}

// helpers
type jsonErr struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, jsonErr{Error: msg})
}
