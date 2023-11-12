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
	//Delete API
	r.HandleFunc("/data/", handlerDelete).Methods("DELETE")
	//Create API
	r.HandleFunc("/data/", handlerPostWithRequestBody).Methods("POST")
	// Update API
	r.HandleFunc("/data/{country}/", handleUpdate).Methods("PUT")

	//Test API
	r.HandleFunc("/data/hello/", handleHello).Methods("GET")

	http.Handle("/", r)

	http.ListenAndServe(":8089", nil)

}

func handleHello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is Hello Database prgm")

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
	var details ScoreTable
	err := json.NewDecoder(r.Body).Decode(details)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	// Use a prepared statement to insert the value
	stmt, err := db.Prepare("INSERT INTO score (Overs, Runs, country) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Execute the prepared statement with the variable value
	_, err = stmt.Exec(details.Overs, details.Runs, details.Country)
	if err != nil {
		log.Fatal(err)
	}

}
