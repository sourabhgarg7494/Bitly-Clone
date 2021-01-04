package main

import (
	"time"
)

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

type singleNode struct {
	Id          string
	Url         string
	Clicks      int
	DateCreated time.Time
}
type trendDataType struct {
	Data []singleNode
}

var trendDataKey string = "trendData"
var queueArray = [...]string{"CreateLinkQueue", "TrendDataUpdateQueue"}

var controlPanelQueue string = "CreateLinkQueue"
var trendDataUpdateQueue string = "TrendDataUpdateQueue"

//var cacheServerLink string = "http://localhost:9001/api/"

var cacheServerLink string = "http://internal-nosql-cluster-alb-1181897715.us-east-1.elb.amazonaws.com/api/"

type createParams struct {
	Url string
	Id  string
}

var mysqlURL string = "cmpe281:cmpe281@tcp(10.0.1.172:3306)/dwarfurl"
