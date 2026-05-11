package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint
	Name      string
	Email     string
	Role      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Order struct {
	ID        uint
	UserID    uint
	Total     float64
	Status    string
	CreatedAt time.Time
}

type Product struct {
	ID        uint
	Name      string
	Category  string
	Price     float64
	Stock     int
	CreatedAt time.Time
}

type OrderItem struct {
	ID        uint
	OrderID   uint
	ProductID uint
	Quantity  int
	UnitPrice float64
}

type AuditLog struct {
	ID         uint
	UserID     uint
	Action     string
	EntityType string
	EntityID   int
	CreatedAt  time.Time
}

// --- Simple queries (existing) ---

func GetUsers(db *gorm.DB, name string) ([]User, error) {
	var users []User
	db.Where("name = ?", name).Find(&users)
	return users, nil
}

func GetAllOrders(db *gorm.DB) ([]Order, error) {
	var orders []Order
	db.Find(&orders)
	return orders, nil
}

// --- Complex queries to test the analyzer ---

// Multi-condition WHERE without indexes on filtered columns
func GetActiveUsersByRole(db *gorm.DB, role string) ([]User, error) {
	var users []User
	db.Where("role = ? AND status = ?", role, "active").Find(&users)
	return users, nil
}

// JOIN query — orders with user info
func GetOrdersWithUsers(db *gorm.DB) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	db.Table("orders").
		Select("orders.*, users.name as user_name, users.email").
		Joins("JOIN users ON users.id = orders.user_id").
		Find(&results)
	return results, nil
}

// Unindexed column filter + ORDER BY
func GetOrdersByStatus(db *gorm.DB, status string) ([]Order, error) {
	var orders []Order
	db.Where("status = ?", status).Order("created_at DESC").Find(&orders)
	return orders, nil
}

// Raw SQL with JOIN across 3 tables — no pagination
func GetOrderDetails(db *gorm.DB, orderID uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	db.Raw(`
		SELECT oi.quantity, oi.unit_price, p.name as product_name, p.category
		FROM order_items oi
		JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id = ?
	`, orderID).Scan(&results)
	return results, nil
}

// SELECT * on a large table with no WHERE
func GetAllProducts(db *gorm.DB) ([]Product, error) {
	var products []Product
	db.Find(&products)
	return products, nil
}

// Filter on unindexed column
func GetProductsByCategory(db *gorm.DB, category string) ([]Product, error) {
	var products []Product
	db.Where("category = ?", category).Find(&products)
	return products, nil
}

// Aggregate query without index on grouped column
func GetOrderTotalsByUser(db *gorm.DB) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	db.Table("orders").
		Select("user_id, SUM(total) as total_spent, COUNT(*) as order_count").
		Group("user_id").
		Find(&results)
	return results, nil
}

// Raw SQL with subquery
func GetUsersWithHighOrders(db *gorm.DB, minTotal float64) ([]User, error) {
	var users []User
	db.Raw(`
		SELECT u.* FROM users u
		WHERE u.id IN (
			SELECT o.user_id FROM orders o
			WHERE o.total > ?
			GROUP BY o.user_id
		)
	`, minTotal).Scan(&users)
	return users, nil
}

// Audit log query — filter on multiple unindexed columns
func GetAuditLogs(db *gorm.DB, entityType string, userID uint) ([]AuditLog, error) {
	var logs []AuditLog
	db.Where("entity_type = ? AND user_id = ?", entityType, userID).
		Order("created_at DESC").
		Find(&logs)
	return logs, nil
}

// LIKE query — always causes full scan
func SearchUsers(db *gorm.DB, query string) ([]User, error) {
	var users []User
	db.Where("name LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%").Find(&users)
	return users, nil
}