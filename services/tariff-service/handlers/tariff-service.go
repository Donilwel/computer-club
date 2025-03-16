package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"services/utils"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Моковые данные (заменить на БД)
var items = []Item{
	{ID: 1, Name: "Item 1"},
	{ID: 2, Name: "Item 2"},
}

// Получить все элементы
func GetItems(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(items)
}

// Получить элемент по ID
func GetItemByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	for _, item := range items {
		if fmt.Sprintf("%d", item.ID) == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	utils.WriteError(w, http.StatusNotFound, "Item not found")
}

// Создать новый элемент
func CreateItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	json.NewDecoder(r.Body).Decode(&newItem)
	newItem.ID = len(items) + 1
	items = append(items, newItem)
	json.NewEncoder(w).Encode(newItem)
}

// Обновить элемент
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	for i, item := range items {
		if fmt.Sprintf("%d", item.ID) == id {
			json.NewDecoder(r.Body).Decode(&items[i])
			items[i].ID = item.ID
			json.NewEncoder(w).Encode(items[i])
			return
		}
	}
	utils.WriteError(w, http.StatusNotFound, "Item not found")
}

// Удалить элемент
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	for i, item := range items {
		if fmt.Sprintf("%d", item.ID) == id {
			items = append(items[:i], items[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	utils.WriteError(w, http.StatusNotFound, "Item not found")
}
