package main

import "time"

type rabbitMqConfigType struct {
	Host  string
	Port  string
	Queue string
	User  string
	Pass  string
}

var configLocal rabbitMqConfigType = rabbitMqConfigType{
	Host:  "10.0.1.110",
	Port:  "5672",
	Queue: "CreateLinkQueue",
	User:  "guest",
	Pass:  "guest",
}

// var configLocal rabbitMqConfigType = rabbitMqConfigType{
// 	Host:  "localhost",
// 	Port:  "5672",
// 	Queue: "CreateLinkQueue",
// 	User:  "guest",
// 	Pass:  "guest",
// }

type singleNode struct {
	Id          string
	Url         string
	Clicks      int
	DateCreated time.Time
}
type trendDataType struct {
	Data []singleNode
}

//var cacheServerLink string = "http://localhost:9001/api/"
var cacheServerLink string = "http://internal-nosql-cluster-alb-1181897715.us-east-1.elb.amazonaws.com/api/"

var randomStringLength int = 7

type createParams struct {
	Url     string
	BaseUrl string
	Id      string
}

type createParams2 struct {
	Url string
	Id  string
}

type returnType struct {
	Url      string
	ShortUrl string
}

var trendDataKey string = "trendData"

var baseURL string = "dwarfurl.tk/"

//var frontendURL string = "http://localhost:3500"

var frontendURL string = "https://dwarfurl2.herokuapp.com"
var frontendURL2 string = "http://dwarfurl2.herokuapp.com"
