package db

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	ollama "github.com/ollama/ollama/api"
	"github.com/rdawson46/dashboard/jobs"
)

func setupTestDB(t *testing.T) *sqliteRepo {
	// Use in-memory database for testing
	// file::memory:?cache=shared allows multiple connections to the same in-memory DB
	db_url := "file::memory:?cache=shared"
	
	// Migrations are in the same directory as this test
	migration_path := "file://migrations"
	
	logger := log.NewWithOptions(os.Stderr, log.Options{Level: log.DebugLevel})
	repo, err := NewSqliteRepository(logger, db_url, migration_path)
	if err != nil {
		t.Fatalf("Failed to create test DB: %v", err)
	}
	
	return repo
}

func TestUserOperations(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()
	ctx := context.Background()

	t.Run("CreateAndGetUser", func(t *testing.T) {
		username := "testuser"
		password := "Password123!" // Must satisfy checkPassword

		user, err := repo.CreateUser(ctx, username, password)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		if user.Username != username {
			t.Errorf("Expected username %s, got %s", username, user.Username)
		}

		fetchedUser, err := repo.GetUser(ctx, user.ID)
		if err != nil {
			t.Fatalf("Failed to fetch user: %v", err)
		}

		if fetchedUser.Username != username {
			t.Errorf("Expected username %s, got %s", username, fetchedUser.Username)
		}
	})

	t.Run("SignInUser", func(t *testing.T) {
		username := "signinuser"
		password := "Password123!"

		_, err := repo.CreateUser(ctx, username, password)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Success case
		user, err := repo.SignInUser(ctx, username, password)
		if err != nil {
			t.Fatalf("Failed to sign in: %v", err)
		}
		if user.Username != username {
			t.Errorf("Expected username %s, got %s", username, user.Username)
		}

		// Failure case - wrong password
		_, err = repo.SignInUser(ctx, username, "wrongpass")
		if err == nil {
			t.Error("Expected error for wrong password, got nil")
		}

		// Failure case - non-existent user
		_, err = repo.SignInUser(ctx, "nonexistent", password)
		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}
	})

	t.Run("GetUsersAndCount", func(t *testing.T) {
		initialCount, _ := repo.GetUserCount(ctx)
		
		for i := 0; i < 5; i++ {
			repo.CreateUser(ctx, fmt.Sprintf("user%d", i), "Password123!")
		}

		count, err := repo.GetUserCount(ctx)
		if err != nil {
			t.Fatalf("Failed to get count: %v", err)
		}
		if count != initialCount + 5 {
			t.Errorf("Expected count %d, got %d", initialCount + 5, count)
		}

		users, err := repo.GetUsers(ctx, 10, 0)
		if err != nil {
			t.Fatalf("Failed to get users: %v", err)
		}
		if int64(len(users)) < 5 {
			t.Errorf("Expected at least 5 users, got %d", len(users))
		}
	})
}

func TestMessageOperations(t *testing.T) {
	os.Setenv("OLLAMA_URL", "") // Ensure no external calls during generateDesc
	repo := setupTestDB(t)
	defer repo.Close()
	ctx := context.Background()

	// Create a user first
	user, _ := repo.CreateUser(ctx, "msguser", "Password123!")

	t.Run("CreateAndGetMessage", func(t *testing.T) {
		messages := []ollama.Message{
			{Role: "user", Content: "Hello, how are you?"},
		}
		model := "test-model"

		msg, err := repo.CreateMessage(ctx, user.ID, model, messages)
		if err != nil {
			t.Fatalf("Failed to create message: %v", err)
		}

		if msg.Id == "" {
			t.Error("Expected non-empty message ID")
		}

		fetchedMessages, err := repo.GetMessage(ctx, user.ID, msg.Id)
		if err != nil {
			t.Fatalf("Failed to fetch message: %v", err)
		}

		if len(fetchedMessages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(fetchedMessages))
		}
		if fetchedMessages[0].Content != messages[0].Content {
			t.Errorf("Expected content %s, got %s", messages[0].Content, fetchedMessages[0].Content)
		}
	})

	t.Run("AddMessage", func(t *testing.T) {
		messages := []ollama.Message{{Role: "user", Content: "Init"}}
		msg, _ := repo.CreateMessage(ctx, user.ID, "test-model", messages)

		newMessages := append(messages, ollama.Message{Role: "assistant", Content: "Hello!"})
		success, err := repo.AddMessage(ctx, msg.Id, user.ID, "test-model", newMessages)
		if err != nil {
			t.Fatalf("Failed to add message: %v", err)
		}
		if !success {
			t.Error("Expected success=true for AddMessage")
		}

		fetched, _ := repo.GetMessage(ctx, user.ID, msg.Id)
		if len(fetched) != 2 {
			t.Errorf("Expected 2 messages, got %d", len(fetched))
		}
	})

	t.Run("GetDescriptions", func(t *testing.T) {
		repo.CreateMessage(ctx, user.ID, "m", []ollama.Message{{Role: "user", Content: "Chat 1"}})
		repo.CreateMessage(ctx, user.ID, "m", []ollama.Message{{Role: "user", Content: "Chat 2"}})

		descs, err := repo.GetDescriptions(ctx, user.ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to get descriptions: %v", err)
		}
		if len(descs) < 2 {
			t.Errorf("Expected at least 2 descriptions, got %d", len(descs))
		}
	})

	t.Run("DeleteMessage", func(t *testing.T) {
		msg, _ := repo.CreateMessage(ctx, user.ID, "m", []ollama.Message{{Role: "user", Content: "To Delete"}})
		
		success, err := repo.DeleteMessage(ctx, msg.Id, user.ID)
		if err != nil {
			t.Fatalf("Failed to delete message: %v", err)
		}
		if !success {
			t.Error("Expected success=true for DeleteMessage")
		}

		_, err = repo.GetMessage(ctx, user.ID, msg.Id)
		if err == nil {
			t.Error("Expected error fetching deleted message, got nil")
		}
	})
}

