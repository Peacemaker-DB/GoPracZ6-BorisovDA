package main

import (
	"log"
	"net/http"

	"github.com/Peacemaker-DB/GoPracZ6-BorisovDA/internal/db"
	"github.com/Peacemaker-DB/GoPracZ6-BorisovDA/internal/httpapi"
	"github.com/Peacemaker-DB/GoPracZ6-BorisovDA/internal/models"
)

func main() {
	d := db.Connect()

	// Автоматически создаст (или обновит) таблицы под наши модели
	if err := d.AutoMigrate(&models.User{}, &models.Note{}, &models.Tag{}); err != nil {
		log.Fatal("migrate:", err)
	}

	r := httpapi.BuildRouter(d)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
