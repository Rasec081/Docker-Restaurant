package repository

import (
	"errors"
	"restaurant-backend/internal/models"
)

// MockUserRepository
type MockUserRepository struct {
	Users           map[int]*models.User
	UsersByUsername map[string]*models.User
	CreateError     error
	UpdateError     error
	DeleteError     error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users:           make(map[int]*models.User),
		UsersByUsername: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) GetByID(id int) (*models.User, error) {
	if user, ok := m.Users[id]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	if user, ok := m.UsersByUsername[username]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) Create(user *models.User) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	user.ID = len(m.Users) + 1
	m.Users[user.ID] = user
	m.UsersByUsername[user.Nombre] = user
	return nil
}

func (m *MockUserRepository) Update(user *models.User) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	if _, ok := m.Users[user.ID]; !ok {
		return errors.New("user not found")
	}
	m.Users[user.ID] = user
	m.UsersByUsername[user.Nombre] = user
	return nil
}

func (m *MockUserRepository) Delete(id int) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	if _, ok := m.Users[id]; !ok {
		return errors.New("user not found")
	}
	delete(m.Users, id)
	// También eliminar del map por username (requiere búsqueda)
	for username, user := range m.UsersByUsername {
		if user.ID == id {
			delete(m.UsersByUsername, username)
			break
		}
	}
	return nil
}

// MockMenuRepository
type MockMenuRepository struct {
	Menus       map[int]*models.Menu
	CreateError error
	UpdateError error
	DeleteError error
}

func NewMockMenuRepository() *MockMenuRepository {
	return &MockMenuRepository{
		Menus: make(map[int]*models.Menu),
	}
}

func (m *MockMenuRepository) GetByID(id int) (*models.Menu, error) {
	if menu, ok := m.Menus[id]; ok {
		return menu, nil
	}
	return nil, errors.New("menu not found")
}

func (m *MockMenuRepository) Create(menu *models.Menu) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	menu.ID = len(m.Menus) + 1
	m.Menus[menu.ID] = menu
	return nil
}

func (m *MockMenuRepository) Update(menu *models.Menu) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	if _, ok := m.Menus[menu.ID]; !ok {
		return errors.New("menu not found")
	}
	m.Menus[menu.ID] = menu
	return nil
}

func (m *MockMenuRepository) Delete(id int) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	if _, ok := m.Menus[id]; !ok {
		return errors.New("menu not found")
	}
	delete(m.Menus, id)
	return nil
}

// MockRestaurantRepository
type MockRestaurantRepository struct {
	Restaurants map[int]*models.Restaurant
	GetAllError error
	CreateError error
	UpdateError error
	DeleteError error
}

func NewMockRestaurantRepository() *MockRestaurantRepository {
	return &MockRestaurantRepository{
		Restaurants: make(map[int]*models.Restaurant),
	}
}

func (m *MockRestaurantRepository) GetAll() ([]models.Restaurant, error) {
	if m.GetAllError != nil {
		return nil, m.GetAllError
	}
	restaurants := make([]models.Restaurant, 0, len(m.Restaurants))
	for _, r := range m.Restaurants {
		restaurants = append(restaurants, *r)
	}
	return restaurants, nil
}

func (m *MockRestaurantRepository) GetByID(id int) (*models.Restaurant, error) {
	if restaurant, ok := m.Restaurants[id]; ok {
		return restaurant, nil
	}
	return nil, errors.New("restaurant not found")
}

func (m *MockRestaurantRepository) Create(restaurant *models.Restaurant) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	restaurant.ID = len(m.Restaurants) + 1
	m.Restaurants[restaurant.ID] = restaurant
	return nil
}

func (m *MockRestaurantRepository) Update(restaurant *models.Restaurant) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	if _, ok := m.Restaurants[restaurant.ID]; !ok {
		return errors.New("restaurant not found")
	}
	m.Restaurants[restaurant.ID] = restaurant
	return nil
}

func (m *MockRestaurantRepository) Delete(id int) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	if _, ok := m.Restaurants[id]; !ok {
		return errors.New("restaurant not found")
	}
	delete(m.Restaurants, id)
	return nil
}

// MockReservationRepository
type MockReservationRepository struct {
	Reservations      map[int]*models.Reservation
	CreateError       error
	DeleteError       error
	AvailabilityError error
	NextID            int
}

func NewMockReservationRepository() *MockReservationRepository {
	return &MockReservationRepository{
		Reservations: make(map[int]*models.Reservation),
		NextID:       1,
	}
}

func (m *MockReservationRepository) Create(reservation *models.Reservation) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	reservation.ID = m.NextID
	m.NextID++
	m.Reservations[reservation.ID] = reservation
	return nil
}

func (m *MockReservationRepository) Delete(id int) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	if _, ok := m.Reservations[id]; !ok {
		return errors.New("reservation not found")
	}
	delete(m.Reservations, id)
	return nil
}

func (m *MockReservationRepository) GetByID(id int) (*models.Reservation, error) {
	if reservation, ok := m.Reservations[id]; ok {
		return reservation, nil
	}
	return nil, errors.New("reservation not found")
}

func (m *MockReservationRepository) IsTableAvailable(tableID int, fecha string) (bool, error) {
	if m.AvailabilityError != nil {
		return false, m.AvailabilityError
	}
	for _, reservation := range m.Reservations {
		if reservation.TableID == tableID && reservation.Fecha == fecha && reservation.Estado == 1 {
			return false, nil
		}
	}
	return true, nil
}

// MockOrderRepository
type MockOrderRepository struct {
	Orders      map[int]*models.Order
	CreateError error
	NextID      int
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		Orders: make(map[int]*models.Order),
		NextID: 1,
	}
}

func (m *MockOrderRepository) GetByID(id int) (*models.Order, error) {
	if order, ok := m.Orders[id]; ok {
		return order, nil
	}
	return nil, errors.New("order not found")
}

func (m *MockOrderRepository) Create(order *models.Order) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	order.ID = m.NextID
	m.NextID++
	m.Orders[order.ID] = order
	return nil
}
