package main

import (
	"bufio"
	"encoding/csv"
	"log"

	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	. "github.com/jdejoya17/quiz/util"
	lg "github.com/jdejoya17/quiz/util/logger"
)

type QuizResults struct {
	total     int
	score     int
	completed chan bool
}

// global variable
type Global struct {
	problems     string
	time_limit   int
	quiz_results QuizResults
}

var g Global

// flag definitions
var (
	rootCmd = &cobra.Command{
		Use:   "quiz",
		Short: "Runs a timed quiz",
		Long: `Runs a timed quiz. The problems are read from
an input file.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Run a Quiz!!!")
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&g.problems,
		"problems",
		"p",
		"problems.csv",
		"csv file containing a list of problems and asnwers",
	)

	rootCmd.PersistentFlags().IntVarP(
		&g.time_limit,
		"time_limit",
		"t",
		30,
		"time limit of the quiz",
	)

	helpFunc := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		helpFunc(rootCmd, []string{})
		os.Exit(0)
	})
}

func main() {
	// parse cmdline arguments
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
	// rootCmd.HelpFunc()

	// open problem list from csv file
	data, err := parseProblems()
	if err != nil {
		lg.Error.Println(ErrorProblemLoad)
		ErrMsg(ErrorProblemLoad)
	}
	g.quiz_results = QuizResults{
		total:     len(data),
		score:     0,
		completed: make(chan bool),
	}

	// start quiz timer
	fmt.Printf("Press Enter Key to start quiz\n")
	fmt.Printf("Time Limit: %v\n", g.time_limit)
	fmt.Scanln()
	tmr := time.NewTimer(
		time.Duration(g.time_limit * int(time.Second)),
	)

	// start quiz
	go func() {
		err := startQuiz(data)
		if err != nil {
			lg.Error.Println(ErrorScanAnswer)
			ErrMsg(ErrorScanAnswer)
		}
	}()

	// wait for timer
	for {
		select {
		case <-tmr.C:
			fmt.Println("Time is up!!")
			close(g.quiz_results.completed)
			goto end_quiz
		case <-g.quiz_results.completed:
			fmt.Println("Quiz complete!!")
			tmr.Stop()
			goto end_quiz
		}
	}
end_quiz:

	// tally results
	fmt.Printf("------------------\n")
	fmt.Printf("RESULTS:\n")
	fmt.Printf("Total: %v\n", g.quiz_results.total)
	fmt.Printf("Score: %v\n", g.quiz_results.score)
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

func startQuiz(problems [][]string) error {
	scanner := bufio.NewScanner(os.Stdin)
	var problem, answer, userAnswer string
	for _, val := range problems {
		problem, answer = val[0], val[1]

		// prompt for answer
		fmt.Printf("Problem: %v\n", problem)
		fmt.Printf("Answer: ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			lg.Error.Println(err)
			return err
		} else {
			userAnswer = scanner.Text()
		}

		// validate answer
		if answer == userAnswer {
			g.quiz_results.score += 1
		}
		fmt.Printf("\n")
	}

	g.quiz_results.completed <- true
	return nil
}
