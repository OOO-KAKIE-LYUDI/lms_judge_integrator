package model

import (
	"database/sql"
	"time"
)

type CodeJudge struct {
	ID            string
	SubmissionId  int
	Language      int
	SourceCode    string
	TestArguments string
	TestResults   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Status        string
	Token         sql.NullString
	Result        sql.NullString
	ResultMessage sql.NullString
}

func NewCodeJudge(ID string, submissionId int, language int, sourceCode string, testArguments string, testResults string, createdAt time.Time, status string) *CodeJudge {
	return &CodeJudge{
		ID:            ID,
		SubmissionId:  submissionId,
		Language:      language,
		SourceCode:    sourceCode,
		TestArguments: testArguments,
		TestResults:   testResults,
		CreatedAt:     createdAt,
		Status:        status,
	}
}

type CodeJudgeStatus int

const (
	New CodeJudgeStatus = iota + 1
	Wait
	Done
)

func GetCodeJudgeStatusName(status CodeJudgeStatus) string {
	switch status {
	case New:
		return "NEW"
	case Wait:
		return "WAIT"
	case Done:
		return "DONE"
	default:
		return ""
	}
}

type CodeJudgeResult int

const (
	InQueue CodeJudgeResult = iota + 1
	Processing
	Accepted
	WrongAnswer
	TimeLimitExceeded
	CompilationError
	SIGSEGV
	SIGXFSZ
	SIGFPE
	SIGABRT
	NZEC
	Other
	InternalError
	ExecFormatError
)

func GetCodeJudgeResultFromInt(value int) CodeJudgeResult {
	switch value {
	case 1:
		return InQueue
	case 2:
		return Processing
	case 3:
		return Accepted
	case 4:
		return WrongAnswer
	case 5:
		return TimeLimitExceeded
	case 6:
		return CompilationError
	case 7:
		return SIGSEGV
	case 8:
		return SIGXFSZ
	case 9:
		return SIGFPE
	case 10:
		return SIGABRT
	case 11:
		return NZEC
	case 12:
		return Other
	case 13:
		return InternalError
	case 14:
		return ExecFormatError
	default:
		return 0
	}
}

type InternalCodeJudgeResult int

const (
	Ok InternalCodeJudgeResult = iota
	InProcess
	Failed
	RuntimeException
	InternalException
	FormatException
	OtherException
)

func GetInternalCodeJudgeResult(result CodeJudgeResult) InternalCodeJudgeResult {
	switch result {
	case Accepted:
		return Ok
	case InQueue, Processing:
		return InProcess
	case WrongAnswer, TimeLimitExceeded, CompilationError:
		return Failed
	case SIGSEGV, SIGXFSZ, SIGFPE, SIGABRT, NZEC, Other:
		return RuntimeException
	case InternalError:
		return InternalException
	case ExecFormatError:
		return FormatException
	default:
		return OtherException
	}
}

func IsCodeJudgeFinished(result InternalCodeJudgeResult) bool {
	if result != InProcess {
		return true
	}

	return false
}

type CreateSubmissionDto struct {
	SubmissionId  int
	Language      int
	SourceCode    string
	TestArguments string
	TestResults   string
}

func NewCreateSubmissionDto(submissionId, language int, sourceCode, testArguments, testResults string) *CreateSubmissionDto {
	return &CreateSubmissionDto{SubmissionId: submissionId, Language: language, SourceCode: sourceCode, TestArguments: testArguments, TestResults: testResults}
}
