package model

import (
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type User struct {
	UserId                 int64        `db:"user_id"`
	Phone                  string       `db:"phone"`
	Nickname               string       `db:"nickname"`
	Avatar                 string       `db:"avatar"`
	SubscriptionType       int32        `db:"subscription_type"`
	SubscriptionExpireTime sql.NullTime `db:"subscription_expire_time"`
	Status                 int32        `db:"status"`
	CreatedAt              time.Time    `db:"created_at"`
	UpdatedAt              time.Time    `db:"updated_at"`
}

type UserModel struct {
	conn sqlx.SqlConn
}

func NewUserModel(conn sqlx.SqlConn) *UserModel {
	return &UserModel{
		conn: conn,
	}
}

// 根据手机号查询用户
func (m *UserModel) FindByPhone(phone string) (*User, error) {
	var user User
	query := `SELECT user_id, phone, nickname, avatar, subscription_type, subscription_expire_time, status, created_at, updated_at
			  FROM users WHERE phone = ? AND status != 3`
	err := m.conn.QueryRow(&user, query, phone)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 根据用户ID查询用户
func (m *UserModel) FindById(userId int64) (*User, error) {
	var user User
	query := `SELECT user_id, phone, nickname, avatar, subscription_type, subscription_expire_time, status, created_at, updated_at
			  FROM users WHERE user_id = ?`
	err := m.conn.QueryRow(&user, query, userId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 创建用户
func (m *UserModel) Insert(user *User) error {
	query := `INSERT INTO users (user_id, phone, nickname, avatar, subscription_type, status)
			  VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.conn.Exec(query, user.UserId, user.Phone, user.Nickname, user.Avatar, user.SubscriptionType, user.Status)
	return err
}

// 更新订阅信息
func (m *UserModel) UpdateSubscription(userId int64, subscriptionType int32, expireTime time.Time) error {
	query := `UPDATE users SET subscription_type = ?, subscription_expire_time = ? WHERE user_id = ?`
	_, err := m.conn.Exec(query, subscriptionType, expireTime, userId)
	return err
}

// 更新用户状态
func (m *UserModel) UpdateStatus(userId int64, status int32) error {
	query := `UPDATE users SET status = ? WHERE user_id = ?`
	_, err := m.conn.Exec(query, status, userId)
	return err
}

// 更新用户信息（昵称、头像）
func (m *UserModel) UpdateInfo(userId int64, nickname, avatar string) error {
	query := `UPDATE users SET nickname = ?, avatar = ? WHERE user_id = ?`
	_, err := m.conn.Exec(query, nickname, avatar, userId)
	return err
}
