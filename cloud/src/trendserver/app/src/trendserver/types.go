package main

type bitly struct {
	Original_url  string
	Short_url  string
	Redirect_url string
	Access_count int
}

type nosql_doc struct {
	Key string `json: "key"`
	Record string `json: "record"`
	Json string `json: "json"`
	Vclock [5]string `json: "vclock"`
	Message string `json: "message"`
}


var links map[string] bitly 

