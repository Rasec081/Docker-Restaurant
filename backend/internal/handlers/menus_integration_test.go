package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"restaurant-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func TestMenusIntegration(t *testing.T) {
	setupTestData(t)
	initIntegrationHandlers()

	router := gin.Default()
	router.POST("/menus", CreateMenu)
	router.GET("/menus/:id", GetMenu)
	router.PUT("/menus/:id", UpdateMenu)
	router.DELETE("/menus/:id", DeleteMenu)

	menu := models.Menu{Nombre: "Integration Dish", Precio: 11.5, RestaurantID: 1}
	body, _ := json.Marshal(menu)
	createReq, _ := http.NewRequest("POST", "/menus", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	if createW.Code != http.StatusCreated {
		t.Fatalf("Esperaba %d, obtuvo %d", http.StatusCreated, createW.Code)
	}

	var createResp map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createResp)
	idFloat, ok := createResp["id"].(float64)
	if !ok {
		t.Fatal("respuesta sin id")
	}
	menuID := int(idFloat)

	getReq, _ := http.NewRequest("GET", "/menus/"+itoa(menuID), nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)
	if getW.Code != http.StatusOK {
		t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, getW.Code)
	}

	menu.Nombre = "Integration Dish Updated"
	body, _ = json.Marshal(menu)
	updateReq, _ := http.NewRequest("PUT", "/menus/"+itoa(menuID), bytes.NewBuffer(body))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	router.ServeHTTP(updateW, updateReq)
	if updateW.Code != http.StatusOK {
		t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, updateW.Code)
	}

	deleteReq, _ := http.NewRequest("DELETE", "/menus/"+itoa(menuID), nil)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)
	if deleteW.Code != http.StatusOK {
		t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, deleteW.Code)
	}
}
