package main

import (
	"bufio"
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
)

var (
	fileName string
	shuffle  bool
)

func init() {
	flag.StringVar(&fileName, "n", "problems", "set the name of the file with the problems")
	flag.BoolVar(&shuffle, "s", false, "change the order of problems")
}

func main() {
	flag.Parse()

	q := quiz{fileName: fileName}
	q.makeProblemList(shuffle)
	q.run()
}

type problem struct {
	question string
	answer   string
}

type quiz struct {
	fileName      string
	numCorrectAns uint8
	problems      []problem
}

func (q *quiz) makeProblemList(flag bool) {
	file, err := os.Open(fmt.Sprintf("%s.csv", q.fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	switch d := formatData(data); flag {
	case true:
		rand.Shuffle(len(d), func(i, j int) {
			d[i], d[j] = d[j], d[i]
		})
		q.problems = d
	default:
		q.problems = d
	}
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
		clear()
		q.print(i)
		a, _ := strconv.Atoi(q.problems[i].answer)
		ua, err := strconv.Atoi(handleUserAns())
		if err != nil {
			continue
		}

		if a != ua {
			continue
		}

		q.numCorrectAns++
	}

	fmt.Fprintf(os.Stdout, "Number of correct answers: %d\nTotal number of problems: %d\n", q.numCorrectAns, len(q.problems))
}

func (q quiz) print(idx int) {
	var str strings.Builder
	str.WriteString(fmt.Sprintf("Problem: %d\nQ: %s\nA: ", idx+1, q.problems[idx].question))
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
