package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chr4/pwgen"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/streadway/amqp"
	"github.com/unrolled/render"
)

func NewServer() *negroni.Negroni {
	fmt.Println("Inside New Server")
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	fmt.Println("Router Created")
	initRoutes(mx, formatter)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{frontendURL, frontendURL2},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"},
	})
	handler := c.Handler(mx)
	n.UseHandler(handler)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/cpping", pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/create", createURLHandler(formatter)).Methods("POST")
	mx.HandleFunc("/trendData", getTrendDataHandler(formatter)).Methods("GET")
}

func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// setupResponse(&w, req)

		// if (*req).Method == "OPTIONS" {
		// 	fmt.Println("PREFLIGHT Request")
		// 	return
		// }
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Control Panel Ping API"})
	}
}

// func setupResponse(w *http.ResponseWriter, req *http.Request) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", frontendURL2)
// 	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
// 	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// }

func createURLHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// setupResponse(&w, req)

		// if (*req).Method == "OPTIONS" {
		// 	fmt.Println("PREFLIGHT Request")
		// 	return
		// }
		id := generateRandomID(randomStringLength)
		var body createParams
		_ = json.NewDecoder(req.Body).Decode(&body)
		if body.Url == "" {
			formatter.JSON(w, http.StatusBadRequest, struct{ Error string }{"No URL Passed"})
		} else {
			body.Id = id
			var result returnType
			result.ShortUrl = body.BaseUrl + id
			result.Url = body.Url

			var body2 createParams2
			body2.Id = id
			body2.Url = body.Url

			publishMessage(body2)
			putLinkInCache(body2)
			formatter.JSON(w, http.StatusOK, result)
		}
	}
}

func getTrendDataHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// setupResponse(&w, req)

		// if (*req).Method == "OPTIONS" {
		// 	fmt.Println("PREFLIGHT Request")
		// 	return
		// }
		var trendData trendDataType
		response, err := http.Get(cacheServerLink + trendDataKey)
		if err != nil {
			failOnError(err, "Trend data get error")
			formatter.JSON(w, http.StatusBadRequest, err)
		} else if response.StatusCode == 200 {
			_ = json.NewDecoder(response.Body).Decode(&trendData)
			formatter.JSON(w, http.StatusOK, trendData)
		} else {
			var result map[string]interface{}
			resString, _ := ioutil.ReadAll(response.Body)
			json.Unmarshal([]byte(resString), &result)
			log.Printf("Some error occured in API call Or Document Not Found: %s", resString)
			formatter.JSON(w, http.StatusBadRequest, result)
		}
	}
}

func generateRandomID(length int) string {
	str := pwgen.AlphaNum(length)
	return str
}

func putLinkInCache(body createParams2) {

	jsonData, _ := json.Marshal(body)
	id := body.Id
	response, err := http.Post(cacheServerLink+id, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Cache Data post request failed: %s\n", err)
	} else if response.StatusCode == 200 {
		data, _ := ioutil.ReadAll(response.Body)
		log.Printf("Success from Cache: %s", string(data))
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		log.Printf("Some Error Occured: %s", string(data))
	}
	//Code to put the Link in Chache Node.
}

func publishMessage(body createParams2) {
	conn, err := amqp.Dial("amqp://" + configLocal.User + ":" + configLocal.Pass + "@" + configLocal.Host + ":" + configLocal.Port + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		configLocal.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	strbody, err := json.Marshal(body)
	failOnError(err, "Failed to encode Body JSON")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        strbody,
		})
	log.Printf(" [x] Sent %s", strbody)
	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Sprintf("%s: %s", msg, err)
	}
}
