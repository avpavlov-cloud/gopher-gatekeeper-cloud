package handlers

import (
	"net/http"
)

func GetAllNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"notes": ["Сходить за хлебом", "Выучить Go", "Настроить Keycloak"]}`))
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status": "Заметка успешно создана"}`))
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	// Сюда попадет только пользователь с ролью 'admin' благодаря middleware
	w.Write([]byte(`{"status": "Заметка удалена администратором"}`))
}
