CREATE TABLE IF NOT EXISTS code_judge (
                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            submission_id BIGINT NOT NULL,
                            language int NOT NULL,
                            source_code TEXT NOT NULL,
                            test_arguments TEXT,
                            test_results TEXT,
                            created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                            updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                            status VARCHAR(20) NOT NULL CHECK (status IN ('NEW', 'WAIT', 'DONE')),
                            token VARCHAR(255),
                            result VARCHAR(50),
                            result_message TEXT
);

CREATE INDEX IF NOT EXISTS idx_code_judge_submission_id ON code_judge (id);
CREATE INDEX IF NOT EXISTS idx_code_judge_status ON code_judge (status);
CREATE INDEX IF NOT EXISTS idx_code_judge_token ON code_judge (token);