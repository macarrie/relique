package job_event

import (
	"database/sql"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"

	"github.com/macarrie/relique/internal/db"
)

func (e *JobEvent) Save() (int64, error) {
	tx, err := db.Handler().Begin()
	if err != nil {
		return 0, fmt.Errorf("cannot start transaction to save job event: %w", err)
	}
	// Defers are stacked, defer are executed in reverse order of stacking
	defer func() {
		if err != nil {
			e.GetLog().With(
				slog.Any("error", err),
			).Debug("Rollback job event save")
			tx.Rollback()
		}
	}()

	query, args, err := sq.Insert("job_events").SetMap(sq.Eq{
		"job_uuid":    e.JobUuid,
		"time":        e.Time,
		"type":        e.Type.Type,
		"description": e.Description,
	}).ToSql()
	if err != nil {
		return 0, fmt.Errorf("cannot build sql query: %w", err)
	}

	result, err := tx.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("cannot save job event into db: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("cannot get last insert ID: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("cannot commit job event save transaction: %w", err)
	}

	e.ID = id
	return id, nil
}

// ListByJobUuid returns the events of a job in the order they happened
func ListByJobUuid(uuid string) ([]JobEvent, error) {
	events := make([]JobEvent, 0)

	query, args, err := sq.Select(
		"id",
		"job_uuid",
		"time",
		"type",
		"description",
	).From(
		"job_events",
	).Where(
		"job_uuid = ?", uuid,
	).OrderBy(
		"time ASC",
		"id ASC",
	).ToSql()
	if err != nil {
		return events, fmt.Errorf("cannot build sql query: %w", err)
	}

	rows, err := db.Handler().Query(query, args...)
	if err != nil {
		return events, fmt.Errorf("cannot query job events from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var e JobEvent
		var description sql.NullString
		if err := rows.Scan(&e.ID, &e.JobUuid, &e.Time, &e.Type.Type, &description); err != nil {
			return events, fmt.Errorf("cannot parse job event from db: %w", err)
		}
		e.Description = description.String
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return events, fmt.Errorf("cannot read job events from db: %w", err)
	}

	return events, nil
}

// DeleteByJobUuid removes all events of a job inside the caller's transaction
func DeleteByJobUuid(tx *sql.Tx, uuid string) error {
	query, args, err := sq.Delete("job_events").Where("job_uuid = ?", uuid).ToSql()
	if err != nil {
		return fmt.Errorf("cannot build sql query: %w", err)
	}

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("cannot delete job events from db: %w", err)
	}

	return nil
}
