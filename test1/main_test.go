package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-playground/assert"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Test_getAllCustomers(t *testing.T) {
	//สร้างDB สำหรับ test
	mockDB_ForTest()

	//เชื่อมต่อ Customers_test.db เพื่อทำการ test
	db, err := gorm.Open(sqlite.Open("Customers_test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	//สร้างตัวแปรสำหรับ test
	type Customer struct {
		Name string `json:"name"`
		Age  int    `json:"age" validate:"numeric"`
	}
	tests := []struct {
		description  string
		requestBody  []Customer
		expectStatus int
	}{
		{
			description:  "case : correct input",
			requestBody:  []Customer{{"april", 25}, {"jeff", 23}},
			expectStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			router := Routersetup(db)
			req := httptest.NewRequest("GET", "/customers", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			//แสดงผลลัพธ์ที่ได้มา
			//fmt.Println("resp >>>>", resp.Body.String())
			//fmt.Println("req >>>>", tt.requestBody)

			//ทำให้ resp ที่ได้รับมา ให้อยู่ในรูปแบบเดี๋ยวกับ data ที่ mock ไว้
			var respBody []Customer
			json.Unmarshal([]byte(resp.Body.Bytes()), &respBody)

			// ตรวจสอบว่าข้อมูลใน responseBody ตรงกับ requestBody หรือไม่
			for i, respCustomer := range respBody {
				fmt.Println(respCustomer, tt.requestBody[i])
				assert.Equal(t, respCustomer, tt.requestBody[i])
			}

		})
	}
}
func Test_getCustomerByID(t *testing.T) {
	//เชื่อมต่อ Customers_test.db เพื่อทำการ test
	db, err := gorm.Open(sqlite.Open("Customers_test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	//สร้างตัวแปรสำหรับ test
	type Customer struct {
		Name string `json:"name"`
		Age  int    `json:"age" validate:"numeric"`
	}

	tests := []struct {
		description  string
		requestBody  []Customer
		path         string
		id           string
		expectStatus int
	}{
		{
			description:  "case correct input1",
			requestBody:  []Customer{{"april", 25}},
			path:         "/customers/",
			id:           "1",
			expectStatus: http.StatusOK,
		},
		{
			description:  "case correct input2",
			requestBody:  []Customer{{"jeff", 23}},
			path:         "/customers/",
			id:           "2",
			expectStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			router := Routersetup(db)

			//สร้าง path ที่ต้องการจะทดสอบ เช่น /customer/ + id
			pathApi := tt.path + tt.id
			req := httptest.NewRequest("GET", pathApi, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			//fmt.Println("resp >>>>", resp.Body.String())

			//ทำให้ resp ที่ได้รับมา ให้อยู่ในรูปแบบเดี๋ยวกับ data ที่ mock ไว้
			var respBody []Customer
			json.Unmarshal([]byte(resp.Body.Bytes()), &respBody)
			//fmt.Println("<>", respBody)

			// ตรวจสอบว่าข้อมูลว่าข้อมูลที่่ได้รับ ตรงกับ data ที่ กำหนดไว้หรือไม่
			for i, respCustomer := range respBody {
				fmt.Println(respCustomer, tt.requestBody[i])
				assert.Equal(t, respCustomer, tt.requestBody[i])
			}
		})
	}
}

func Test_addCustomer(t *testing.T) {

	mockDB, _ := gorm.Open(sqlite.Open("Customers_test.db"), &gorm.Config{SkipDefaultTransaction: true})

	tests := []struct {
		description  string
		requestBody  Customer
		expectStatus int
	}{
		{
			description:  "input",
			requestBody:  Customer{"Nuy", 55},
			expectStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Create a new router using the mock DB
			router := Routersetup(mockDB)

			// Convert request body to JSON
			reqBody, _ := json.Marshal(tt.requestBody)

			// Create a new HTTP request
			req := httptest.NewRequest("POST", "/customers/add", bytes.NewReader(reqBody))
			resp := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(resp, req)
			// Print the response body for debugging
			fmt.Println("resp >>>>", resp.Body.String())

			// Assert the response status code
			assert.Equal(t, tt.expectStatus, resp.Code)

		})
	}
}

func Test_deleteCustomer(t *testing.T) {

	mockDB, _ := gorm.Open(sqlite.Open("Customers_test.db"), &gorm.Config{SkipDefaultTransaction: true})

	tests := []struct {
		description  string
		id           string
		expectStatus int
	}{
		{
			description:  "input correc2",
			id:           "3",
			expectStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Create a new router using the mock DB
			router := Routersetup(mockDB)

			// Create a new HTTP request
			req := httptest.NewRequest("DELETE", "/customers/id/"+tt.id, nil)
			resp := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(resp, req)
			// Print the response body for debugging
			//fmt.Println("resp >>>>", resp)

			// Assert the response status code
			assert.Equal(t, tt.expectStatus, resp.Code)

		})
	}
}

func mockDB_ForTest() {
	err := os.Remove("Customers_test.db")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("File deleted successfully")
	// Open or create the SQLite database file
	db, err := sql.Open("sqlite3", "Customers_test.db")
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
