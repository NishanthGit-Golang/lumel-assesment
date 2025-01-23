package services

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"revenue/database"
	"revenue/models"

	"gorm.io/gorm"
)

const BatchSize = 1000

func LoadData(filePath string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV data: %w", err)
	}

	for i, record := range records {
		fmt.Printf("Processing row %d: %v\n", i+1, record)

		if i == 0 {
			continue
		}

		if len(record) != 15 {
			fmt.Printf("Skipping malformed row %d (expected 15 fields, got %d): %v\n", i+1, len(record), record)
			continue
		}

		quantitySold, err := strconv.Atoi(record[7])
		if err != nil {
			fmt.Printf("Invalid Quantity Sold at row %d: %v\n", i+1, err)
			continue
		}

		unitPrice, err := strconv.ParseFloat(record[8], 64)
		if err != nil {
			fmt.Printf("Invalid Unit Price at row %d: %v\n", i+1, err)
			continue
		}

		discount, err := strconv.ParseFloat(record[9], 64)
		if err != nil {
			fmt.Printf("Invalid Discount at row %d: %v\n", i+1, err)
			continue
		}

		shippingCost, err := strconv.ParseFloat(record[10], 64)
		if err != nil {
			fmt.Printf("Invalid Shipping Cost at row %d: %v\n", i+1, err)
			continue
		}

		dateOfSale, err := time.Parse("2006-01-02", record[6])
		if err != nil {
			fmt.Printf("Invalid Date of Sale at row %d: %v\n", i+1, err)
			continue
		}

		order := models.Order{
			ID:            record[0],
			CustomerID:    record[2],
			DateOfSale:    dateOfSale,
			Region:        record[5],
			PaymentMethod: record[11],
			ShippingCost:  shippingCost,
		}

		product := models.Product{
			ID:           record[1],
			Name:         record[3],
			Category:     record[4],
			UnitPrice:    unitPrice,
			QuantitySold: quantitySold,
			Discount:     discount,
		}

		customer := models.Customer{
			ID:      record[2],
			Name:    record[12],
			Email:   record[13],
			Address: record[14],
		}

		if err := saveToDatabase(order, product, customer); err != nil {
			fmt.Printf("Error saving row %d: %v\n", i+1, err)
		}
	}

	return nil
}

func saveToDatabase(order models.Order, product models.Product, customer models.Customer) error {
	tx := database.DB.Begin()

	// Check if the customer exists
	var existingCustomer models.Customer
	if err := tx.Where("id = ?", customer.ID).First(&existingCustomer).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return fmt.Errorf("failed to check if customer exists: %w", err)
	}

	if existingCustomer.ID == "" { // Customer does not exist, insert new customer
		if err := tx.Create(&customer).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to save customer: %w", err)
		}
	} else { // Customer exists, skip insert or update customer details if needed
		// Optionally update customer details if necessary
		// e.g., existingCustomer.Name = customer.Name
		if err := tx.Save(&existingCustomer).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update customer: %w", err)
		}
	}

	// Check if the product exists
	var existingProduct models.Product
	if err := tx.Where("id = ?", product.ID).First(&existingProduct).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return fmt.Errorf("failed to check if product exists: %w", err)
	}

	if existingProduct.ID == "" { // Product does not exist, insert new product
		if err := tx.Create(&product).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to save product: %w", err)
		}
	} else { // Product exists, update it
		product.QuantitySold += existingProduct.QuantitySold // Example of updating quantity_sold
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update product: %w", err)
		}
	}

	// Save the order
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save order: %w", err)
	}

	return tx.Commit().Error
}
