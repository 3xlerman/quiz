package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func handleErr(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return ret
}

func main() {

	csvFilename := flag.String("csv", "problems.csv", "a .csv file in the format of 'question, answer'")
	timeLimit := flag.Int("limit", 10, "the time limit for each question in seconds")

	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		handleErr("Failed to open the .csv file: " + *csvFilename)
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		handleErr("Failed to parse the .csv file")
	}

	problems := parseLines(lines)

	// Counting seconds
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0

problemLoop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, problem.question)

		answerCh := make(chan string)

		// Waiting an answer, when main goroutine is counting seconds
		go func() {
			var answer string
			_, err := fmt.Scanf("%s\n", &answer)
			if err != nil {
				handleErr("Failed to handle user input")
			}
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("Time is over.")
			break problemLoop
		case answer := <-answerCh:
			if answer == problem.answer {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d of %d.\n", correct, len(problems))
}
