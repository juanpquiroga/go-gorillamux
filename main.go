package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:ID`
	Name    string `json:Name`
	Content string `json:Content`
}

type allTasks []task

var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task One",
		Content: "Some content",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a valid task")
	}

	json.Unmarshal(reqBody, &newTask)

	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for _, task := range tasks {
		if task.ID == taskID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Invalid ID")
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")

	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "The task with ID %v has been removed successfully", taskID)
			return
		}
	}
	fmt.Fprintf(w, "Invalid ID")
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")

	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	var updatedTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Enter a valid task")
		return
	}

	json.Unmarshal(reqBody, &updatedTask)
	updatedTask.ID = taskID

	for i, task := range tasks {
		if task.ID == taskID {
			// Eliminar la tarea
			tasks = append(tasks[:i], tasks[i+1:]...)
			// Crea nuevamente la tarea
			tasks = append(tasks, updatedTask)
			json.NewEncoder(w).Encode(updatedTask)
			return
		}
	}
	fmt.Fprintf(w, "Invalid ID")
}

func searchTask(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	key := r.FormValue("key")
	fmt.Fprintf(w, key)
}

func headerTask(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	key := r.Header.Get("X-Liftit")
	fmt.Fprintf(w, key)
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my API en GO")
}

func main() {
	fmt.Println("Hello World")
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	router.Path("/task/search").Queries("key", "{key}").HandlerFunc(searchTask).Name("searchTask").Methods("GET")
	router.Path("/task/search").HandlerFunc(searchTask).Methods("GET")
	router.Path("/task/header").HandlerFunc(headerTask).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", router))
}
