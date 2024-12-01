package botmock

import (
	"context"

	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
)

type WebTgClientMock struct{}

func (m *WebTgClientMock) SendRespond(ctx context.Context, tgCreaterID int64, taskID string, worker *usermodel.User) error {

	return nil
}
func (m *WebTgClientMock) WaitFiles(ctx context.Context, tgCreaterID int64) error {
	return nil
}

func (m *WebTgClientMock) SelectWorker(ctx context.Context, tgWorkerID int64, taskID string) error {
	return nil
}
