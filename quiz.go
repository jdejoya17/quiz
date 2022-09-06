package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	. "github.com/jdejoya17/quiz/util"
	lg "github.com/jdejoya17/quiz/util/logger"
)

type QuizResults struct {
	total int
	score int
}

// global variable
type Global struct {
	problems   string
	time_limit int
}

var g Global

func init() {
	lg.Info.Println("initializing")

	// long flag name
	flag.StringVar(
		&g.problems,
		"problems",
		"problems.csv",
		"CSV file containing a list of problems and answers",
	)
	// short flag name
	flag.StringVar(
		&g.problems,
		"p",
		"problems.csv",
		"CSV file containing a list of problems and answers",
	)

	// long flag name
	flag.IntVar(
		&g.time_limit,
		"timelimit",
		2,
		"Time limit for the quiz",
	)
}

func main() {
	// parse cmdline arguments
	flag.Parse()

	// open problem list from csv file
	data, err := parseProblems()
	if err != nil {
		lg.Error.Println(ErrorProblemLoad)
		ErrMsg(ErrorProblemLoad)
	}

	// start quiz
	results, err := startQuiz(data)
	if err != nil {
		lg.Error.Println(ErrorScanAnswer)
		ErrMsg(ErrorScanAnswer)
	}

	// tally results
	fmt.Printf("------------------\n")
	fmt.Printf("RESULTS:\n")
	fmt.Printf("Total: %v\n", results.total)
	fmt.Printf("Score: %v\n", results.score)
}

func parseProblems() (data [][]string, err error) {
	// open csv file
	f, err := os.Open(g.problems)
	if err != nil {
		lg.Error.Println(err)
		return data, err
	}
	defer f.Close()

	// parse csv file
	csvReader := csv.NewReader(f)
	data, err = csvReader.ReadAll()
	if err != nil {
		lg.Error.Println(err)
		return data, err
	}

	return data, err
}

func startQuiz(problems [][]string) (QuizResults, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var problem, answer, userAnswer string
	results := QuizResults{total: len(problems), score: 0}
	for _, val := range problems {
		problem, answer = val[0], val[1]

		// prompt for answer
		fmt.Printf("Problem: %v\n", problem)
		fmt.Printf("Answer: ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			lg.Error.Println(err)
			return results, err
		} else {
			userAnswer = scanner.Text()
		}

		// validate answer
		if answer == userAnswer {
			results.score += 1
		}
		fmt.Printf("\n")
	}

	return results, nil
}
