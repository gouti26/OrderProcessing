package handler

import (
	"OrderProcessing/common"
	"OrderProcessing/models"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB    *gorm.DB
	Queue chan common.OrderQueueRequest
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		DB:    db,
		Queue: make(chan common.OrderQueueRequest, 1000),
	}
}

// func (h *Handler) getOrderAmount(orderRequest common.OrderRequest) float64 {
// 	itemIds := []uint64{}
// 	for _, itemInfo := range orderRequest.ItemInfo {
// 		itemIds = append(itemIds, itemInfo.ItemID)
// 	}

// 	itemPriceInformation := []common.ItemPriceInfo{}
// 	err := h.DB.Table("items").Where("item_id IN (?)", itemIds).Find((&itemPriceInformation)).Error
// 	if err != nil {
// 		errMsg := fmt.Sprintf("item price information db error: %+v", err)
// 		log.Fatal(errMsg)
// 	}

// 	itemPriceMap := make(map[uint64]float64)
// 	for _, itemPriceInfo := range itemPriceInformation {
// 		itemPriceMap[itemPriceInfo.ItemID] = itemPriceInfo.Price
// 	}

// 	totalOrderAmount := 0.0
// 	for _, itemInfo := range orderRequest.ItemInfo {
// 		totalOrderAmount += itemPriceMap[itemInfo.ItemID] * float64(itemInfo.Quantity)
// 	}
// 	return totalOrderAmount

// }

func (h *Handler) CreateOrder(c *gin.Context) {
	var orderRequest common.OrderRequest
	err := c.ShouldBindJSON(&orderRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// call getOrderAmount() function to fetch totalOrderAmount once items are added in the table
	totalOrderAmount := 100.00
	presentTime := time.Now()
	order := models.Order{
		OrderID:     orderRequest.OrderID,
		UserID:      orderRequest.UserID,
		Status:      "Pending",
		CreatedAt:   presentTime,
		TotalAmount: totalOrderAmount,
	}

	err = h.DB.Table("orders").Create(&order).Error
	if err != nil {
		errMsg := fmt.Sprintf("Order creation db error: %+v", err)
		log.Fatal(errMsg)
		return
	}
	fmt.Printf("total items %+v", len(orderRequest.ItemInfo))
	for _, itemInfo := range orderRequest.ItemInfo {
		orderInfo := models.OrderInformation{
			OrderID:  orderRequest.OrderID,
			ItemID:   itemInfo.ItemID,
			Quantity: itemInfo.Quantity,
		}
		err = h.DB.Table("order_informations").Create(&orderInfo).Error
		if err != nil {
			errMsg := fmt.Sprintf("OrderInformation creation db error: %+v", err)
			log.Fatal(errMsg)
			return
		}
	}

	orderQueueRequest := common.OrderQueueRequest{
		OrderId: orderRequest.OrderID,
		OrderRequest: common.OrderRequest{
			UserID:   orderRequest.UserID,
			ItemInfo: orderRequest.ItemInfo,
		},
	}
	h.Queue <- orderQueueRequest
	c.JSON(http.StatusOK, gin.H{"order_id": orderRequest.OrderID})
}

// GetOrderStatus API Handler
func (h *Handler) GetOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var order models.Order
	err := h.DB.First(&order, "order_id = ?", orderID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": order.Status})
}

// GetMetrics API Handler
func (h *Handler) GetMetrics(c *gin.Context) {
	var totalOrders int64
	var avgProcessingTime float64
	var pendingCount, processingCount, completedCount int64

	h.DB.Model(&models.Order{}).Count(&totalOrders)
	h.DB.Model(&models.Order{}).Where("status = ?", "Pending").Count(&pendingCount)
	h.DB.Model(&models.Order{}).Where("status = ?", "Processing").Count(&processingCount)
	h.DB.Model(&models.Order{}).Where("status = ?", "Completed").Count(&completedCount)

	c.JSON(http.StatusOK, gin.H{
		"total_orders":        totalOrders,
		"avg_processing_time": avgProcessingTime,
		"Pending":             pendingCount,
		"Processing":          processingCount,
		"Completed":           completedCount,
	})
}
