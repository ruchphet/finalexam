package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/ruchphet/finalexam/db"
	entity "github.com/ruchphet/finalexam/entity"
)

//TOKEN Authen
var TOKEN = "token2019"

func authenticater(c *gin.Context) bool {
	token := c.GetHeader("Authorization")
	return token == TOKEN
}

func authMiddleware(c *gin.Context) {
	if !authenticater(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error"})
		c.Abort()
	}
	c.Next()
}

func createCustomer(c *gin.Context) {
	var customer entity.Customer
	err := c.ShouldBindJSON(&customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ERROR"})
		return
	}
	custID, err := db.InsertCustomer(customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	customerAdded, err := db.GetCustomerByID(custID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, customerAdded)

}

func getCustomerByID(c *gin.Context) {
	paramID := c.Param("id")
	if paramID != "" {
		intID, _ := strconv.Atoi(paramID)
		customer, err := db.GetCustomerByID(intID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, customer)
	}
}

func getAllCustomer(c *gin.Context) {
	custList, err := db.GetAllCustomer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, custList)
}

func updateCustomer(c *gin.Context) {
	paramID := c.Param("id")
	var customer entity.Customer
	err := c.ShouldBindJSON(&customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
		return
	}
	if paramID != "" {
		updatedCust, err := db.UpdateCustomer(customer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
			return
		}
		c.JSON(http.StatusOK, updatedCust)
	}
}

func deleteCustomer(c *gin.Context) {
	paramID := c.Param("id")
	if paramID != "" {
		intID, _ := strconv.Atoi(paramID)
		err := db.DeleteCustomerByID(intID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
	}
}

func main() {
	r := gin.Default()
	r.Use(authMiddleware)
	db.CreateCustTable()

	r.POST("/customers", createCustomer)
	r.GET("customers/:id", getCustomerByID)
	r.GET("/customers", getAllCustomer)
	r.PUT("customers/:id", updateCustomer)
	r.DELETE("customers/:id", deleteCustomer)
	r.Run(":2019")
}
