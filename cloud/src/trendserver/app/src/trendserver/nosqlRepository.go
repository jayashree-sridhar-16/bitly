package main

import (
	"bytes"
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
)


func updateNosql(l bitly) {
	fmt.Println("Updating...")
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(l)
	fmt.Printf("\n%+v", l)

	req,_ := http.NewRequest(
		"PUT",
		nosql_url + "/api/" + l.Short_url,
		b,
	)

	// add a request header
	req.Header.Add( "Content-Type", "application/json; charset=UTF-8" )
	
	// send an HTTP using `req` object
	_, err := http.DefaultClient.Do( req )

	
	if err != nil {
		failOnError(err, "Update to nosql db failed!")
	}

	links[l.Short_url] = l
}

func createNosql(l bitly) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(l)

	_, err := http.Post(nosql_url + "/api/" + l.Short_url, 
		"application/json; charset=UTF-8",
		b)
	if err != nil {
		failOnError(err, "Create to nosql db failed!")
	}

	links[l.Short_url] = l
}


func getLinkFromNosql(url string) bitly {
	res, err := http.Get(nosql_url + "/api/" + url)
	if err != nil {
		failOnError(err, "Link does not exist.")
	} 
	defer res.Body.Close()

	var b bitly
	_ = json.NewDecoder(res.Body).Decode(&b)
	fmt.Printf("\n%+v", b)

	return b
}


func getAllLinksFromNosql() {
	res, err := http.Get(nosql_url + "/api")
	if err != nil {
		failOnError(err, "Failed to get all api from nosql")
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var api []nosql_doc
	json.Unmarshal(body, &api)

	for x := range api {
		//fmt.Printf("%+v\n",api[x])
		var l = getLinkFromNosql(api[x].Key)
		links[api[x].Key] = l
	}
}