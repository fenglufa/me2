package model

import (
	"database/sql"
	"strings"
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
	query := `INSERT INTO avatars (avatar_id, user_id, nickname, avatar_url, gender, birth_date, occupation, marital_status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	return m.conn.Exec(query, avatar.AvatarId, avatar.UserId, avatar.Nickname, avatar.AvatarUrl, avatar.Gender,
		avatar.BirthDate, avatar.Occupation, avatar.MaritalStatus)
}

func (m *AvatarModel) FindByUserId(userId int64) (*Avatar, error) {
	var avatar Avatar
	query := `SELECT id, avatar_id, user_id, nickname, avatar_url, gender, DATE_FORMAT(birth_date, '%Y-%m-%d') as birth_date, occupation,
		marital_status, created_at, updated_at
		FROM avatars WHERE user_id = ?`
	err := m.conn.QueryRow(&avatar, query, userId)
	if err != nil {
		return nil, err
	}
	return &avatar, nil
}

func (m *AvatarModel) FindByAvatarId(avatarId int64) (*Avatar, error) {
	var avatar Avatar
	query := `SELECT id, avatar_id, user_id, nickname, avatar_url, gender, DATE_FORMAT(birth_date, '%Y-%m-%d') as birth_date, occupation,
		marital_status, created_at, updated_at
		FROM avatars WHERE avatar_id = ?`
	err := m.conn.QueryRow(&avatar, query, avatarId)
	if err != nil {
		return nil, err
	}
	return &avatar, nil
}

func (m *AvatarModel) UpdateProfile(avatarId int64, nickname, avatarUrl string, gender int32, birthDate, occupation string, maritalStatus int32) error {
	updates := []string{}
	args := []interface{}{}

	if nickname != "" {
		updates = append(updates, "nickname = ?")
		args = append(args, nickname)
	}
	if avatarUrl != "" {
		updates = append(updates, "avatar_url = ?")
		args = append(args, avatarUrl)
	}
	if gender > 0 {
		updates = append(updates, "gender = ?")
		args = append(args, gender)
	}
	if birthDate != "" {
		updates = append(updates, "birth_date = ?")
		args = append(args, birthDate)
	}
	if occupation != "" {
		updates = append(updates, "occupation = ?")
		args = append(args, occupation)
	}
	if maritalStatus > 0 {
		updates = append(updates, "marital_status = ?")
		args = append(args, maritalStatus)
	}

	if len(updates) == 0 {
		return nil
	}

	args = append(args, avatarId)
	query := "UPDATE avatars SET " + strings.Join(updates, ", ") + " WHERE avatar_id = ?"
	_, err := m.conn.Exec(query, args...)
	return err
}
