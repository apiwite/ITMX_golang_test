package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Customer struct {
	Name string `json:"name"`
	Age  int    `json:"age" validate:"numeric"`
}
type CustomerStatus struct {
	Status string `json:"status"`
}

func DBsetup() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("Customers.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Error opening database:", err)
		return nil
	}
	return db
}

func Routersetup(db *gorm.DB) *gin.Engine {

	router := gin.Default()

	router.GET("/customers", func(c *gin.Context) {
		getAllCustomers(c, db)
	})
	router.POST("/customers/add", func(c *gin.Context) {
		addCustomer(c, db)
	})
	router.GET("/customers/:id", func(c *gin.Context) {
		getCustomerByID(c, db)
	})
	router.PUT("/customers/update/:id", func(c *gin.Context) {
		updateCustomer(c, db)
	})
	router.DELETE("/customers/id/:id", func(c *gin.Context) {
		deleteCustomer(c, db)
	})
	router.GET("/resetDB", func(c *gin.Context) {
		resetDB()
	})

	return router
}

func main() {
	db := DBsetup()
	if db == nil {
		fmt.Println("Failed to connect to the database")
		return
	}

	router := Routersetup(db)
	// รันเซิร์ฟเวอร์ที่ port 8080
	router.Run(":8080")
}

func addCustomer(c *gin.Context, db *gorm.DB) {
	var newCustomer Customer
	var validate = validator.New()
	// รับข้อมูลลูกค้าจาก JSON และแปลงเป็นโครงสร้าง Customer
	if err := c.BindJSON(&newCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	}
	if err := validate.Struct(newCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	}

	// เพิ่มลูกค้าใหม่ลงในฐานข้อมูล
	db.Create(&newCustomer)

	// ส่งข้อมูลลูกค้าที่เพิ่มเข้าไปกลับเป็น JSON
	c.JSON(http.StatusOK, "ADD Success")
}

func getAllCustomers(c *gin.Context, db *gorm.DB) []Customer {
	var customers []Customer
	db.Find(&customers)

	// ส่งข้อมูลลูกค้าทั้งหมดกลับเป็น JSON
	fmt.Println("Customer list", customers)
	c.JSON(http.StatusOK, customers)
	return customers
}

func getCustomerByID(c *gin.Context, db *gorm.DB) {
	// รับค่าพารามิเตอร์ name จาก URL
	id := c.Param("id")

	//ค้นหาโดยใช้ name
	var customers []Customer
	if err := db.Where("id = ?", id).Find(&customers).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error db.Where in getCustomerByID "})
		return
	}
	if len(customers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "ID:" + id + " not found"})
		return
	}

	// ส่งข้อมูลลูกค้าที่พบกลับเป็น JSON
	c.JSON(http.StatusOK, customers)
}

func updateCustomer(c *gin.Context, db *gorm.DB) {
	// รับค่า ID จาก URL
	id := c.Param("id")
	var updatedCustomer Customer
	var customer_status CustomerStatus

	// รับข้อมูลลูกค้าที่ต้องการอัปเดตจาก JSON และแปลงเป็นโครงสร้าง Customer
	if err := c.BindJSON(&updatedCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// อัปเดตข้อมูลลูกค้าในฐานข้อมูลโดยอิงจาก ID
	result := db.Model(&Customer{}).Where("id = ?", id).Updates(&updatedCustomer)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"ERROR": " Customer not found"})
		return
	}

	// กำหนดสถานะการอัปเดตเป็น "updated"
	customer_status.Status = "updated"

	// ส่งข้อมูลลูกค้าที่อัปเดตเรียบร้อยแล้วกลับเป็น JSON
	c.JSON(http.StatusOK, customer_status)
}

func deleteCustomer(c *gin.Context, db *gorm.DB) {
	// รับค่า ID จาก URL
	id := c.Param("id")

	// ลบข้อมูลลูกค้าโดยใช้ SQL Query ที่มีเงื่อนไข
	result := db.Exec("DELETE FROM customers WHERE id = ?", id)

	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer ID not found"})
		return
	}
	fmt.Println("DELETE OK")
	// ส่งข้อความยืนยันการลบกลับเป็น JSON
	c.JSON(http.StatusOK, gin.H{"message": "Customer ID : " + id + " deleted successfully"})
}

func resetDB() {

	db, err := sql.Open("sqlite3", "Customers.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// Create a table in the database if it does not exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS customers (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						name TEXT NOT NULL,
						age INTEGER NOT NULL
					)`)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	// Insert customer data into the table
	_, err = db.Exec("INSERT INTO customers (name, age) VALUES (?, ?)", "april", 25)
	if err != nil {
		fmt.Println("Error inserting data:", err)
		return
	}
	_, err = db.Exec("INSERT INTO customers (name, age) VALUES (?, ?)", "jeff", 23)
	if err != nil {
		fmt.Println("Error inserting data:", err)
		return
	}

	fmt.Println("Database created and initialized successfully!")
}
