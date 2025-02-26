package migration

import (
	"OrderProcessing/handler"
	"OrderProcessing/models"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	DB      *gorm.DB
	Mutex   sync.Mutex
	Handler *handler.Handler
}

func NewServer(db *gorm.DB, handlerService *handler.Handler) *Server {
	return &Server{
		DB:      db,
		Handler: handlerService,
	}
}

func (m *Server) InitializeServer() {
	m.InitMigration()
	m.StartWorkerPool(10) // Start 10 workers for order processing
}

func (s *Server) SetupRouter() *gin.Engine {
	r := gin.Default()
	handlerService := s.Handler
	r.POST("/orders", handlerService.CreateOrder)
	r.GET("/orders/:id", handlerService.GetOrderStatus)
	r.GET("/orders/metrics", handlerService.GetMetrics)
	return r
}

func (s *Server) migrate(entity interface{}) (err error) {
	db := s.DB
	if db == nil {
		log.Fatal("DB NIL")
	}

	err = db.AutoMigrate(entity)
	if err != nil {
		log.Fatal("Migration error: ", err)
	}
	return nil
}

func (m *Server) InitMigration() {
	log.Printf("DB table migration started...")

	dbTables := map[string]interface{}{
		"orders":             &models.Order{},
		"items":              &models.Item{},
		"order_informations": &models.OrderInformation{},
		"users":              &models.User{},
	}

	for key, table := range dbTables {
		err := m.migrate(table)
		if err != nil {
			errMsg := fmt.Sprintf("Table: %+v with migration error: %+v", key, err)
			log.Fatal(errMsg)
		}
	}
	log.Printf("DB table migration concluded...")
}

// StartWorkerPool starts background workers
func (s *Server) StartWorkerPool(workerCount int) {
	fmt.Printf("\nReady for order processing\n")
	for range workerCount {
		go s.processOrders()
	}
}

// Worker function to process orders
func (s *Server) processOrders() {
	for order := range s.Handler.Queue {
		time.Sleep(1 * time.Second) // Simulate processing time
		s.Mutex.Lock()
		s.DB.Table("orders").Where("order_id = ?", order.OrderId).Update("status", "Completed")
		s.Mutex.Unlock()
		fmt.Println()
		fmt.Printf("\nOrderID %+v Processed", order.OrderId)
	}
}
