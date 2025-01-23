package models

import "time"

type Customer struct {
	ID      string `gorm:"primaryKey;size:191" json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type Product struct {
	ID           string  `gorm:"primaryKey;size:191" json:"id"`
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	UnitPrice    float64 `json:"unit_price"`
	QuantitySold int     `json:"quantity_sold"`
	Discount     float64 `json:"discount"`
}

type Order struct {
	ID            string    `gorm:"primaryKey;size:191" json:"id"`
	CustomerID    string    `gorm:"size:191;not null" json:"customer_id"`
	Customer      Customer  `gorm:"foreignKey:CustomerID"`
	DateOfSale    time.Time `json:"date_of_sale"`
	Region        string    `json:"region"`
	PaymentMethod string    `json:"payment_method"`
	ShippingCost  float64   `json:"shipping_cost"`
}
