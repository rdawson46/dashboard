package jobs

import (
	"testing"
)

func TestTask(t *testing.T) {
	t.Run("NewTask_Valid", func(t *testing.T) {
		task := NewTask("LLM")
		if task == nil {
			t.Fatal("expected task, got nil")
		}
		if task.T != llm {
			t.Errorf("expected task type %v, got %v", llm, task.T)
		}
	})

	t.Run("NewTask_Invalid", func(t *testing.T) {
		task := NewTask("Invalid")
		if task != nil {
			t.Error("expected nil for invalid task type")
		}
	})

	t.Run("isValidTask", func(t *testing.T) {
		if !llm.isValidTask() {
			t.Error("LLM should be valid")
		}
		if !tool.isValidTask() {
			t.Error("Tool should be valid")
		}
		if taskType("other").isValidTask() {
			t.Error("other should be invalid")
		}
	})
}

func TestJob(t *testing.T) {
	t.Run("NewJob", func(t *testing.T) {
		task := NewTask("LLM")
		job := NewJob(task, "Test Job", "daily")
		if job == nil {
			t.Fatal("expected job, got nil")
		}
		if job.Name != "Test Job" {
			t.Errorf("expected name Test Job, got %s", job.Name)
		}
		if job.Status != StatusPending {
			t.Errorf("expected status Pending, got %s", job.Status)
		}
	})

	t.Run("isValidStatus", func(t *testing.T) {
		if !StatusPending.isValidStatus() {
			t.Error("Pending should be valid")
		}
		if !StatusRunning.isValidStatus() {
			t.Error("Running should be valid")
		}
		if !StatusCompleted.isValidStatus() {
			t.Error("Completed should be valid")
		}
		if !StatusFailed.isValidStatus() {
			t.Error("Failed should be valid")
		}
		if Status("other").isValidStatus() {
			t.Error("other should be invalid")
		}
	})
}
