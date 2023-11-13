package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// mysql> select * from score;
// +-------+------+---------+
// | Overs | Runs | country |
// +-------+------+---------+
// |  45.9 |  324 | NULL    |
// +-------+------+---------+
// 1 row in set (0.00 sec)

// mysql> update score set country='India' where Runs=324;
// Query OK, 1 row affected (0.00 sec)
// Rows matched: 1  Changed: 1  Warnings: 0

// mysql> select * from score;
// +-------+------+---------+
// | Overs | Runs | country |
// +-------+------+---------+
// |  45.9 |  324 | India   |
// +-------+------+---------+
// 1 row in set (0.00 sec)

// mysql> insert into score (Overs, Runs, country) values (40.9, 23, 'Pakisthan')
//     -> ^C
// mysql> insert into score (Overs, Runs, country) values (40.9, 23, 'Pakisthan');
// Query OK, 1 row affected (0.01 sec)

// mysql> insert into score (Overs, Runs, country) values (30, 256, 'SriLanka');
// Query OK, 1 row affected (0.00 sec)

// mysql> select * from score;
// +-------+------+-----------+
// | Overs | Runs | country   |
// +-------+------+-----------+
// |  45.9 |  324 | India     |
// |  40.9 |   23 | Pakisthan |
// |    30 |  256 | SriLanka  |
// +-------+------+-----------+
// 3 rows in set (0.01 sec)

var db *sql.DB

type ScoreTable struct {
	Country string  `json: Country`
	Overs   float32 `json: Overs`
	Runs    int     `json: Runs`
}

type responseBody struct {
	Status string `json: status`
}

func main() {
	// mysql -u root -p mysqlpass
	fmt.Println("This is my first mysql")

	r := mux.NewRouter()
	//Read API
	r.HandleFunc("/data/", handleGetWithQueryParam).Methods("GET")
	//Delete API with country value as query parameter
	r.HandleFunc("/data/", handlerDelete).Methods("DELETE")
	//Create API
	r.HandleFunc("/data/", handlerPostWithRequestBody).Methods("POST")
	// Update API with country value as path variable
	r.HandleFunc("/data/{country}", handleUpdate).Methods("PUT")

	//Test API
	r.HandleFunc("/data/hello/", handleHello).Methods("GET")

	//GET all details
	r.HandleFunc("/data/Getall", handlerGetAll).Methods("GET")

	http.Handle("/", r)

	http.ListenAndServe(":8089", nil)

}

func handleHello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is Hello Database prgm")

}

func handlerGetAll(w http.ResponseWriter, r *http.Request) {

	// curl -X GET "http://127.0.0.1:8089/data/Getall"

	db, err := sql.Open("mysql", "root:mysqlpass@tcp(localhost:3306)/cricket")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Use a prepared statement to insert the value
	rows, err := db.Query("SELECT * FROM score")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var resp1 responseBody
	for rows.Next() {
		var details ScoreTable
		err := rows.Scan(&details.Overs, &details.Runs, &details.Country)
		if err != nil {

			log.Fatal(err)

		}

		fmt.Println(details.Country, details.Overs, details.Runs)
		resp1.Status = resp1.Status + " " + fmt.Sprintf("Country: %s, Overs: %f, Runs: %d}", details.Country, details.Runs, details.Overs)

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp1)

}

func handleGetWithQueryParam(w http.ResponseWriter, r *http.Request) {

	//curl -X GET "http://127.0.0.1:8089/data/?country="Iran""

	//curl -X GET "http://127.0.0.1:8089/data/?country="Pakisthan""
	QueryParams := r.URL.Query()

	name := QueryParams.Get("country")

	db, err := sql.Open("mysql", "root:mysqlpass@tcp(localhost:3306)/cricket")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Use a prepared statement to insert the value
	stmt, err := db.Prepare("SELECT * FROM score WHERE country = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(name)
	if err != nil {
		log.Fatal(err)

	}
	var resp1 responseBody
	resp1.Status = "NO DATA FOUND"
	for rows.Next() {
		var details ScoreTable
		err := rows.Scan(&details.Overs, &details.Runs, &details.Country)
		if err != nil {

			log.Fatal(err)

		}

		fmt.Println("Get details of country: ", name)
		fmt.Println(details.Country, details.Overs, details.Runs)
		resp1.Status = fmt.Sprintf("Country: %s, Overs: %d, Runs: %d}", details.Country, details.Runs, details.Overs)

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp1)

}

func handleUpdate(w http.ResponseWriter, r *http.Request) {

	// curl -X PUT "http://127.0.0.1:8089/data/SriLanka" -d '{"Runs": 278, "Overs": 50.0}'

	var details ScoreTable
	// Extract the path variable "country" from the URL
	vars := mux.Vars(r)
	name := vars["country"]

	err := json.NewDecoder(r.Body).Decode(&details)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:mysqlpass@tcp(localhost:3306)/cricket")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Use a prepared statement to insert the value
	stmt, err := db.Prepare("UPDATE score SET Overs = ?,Runs = ? WHERE country = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var resp1 responseBody
	rows, err := stmt.Exec(details.Overs, details.Runs, name)
	fmt.Println(details.Overs, details.Runs, name)
	result, err := rows.RowsAffected()
	if result != 0 {
		resp1.Status = "Updated"

	} else {
		resp1.Status = "Details not found"

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp1)

	fmt.Println("Updated details of country: ", name)

}

func handlerDelete(w http.ResponseWriter, r *http.Request) {

	// curl -X DELETE "http://127.0.0.1:8089/data/?country="Pakisthan""
	// {"Status":"Details not found"}

	// curl -X DELETE "http://127.0.0.1:8089/data/?country="Pakisthan""
	// {"Status":"Deleted"}

	QueryParams := r.URL.Query()

	name := QueryParams.Get("country")

	db, err := sql.Open("mysql", "root:mysqlpass@tcp(localhost:3306)/cricket")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Use a prepared statement to insert the value
	stmt, err := db.Prepare("DELETE FROM score where country=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var resp1 responseBody
	rows, err := stmt.Exec(name)
	result, err := rows.RowsAffected()
	if result != 0 {
		resp1.Status = "Deleted"

	} else {
		resp1.Status = "Details not found"

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp1)

	fmt.Println("Deleted details of country: ", name)

}

func handlerPostWithRequestBody(w http.ResponseWriter, r *http.Request) {

	// curl -X POST "http://127.0.0.1:8089/data/" -d '{"Country":"Newzealand", "Runs": 278, "Overs": 50.0}'

	// curl -X POST "http://127.0.0.1:8089/data/" -d '{"Country":"Iraq", "Runs": 328, "Overs": 30.0}'

	var details ScoreTable
	err := json.NewDecoder(r.Body).Decode(&details)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:mysqlpass@tcp(localhost:3306)/cricket")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Use a prepared statement to insert the value
	stmt, err := db.Prepare("INSERT INTO score (Overs, Runs, country) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var resp1 responseBody
	// Execute the prepared statement with the variable value
	_, err = stmt.Exec(details.Overs, details.Runs, details.Country)
	if err != nil {
		resp1.Status = "ERROR"

	} else {
		resp1.Status = "OK"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp1)

	fmt.Println("Added details of country: ", details.Country)

}
