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

func main() {
	// mysql -u root -p mysqlpass
	fmt.Println("This is my first mysql")

	r := mux.NewRouter()
	//Read API
	r.HandleFunc("/data/{country}", handleGetWithQueryParam).Methods("GET")
	//Delete API
	r.HandleFunc("/data/{country}", handlerDelete).Methods("DELETE")
	//Create API
	r.HandleFunc("/data/", handlerPostWithRequestBody).Methods("POST")
	// Update API
	r.HandleFunc("/data/{country}/", handleUpdate).Methods("PUT")

	http.Handle("/", r)

	http.ListenAndServe(":8089", nil)

	db, err := sql.Open("mysql", "root:mysqlpass@tcp(localhost:3306)/cricket")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func handleGetWithQueryParam(w http.ResponseWriter, r *http.Request) {
	QueryParams := r.URL.Query()

	name := QueryParams.Get("country")

	rows, err := db.Query("SELECT * FROM score where country=", name)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var country string
		var Runs int
		var Overs float32
		err := rows.Scan(&Overs, &Runs, &country)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Get details of country: ", name)
		fmt.Println(country, Runs, Overs)
	}

}

func handleUpdate(w http.ResponseWriter, r *http.Request) {

}

func handlerDelete(w http.ResponseWriter, r *http.Request) {

	QueryParams := r.URL.Query()

	name := QueryParams.Get("country")

	// Use a prepared statement to insert the value
	stmt, err := db.Prepare("DELETE FROM score where country=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Execute the prepared statement with the variable value
	_, err = stmt.Exec(name)
	if err != nil {
		log.Fatal(err)
	}

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
