package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"lms_judge_integrator/internal/model"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const (
	codeJudgeTableName     = "code_judge"
	idFieldName            = "id"
	assignmentIdFieldName  = "assignment_id"
	languageFieldName      = "language"
	sourceCodeFieldName    = "source_code"
	testArgumentsFieldName = "test_arguments"
	testResultsFieldName   = "test_results"
	statusFieldName        = "status"
	createdAtFieldName     = "created_at"
	updatedAtFieldName     = "updated_at"
	resultFieldName        = "result"
	resultMessageFieldName = "result_message"
	tokenFieldName         = "token"
	allFields              = "*"
)

type PostgresRepository struct {
	db *pgxpool.Pool
	sb sq.StatementBuilderType
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PostgresRepository) CreateSubmission(ctx context.Context, submissionDto *model.CreateSubmissionDto) (*model.CodeJudge, error) {
	now := time.Now()
	query, args, err := r.sb.Insert(codeJudgeTableName).
		Columns(
			assignmentIdFieldName,
			languageFieldName,
			sourceCodeFieldName,
			testArgumentsFieldName,
			testResultsFieldName,
			statusFieldName,
			createdAtFieldName,
		).
		Values(
			submissionDto.AssignmentId,
			submissionDto.Language,
			submissionDto.SourceCode,
			submissionDto.TestArguments,
			submissionDto.TestResults,
			model.GetCodeJudgeStatusName(model.New),
			now,
		).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var id string
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create submission: %w", err)
	}

	return model.NewCodeJudge(
		id,
		submissionDto.AssignmentId,
		submissionDto.Language,
		submissionDto.SourceCode,
		submissionDto.TestArguments,
		submissionDto.TestResults,
		now,
		model.GetCodeJudgeStatusName(model.New)), nil
}

func (r *PostgresRepository) GetPendingSubmissions(ctx context.Context) ([]model.CodeJudge, error) {
	query, args, err := r.sb.Select(allFields).
		From(codeJudgeTableName).
		Where(sq.NotEq{statusFieldName: model.GetCodeJudgeStatusName(model.Done)}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query submissions: %w", err)
	}
	defer rows.Close()

	var submissions []model.CodeJudge
	for rows.Next() {
		var sub model.CodeJudge
		err := rows.Scan(
			&sub.ID,
			&sub.AssignmentId,
			&sub.Language,
			&sub.SourceCode,
			&sub.TestArguments,
			&sub.TestResults,
			&sub.CreatedAt,
			&sub.UpdatedAt,
			&sub.Status,
			&sub.Token,
			&sub.Result,
			&sub.ResultMessage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan submission: %w", err)
		}
		submissions = append(submissions, sub)
	}
	return submissions, nil
}

func (r *PostgresRepository) UpdateSubmissionToken(ctx context.Context, id string, token string) error {
	query, args, err := r.sb.Update(codeJudgeTableName).
		Set(tokenFieldName, token).
		Set(statusFieldName, model.GetCodeJudgeStatusName(model.Wait)).
		Set(updatedAtFieldName, time.Now()).
		Where(sq.Eq{idFieldName: id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return checkRowsAffected(result)
}

func (r *PostgresRepository) UpdateSubmissionStatus(ctx context.Context, id string, status string) error {
	query, args, err := r.sb.Update(codeJudgeTableName).
		Set(statusFieldName, status).
		Set(updatedAtFieldName, time.Now()).
		Where(sq.Eq{idFieldName: id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return checkRowsAffected(result)
}

func (r *PostgresRepository) UpdateSubmissionResult(ctx context.Context, id string, judgeResult map[string]interface{}) error {
	var (
		resultStatus  sql.NullString
		resultMessage sql.NullString
		testOutput    sql.NullString
	)

	if status, ok := judgeResult[statusFieldName].(map[string]interface{}); ok {
		resultStatus.String = fmt.Sprintf("%v", status[idFieldName])
		resultStatus.Valid = true

		if desc, ok := status["description"].(string); ok {
			resultMessage.String = desc
			resultMessage.Valid = true
		}
	}

	if stdout, ok := judgeResult["stdout"].(string); ok {
		testOutput.String = stdout
		testOutput.Valid = true
	} else if stderr, ok := judgeResult["stderr"].(string); ok {
		testOutput.String = stderr
		testOutput.Valid = true
	}

	query, args, err := r.sb.Update(codeJudgeTableName).
		SetMap(map[string]interface{}{
			statusFieldName:        model.GetCodeJudgeStatusName(model.Done),
			resultFieldName:        resultStatus,
			resultMessageFieldName: resultMessage,
			testResultsFieldName:   testOutput,
			updatedAtFieldName:     time.Now(),
		}).
		Where(sq.Eq{idFieldName: id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return checkRowsAffected(result)
}

func checkRowsAffected(result pgconn.CommandTag) error {
	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}
