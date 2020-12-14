package main;

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"sort"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var trendServer_url = "http://10.0.1.224/links" 
//var trendServer_url = "http://localhost:3002/links" 

func init() {
	testTrendServer()
}

func testTrendServer() {
	_, err := http.Get(trendServer_url)
	if err != nil {
	   log.Fatalln(err)
	}
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
	mx.HandleFunc("/redirect/{url}", redirectToUrl).Methods("GET", http.MethodOptions)
	mx.HandleFunc("/links", getLinks(formatter)).Methods("GET", http.MethodOptions)
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Link Redirect server active"})
	}
}

func redirectToUrl(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(req)
	var url string = params["url"]
	
	res, err := http.Get(trendServer_url + "/" + url)
	if err != nil {
		failOnError(err, "Url not found.")
	}
	defer res.Body.Close()

	var b bitly
	_ = json.NewDecoder(res.Body).Decode(&b)
	fmt.Printf("%+v", b)
	

	req.Host = b.Original_url
	http.Redirect(writer, req, req.Host, 302)
}

// API Ping Handler
func getLinks(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		res, err := http.Get(trendServer_url)
		if err != nil {
			failOnError(err, "Url not found.")
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		data := make([]bitly, 0)
		json.Unmarshal(body, &data)

		sort.SliceStable(data, func(i, j int) bool {
		    return data[i].Access_count > data[j].Access_count
		})

		formatter.JSON(w, http.StatusOK, data)
	}
}


