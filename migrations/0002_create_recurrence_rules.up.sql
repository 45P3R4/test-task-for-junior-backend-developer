CREATE TABLE IF NOT EXISTS recurrence_rules (
    task_id BIGINT PRIMARY KEY,
    recurrence_type TEXT NOT NULL,
    recurrence_modifiers TEXT[] DEFAULT '{}',
    end_date TIMESTAMPTZ DEFAULT NULL,
    max_occurrences INT DEFAULT NULL,
    interval INT DEFAULT NULL,
    days INT[] DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_recurrence_task 
        FOREIGN KEY (task_id) 
        REFERENCES tasks(id) 
        ON DELETE CASCADE,
);

CREATE INDEX IF NOT EXISTS idx_recurrence_task_id ON recurrence_rules (task_id);
CREATE INDEX IF NOT EXISTS idx_recurrence_type ON recurrence_rules (recurrence_type);
CREATE INDEX IF NOT EXISTS idx_recurrence_end_date ON recurrence_rules (end_date) WHERE end_date IS NOT NULL;

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_tasks_updated_at
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_recurrence_rules_updated_at
    BEFORE UPDATE ON recurrence_rules
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
