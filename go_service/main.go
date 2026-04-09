package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Address is a nested struct inside Order.
type Address struct {
	Street  string `json:"street"  binding:"required"`
	City    string `json:"city"    binding:"required"`
	Country string `json:"country" binding:"required"`
	Zip     string `json:"zip"     binding:"required"`
}

// OrderItem represents one line in an order.
type OrderItem struct {
	ProductID   int     `json:"product_id"   binding:"required,min=1"`
	ProductName string  `json:"product_name" binding:"required"`
	Quantity    int     `json:"quantity"     binding:"required,min=1"`
	UnitPrice   float64 `json:"unit_price"   binding:"required,gt=0"`
}

// Order is the top-level complex structure.
type Order struct {
	ID          int         `json:"id"`
	CustomerID  int         `json:"customer_id"  binding:"required,min=1"`
	Items       []OrderItem `json:"items"        binding:"required,min=1,dive"`
	ShipTo      Address     `json:"ship_to"      binding:"required"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
}

var orders []Order
var nextID = 1

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "go-service"})
	})
	r.POST("/orders", createOrderHandler)
	r.GET("/orders/:id", getOrderHandler)
	r.GET("/orders", listOrdersHandler)
	return r
}

func createOrderHandler(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var total float64
	for _, item := range order.Items {
		total += float64(item.Quantity) * item.UnitPrice
	}

	order.ID = nextID
	nextID++
	order.TotalAmount = total
	order.Status = "pending"
	order.CreatedAt = time.Now().UTC()

	orders = append(orders, order)
	c.JSON(http.StatusCreated, order)
}

func getOrderHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	for _, order := range orders {
		if order.ID == id {
			c.JSON(http.StatusOK, order)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
}

func listOrdersHandler(c *gin.Context) {
	c.JSON(http.StatusOK, orders)
}

func main() {
	r := SetupRouter()
	r.Run(":8082")
}
