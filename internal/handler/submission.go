package handler

import (
	"encoding/json"
	"lms_judge_integrator/internal/model"
	"lms_judge_integrator/internal/service"
	"net/http"
)

type SubmissionHandler struct {
	service *service.JudgeService
}

func NewSubmissionHandler(s *service.JudgeService) *SubmissionHandler {
	return &SubmissionHandler{service: s}
}

func (h *SubmissionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SubmissionId  int    `json:"submissionId"`
		Language      int    `json:"language"`
		SourceCode    string `json:"sourceCode"`
		TestArguments string `json:"testArguments"`
		TestResults   string `json:"testResults"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createSubmissionDto := model.NewCreateSubmissionDto(req.SubmissionId, req.Language, req.SourceCode, req.TestArguments, req.TestResults)
	if _, err := h.service.CreateSubmission(r.Context(), createSubmissionDto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
