package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {

	counts := make(map[string]int)

	csvFile := flag.String("csv", "problems.csv", "CSV Problem File")
	timerInt := flag.Int("timer", 3, "time in seconds to answer question")

	if *csvFile == "" {

		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()

	file, err := os.Open(*csvFile)

	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(file)
	input := bufio.NewReader(os.Stdin)
	timer := time.Duration(*timerInt) * time.Second
	c1 := make(chan string, 1)

	fmt.Println("Timer between questions set to: ", timer)

	for {

		record, err := reader.Read()

		if err == io.EOF {

			break
		}

		if err != nil {
			log.Fatal(err)
		}

		go func() {
			fmt.Println(record[0], "=?")
			userInput, _ := input.ReadString('\n')
			userInput = strings.Replace(userInput, "\n", "", -1)
			c1 <- userInput
		}()

		select {

		case userInput := <-c1:

			if userInput == "exit" {
				os.Exit(1)
			}

			if userInput == record[1] {
				counts["Correct"]++
			} else {
				counts["Incorrect"]++
			}
		case <-time.After(timer):
			counts["Incorrect"]++
			fmt.Println("Time expired. Answer marked as incorrect")
		}

	}

	fmt.Println("Correct answers: \t", counts["Correct"], "\nIncorrect answers: \t", counts["Incorrect"])

}