func TestFileOperations(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()
	ctx := context.Background()

	user, _ := repo.CreateUser(ctx, "fileuser", "Password123!")

	t.Run("SaveAndGetFile", func(t *testing.T) {
		fileName := "test.txt"
		contentType := "text/plain"
		content := "Hello World"

		id, err := repo.SaveFile(ctx, user.ID, fileName, contentType, content)
		if err != nil {
			t.Fatalf("Failed to save file: %v", err)
		}

		file, err := repo.GetFile(ctx, id, user.ID)
		if err != nil {
			t.Fatalf("Failed to get file: %v", err)
		}

		if file.FileName != fileName {
			t.Errorf("Expected filename %s, got %s", fileName, file.FileName)
		}
		if file.Content != content {
			t.Errorf("Expected content %s, got %s", content, file.Content)
		}
	})

	t.Run("UpdateFile", func(t *testing.T) {
		id, _ := repo.SaveFile(ctx, user.ID, "old.txt", "text/plain", "old")
		
		err := repo.UpdateFile(ctx, id, user.ID, "new.txt", "text/plain", "new")
		if err != nil {
			t.Fatalf("Failed to update file: %v", err)
		}

		file, _ := repo.GetFile(ctx, id, user.ID)
		if file.FileName != "new.txt" {
			t.Errorf("Expected new filename, got %s", file.FileName)
		}
	})

	t.Run("GetFilesAndCount", func(t *testing.T) {
		initialCount, _ := repo.GetFilesCount(ctx, user.ID)
		repo.SaveFile(ctx, user.ID, "f1.txt", "t", "c")
		repo.SaveFile(ctx, user.ID, "f2.txt", "t", "c")

		count, err := repo.GetFilesCount(ctx, user.ID)
		if err != nil {
			t.Fatalf("Failed to get count: %v", err)
		}
		if count != initialCount + 2 {
			t.Errorf("Expected count %d, got %d", initialCount + 2, count)
		}

		files, err := repo.GetFiles(ctx, user.ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to get files: %v", err)
		}
		if len(files) < 2 {
			t.Errorf("Expected at least 2 files, got %d", len(files))
		}
	})

	t.Run("DeleteFile", func(t *testing.T) {
		id, _ := repo.SaveFile(ctx, user.ID, "del.txt", "t", "c")
		err := repo.DeleteFile(ctx, id, user.ID)
		if err != nil {
			t.Fatalf("Failed to delete file: %v", err)
		}

		_, err = repo.GetFile(ctx, id, user.ID)
		if err == nil {
			t.Error("Expected error fetching deleted file, got nil")
		}
	})
}

func TestJobOperations(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()
	ctx := context.Background()

	user, _ := repo.CreateUser(ctx, "jobuser", "Password123!")

	t.Run("CreateAndGetJob", func(t *testing.T) {
		job := &jobs.Job{
			Name: "Test Job",
			Task: jobs.NewTask("LLM"),
			Model: "test-model",
		}
		job.Task.Query = "test-query"

		createdJob, err := repo.CreateJob(ctx, user.ID, job)
		if err != nil {
			t.Fatalf("Failed to create job: %v", err)
		}

		if createdJob.Id == "" {
			t.Error("Expected non-empty job ID")
		}

		fetchedJob, err := repo.GetJob(ctx, createdJob.Id, user.ID)
		if err != nil {
			t.Fatalf("Failed to fetch job: %v", err)
		}

		if fetchedJob.Model != job.Model {
			t.Errorf("Expected model %s, got %s", job.Model, fetchedJob.Model)
		}
	})

	t.Run("UpdateAndPeekJob", func(t *testing.T) {
		job := &jobs.Job{
			Name: "Peek Job",
			Task: jobs.NewTask("LLM"),
			Status: jobs.StatusPending,
			Time: time.Now().UTC().Add(-time.Minute).Format(time.RFC3339),
		}
		job.Task.Query = "peek-query"
		repo.CreateJob(ctx, user.ID, job)

		peeked, err := repo.Peek(ctx)
		if err != nil {
			t.Fatalf("Failed to peek job: %v", err)
		}
		if peeked == nil {
			t.Fatal("Expected to peek a job, got nil")
		}
		if peeked.Status != jobs.StatusRunning {
			t.Errorf("Expected status 'Running' after peek, got %s", peeked.Status)
		}
	})

	t.Run("GetJobsAndCount", func(t *testing.T) {
		initialCount, _ := repo.GetJobCount(ctx, user.ID)
		j1 := &jobs.Job{Name: "j1", Task: jobs.NewTask("LLM")}
		j1.Task.Query = "t"
		repo.CreateJob(ctx, user.ID, j1)
		j2 := &jobs.Job{Name: "j2", Task: jobs.NewTask("LLM")}
		j2.Task.Query = "t"
		repo.CreateJob(ctx, user.ID, j2)



		count, err := repo.GetJobCount(ctx, user.ID)
		if err != nil {
			t.Fatalf("Failed to get count: %v", err)
		}
		if count < initialCount + 2 {
			t.Errorf("Expected count to increase by at least 2, got %d", count)
		}

		jobList, err := repo.GetJobs(ctx, user.ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to get jobs: %v", err)
		}
		if len(jobList) < 2 {
			t.Errorf("Expected at least 2 jobs, got %d", len(jobList))
		}
	})
}
