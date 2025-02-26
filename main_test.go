package main

import (
	"OrderProcessing/common"
	"OrderProcessing/handler"
	"OrderProcessing/migration"
	"OrderProcessing/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Setup test database
func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Order{}, &models.OrderInformation{})
	return db
}

func setupTestServer() *migration.Server {
	db := setupTestDB()
	handlerService := handler.NewHandler(db)
	server := migration.NewServer(db, handlerService)
	go server.StartWorkerPool(2) // Start 2 workers
	return server
}

// Test creating an order (POST /orders)
func TestCreateOrder(t *testing.T) {
	server := setupTestServer()
	router := server.SetupRouter()

	itemInformation := []common.ItemInfo{
		{
			ItemID:   101,
			Quantity: 2,
		}, {
			ItemID:   102,
			Quantity: 1,
		}}

	order := common.OrderRequest{
		UserID:   1,
		OrderID:  1000,
		ItemInfo: itemInformation,
	}

	body, _ := json.Marshal(order)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NotNil(t, response["order_id"])
}

// Test fetching an order status (GET /orders/:id)
func TestGetOrderStatus(t *testing.T) {
	server := setupTestServer()
	router := server.SetupRouter()

	itemInformation := []common.ItemInfo{
		{
			ItemID:   101,
			Quantity: 2,
		}, {
			ItemID:   102,
			Quantity: 1,
		}}

	order := common.OrderRequest{
		UserID:   1,
		OrderID:  1001,
		ItemInfo: itemInformation,
	}

	body, _ := json.Marshal(order)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var orderResponse common.OrderCreationResponse
	body, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &orderResponse)

	url := "/orders/" + fmt.Sprint(orderResponse.OrderId)
	req, _ = http.NewRequest("GET", url, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)

	assert.Equal(t, "Pending", response["status"])
}

// Test order processing (Queue + Worker)
func TestOrderProcessing(t *testing.T) {
	server := setupTestServer()
	router := server.SetupRouter()
	itemInformation := []common.ItemInfo{
		{
			ItemID:   101,
			Quantity: 2,
		}, {
			ItemID:   102,
			Quantity: 1,
		}}

	order := common.OrderRequest{
		UserID:   1,
		OrderID:  1003,
		ItemInfo: itemInformation,
	}

	body, _ := json.Marshal(order)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var orderResponse common.OrderCreationResponse
	body, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &orderResponse)

	orderQueueRequest := common.OrderQueueRequest{
		OrderRequest: common.OrderRequest{
			UserID:   1,
			ItemInfo: itemInformation,
		},
		OrderId: orderResponse.OrderId,
	}
	// Push to queue and let worker process it
	server.Handler.Queue <- orderQueueRequest
	time.Sleep(2 * time.Second) // Wait for worker to process
	var updatedOrder models.Order
	server.DB.First(&updatedOrder, orderResponse.OrderId)
	assert.Equal(t, "Completed", updatedOrder.Status)
}

// Test fetching metrics (GET /orders/metrics)
func TestGetMetrics(t *testing.T) {
	server := setupTestServer()
	router := server.SetupRouter()

	// Creating some test data
	server.DB.Create(&models.Order{
		UserID:      1,
		OrderID:     1004,
		TotalAmount: 50,
		Status:      "Completed",
	})

	server.DB.Create(&models.Order{
		UserID:      2,
		OrderID:     1005,
		TotalAmount: 30,
		Status:      "Pending",
	})

	req, _ := http.NewRequest("GET", "/orders/metrics", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)

	assert.Equal(t, 1.0, response["Completed"]) // Ensure at least 1 order is Completed
	assert.Equal(t, 1.0, response["Pending"])   // Ensure at least 1 order is Pending
}
