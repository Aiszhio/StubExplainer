package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/Aiszhio/StubExplainer/internal/model"
	"github.com/Aiszhio/StubExplainer/pkg/config"
)

// ClickHouseStorage stores explanations in ClickHouse.
type ClickHouseStorage struct {
	conn clickhouse.Conn
}

func NewClickHouseStorage(ctx context.Context, cfg config.Config) (*ClickHouseStorage, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.ClickHouseAddr},
		Auth: clickhouse.Auth{
			Database: cfg.ClickHouseDatabase,
			Username: cfg.ClickHouseUsername,
			Password: cfg.ClickHousePassword,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: time.Second * 10,
	})
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return &ClickHouseStorage{conn: conn}, nil
}

func (s *ClickHouseStorage) SaveBatch(items []model.Explanation) error {
	if len(items) == 0 {
		return nil
	}

	ctx := context.Background()
	batch, err := s.conn.PrepareBatch(ctx, "INSERT INTO explanations (incident_id, user_id, decision, score, features_json, created_at)")
	if err != nil {
		return err
	}

	for _, item := range items {
		features, err := json.Marshal(item.Features)
		if err != nil {
			return err
		}

		if err := batch.Append(
			item.IncidentID,
			item.UserID,
			item.Decision,
			item.Score,
			string(features),
			item.CreatedAt,
		); err != nil {
			return err
		}
	}

	return batch.Send()
}

func (s *ClickHouseStorage) GetByIncidentID(incidentID string) (model.Explanation, error) {
	ctx := context.Background()

	var item model.Explanation
	var featuresJSON string

	row := s.conn.QueryRow(ctx, `
		SELECT incident_id, user_id, decision, score, features_json, created_at
		FROM explanations
		WHERE incident_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`, incidentID)

	if err := row.Scan(
		&item.IncidentID,
		&item.UserID,
		&item.Decision,
		&item.Score,
		&featuresJSON,
		&item.CreatedAt,
	); err != nil {
		return model.Explanation{}, errors.New("explanation not found")
	}

	if err := json.Unmarshal([]byte(featuresJSON), &item.Features); err != nil {
		return model.Explanation{}, err
	}

	return item, nil
}

func (s *ClickHouseStorage) Close() error {
	return s.conn.Close()
}
