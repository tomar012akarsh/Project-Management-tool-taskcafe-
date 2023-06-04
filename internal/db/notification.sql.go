// Code generated by sqlc. DO NOT EDIT.
// source: notification.sql

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createNotification = `-- name: CreateNotification :one
INSERT INTO notification (caused_by, data, action_type, created_on)
  VALUES ($1, $2, $3, $4) RETURNING notification_id, caused_by, action_type, data, created_on
`

type CreateNotificationParams struct {
	CausedBy   uuid.UUID       `json:"caused_by"`
	Data       json.RawMessage `json:"data"`
	ActionType string          `json:"action_type"`
	CreatedOn  time.Time       `json:"created_on"`
}

func (q *Queries) CreateNotification(ctx context.Context, arg CreateNotificationParams) (Notification, error) {
	row := q.db.QueryRowContext(ctx, createNotification,
		arg.CausedBy,
		arg.Data,
		arg.ActionType,
		arg.CreatedOn,
	)
	var i Notification
	err := row.Scan(
		&i.NotificationID,
		&i.CausedBy,
		&i.ActionType,
		&i.Data,
		&i.CreatedOn,
	)
	return i, err
}

const createNotificationNotifed = `-- name: CreateNotificationNotifed :one
INSERT INTO notification_notified (notification_id, user_id) VALUES ($1, $2) RETURNING notified_id, notification_id, user_id, read, read_at
`

