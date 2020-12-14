package main

type bitly struct {
	Original_url  string
	Short_url  string
	Redirect_url string
	Access_count int
}


var links map[string] bitly 

