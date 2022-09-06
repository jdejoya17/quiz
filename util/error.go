package error

import (
	"fmt"
	"os"
)

// base error definition
type QuizBaseError struct {
	rc  int
	msg string
}

func (e *QuizBaseError) New(rc int, msg string) {
	e.rc = rc
	e.msg = msg
}

func (e *QuizBaseError) Error() string {
	return e.msg
}

// error codes
const (
	ErrorProblemLoadStatus int = 10
	ErrorScanAnswerStatus      = 11
)

// error message
const (
	ErrorProblemLoadMsg string = "Unable to load problems from input file."
	ErrorScanAnswerMsg         = "Unable to scan problem quiz answer"
)

// exportable error definitions
var (
	ErrorProblemLoad QuizBaseError // unable to load quiz file
	ErrorScanAnswer  QuizBaseError // failure to scan user answer
)

// error implementations
func init() {

	ErrorProblemLoad.New(
		ErrorProblemLoadStatus,
		ErrorProblemLoadMsg,
	)

	ErrorScanAnswer.New(
		ErrorScanAnswerStatus,
		ErrorScanAnswerMsg,
	)
}

func ErrMsg(e QuizBaseError) {
	fmt.Printf("quiz: %v\n", e.msg)
	os.Exit(e.rc)
}
