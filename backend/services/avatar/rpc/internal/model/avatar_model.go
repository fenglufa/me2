package model

import (
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Avatar struct {
	Id            int64     `db:"id"`
	AvatarId      int64     `db:"avatar_id"`
	UserId        int64     `db:"user_id"`
	Nickname      string    `db:"nickname"`
	AvatarUrl     string    `db:"avatar_url"`
	Gender        int32     `db:"gender"`
	BirthDate     string    `db:"birth_date"`
	Occupation    string    `db:"occupation"`
	MaritalStatus int32     `db:"marital_status"`
	Warmth        int32     `db:"warmth"`
	Adventurous   int32     `db:"adventurous"`
	Social        int32     `db:"social"`
	Creative      int32     `db:"creative"`
	Calm          int32     `db:"calm"`
	Energetic     int32     `db:"energetic"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type AvatarModel struct {
	conn sqlx.SqlConn
}

func NewAvatarModel(conn sqlx.SqlConn) *AvatarModel {
	return &AvatarModel{conn: conn}
}

func (m *AvatarModel) Insert(avatar *Avatar) (sql.Result, error) {
	query := `INSERT INTO avatars (avatar_id, user_id, nickname, avatar_url, gender, birth_date, occupation,
		marital_status, warmth, adventurous, social, creative, calm, energetic)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return m.conn.Exec(query, avatar.AvatarId, avatar.UserId, avatar.Nickname, avatar.AvatarUrl, avatar.Gender,
		avatar.BirthDate, avatar.Occupation, avatar.MaritalStatus, avatar.Warmth, avatar.Adventurous,
		avatar.Social, avatar.Creative, avatar.Calm, avatar.Energetic)
}

func (m *AvatarModel) FindByUserId(userId int64) (*Avatar, error) {
	var avatar Avatar
	query := `SELECT id, avatar_id, user_id, nickname, avatar_url, gender, birth_date, occupation,
		marital_status, warmth, adventurous, social, creative, calm, energetic, created_at, updated_at
		FROM avatars WHERE user_id = ?`
	err := m.conn.QueryRow(&avatar, query, userId)
	if err != nil {
		return nil, err
	}
	return &avatar, nil
}

func (m *AvatarModel) FindByAvatarId(avatarId int64) (*Avatar, error) {
	var avatar Avatar
	query := `SELECT id, avatar_id, user_id, nickname, avatar_url, gender, birth_date, occupation,
		marital_status, warmth, adventurous, social, creative, calm, energetic, created_at, updated_at
		FROM avatars WHERE avatar_id = ?`
	err := m.conn.QueryRow(&avatar, query, avatarId)
	if err != nil {
		return nil, err
	}
	return &avatar, nil
}

func (m *AvatarModel) UpdateProfile(avatarId int64, nickname, avatarUrl string) error {
	query := `UPDATE avatars SET nickname = ?, avatar_url = ? WHERE avatar_id = ?`
	_, err := m.conn.Exec(query, nickname, avatarUrl, avatarId)
	return err
}
