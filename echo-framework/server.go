package main

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type todo struct {
	ID        string `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

var todos = []todo{
	{ID: "1", Item: "Clean Room", Completed: false},
	{ID: "2", Item: "Read Book", Completed: false},
	{ID: "3", Item: "Record Video", Completed: false},
}

func getTodos(context echo.Context) error {
	return context.JSON(http.StatusOK, todos)
}

func addTodo(context echo.Context) error {
	var newTodo todo

	if err := context.Bind(&newTodo); err != nil {
		return err
	}

	todos = append(todos, newTodo)
	return context.JSON(http.StatusCreated, newTodo)
}

func getTodoById(id string) (*todo, error) {
	for i, t := range todos {
		if t.ID == id {
			return &todos[i], nil
		}
	}

	return nil, errors.New("todo not found")
}

func getTodo(context echo.Context) error {
	id := context.Param("id")
	todo, err := getTodoById(id)

	if err != nil {
		return context.JSON(http.StatusNotFound, map[string]string{"message": "Todo not found"})
	}

	return context.JSON(http.StatusOK, todo)
}

func main() {
	e := echo.New()
	e.GET("/todos", getTodos)
	e.GET("/todos/:id", getTodo)
	e.POST("/todos", addTodo)
	e.Start("localhost:9090")
}
