package worker

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"lms_judge_integrator/internal/model"
	"lms_judge_integrator/internal/repository"
	"lms_judge_integrator/internal/service"
	"time"
)

type CheckerWorker struct {
	repo     *repository.PostgresRepository
	service  *service.JudgeService
	producer sarama.SyncProducer
	interval time.Duration
	stopChan chan struct{}
}

func NewCheckerWorker(repo *repository.PostgresRepository, service *service.JudgeService, producer sarama.SyncProducer, interval time.Duration) *CheckerWorker {
	return &CheckerWorker{
		repo:     repo,
		service:  service,
		producer: producer,
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

func (w *CheckerWorker) Run() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.processSubmissions()
		case <-w.stopChan:
			return
		}
	}
}

func (w *CheckerWorker) Stop() {
	close(w.stopChan)
}

func (w *CheckerWorker) processSubmissions() {
	ctx := context.Background()
	submissions, err := w.repo.GetPendingSubmissions(ctx)
	if err != nil {
		return
	}

	for _, sub := range submissions {
		switch sub.Status {
		case "NEW":
			w.service.SubmitToJudge0(ctx, &sub)
			w.repo.UpdateSubmissionStatus(ctx, sub.ID, "WAIT")
		case "WAIT":
			statusID, result, err := w.service.CheckSubmissionStatus(ctx, sub.Token.String)
			if err != nil {
				continue
			}

			if model.IsCodeJudgeFinished(statusID) {
				w.repo.UpdateSubmissionResult(ctx, sub.ID, result)

				msg, _ := json.Marshal(map[string]interface{}{
					"submissionId":  sub.SubmissionId,
					"codeJudgeId":   sub.ID,
					"status":        model.GetCodeJudgeStatusName(model.Done),
					"result":        statusID,
					"resultMessage": result["message"],
				})

				w.producer.SendMessage(&sarama.ProducerMessage{
					Topic: "code-judge-results",
					Value: sarama.ByteEncoder(msg),
				})
			}
		}
	}
}
