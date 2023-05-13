package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type Problem struct {
	question string
	answer   string
}

func main() {
	// read the file name
	csvFile := flag.String("problems", "problems.csv", "csv file have questions")
	timeLimit := flag.Int("limit", 0, "the time limit for the quiz in seconds")
	flag.Parse()
	file, err := os.Open(*csvFile)
	if err != nil {
		exit(fmt.Sprintf("can't open the file: %s", *csvFile))
	}
	ex := csv.NewReader(file)
	content, err := ex.ReadAll()
	if err != nil {
		exit("Failed to parse the file")
	}
	quizzes := ParseFileContent(content)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	<-timer.C
	correct := 0
	wrong := 0
ProblemLoop:
	for i, p := range quizzes {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)
		answerChan := make(chan string)
		go func() {
			var input string
			fmt.Scanln(&input)
			answerChan <- input

		}()
		select {
		case <-timer.C:
			fmt.Println()
			break ProblemLoop
		case input := <-answerChan:
			if input == p.answer {
				correct++
			} else {
				wrong--
			}
		}
	}
	total := len(content)
	fmt.Printf("You scored %d out of %d.\n", correct, total)
}
func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
func ParseFileContent(content [][]string) []Problem {
	problems := make([]Problem, len(content))
	for i, line := range content {
		problems[i] = Problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return problems
}
