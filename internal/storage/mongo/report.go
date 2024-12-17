package mongoStore

import (
	"context"
	"errors"
	"time"

	reportmodel "github.com/Negat1v9/work-marketplace/model/report"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNoReport = errors.New("no report")
)

type reportRepository struct {
	c *mongo.Collection
}

func newReportRepo(c *mongo.Collection) *reportRepository {
	return &reportRepository{
		c: c,
	}
}

func (r *reportRepository) createIndexes(ctx context.Context) ([]string, error) {
	indexed := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "reported_by.id", Value: 1}},
			Options: options.Index().SetName("report_reported_by_id_idx"),
		},
		{
			Keys:    bson.D{{Key: "reported_user.id", Value: 1}},
			Options: options.Index().SetName("report_reported_user_id_idx"),
		},
		{
			Keys:    bson.D{{Key: "task_id", Value: 1}},
			Options: options.Index().SetName("report_task_id_idx"),
		},
	}
	res, err := r.c.Indexes().CreateMany(ctx, indexed)
	return res, err
}
func (r *reportRepository) Create(ctx context.Context, report *reportmodel.Report) (string, error) {
	report.CreatedAt = time.Now().UTC()
	res, err := r.c.InsertOne(ctx, report)
	if err != nil {
		return "", err
	}

	id := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), nil
}
func (r *reportRepository) FindOne(ctx context.Context, filter bson.D) (*reportmodel.Report, error) {
	var report reportmodel.Report
	err := r.c.FindOne(ctx, filter).Decode(&report)
	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ErrNoReport
	case err != nil:
		return nil, err
	}

	return &report, nil
}
