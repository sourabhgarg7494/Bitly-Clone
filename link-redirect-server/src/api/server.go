package main

import (
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
	fmt.Println("Link Redirect Server")
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
	mx.HandleFunc("/lrping", pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/{id}", redirectURLHandler(formatter)).Methods("GET")
}

func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Link Redirect Ping API"})
	}
}

func redirectURLHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		var id string = params["id"]

		response, err := http.Get(cacheServerLink + id)
		if err != nil {
			log.Printf("Cache get error: %s", err)
		} else if response.StatusCode == 200 {
			var body createParams
			_ = json.NewDecoder(response.Body).Decode(&body)
			log.Println("Before Publish Message")
			publishMessage(body)
			log.Println("After publish Message")
			http.Redirect(w, req, body.Url, 302)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			log.Printf("Some Error Occured: %s", string(data))
		}
	}
}

func generateRandomID(length int) string {
	str := pwgen.AlphaNum(length)
	return str
}

func publishMessage(body createParams) {
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
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
