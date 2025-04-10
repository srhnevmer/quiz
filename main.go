package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	fileName := flag.String("n", "problems", "set the name of the file with the problems")
	reorder := flag.Bool("r", false, "change the order of problems")
	timer := flag.Int64("t", 30, "time limit to complete the quiz")
	flag.Parse()

	q := quiz{
		fileName:  *fileName,
		reorder:   *reorder,
		timeLimit: time.Duration(*timer),
	}
	q.makeProblemList()
	q.run()
}

type problem struct {
	question string
	answer   string
}

type quiz struct {
	fileName      string
	reorder       bool
	numCorrectAns uint8
	problems      []problem
	timeLimit     time.Duration
}

func (q *quiz) makeProblemList() {
	file, err := os.Open(fmt.Sprintf("%s.csv", q.fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	d := formatData(data)
	if q.reorder {
		rand.Shuffle(len(d), func(i, j int) {
			d[i], d[j] = d[j], d[i]
		})
	}
	q.problems = d
}

func formatData(data [][]string) []problem {
	list := make([]problem, len(data))
	for i, line := range data {
		list[i].question = line[0]
		list[i].answer = line[1]
	}
	return list
}

func (q *quiz) run() {
Loop:
	for {
		clear()
		fmt.Fprintf(os.Stdout, "This quiz has a time limit\nYou have %d seconds to complete this quiz\nPress [y]es to continue or [n]o to exit: ", q.timeLimit)
		switch in := userInput(); {
		case bytes.Equal(in, []byte{121}), bytes.Equal(in, []byte{121, 101, 115}):
			break Loop
		case bytes.Equal(in, []byte{110}), bytes.Equal(in, []byte{110, 111}):
			os.Exit(0)
		}
	}

	timer := time.NewTimer(q.timeLimit * time.Second)
	go func() {
		<-timer.C
		resultNotification(q.numCorrectAns, len(q.problems))
		os.Exit(0)
	}()

	for i := range q.problems {
		clear()
		q.print(i)
		a, _ := strconv.Atoi(q.problems[i].answer)
		ua, err := strconv.Atoi(string(userInput()))
		if err != nil {
			continue
		}

		if a != ua {
			continue
		}

		q.numCorrectAns++
	}

	resultNotification(q.numCorrectAns, len(q.problems))
}

func (q quiz) print(idx int) {
	var str strings.Builder
	str.WriteString(fmt.Sprintf("Problem: %d\nQ: %s\nA: ", idx+1, q.problems[idx].question))
	fmt.Fprint(os.Stdout, str.String())
}

func clear() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd = exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func userInput() []byte {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	err := s.Err()
	if err != nil {
		log.Fatal(err)
	}
	return bytes.ToLower(bytes.TrimSpace(s.Bytes()))
}

func resultNotification(quantity uint8, total int) {
	clear()
	fmt.Fprintf(os.Stdout, "Number of correct answers: %d\nTotal number of problems: %d\n", quantity, total)
}
