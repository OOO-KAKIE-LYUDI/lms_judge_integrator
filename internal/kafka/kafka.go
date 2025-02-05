package kafka

import (
	"github.com/IBM/sarama"
	"os"
)

const kafkaUrlVar = "kafka.broker.url"

type JudgeIntegratorKafka struct {
	kafkaBrokers []string
	config       *sarama.Config
}

func (j *JudgeIntegratorKafka) KafkaBrokers() []string {
	return j.kafkaBrokers
}

func (j *JudgeIntegratorKafka) Config() *sarama.Config {
	return j.config
}

func NewJudgeIntegratorKafka(kafkaBrokers []string, config *sarama.Config) *JudgeIntegratorKafka {
	return &JudgeIntegratorKafka{kafkaBrokers: kafkaBrokers, config: config}
}

func InitKafka() *JudgeIntegratorKafka {
	kafkaBrokers := []string{os.Getenv(kafkaUrlVar)}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	return NewJudgeIntegratorKafka(kafkaBrokers, config)
}
