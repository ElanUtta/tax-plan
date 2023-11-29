package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/elanutta/go-intensivo/internal/usecase"
	"github.com/elanutta/go-intensivo/pkg/rabbitmq"
	_ "github.com/mattn/go-sqlite3"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// orderRepository := database.NewOrderRepository(db)

	// uc := usecase.NewCalculateFinalPrice(orderRepository)
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}

	defer ch.Close()

	// msgRabbitmqChannel := make(chan amqp.Delivery)

	canal := make(chan usecase.OrderInput)

	go func() {
		for i := 4; i < 1000000; i++ {
			my_str := strconv.Itoa(i)
			canal <- usecase.OrderInput{my_str, 3.0, 3.1}
		}
	}()

	// for n := range canal {
	// 	fmt.Println(n)
	// }

	// time.Sleep(time.Second * 10)

	// println(canal)

	start := time.Now()

	for n := range canal {
		go rabbitmq.Publisher(ch, n)
		go rabbitmq.Publisher(ch, n)
		go rabbitmq.Publisher(ch, n)
		rabbitmq.Publisher(ch, n)
		// rabbitmq.Publisher(ch, n)
	}

	println("time duration", time.Since(start))

	// go rabbitmq.Consume(ch, msgRabbitmqChannel)
	// rabbitmqWorker(msgRabbitmqChannel, uc)

}

func rabbitmqWorker(msgChan chan amqp.Delivery, uc *usecase.CalculateFinalprice) {
	fmt.Println("Starting rabbitmq")
	for msg := range msgChan {
		var input usecase.OrderInput
		err := json.Unmarshal(msg.Body, &input)
		if err != nil {
			panic(err)
		}

		output, err := uc.Execute(input)
		if err != nil {
			panic(err)
		}

		msg.Ack(false)
		fmt.Println("Procced massage and saved on BD", output)
	}
}
