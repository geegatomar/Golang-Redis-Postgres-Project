// Note that these tests must be run only when the main server is up and running
// These tests are only to check if the apis are working end-to-end with our current database
// and cache. The tests involving a mock db to only test a part (ex- api layer) has not been done here.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type Response struct {
	ID              uint           `json:"ID"`
	CreatedAt       time.Time      `json:"CreatedAt"`
	UpdatedAt       time.Time      `json:"UpdatedAt"`
	DeletedAt       gorm.DeletedAt `json:"DeletedAt"`
	TaskId          string         `json:"taskId"`
	TaskDescription string         `json:"taskDescription"`
}

func createPostBody(id string, desc string) *bytes.Buffer {
	postBody, _ := json.Marshal(map[string]string{
		"taskId":          id,
		"taskDescription": desc})
	return bytes.NewBuffer(postBody)
}

func createTodoTask(id string, desc string) {
	resp, err := http.Post("http://localhost:8081/todos", "application/json", createPostBody(id, desc))
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
}

func setupSuite() {
	createTodoTask("T001", "First Task: Eat breakfast")
}

func deleteTodoTask(id string) {
	client := &http.Client{}
	// Create request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8081/todo/%s", id), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

}

func tearDownSuite() {
	deleteTodoTask("T001")
}

func TestHandleGetTodo(t *testing.T) {
	// In the setup suite, we are adding certain todos to the database, so that we can
	// demonstrate each of these apis. And in the teardown suite, we remove them from the db.
	setupSuite()
	defer tearDownSuite()
	response, err := http.Get("http://localhost:8081/todo/T001")
	if err != nil {
		panic("Get request to /todo/taskId failed")
	}
	body, _ := ioutil.ReadAll(response.Body)
	r := new(Response)
	json.Unmarshal(body, &r)
	assert.Equal(t, response.StatusCode, 200, "Asserting status code")
	assert.Equal(t, r.TaskId, "T001", "Asserting taskId")
}

func TestHandleGetAllTodos(t *testing.T) {
	setupSuite()
	defer tearDownSuite()
	response, err := http.Get("http://localhost:8081/todos")
	if err != nil {
		panic("Get request to /todos failed")
	}
	body, _ := ioutil.ReadAll(response.Body)
	var r []Response
	json.Unmarshal(body, &r)

	found := false
	for _, resp := range r {
		if resp.TaskId == "T001" {
			found = true
		}
	}
	assert.Equal(t, response.StatusCode, 200, "Asserting status code")
	assert.Equal(t, found, true, "Asserting taskId present in list of task elements")
}

func TestHandleCreateTodo(t *testing.T) {
	response, err := http.Post("http://localhost:8081/todos", "application/json", createPostBody("T003", "Get an internship!"))
	if err != nil {
		panic("Post request to /todos failed")
	}
	body, _ := ioutil.ReadAll(response.Body)
	r := new(Response)
	json.Unmarshal(body, &r)
	assert.Equal(t, response.StatusCode, 200, "Asserting status code")
	assert.Equal(t, r.TaskId, "T003", "Asserting taskId")
}

func TestHandleUpdateTodo(t *testing.T) {
	setupSuite()
	defer tearDownSuite()
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPut, "http://localhost:8081/todo/T001", createPostBody("T001", "Updated task"))
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	r := new(Response)
	json.Unmarshal(body, &r)
	assert.Equal(t, resp.StatusCode, 200, "Asserting status code")
	assert.Equal(t, r.TaskDescription, "Updated task", "Asserting taskDescription")
}

func TestHandleDeleteTodo(t *testing.T) {
	setupSuite()
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8081/todo/T001"), nil)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, resp.StatusCode, 200, "Asserting status code")
	assert.Equal(t, string(body), "\"The user with taskId T001 is deleted\"\n", "Asserting response body")
}