type CreateNotificationNotifedParams struct {
	NotificationID uuid.UUID `json:"notification_id"`
	UserID         uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateNotificationNotifed(ctx context.Context, arg CreateNotificationNotifedParams) (NotificationNotified, error) {
	row := q.db.QueryRowContext(ctx, createNotificationNotifed, arg.NotificationID, arg.UserID)
	var i NotificationNotified
	err := row.Scan(
		&i.NotifiedID,
		&i.NotificationID,
		&i.UserID,
		&i.Read,
		&i.ReadAt,
	)
	return i, err
}

const getAllNotificationsForUserID = `-- name: GetAllNotificationsForUserID :many
SELECT notified_id, nn.notification_id, user_id, read, read_at, n.notification_id, caused_by, action_type, data, created_on FROM notification_notified AS nn
  INNER JOIN notification AS n ON n.notification_id = nn.notification_id
  WHERE nn.user_id = $1
`

type GetAllNotificationsForUserIDRow struct {
	NotifiedID       uuid.UUID       `json:"notified_id"`
	NotificationID   uuid.UUID       `json:"notification_id"`
	UserID           uuid.UUID       `json:"user_id"`
	Read             bool            `json:"read"`
	ReadAt           sql.NullTime    `json:"read_at"`
	NotificationID_2 uuid.UUID       `json:"notification_id_2"`
	CausedBy         uuid.UUID       `json:"caused_by"`
	ActionType       string          `json:"action_type"`
	Data             json.RawMessage `json:"data"`
	CreatedOn        time.Time       `json:"created_on"`
}

func (q *Queries) GetAllNotificationsForUserID(ctx context.Context, userID uuid.UUID) ([]GetAllNotificationsForUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllNotificationsForUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllNotificationsForUserIDRow
	for rows.Next() {
		var i GetAllNotificationsForUserIDRow
		if err := rows.Scan(
			&i.NotifiedID,
			&i.NotificationID,
			&i.UserID,
			&i.Read,
			&i.ReadAt,
			&i.NotificationID_2,
			&i.CausedBy,
			&i.ActionType,
			&i.Data,
			&i.CreatedOn,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNotificationByID = `-- name: GetNotificationByID :one
SELECT notification_id, caused_by, action_type, data, created_on FROM notification WHERE notification_id = $1
`

func (q *Queries) GetNotificationByID(ctx context.Context, notificationID uuid.UUID) (Notification, error) {
	row := q.db.QueryRowContext(ctx, getNotificationByID, notificationID)
	var i Notification
	err := row.Scan(
		&i.NotificationID,
		&i.CausedBy,
		&i.ActionType,
		&i.Data,
		&i.CreatedOn,
	)
	return i, err
}

const getNotificationsForUserIDCursor = `-- name: GetNotificationsForUserIDCursor :many
SELECT n.notification_id, n.caused_by, n.action_type, n.data, n.created_on, nn.notified_id, nn.notification_id, nn.user_id, nn.read, nn.read_at FROM notification_notified AS nn
  INNER JOIN notification AS n ON n.notification_id = nn.notification_id
  WHERE (n.created_on, n.notification_id) < ($1::timestamptz, $2::uuid)
  AND nn.user_id = $3::uuid
  AND ($4::boolean = false OR nn.read = false)
  AND ($5::boolean = false OR n.action_type = ANY($6::text[]))
  ORDER BY n.created_on DESC
  LIMIT $7::int
`

type GetNotificationsForUserIDCursorParams struct {
	CreatedOn        time.Time `json:"created_on"`
	NotificationID   uuid.UUID `json:"notification_id"`
	UserID           uuid.UUID `json:"user_id"`
	EnableUnread     bool      `json:"enable_unread"`
	EnableActionType bool      `json:"enable_action_type"`
	ActionType       []string  `json:"action_type"`
	LimitRows        int32     `json:"limit_rows"`
}

type GetNotificationsForUserIDCursorRow struct {
	NotificationID   uuid.UUID       `json:"notification_id"`
	CausedBy         uuid.UUID       `json:"caused_by"`
	ActionType       string          `json:"action_type"`
	Data             json.RawMessage `json:"data"`
	CreatedOn        time.Time       `json:"created_on"`
	NotifiedID       uuid.UUID       `json:"notified_id"`
	NotificationID_2 uuid.UUID       `json:"notification_id_2"`
	UserID           uuid.UUID       `json:"user_id"`
	Read             bool            `json:"read"`
	ReadAt           sql.NullTime    `json:"read_at"`
}

func (q *Queries) GetNotificationsForUserIDCursor(ctx context.Context, arg GetNotificationsForUserIDCursorParams) ([]GetNotificationsForUserIDCursorRow, error) {
	rows, err := q.db.QueryContext(ctx, getNotificationsForUserIDCursor,
		arg.CreatedOn,
		arg.NotificationID,
		arg.UserID,
		arg.EnableUnread,
		arg.EnableActionType,
		pq.Array(arg.ActionType),
		arg.LimitRows,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetNotificationsForUserIDCursorRow
	for rows.Next() {
		var i GetNotificationsForUserIDCursorRow
		if err := rows.Scan(
			&i.NotificationID,
			&i.CausedBy,
			&i.ActionType,
			&i.Data,
			&i.CreatedOn,
			&i.NotifiedID,
			&i.NotificationID_2,
			&i.UserID,
			&i.Read,
			&i.ReadAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNotificationsForUserIDPaged = `-- name: GetNotificationsForUserIDPaged :many
SELECT n.notification_id, n.caused_by, n.action_type, n.data, n.created_on, nn.notified_id, nn.notification_id, nn.user_id, nn.read, nn.read_at FROM notification_notified AS nn
  INNER JOIN notification AS n ON n.notification_id = nn.notification_id
  WHERE nn.user_id = $1::uuid
  AND ($2::boolean = false OR nn.read = false)
  AND ($3::boolean = false OR n.action_type = ANY($4::text[]))
  ORDER BY n.created_on DESC
  LIMIT $5::int
`

type GetNotificationsForUserIDPagedParams struct {
	UserID           uuid.UUID `json:"user_id"`
	EnableUnread     bool      `json:"enable_unread"`
	EnableActionType bool      `json:"enable_action_type"`
	ActionType       []string  `json:"action_type"`
	LimitRows        int32     `json:"limit_rows"`
}

type GetNotificationsForUserIDPagedRow struct {
	NotificationID   uuid.UUID       `json:"notification_id"`
	CausedBy         uuid.UUID       `json:"caused_by"`
	ActionType       string          `json:"action_type"`
	Data             json.RawMessage `json:"data"`
	CreatedOn        time.Time       `json:"created_on"`
	NotifiedID       uuid.UUID       `json:"notified_id"`
	NotificationID_2 uuid.UUID       `json:"notification_id_2"`
	UserID           uuid.UUID       `json:"user_id"`
	Read             bool            `json:"read"`
	ReadAt           sql.NullTime    `json:"read_at"`
}

func (q *Queries) GetNotificationsForUserIDPaged(ctx context.Context, arg GetNotificationsForUserIDPagedParams) ([]GetNotificationsForUserIDPagedRow, error) {
	rows, err := q.db.QueryContext(ctx, getNotificationsForUserIDPaged,
		arg.UserID,
		arg.EnableUnread,
		arg.EnableActionType,
		pq.Array(arg.ActionType),
		arg.LimitRows,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetNotificationsForUserIDPagedRow
	for rows.Next() {
		var i GetNotificationsForUserIDPagedRow
		if err := rows.Scan(
			&i.NotificationID,
			&i.CausedBy,
			&i.ActionType,
			&i.Data,
			&i.CreatedOn,
			&i.NotifiedID,
			&i.NotificationID_2,
			&i.UserID,
			&i.Read,
			&i.ReadAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNotifiedByID = `-- name: GetNotifiedByID :one
SELECT notified_id, nn.notification_id, user_id, read, read_at, n.notification_id, caused_by, action_type, data, created_on FROM notification_notified as nn
  INNER JOIN notification AS n ON n.notification_id = nn.notification_id
  WHERE notified_id = $1
`

type GetNotifiedByIDRow struct {
	NotifiedID       uuid.UUID       `json:"notified_id"`
	NotificationID   uuid.UUID       `json:"notification_id"`
	UserID           uuid.UUID       `json:"user_id"`
	Read             bool            `json:"read"`
	ReadAt           sql.NullTime    `json:"read_at"`
	NotificationID_2 uuid.UUID       `json:"notification_id_2"`
	CausedBy         uuid.UUID       `json:"caused_by"`
	ActionType       string          `json:"action_type"`
	Data             json.RawMessage `json:"data"`
	CreatedOn        time.Time       `json:"created_on"`
}

func (q *Queries) GetNotifiedByID(ctx context.Context, notifiedID uuid.UUID) (GetNotifiedByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getNotifiedByID, notifiedID)
	var i GetNotifiedByIDRow
	err := row.Scan(
		&i.NotifiedID,
		&i.NotificationID,
		&i.UserID,
		&i.Read,
		&i.ReadAt,
		&i.NotificationID_2,
		&i.CausedBy,
		&i.ActionType,
		&i.Data,
		&i.CreatedOn,
	)
	return i, err
}

const getNotifiedByIDNoExtra = `-- name: GetNotifiedByIDNoExtra :one
SELECT notified_id, notification_id, user_id, read, read_at FROM notification_notified as nn WHERE nn.notified_id = $1
`

func (q *Queries) GetNotifiedByIDNoExtra(ctx context.Context, notifiedID uuid.UUID) (NotificationNotified, error) {
	row := q.db.QueryRowContext(ctx, getNotifiedByIDNoExtra, notifiedID)
	var i NotificationNotified
	err := row.Scan(
		&i.NotifiedID,
		&i.NotificationID,
		&i.UserID,
		&i.Read,
		&i.ReadAt,
	)
	return i, err
}

const hasUnreadNotification = `-- name: HasUnreadNotification :one
SELECT EXISTS (SELECT 1 FROM notification_notified WHERE read = false AND user_id = $1)
`

func (q *Queries) HasUnreadNotification(ctx context.Context, userID uuid.UUID) (bool, error) {
	row := q.db.QueryRowContext(ctx, hasUnreadNotification, userID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const markAllNotificationsRead = `-- name: MarkAllNotificationsRead :exec
UPDATE notification_notified SET read = true, read_at = $2 WHERE user_id = $1
`

type MarkAllNotificationsReadParams struct {
	UserID uuid.UUID    `json:"user_id"`
	ReadAt sql.NullTime `json:"read_at"`
}

func (q *Queries) MarkAllNotificationsRead(ctx context.Context, arg MarkAllNotificationsReadParams) error {
	_, err := q.db.ExecContext(ctx, markAllNotificationsRead, arg.UserID, arg.ReadAt)
	return err
}

const markNotificationAsRead = `-- name: MarkNotificationAsRead :exec
UPDATE notification_notified SET read = $3, read_at = $2 WHERE user_id = $1 AND notified_id = $4
`

type MarkNotificationAsReadParams struct {
	UserID     uuid.UUID    `json:"user_id"`
	ReadAt     sql.NullTime `json:"read_at"`
	Read       bool         `json:"read"`
	NotifiedID uuid.UUID    `json:"notified_id"`
}

func (q *Queries) MarkNotificationAsRead(ctx context.Context, arg MarkNotificationAsReadParams) error {
	_, err := q.db.ExecContext(ctx, markNotificationAsRead,
		arg.UserID,
		arg.ReadAt,
		arg.Read,
		arg.NotifiedID,
	)
	return err
}