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
)

func main() {

	counts := make(map[string]int)

	csvFile := flag.String("csv", "problems.csv", "CSV Problem File")

	if *csvFile == "" {

		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()

	file, _ := os.Open(*csvFile)
	reader := csv.NewReader(file)
	input := bufio.NewReader(os.Stdin)

	for {

		record, err := reader.Read()

		if err == io.EOF {

			break
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(record[0], "=?")

		userInput, _ := input.ReadString('\n')
		userInput = strings.Replace(userInput, "\n", "", -1)

		if userInput == record[1] {
			counts["Correct"]++

		} else {

			counts["Incorrect"]++
		}

	}

<<<<<<< HEAD
	fmt.Println("Correct answers: \t", counts["Correct"], "\nIncorrect answers: \t", counts["Incorrect"])
=======
	fmt.Println("Correct answers: ", counts["Correct"])
>>>>>>> fe5991875c61c8130bfd17791a3eab4f367cfeaa

}
