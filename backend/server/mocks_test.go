package server

import (
	"context"

	ollama "github.com/ollama/ollama/api"
	"github.com/rdawson46/dashboard/db"
	"github.com/rdawson46/dashboard/jobs"
)

type mockRepository struct {
	SignInUserFunc       func(ctx context.Context, username, password string) (*db.User_db, error)
	CreateUserFunc       func(ctx context.Context, username, password string) (*db.User_db, error)
	CreateMessageFunc     func(ctx context.Context, userId, model string, message []ollama.Message) (db.CreateMessage, error)
	AddMessageFunc        func(ctx context.Context, messageId, userId, model string, message []ollama.Message) (bool, error)
	GetMessageFunc        func(ctx context.Context, userId, messageId string) ([]ollama.Message, error)
	GetDescriptionsFunc   func(ctx context.Context, userId string, limit, offset int) (db.Descriptions, error)
	DeleteMessageFunc     func(ctx context.Context, id string, user_id string) (bool, error)
	SaveFileFunc          func(ctx context.Context, userId, fileName, contentType, content string) (string, error)
	GetFileFunc           func(ctx context.Context, fileId, userId string) (*db.File, error)
	GetFilesFunc          func(ctx context.Context, userId string, limit, offset int) ([]*db.File, error)
	GetFilesCountFunc     func(ctx context.Context, userId string) (int, error)
	DeleteFileFunc        func(ctx context.Context, fileId, userId string) error
	UpdateFileFunc        func(ctx context.Context, fileId, userId, fileName, contentType, content string) error
	CreateJobFunc         func(ctx context.Context, userId string, job *jobs.Job) (*jobs.Job, error)
	GetJobFunc            func(ctx context.Context, jobId string, userId string) (*jobs.Job, error)
	GetJobsFunc           func(ctx context.Context, userId string, limit, offset int) (jobs.Jobs, error)
	UpdateJobFunc         func(ctx context.Context, job jobs.Job, userId string) (*jobs.Job, error)
	DeleteJobFunc         func(ctx context.Context, jobId string, userId string) error
	GetPerferredModelFunc func(ctx context.Context, userId string) (string, error)
}

func (m *mockRepository) CreateUser(ctx context.Context, username, password string) (*db.User_db, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, username, password)
	}
	return &db.User_db{ID: "test-id", Username: username}, nil
}

func (m *mockRepository) SignInUser(ctx context.Context, username, password string) (*db.User_db, error) {
	if m.SignInUserFunc != nil {
		return m.SignInUserFunc(ctx, username, password)
	}
	return &db.User_db{ID: "test-id", Username: username}, nil
}

func (m *mockRepository) CreateMessage(ctx context.Context, userId, model string, message []ollama.Message) (db.CreateMessage, error) {
	if m.CreateMessageFunc != nil {
		return m.CreateMessageFunc(ctx, userId, model, message)
	}
	return db.CreateMessage{Id: "new-message-id", New: true}, nil
}

func (m *mockRepository) AddMessage(ctx context.Context, messageId, userId, model string, message []ollama.Message) (bool, error) {
	if m.AddMessageFunc != nil {
		return m.AddMessageFunc(ctx, messageId, userId, model, message)
	}
	return true, nil
}

func (m *mockRepository) GetMessage(ctx context.Context, userId, messageId string) ([]ollama.Message, error) {
	if m.GetMessageFunc != nil {
		return m.GetMessageFunc(ctx, userId, messageId)
	}
	return nil, nil
}

func (m *mockRepository) GetDescriptions(ctx context.Context, userId string, limit, offset int) (db.Descriptions, error) {
	if m.GetDescriptionsFunc != nil {
		return m.GetDescriptionsFunc(ctx, userId, limit, offset)
	}
	return db.Descriptions{}, nil
}

func (m *mockRepository) DeleteMessage(ctx context.Context, id string, user_id string) (bool, error) {
	if m.DeleteMessageFunc != nil {
		return m.DeleteMessageFunc(ctx, id, user_id)
	}
	return true, nil
}

