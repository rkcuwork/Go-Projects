package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-postgres/models"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_"github.com/lib/pq"
)

type response struct {
	ID      int64 `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("middleware :: middleware :: CreateConnection :: Error loading.env file :: %v", err)
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		// log.Fatalf("middleware :: middleware :: CreateConnection :: sql.open not working")
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatalf("middleware :: middleware :: CreateConnection :: db.Ping not working")
		panic(err)
	}

	fmt.Println("Successfully connected to postgres")

	return db

}


func CreateStock(w http.ResponseWriter, r *http.Request)  {
	var stock models.Stock
	
	err := json.NewDecoder(r.Body).Decode(&stock)
	if err!= nil {
        log.Fatalf("middleware :: middleware :: CreateStock :: Error decoding request body %v", err)
    }

	insertID := insertStock(stock)

	res := response{
		ID:insertID,
        Message: "Stock created Successfully",
	}

	json.NewEncoder(w).Encode(res)
}


func GetAllStock(w http.ResponseWriter, r *http.Request)  {
	stocks, err := getAllStock()

	if err!= nil {
		log.Fatalf("middleware :: middleware :: GetAllStock Error getting stocks %v", err)
     
    }

	json.NewEncoder(w).Encode(stocks)

}

func GetStock(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err!= nil {
        log.Fatalf("middleware :: middleware :: GetStock :: Error parsing id %v", err)
    }

	stock, err := getStock(int64(id))
	if err!= nil {
        log.Fatalf("middleware :: middleware :: GetStock :: Error getting stock %v", err)
    }

	json.NewEncoder(w).Encode(stock)
}

func UpdateStock(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err!= nil {
        log.Fatalf("middleware :: middleware :: UpdateStocks :: Error parsing id %v", err)
    }

	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)
	if err!= nil {
        log.Fatalf("middleware :: middleware :: UpdateStocks :: Error decoding request body %v", err)
    }

	updatedRows, err := updateStock(int64(id), stock)
	if err!= nil {
        log.Fatalf("middleware :: middleware :: UpdateStocks :: Error updating stock %v", err)
    }

	msg := fmt.Sprintf("stock updated successfully. Total rows affected: %v", updatedRows)

	res := response{
        ID:int64(id),
        Message: msg,
    }

	json.NewEncoder(w).Encode(res)

}

func DeleteStock(w http.ResponseWriter, r *http.Request)  {
    params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err!= nil {
        log.Fatalf("middleware :: middleware :: DeleteStock :: Error parsing id %v", err)
    }
	deletedRows, err := deleteStock(int64(id))
	if err!= nil {
        log.Fatalf("middleware :: middleware :: DeleteStock :: Error deleting stock %v", err)
    }
	msg := fmt.Sprintf("stock deleted successfully. Total rows affected: %v", deletedRows)
	res := response{
        ID:int64(id),
        Message: msg,
    }
	json.NewEncoder(w).Encode(res)
}

func insertStock(stock models.Stock) int64 {				
	db := CreateConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO stocks(name, price, company) VALUES($1, $2, $3) RETURNING stockid`

	var id int64

	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err!= nil {
        log.Fatalf("middleware :: middleware :: insertStock :: Error inserting stock %v", err)
    }

	return id
}

func getAllStock() ([]models.Stock, error) {
	db := CreateConnection()
	defer db.Close()

	var stocks []models.Stock
	sqlStatement := `SELECT * FROM stocks`
	rows, err := db.Query(sqlStatement)

	if err!= nil {
        log.Fatalf("middleware :: middleware :: getAllStock :: Error getting stocks %v", err)
    }

	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
        err := rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

        if err!= nil {
            log.Fatalf("middleware :: middleware :: getAllStock :: Error getting stocks %v", err)
        }

        stocks = append(stocks, stock)
	}

	return stocks, err
}

func getStock(id int64) (models.Stock, error) {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := `SELECT * FROM stocks where stockid = $1`
	var stock models.Stock
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	switch err{
	case sql.ErrNoRows:
		fmt.Println("middleware :: middleware :: getStock :: No rows returned")
	case nil:
		return stock, nil
	default:
		log.Fatalf("middleware :: middleware :: getStock :: unable to scan rows %v", err)
	}

	return stock, nil
}

func updateStock(id int64, stock models.Stock) (int64, error) {
    // create the postgres db connection
	db := CreateConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)

	if err != nil {
		log.Fatalf("middleware :: middleware :: updateStock :: Error updating stock %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("middleware :: middleware :: updateStock :: Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected,err
}

func deleteStock(id int64) (int64,error) {
	db := CreateConnection()
	defer db.Close()


	sqlStatement := `DELETE FROM stocks WHERE stockid=$1`
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("middleware :: middleware :: deleteStock :: Error deleting stock %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("middleware :: middleware :: deleteStock :: Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected,err
}