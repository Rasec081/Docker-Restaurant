package models

// User representa un usuario del sistema
type User struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	RoleID int    `json:"role_id"`
}

// Restaurant representa un restaurante
type Restaurant struct {
	ID      int    `json:"id"`
	Nombre  string `json:"nombre"`
	Estado  int    `json:"estado"`
	AdminID int    `json:"admin_id"`
}

// Menu representa un plato en el menú
type Menu struct {
	ID           int     `json:"id"`
	Nombre       string  `json:"nombre"`
	Precio       float64 `json:"precio"`
	RestaurantID int     `json:"restaurant_id"`
}

// Reservation representa una reserva
type Reservation struct {
	ID       int    `json:"id"`
	TableID  int    `json:"table_id"`
	ClientID int    `json:"client_id"`
	Fecha    string `json:"fecha"`
	Estado   int    `json:"estado"`
}

// Order representa una orden
type Order struct {
	ID           int    `json:"id"`
	TableID      *int   `json:"table_id,omitempty"`
	ClientID     int    `json:"client_id"`
	OrdersType   string `json:"orders_type"`
	RestaurantID int    `json:"restaurant_id"`
}

// LoginRequest representa credenciales de login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest representa datos de registro
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
