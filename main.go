package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
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

func main() {
	// mysql -u root -p mysqlpass
	fmt.Println("This is my first mysql")

	db, err := sql.Open("mysql", "root:mysqlpass@tcp(localhost:3306)/cricket")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT * FROM score")
	fmt.Println("This is Teju db")
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
		fmt.Println(country, Runs, Overs)
	}
	defer db.Close()
}
