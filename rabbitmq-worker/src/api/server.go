package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
)

func initializeAllWorkers() {
	stop := make(chan bool)
	for _, element := range queueArray {
		go worker(element)
	}
	<-stop
}

func worker(queue string) {
	fmt.Printf("Worker started on %s\n", queue)
	conn, err := amqp.Dial("amqp://" + configLocal.User + ":" + configLocal.Pass + "@" + configLocal.Host + ":" + configLocal.Port + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	messages, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to Consume from channel")

	stop := make(chan bool)
	go func() {
		for d := range messages {
			if queue == controlPanelQueue {
				CreateLinkQueueHandlerMysql(d)
				InsertTrendDataHandler(d)
			} else if queue == trendDataUpdateQueue {
				UpdateTrendDataHandler(d)
			}
		}
	}()
	<-stop
}

//CreateLinkQueueHandlerMysql : Handler for the mesages received from CreateLinkQueue
func CreateLinkQueueHandlerMysql(d amqp.Delivery) {
	log.Printf("Received a message: %s", d.Body)
	data := &createParams{}
	//var data createParams
	err := json.Unmarshal(d.Body, data)
	if err != nil {
		failOnError(err, "Error in decoding Data")
	}

	db, err := sql.Open("mysql", mysqlURL)

	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO urls (original_url,uniqueId,clicks) VALUES (?,?,0)")
	if err != nil {
		failOnError(err, "Error in generating SQL statement")
	}
	_, err = stmt.Exec(data.Url, data.Id)
	if err != nil {
		failOnError(err, "Error in executing SQL statement")
	}
	log.Println("Data Successfully Inserted!!")
}

//InsertTrendDataHandler : Handler for inserting the data in Trend Document
func InsertTrendDataHandler(d amqp.Delivery) {
	log.Printf("Insert Trend Data: %s", d.Body)
	loc, _ := time.LoadLocation("America/Los_Angeles")
	data := &createParams{}
	err := json.Unmarshal(d.Body, data)
	var trendData trendDataType
	var newNode singleNode
	newNode.Clicks = 0
	newNode.DateCreated = time.Now().In(loc)
	newNode.Id = data.Id
	newNode.Url = data.Url
	//var data createParams
	if err != nil {
		failOnError(err, "Error in decoding Data")
	}
	response, err := http.Get(cacheServerLink + trendDataKey)
	if err != nil {
		failOnError(err, "Trend data get error")
	} else if response.StatusCode == 200 {
		_ = json.NewDecoder(response.Body).Decode(&trendData)
		trendData.Data = append(trendData.Data, newNode)
		TrendDataUpdateRequestCall(trendData)
	} else {
		var result map[string]interface{}
		resString, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal([]byte(resString), &result)
		//failOnError(result, "Some error occured in API call Or Document Not Found")
		log.Printf("Some error occured in API call Or Document Not Found: %s", resString)
		trendData.Data = append(trendData.Data, newNode)
		TrendDataPostRequestCall(trendData)
	}
}

//TrendDataUpdateRequestCall : Handler for updateing the Trend Data Request Call
func TrendDataUpdateRequestCall(trendData trendDataType) {
	jsonData, _ := json.Marshal(trendData)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, cacheServerLink+trendDataKey, bytes.NewBuffer(jsonData))
	if err != nil {
		failOnError(err, "Cache Data Put request failed")
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	response, err := client.Do(req)
	if err != nil {
		failOnError(err, "Error in Cache Put Request")
	} else if response.StatusCode == 200 {
		data, _ := ioutil.ReadAll(response.Body)
		log.Printf("Success from Cache: %s\n", string(data))
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		log.Printf("Some Error Occured: %s", string(data))
	}
}

//TrendDataPostRequestCall : Handler function for Posting the Object into Trend Data
func TrendDataPostRequestCall(trendData trendDataType) {
	jsonData, _ := json.Marshal(trendData)
	res, err := http.Post(cacheServerLink+trendDataKey, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Cache Data post request failed", err)
	} else if res.StatusCode == 200 {
		data, _ := ioutil.ReadAll(res.Body)
		log.Printf("Success from Cache: %s", string(data))
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		log.Printf("Some Error Occured: %s", string(data))
	}
}

//UpdateTrendDataHandler : Handler for updating the Trend Data
func UpdateTrendDataHandler(d amqp.Delivery) {
	log.Printf("Update Trend Data: %s", d.Body)

	data := &createParams{}
	//var data createParams
	err := json.Unmarshal(d.Body, data)
	if err != nil {
		failOnError(err, "Error in decoding Data")
	}

	response, err := http.Get(cacheServerLink + trendDataKey)
	if err != nil {
		failOnError(err, "Trend data get error")
	} else if response.StatusCode == 200 {
		var trendData trendDataType
		_ = json.NewDecoder(response.Body).Decode(&trendData)
		for i := range trendData.Data {
			if trendData.Data[i].Id == data.Id {
				trendData.Data[i].Clicks++
			}
		}
		log.Printf("Updated Data : %s", trendData)
		db, err := sql.Open("mysql", mysqlURL)

		defer db.Close()
		stmt, err := db.Prepare("UPDATE urls SET clicks =clicks+1 WHERE uniqueId = ?")
		if err != nil {
			failOnError(err, "Error in generating SQL statement")
		}
		_, err = stmt.Exec(data.Id)
		if err != nil {
			failOnError(err, "Error in executing SQL statement")
		}
		log.Println("Data Successfully Inserted!!")
		TrendDataUpdateRequestCall(trendData)
	} else {
		var result map[string]interface{}
		resString, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal([]byte(resString), &result)
		log.Printf("Some error occured in API call Or Document Not Found: %s", resString)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
	}
}
