package mongo_mock

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockSession struct {
	context.Context
	mongo.Session
	ActiveTransaction  bool
	ClusterTimeValue   bson.Raw
	OperationTimeValue *primitive.Timestamp
}

func (m *MockSession) StartTransaction(opts ...*options.TransactionOptions) error {

	if m.ActiveTransaction {
		return errors.New("transaction already started")
	}
	m.ActiveTransaction = true
	return nil
}

func (m *MockSession) AbortTransaction(ctx context.Context) error {
	if !m.ActiveTransaction {
		return errors.New("no active transaction to abort")
	}
	m.ActiveTransaction = false
	return nil
}

func (m *MockSession) CommitTransaction(ctx context.Context) error {

	if !m.ActiveTransaction {
		return errors.New("no active transaction to commit")
	}
	m.ActiveTransaction = false
	return nil
}

func (m *MockSession) WithTransaction(ctx context.Context, fn func(ctx mongo.SessionContext) (interface{}, error),
	opts ...*options.TransactionOptions) (interface{}, error) {

	if err := m.StartTransaction(opts...); err != nil {
		return nil, err
	}
	defer func() {
		_ = m.AbortTransaction(ctx)
	}()

	result, err := fn(ctx.(mongo.SessionContext))
	if err != nil {
		return nil, err
	}

	if err := m.CommitTransaction(ctx); err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MockSession) EndSession(ctx context.Context) {
	m.ActiveTransaction = false
}

func (m *MockSession) ClusterTime() bson.Raw {
	return nil
}

func (m *MockSession) OperationTime() *primitive.Timestamp {
	return nil
}

func (m *MockSession) Client() *mongo.Client {
	return nil
}

func (m *MockSession) ID() bson.Raw {
	return nil
}

func (m *MockSession) AdvanceClusterTime(bson.Raw) error {
	return nil
}

func (m *MockSession) AdvanceOperationTime(ts *primitive.Timestamp) error {
	return nil
}
