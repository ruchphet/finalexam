package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	entity "github.com/ruchphet/finalexam/entity"
)

//GetAllCustomer -> get all customer profiles
func GetAllCustomer() (custList []entity.Customer, e error) {
	custList = []entity.Customer{}
	db, err := CreateConnection()
	defer db.Close()
	statement, err := db.Prepare("SELECT id, name, email, status FROM customers")
	if err != nil {
		log.Printf("[Query All] Error can not prepare statement : ", err)
		return custList, err
	}
	rows, err := statement.Query()
	if err != nil {
		log.Printf("[Query All] Error can not query statement : ", err)
		return custList, err
	}

	for rows.Next() {
		var id int
		var name, email, status string
		err := rows.Scan(&id, &name, &email, &status)
		if err != nil {
			log.Printf("[Query Todos] Error Can not query data : %s\n", err)
		}
		customer := entity.Customer{
			ID:     id,
			Name:   name,
			Email:  email,
			Status: status,
		}
		custList = append(custList, customer)
	}

	return custList, err

}

//CreateConnection ...
func CreateConnection() (db *sql.DB, e error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL")) //os.Getenv("DAYABASE_URL")
	log.Printf("Connected to %s\n", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println("Connect DB Error : ", err)
		return db, err
	}
	return db, err
}

//InsertCustomer for insert a  customer
func InsertCustomer(customer entity.Customer) (idResult int, e error) {
	db, _ := CreateConnection()
	defer db.Close()
	row := db.QueryRow("INSERT INTO customers (name, email, status) VALUES ($1, $2, $3) RETURNING id ", customer.Name, customer.Email, customer.Status)
	var id int
	err := row.Scan(&id)
	if err != nil {
		log.Printf("Error Can not insert data : %s\n", err)
		return 0, err
	}
	log.Println("Insert customer completed.")
	return id, err
}

//GetCustomerByID ...
func GetCustomerByID(id int) (customer entity.Customer, e error) {
	db, _ := CreateConnection()
	defer db.Close()
	statement, err := db.Prepare("SELECT id, name, email, status FROM customers WHERE id=$1")
	if err != nil {
		log.Printf("Error Can not Prepare statement : %s\n", err)
		return entity.Customer{}, err
	}
	row := statement.QueryRow(id)
	var returnID int
	var name, email, status string
	err = row.Scan(&returnID, &name, &email, &status)

	if err != nil {
		log.Printf("[Query Todos] Error Can not query data : %s\n", err)
		return entity.Customer{}, err
	}
	return entity.Customer{
		ID:     returnID,
		Name:   name,
		Email:  email,
		Status: status,
	}, nil
}

//UpdateCustomer ...
func UpdateCustomer(customer entity.Customer) (entity.Customer, error) {
	db, _ := CreateConnection()
	defer db.Close()
	var cust entity.Customer
	statement, err := db.Prepare("UPDATE customers SET name=$2 , email=$3, status=$4 WHERE id=$1")
	if err != nil {
		log.Printf("[Update Todos] Error Can not Prepare statement : %s\n", err)
		return cust, err
	}
	if _, err = statement.Exec(customer.ID, customer.Name, customer.Email, customer.Status); err != nil {
		log.Printf("[Update Todos] Error Can not Exec statement : %s\n", err)
		return cust, err
	}

	return GetCustomerByID(customer.ID)
}

//DeleteCustomerByID ...
func DeleteCustomerByID(id int) error {
	db, err := CreateConnection()
	if err != nil {
		log.Printf("[Delete Customer] Error Can not Prepare statement : %s\n", err)
		return err
	}
	defer db.Close()
	statement, err := db.Prepare("DELETE FROM customers WHERE id=$1")
	if err != nil {
		log.Printf("[Delete Customer] Error Can not Prepare statement : %s\n", err)
		return err
	}
	if _, err = statement.Exec(id); err != nil {
		log.Printf("[Delete Customer] Error Can not Exec statement : %s\n", err)
		return err
	}
	return err
}

//CreateCustTable ...
func CreateCustTable() {
	db, _ := CreateConnection()
	createTB := `CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);`
	_, err := db.Exec(createTB)
	if err != nil {
		log.Println("Can not create table : ", err)
	}
	log.Println("Create Table success.")
	defer db.Close()
}
