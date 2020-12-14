package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/streadway/amqp"
)

//var lrs_host = "localhost" // "10.0.1.51", localhost
//var lrs_port = "3001" // 3001 , 80
//var redirect_url = "http://" + lrs_host + ":" + lrs_port + "/redirect/"
var redirect_url = "https://3cxqoqa6ci.execute-api.us-east-1.amazonaws.com/prod" + "/redirect/"

//replace with correct username password and host
var mysql_connect = "cmpe281:cmpe281bitly@tcp(mysql.cjnv9mdkfify.us-east-1.rds.amazonaws.com:3306)/bitly"
//var mysql_connect = "cmpe281:cmpe281@tcp(localhost:3306)/bitly"

// RabbitMQ Config
var rabbitmq_server = "10.0.1.224" // "10.0.1.224", localhost
var rabbitmq_port = "5672"
var rabbitmq_queue = "createlink"
var rabbitmq_user = "guest"
var rabbitmq_pass = "guest"

func init() {
	testSQL()
	testRabbitMQ()
}

//Check if mysql backend is reachable
func testSQL() {
	db, err := sql.Open("mysql", mysql_connect)
	if err != nil {
		log.Fatal(err)
	} else {
		var (
			original_link string
			short_link    string
		)
		rows, err := db.Query("select original_link, short_link from link")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&original_link, &short_link)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(original_link, short_link)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
	defer db.Close()
}

//Check if rabbitmq is reachable
func testRabbitMQ() {
	conn, err := amqp.Dial("amqp://"+rabbitmq_user+":"+rabbitmq_pass+"@"+rabbitmq_server+":"+rabbitmq_port+"/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
}


// Helper Functions
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	mx.Use(mux.CORSMethodMiddleware(mx))
	n.UseHandler(mx)
	return n
}

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/", pingHandler(formatter)).Methods("GET", http.MethodOptions)
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET", http.MethodOptions)
	mx.HandleFunc("/testencode", encodeHandler(formatter)).Methods("GET", http.MethodOptions)
	mx.HandleFunc("/links/create", createShortLink(formatter)).Methods("POST", http.MethodOptions)
	mx.HandleFunc("/links/{url}", getShortLink(formatter)).Methods("GET", http.MethodOptions)
	mx.HandleFunc("/links/{url}", deleteShortlink(formatter)).Methods("DELETE", http.MethodOptions)
	mx.HandleFunc("/links", getShortLink(formatter)).Methods("GET", http.MethodOptions)
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Control panel active!"})
	}
}

func encodeHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		url := encode_url("https://github.com/nguyensjsu/cmpe281-jayashree-sridhar-16/tree/bitly/cloud")
		formatter.JSON(w, http.StatusOK, url)
	}
}

func createShortLink(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if req.Method == http.MethodOptions {
	        return
	    }


		var l bitly
		_ = json.NewDecoder(req.Body).Decode(&l)
		fmt.Println("Shortening the url: ", l.Original_url)

		l.Short_url = encode_url(l.Original_url)
		fmt.Println("short url:", l.Short_url)

		var (
			original_url string
			short_url    string
		)
		db, err := sql.Open("mysql", mysql_connect)
		defer db.Close()
		stmt, err := db.Prepare("insert into link(original_link, short_link) values(?,?)")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(l.Original_url, l.Short_url)
		if err != nil {
			log.Fatal(err)
		}

		if err != nil {
			log.Fatal(err)
		} else {
			rows, _ := db.Query("select original_link, short_link from link where original_link = ?", l.Original_url)
			defer rows.Close()
			for rows.Next() {
				rows.Scan(&original_url, &short_url)
				log.Println(original_url, short_url)
			}
		}

		result := bitly{
			Original_url: original_url,
			Short_url:    short_url,
			Redirect_url: redirect_url + short_url,
		}

		queue_send(result)

		fmt.Println("Bitly link:", result)
		formatter.JSON(w, http.StatusOK, result)

	}
}

func getShortLink(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		params := mux.Vars(req)
		var url string = params["url"]

		//Connect to mysql db
		db, _ := sql.Open("mysql", mysql_connect)
		defer db.Close()

		var (
			original_url string
			short_url    string
		)

		if url == "" {
			fmt.Println("Stored links:", links)
			links_array := make([]bitly, 0)
			rows, err := db.Query("select * from link")
			defer rows.Close()
			if err != nil {
				log.Fatal(err)
			} else {
				for rows.Next() {
					rows.Scan(&original_url, &short_url)
					log.Println(original_url, short_url)
					var l = bitly{
						Original_url: original_url,
						Short_url:    short_url,
						Redirect_url: redirect_url + short_url,
					}
					links_array = append(links_array, l)
				}
			}
			formatter.JSON(w, http.StatusOK, links_array)
		} else {
			rows, err := db.Query("select * from link where short_link=?", url)
			defer rows.Close()
			if err != nil {
				log.Fatal(err)
			} else {				
				if rows.Next() {
					rows.Scan(&original_url, &short_url)
					log.Println(original_url, short_url)
					var l = bitly{
						Original_url: original_url,
						Short_url:    short_url,
						Redirect_url: redirect_url + short_url,
					}
					fmt.Println("Bitly link: ", l)
					formatter.JSON(w, http.StatusOK, l)
				} else {
					formatter.JSON(w, http.StatusNotFound, struct{ Message string }{"Link not found!"})
				}
			}

			
		}

	}
}


func deleteShortlink(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		params := mux.Vars(req)
		var url string = params["url"]

		//Connect to mysql db
		db, _ := sql.Open("mysql", mysql_connect)
		defer db.Close()

		if url != "" {
			stmt, err := db.Prepare("delete from link where short_link =?")
			if err != nil {
				log.Fatal(err)
			}
			_, err = stmt.Exec(url)
			if err != nil {
				log.Fatal(err)
			}
			formatter.JSON(w, http.StatusNoContent, struct{ Test string }{"Url deleted!"})
		}
	}
}


// Send Order to Queue for Processing
func queue_send(message bitly) {
	conn, err := amqp.Dial("amqp://"+rabbitmq_user+":"+rabbitmq_pass+"@"+rabbitmq_server+":"+rabbitmq_port+"/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		rabbitmq_queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body, _ := json.Marshal(message)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}

