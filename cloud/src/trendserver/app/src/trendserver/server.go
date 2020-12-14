package main

import (	
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"strings"
	//"bytes"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/streadway/amqp"
)


// RabbitMQ Config
var rabbitmq_server = "10.0.1.224" // "10.0.1.224", "localhost"
var rabbitmq_port = "5672"
var rabbitmq_queue = "createlink"
var rabbitmq_user = "guest"
var rabbitmq_pass = "guest"

//NoSQL configuration
//var nosql_url = "http://nosql-elb-a35c543da7694aea.elb.us-east-1.amazonaws.com" // replace with load balancer url
var nosql_url = "http://10.0.3.234:9090";

func init() {
	testNosql()
	testRabbitMQ()

	if links == nil {
		links = make(map[string]bitly)
	}
}

func testNosql() {
	_, err := http.Get(nosql_url)
	if err != nil {
	   log.Fatalln(err)
	}
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
	mx.HandleFunc("/links", createLink(formatter)).Methods("GET", http.MethodOptions)
	mx.HandleFunc("/links/{url}", getLink(formatter)).Methods("GET", http.MethodOptions)
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Trend server active!"})
	}
}

//Create short link document in nosql
func createLink(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		if links == nil {
			links = make(map[string]bitly)
		}

		var urls []bitly = queue_receive()

		for _,x := range urls {
			fmt.Printf("%+v\n",x)
			fmt.Printf("%s\n", x.Short_url)
		}

		for _,i := range urls {
			createNosql(i)
		}
		getAllLinksFromNosql()
		//var data []bitly
		data := make([]bitly, 0)
		for _, element := range links {
			fmt.Printf("%+v", element)
			data = append(data, element)
		}
		formatter.JSON(w, http.StatusOK, data)
	}
}


//Get link from cache or nosql db
func getLink(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		params := mux.Vars(req)
		var url string = params["url"]

		l, ok := links[url]
		if ok == false {
			b := getLinkFromNosql(url)
			incrementAccessCount(b)
			formatter.JSON(w, http.StatusOK, b)
		} else {
			incrementAccessCount(l)
			formatter.JSON(w, http.StatusOK, l)
		}
	}
}


func incrementAccessCount(a bitly) {
	a.Access_count += 1
	updateNosql(a)
}


// Receive Order from Queue to Process
func queue_receive() []bitly {
	conn, err := amqp.Dial("amqp://"+rabbitmq_user+":"+rabbitmq_pass+"@"+rabbitmq_server+":"+rabbitmq_port+"/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		rabbitmq_queue, // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"trs",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	create_channel := make(chan bitly)
	
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var link bitly
			if err := json.Unmarshal(d.Body, &link); err != nil {
				panic(err)
			} 
			create_channel <- link
		}
		close(create_channel)
	}()

	err = ch.Cancel("trs", false)
	if err != nil {
	    log.Fatalf("basic.cancel: %v", err)
	}

	var links_array []bitly
	for n := range create_channel {
    	links_array = append(links_array, n)
    }	

    return links_array
}
