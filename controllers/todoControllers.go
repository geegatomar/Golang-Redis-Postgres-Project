package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/geegatomar/todo/models"
	"github.com/gorilla/mux"
)

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todos []models.ToDo
	models.DB.Find(&todos)
	json.NewEncoder(w).Encode(todos)
}

// Get the ToDo task based on the taskId specified in the url params
func GetTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var todo models.ToDo

	// If the element is found in the redis cache, directly return it
	res := models.GetFromCache(models.REDIS, params["taskId"])
	if res != nil {
		fmt.Println("Obtained from redis cache")
		io.WriteString(w, res.(string))
		return
	}
	fmt.Println("Element not found in redis cache")
	models.DB.First(&todo, "task_id = ?", params["taskId"])
	fmt.Println("Setting element in cache")

	// Set element in the redis cache before returning the result
	models.SetInCache(models.REDIS, todo.TaskId, todo)
	json.NewEncoder(w).Encode(todo)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo models.ToDo
	json.NewDecoder(r.Body).Decode(&todo)
	fmt.Print(todo)
	models.DB.Create(&todo)

	// Set the element in the redis cache when user creates a new entry
	models.SetInCache(models.REDIS, todo.TaskId, todo)
	json.NewEncoder(w).Encode(todo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var todo models.ToDo
	models.DB.First(&todo, "task_id = ?", params["taskId"])
	json.NewDecoder(r.Body).Decode(&todo)
	models.DB.Save(&todo)

	// When element gets updated, we make sure to update the redis cache entry as well
	models.SetInCache(models.REDIS, todo.TaskId, todo)
	json.NewEncoder(w).Encode(todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var todo models.ToDo
	models.DB.Delete(&todo, "task_id = ?", params["taskId"])

	// Deleting the element from the redis cache as well
	models.DeleteFromCache(models.REDIS, params["taskId"])
	json.NewEncoder(w).Encode(fmt.Sprintf("The user with taskId %s is deleted", params["taskId"]))
}
