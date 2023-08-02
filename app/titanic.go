package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

type Person struct {
	Name                    string  `json:"name"`
	Sex                     string  `json:"sex"`
	Age                     int     `json:"age"`
	PassengerClass          int     `json:"passengerClass"`
	SiblingsOrSpousesAboard int     `json:"siblingsOrSpousesAboard"`
	ParentsOrChildrenAboard int     `json:"parentsOrChildrenAboard"`
	Fare                    float64 `json:"fare"`
	Survived                bool    `json:"survived"`
	Uuid                    string  `json:"uuid"`
}

func main() {
	initDbSchema()

	r := gin.Default()

	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetSlowTime(10)
	m.Use(r)

	r.GET("/people", getPeople)
	r.GET("/people/:uuid", getPerson)
	r.POST("/people", createPerson)
	r.PUT("/people/:uuid", updatePerson)
	r.DELETE("/people/:uuid", deletePerson)

	r.Run()
	log.Println("Starting Gin server")
}

func getPeople(c *gin.Context) {
	var pp []Person
	db := getDbCon()
	defer db.Close()

	if err := db.Select(&pp, `SELECT * FROM people;`); err != nil {
		log.Fatal(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	log.Printf("Retrieved %d people", len(pp))
	c.JSON(http.StatusOK, pp)
}

func getPerson(c *gin.Context) {
	uuid := c.Param("uuid")

	var p Person

	db := getDbCon()
	defer db.Close()
	statement := `SELECT * FROM people WHERE uuid = $1;`
	if err := db.Get(&p, statement, uuid); err != nil {
		log.Println(err)
		c.Status(http.StatusNotFound)
		return
	}
	//Returns the first person with the given uuid in the resultset
	c.JSON(http.StatusOK, p)
}

func createPerson(c *gin.Context) {
	var p Person
	//Generate uuid in the app instead of a db plugin to allow switching to a postgres-compatible db more easily
	p.Uuid = uuid.New().String()
	if err := c.BindJSON(&p); err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	db := getDbCon()
	defer db.Close()

	statement := `INSERT INTO people (name, sex, age, passengerClass, siblingsOrSpousesAboard, parentsOrChildrenAboard, fare, survived, uuid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
	_, err := db.Exec(statement, p.Name, p.Sex, p.Age, p.PassengerClass, p.SiblingsOrSpousesAboard, p.ParentsOrChildrenAboard, p.Fare, p.Survived, p.Uuid)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	log.Printf("Created person with uuid %s.\n", p.Uuid)
	c.JSON(http.StatusOK, p)
}

func updatePerson(c *gin.Context) {
	uuid := c.Param("uuid")
	var p Person
	if err := c.BindJSON(&p); err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}
	p.Uuid = uuid

	db := getDbCon()
	defer db.Close()

	statement := `UPDATE people SET (name, sex, age, passengerClass, siblingsOrSpousesAboard, parentsOrChildrenAboard, fare, survived, uuid) = ($1, $2, $3, $4, $5, $6, $7, $8, $9) WHERE uuid = $9;`
	_, err := db.Exec(statement, p.Name, p.Sex, p.Age, p.PassengerClass, p.SiblingsOrSpousesAboard, p.ParentsOrChildrenAboard, p.Fare, p.Survived, uuid)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusNotFound)
		return
	}
	log.Printf("Updated person with uuid %s.\n", uuid)
	c.JSON(http.StatusOK, p)
}

func deletePerson(c *gin.Context) {
	uuid := c.Param("uuid")

	db := getDbCon()
	defer db.Close()
	statement := `DELETE FROM people WHERE uuid = $1;`
	_, err := db.Exec(statement, uuid)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusNotFound)
		return
	}
	log.Printf("Deleted person with uuid %s.\n", uuid)
	c.Status(http.StatusOK)
}

func getDbCon() *sqlx.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Println(err)
	}

	conParams := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", conParams)
	if err != nil {
		log.Println(err)
	}
	return db
}

func initDbSchema() {
	//Idempotent to avoid failure if table already exsists
	schema := `CREATE TABLE IF NOT EXISTS people (
    name text,
		sex text,
    age integer,
		passengerClass integer,
		siblingsOrSpousesAboard integer,
		parentsOrChildrenAboard integer,
		fare float,
		survived boolean,
		uuid text);`
	//Use exponential backoff for the initial db connection to prevent a race condition with the db when both containers are deployed at the same
	for i := 0; i < 5; i++ {
		db := getDbCon()
		_, err := db.Exec(schema)
		if err == nil {
			log.Println("Connected to db and ensured table exists.")
			return
		}
		if i == 4 {
			log.Fatalf("Initial connection to db failed after 5 retries: %s\n", err)
		}
		log.Printf("Initial connection to db failed, attempt %d", i+1)
		time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
	}
}
