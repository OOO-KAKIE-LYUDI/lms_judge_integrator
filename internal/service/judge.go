package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"lms_judge_integrator/internal/model"
	"lms_judge_integrator/internal/repository"
	"net/http"
)

type JudgeService struct {
	repo      *repository.PostgresRepository
	judge0URL string
}

func NewJudgeService(repo *repository.PostgresRepository, judge0URL string) *JudgeService {
	return &JudgeService{
		repo:      repo,
		judge0URL: judge0URL,
	}
}

func (s *JudgeService) CreateSubmission(ctx context.Context, submissionDto *model.CreateSubmissionDto) (*model.CodeJudge, error) {
	return s.repo.CreateSubmission(ctx, submissionDto)
}

func (s *JudgeService) SubmitToJudge0(ctx context.Context, submission *model.CodeJudge) error {
	type Judge0Request struct {
		SourceCode     string `json:"source_code"`
		LanguageID     int    `json:"language_id"`
		StdIn          string `json:"stdin"`
		ExpectedOutput string `json:"expected_output"`
	}

	reqBody := Judge0Request{
		SourceCode:     submission.SourceCode,
		LanguageID:     submission.Language,
		StdIn:          submission.TestArguments,
		ExpectedOutput: submission.TestResults,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s/submissions", s.judge0URL),
		"application/json", bytes.NewReader(reqBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	return s.repo.UpdateSubmissionToken(ctx, submission.ID, result.Token)
}

func (s *JudgeService) CheckSubmissionStatus(ctx context.Context, token string) (model.InternalCodeJudgeResult, map[string]interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/submissions/%s", s.judge0URL, token))
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, nil, err
	}

	status := result["status"].(map[string]interface{})
	if status == nil {
		return 0, nil, err
	}

	rawID, ok := status["id"].(float64)
	if !ok {
		return 0, nil, fmt.Errorf("invalid id type")
	}

	id := model.CodeJudgeResult(int(rawID))
	return model.GetInternalCodeJudgeResult(id), result, nil
}
