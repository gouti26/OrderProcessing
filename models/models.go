package models

import "time"

type Order struct {
	OrderID     uint64    `gorm:"primaryKey;autoIncrement" json:"order_id"`
	UserID      uint64    `json:"customer_id"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type Item struct {
	ItemID   uint    `gorm:"primaryKey;autoIncrement" json:"item_id"`
	ItemName string  `gorm:"type:varchar(255);not null" json:"item_name"`
	Price    float64 `gorm:"type:decimal(10,2);not null" json:"price"`
}

type OrderInformation struct {
	OrderID  uint64 `gorm:"primaryKey" json:"order_id"`
	ItemID   uint64 `gorm:"primaryKey" json:"item_id"`
	Quantity int    `gorm:"default:1" json:"quantity"`
	Order    Order  `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Item     Item   `gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
}

type User struct {
	ID          uint64     `json:"id" gorm:"primary_key;not null;auto_increment"`
	Name        string     `json:"name" gorm:"type:varchar(200);not null"`
	PhoneNumber string     `json:"phone_number" gorm:"unique;index:phone_number_idx"`
	Address     string     `json:"address"`
	Email       string     `json:"email"`
	LoginTime   time.Time  `json:"login_time"`
	Status      bool       `json:"status" gorm:"default:1;not null"` // 0 -> blocked, 1-> unblocked
	CreatedAt   time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
