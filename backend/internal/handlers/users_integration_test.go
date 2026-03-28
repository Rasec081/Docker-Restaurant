package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUsersIntegration(t *testing.T) {
	setupTestData(t)
	initIntegrationHandlers()

	router := gin.Default()
	router.GET("/users/me", func(c *gin.Context) {
		c.Set("username", "admin")
		GetUserMe(c)
	})
	router.PUT("/users/:id", UpdateUser)
	router.DELETE("/users/:id", DeleteUser)

	getReq, _ := http.NewRequest("GET", "/users/me", nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)
	if getW.Code != http.StatusOK {
		t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, getW.Code)
	}

	updateBody, _ := json.Marshal(map[string]interface{}{"nombre": "Admin Updated", "rol": 1})
	updateReq, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	router.ServeHTTP(updateW, updateReq)
	if updateW.Code != http.StatusOK {
		t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, updateW.Code)
	}

	// Crear usuario extra para eliminar
	var userID int
	testDB.QueryRow("INSERT INTO Users (nombre, role_id) VALUES ('Integration Delete', 2) RETURNING user_id").Scan(&userID)

	deleteReq, _ := http.NewRequest("DELETE", "/users/"+itoa(userID), nil)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)
	if deleteW.Code != http.StatusOK {
		t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, deleteW.Code)
	}
}
