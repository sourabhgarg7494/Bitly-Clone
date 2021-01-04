package main

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
	Queue: "TrendDataUpdateQueue",
	User:  "guest",
	Pass:  "guest",
}

var randomStringLength int = 7

//var cacheServerLink string = "http://localhost:9001/api/"
var cacheServerLink string = "http://internal-nosql-cluster-alb-1181897715.us-east-1.elb.amazonaws.com/api/"

type createParams struct {
	Url string
	Id  string
}

type returnType struct {
	Url      string
	ShortUrl string
}

var frontendURL string = "https://dwarfurl2.herokuapp.com"
var frontendURL2 string = "http://dwarfurl2.herokuapp.com"