func (m *mockRepository) SaveFile(ctx context.Context, userId, fileName, contentType, content string) (string, error) {
	if m.SaveFileFunc != nil {
		return m.SaveFileFunc(ctx, userId, fileName, contentType, content)
	}
	return "file-id", nil
}

func (m *mockRepository) GetFile(ctx context.Context, fileId, userId string) (*db.File, error) {
	if m.GetFileFunc != nil {
		return m.GetFileFunc(ctx, fileId, userId)
	}
	return &db.File{ID: fileId, FileName: "test.txt"}, nil
}

func (m *mockRepository) GetFiles(ctx context.Context, userId string, limit, offset int) ([]*db.File, error) {
	if m.GetFilesFunc != nil {
		return m.GetFilesFunc(ctx, userId, limit, offset)
	}
	return nil, nil
}

func (m *mockRepository) GetFilesCount(ctx context.Context, userId string) (int, error) {
	if m.GetFilesCountFunc != nil {
		return m.GetFilesCountFunc(ctx, userId)
	}
	return 0, nil
}

func (m *mockRepository) DeleteFile(ctx context.Context, fileId, userId string) error {
	if m.DeleteFileFunc != nil {
		return m.DeleteFileFunc(ctx, fileId, userId)
	}
	return nil
}

func (m *mockRepository) UpdateFile(ctx context.Context, fileId, userId, fileName, contentType, content string) error {
	if m.UpdateFileFunc != nil {
		return m.UpdateFileFunc(ctx, fileId, userId, fileName, contentType, content)
	}
	return nil
}

func (m *mockRepository) CreateJob(ctx context.Context, userId string, job *jobs.Job) (*jobs.Job, error) {
	if m.CreateJobFunc != nil {
		return m.CreateJobFunc(ctx, userId, job)
	}
	return job, nil
}

func (m *mockRepository) GetJob(ctx context.Context, jobId string, userId string) (*jobs.Job, error) {
	if m.GetJobFunc != nil {
		return m.GetJobFunc(ctx, jobId, userId)
	}
	return &jobs.Job{Id: jobId}, nil
}

func (m *mockRepository) GetJobs(ctx context.Context, userId string, limit, offset int) (jobs.Jobs, error) {
	if m.GetJobsFunc != nil {
		return m.GetJobsFunc(ctx, userId, limit, offset)
	}
	return nil, nil
}

func (m *mockRepository) UpdateJob(ctx context.Context, job jobs.Job, userId string) (*jobs.Job, error) {
	if m.UpdateJobFunc != nil {
		return m.UpdateJobFunc(ctx, job, userId)
	}
	return &job, nil
}

func (m *mockRepository) DeleteJob(ctx context.Context, jobId string, userId string) error {
	if m.DeleteJobFunc != nil {
		return m.DeleteJobFunc(ctx, jobId, userId)
	}
	return nil
}

func (m *mockRepository) GetPerferredModel(ctx context.Context, userId string) (string, error) {
	if m.GetPerferredModelFunc != nil {
		return m.GetPerferredModelFunc(ctx, userId)
	}
	return "default", nil
}

func (m *mockRepository) SetPerferredModel(ctx context.Context, userId, model string) error { return nil }
func (m *mockRepository) GetUser(ctx context.Context, id string) (*db.User_db, error) { return nil, nil }
func (m *mockRepository) GetUsers(ctx context.Context, limit, offset int64) ([]*db.User_db, error) { return nil, nil }
func (m *mockRepository) GetUserCount(ctx context.Context) (int64, error) { return 0, nil }
func (m *mockRepository) UpdateUser() {}
func (m *mockRepository) Peek(ctx context.Context) (*jobs.Job, error) { return nil, nil }
func (m *mockRepository) GetJobCount(ctx context.Context, userId string) (int, error) { return 0, nil }
func (m *mockRepository) Close() error { return nil }
func (m *mockRepository) GetMessages() {}
func (m *mockRepository) GetMessageCount() {}
func (m *mockRepository) UpdateMessage() {}
