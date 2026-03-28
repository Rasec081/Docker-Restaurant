package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

func TestGetMenu_WithMock(t *testing.T) {
	// Crear mock repository
	mockRepo := repository.NewMockMenuRepository()

	// Agregar datos de prueba
	testMenu := &models.Menu{
		ID:           1,
		Nombre:       "Pizza Napolitana",
		Precio:       7500.00,
		RestaurantID: 1,
	}
	mockRepo.Menus[1] = testMenu

	// Inicializar handler con mock
	InitMenuHandler(mockRepo)

	router := setupTestRouter()
	router.GET("/menus/:id", GetMenu)

	tests := []struct {
		name           string
		menuID         string
		expectedStatus int
		expectedName   string
	}{
		{
			name:           "Get existing menu",
			menuID:         "1",
			expectedStatus: http.StatusOK,
			expectedName:   "Pizza Napolitana",
		},
		{
			name:           "Get non-existent menu",
			menuID:         "999",
			expectedStatus: http.StatusNotFound,
			expectedName:   "",
		},
		{
			name:           "Get menu with invalid ID",
			menuID:         "abc",
			expectedStatus: http.StatusBadRequest,
			expectedName:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/menus/"+tt.menuID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Esperaba %d, obtuvo %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedName != "" {
				var response models.Menu
				json.Unmarshal(w.Body.Bytes(), &response)
				if response.Nombre != tt.expectedName {
					t.Errorf("Esperaba %s, obtuvo %s", tt.expectedName, response.Nombre)
				}
			}
		})
	}
}

func TestCreateMenu_WithMock(t *testing.T) {
	mockRepo := repository.NewMockMenuRepository()
	InitMenuHandler(mockRepo)

	router := setupTestRouter()
	router.POST("/menus", CreateMenu)

	tests := []struct {
		name           string
		menu           models.Menu
		expectedStatus int
	}{
		{
			name: "Create valid menu",
			menu: models.Menu{
				Nombre:       "New Dish",
				Precio:       15.99,
				RestaurantID: 1,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Create menu with empty name",
			menu: models.Menu{
				Nombre:       "",
				Precio:       15.99,
				RestaurantID: 1,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Create menu with zero price",
			menu: models.Menu{
				Nombre:       "Free Dish",
				Precio:       0,
				RestaurantID: 1,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock para cada test
			mockRepo.Menus = make(map[int]*models.Menu)

			jsonBody, _ := json.Marshal(tt.menu)
			req, _ := http.NewRequest("POST", "/menus", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Esperaba %d, obtuvo %d", tt.expectedStatus, w.Code)
			}

			// Verificar que se creó correctamente
			if tt.expectedStatus == http.StatusCreated {
				if len(mockRepo.Menus) != 1 {
					t.Error("El menú no fue guardado en el repositorio")
				}
			}
		})
	}
}

func TestUpdateMenu_WithMock(t *testing.T) {
	mockRepo := repository.NewMockMenuRepository()

	// Crear menú inicial
	initialMenu := &models.Menu{
		ID:           1,
		Nombre:       "Original Dish",
		Precio:       10.00,
		RestaurantID: 1,
	}
	mockRepo.Menus[1] = initialMenu

	InitMenuHandler(mockRepo)

	router := setupTestRouter()
	router.PUT("/menus/:id", UpdateMenu)

	tests := []struct {
		name           string
		menuID         string
		menu           models.Menu
		expectedStatus int
	}{
		{
			name:   "Update existing menu",
			menuID: "1",
			menu: models.Menu{
				Nombre:       "Updated Dish",
				Precio:       20.00,
				RestaurantID: 1,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Update non-existent menu",
			menuID: "999",
			menu: models.Menu{
				Nombre:       "Updated Dish",
				Precio:       20.00,
				RestaurantID: 1,
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "Update with invalid ID",
			menuID: "abc",
			menu: models.Menu{
				Nombre:       "Updated Dish",
				Precio:       20.00,
				RestaurantID: 1,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.menu)
			req, _ := http.NewRequest("PUT", "/menus/"+tt.menuID, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Esperaba %d, obtuvo %d", tt.expectedStatus, w.Code)
			}

			// Verificar actualización
			if tt.expectedStatus == http.StatusOK {
				updated, _ := mockRepo.GetByID(1)
				if updated.Nombre != "Updated Dish" {
					t.Error("El menú no fue actualizado correctamente")
				}
			}
		})
	}

	t.Run("Update with invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/menus/1", bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Esperaba %d, obtuvo %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestDeleteMenu_WithMock(t *testing.T) {
	mockRepo := repository.NewMockMenuRepository()

	// Crear menú para eliminar
	testMenu := &models.Menu{
		ID:           1,
		Nombre:       "To Delete",
		Precio:       5.99,
		RestaurantID: 1,
	}
	mockRepo.Menus[1] = testMenu

	InitMenuHandler(mockRepo)

	router := setupTestRouter()
	router.DELETE("/menus/:id", DeleteMenu)

	tests := []struct {
		name           string
		menuID         string
		expectedStatus int
	}{
		{
			name:           "Delete existing menu",
			menuID:         "1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Delete non-existent menu",
			menuID:         "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Delete with invalid ID",
			menuID:         "abc",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/menus/"+tt.menuID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Esperaba %d, obtuvo %d", tt.expectedStatus, w.Code)
			}

			// Verificar que se eliminó
			if tt.expectedStatus == http.StatusOK {
				if len(mockRepo.Menus) != 0 {
					t.Error("El menú no fue eliminado")
				}
			}
		})
	}
}
