package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var q quiz
	q.makeProblemList()
	q.run()
}

type problem struct {
	question string
	answer   string
}

type quiz struct {
	numCorrectAns uint8
	problems      []problem
}

func (q *quiz) makeProblemList() {
	file, err := os.Open("problems.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	q.problems = formatData(data)
}

func formatData(data [][]string) []problem {
	list := make([]problem, 0, 5)
	for _, line := range data {
		var p problem
		for i, str := range line {
			switch i {
			case 0:
				p.question = str
			default:
				p.answer = str
			}
		}
		list = append(list, p)
	}
	res := make([]problem, len(list))
	copy(res, list)
	return res
}

func (q *quiz) run() {
	for i := range q.problems {
		q.print(i)
		ua, err := strconv.Atoi(handleUserAns())
		if err != nil {
			continue
		}

		a, _ := strconv.Atoi(q.problems[i].answer)
		if a != ua {
			continue
		}

		q.numCorrectAns++
	}

	fmt.Fprintf(os.Stdout, "Number of correct answers: %d\nTotal number of problems: %d\n", q.numCorrectAns, len(q.problems))
}

func (q quiz) print(num int) {
	var str strings.Builder
	str.WriteString(fmt.Sprintf("Problem: %d\nQ: %s\nA: ", num+1, q.problems[num].question))
	fmt.Fprint(os.Stdout, str.String())
}

func handleUserAns() string {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	err := s.Err()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(s.Text())
}
