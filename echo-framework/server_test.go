package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetTodos(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/todos", nil) // request for testing
	res := httptest.NewRecorder()                             // recording response
	c := e.NewContext(req, res)

	if assert.NoError(t, getTodos(c)) {
		assert.Equal(t, http.StatusOK, res.Code) // same response codes

		var responseTodos []todo
		err := json.Unmarshal(res.Body.Bytes(), &responseTodos)
		assert.NoError(t, err)
	}
}

func TestAddTodo(t *testing.T) {
	e := echo.New()

	t.Run("Valid Request", func(t *testing.T) {
		validTodo := todo{ID: "1", Item: "Clean Room", Completed: false}
		reqBody, _ := json.Marshal(validTodo)
		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)

		todos = make([]todo, 0)
		err := addTodo(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, res.Code)

		var responseTodo todo
		err = json.Unmarshal(res.Body.Bytes(), &responseTodo)
		assert.NoError(t, err)
		assert.Equal(t, validTodo, responseTodo)

		assert.Len(t, todos, 1)
		assert.Equal(t, validTodo, todos[0])
	})

	t.Run("Invalid Request", func(t *testing.T) {
		invalidTodo := todo{ID: "1", Completed: false}
		reqBody, _ := json.Marshal(invalidTodo)
		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)

		todos = make([]todo, 0)
		err := addTodo(c)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Len(t, todos, 0)
	})
}

func TestGetTodoById(t *testing.T) {
	// Test case: Valid ID, todo exists
	foundTodo, err := getTodoById("2")
	assert.NoError(t, err)
	assert.NotNil(t, foundTodo)
	assert.Equal(t, "Read Book", foundTodo.Item)

	// Test case: Invalid ID, todo not found
	notFoundTodo, err := getTodoById("5")
	assert.Error(t, err)
	assert.Nil(t, notFoundTodo)
	assert.EqualError(t, err, "todo not found")
}

// func TestGetTodo(t *testing.T) {
// 	e := echo.New()

// 	todos := []todo{
// 		{ID: "1", Item: "Clean Room", Completed: false},
// 		{ID: "2", Item: "Read Book", Completed: false},
// 		{ID: "3", Item: "Record Video", Completed: false},
// 	}

// 	req := httptest.NewRequest(http.MethodGet, "/todos/2", nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	c.SetPath("/todos/:id")

// 	// Set the todos slice to the context as if it's coming from the middleware or a prior handler
// 	c.Set("todos", todos)

// 	err := getTodo(c)

// 	if assert.NoError(t, err) {
// 		// Updated status code assertion to account for the case when the todo is not found
// 		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, rec.Code)

// 		if rec.Code == http.StatusOK {
// 			var responseTodo todo
// 			err := json.Unmarshal(rec.Body.Bytes(), &responseTodo)
// 			assert.NoError(t, err)

// 			assert.Equal(t, "Read Book", responseTodo.Item)
// 		}
// 	}
// }
