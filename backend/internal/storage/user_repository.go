package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/model"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	rdb *redis.Client
}

func NewUserRepository(rdb *redis.Client) *UserRepository {
	return &UserRepository{rdb: rdb}
}

func (r *UserRepository) GetUser(ctx context.Context, id int64) (*model.User, error) {
	key := fmt.Sprintf("user:%d", id)
	val, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		return nil, err
	}

	var user model.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}
	if user.Portfolio == nil {
		user.Portfolio = make(map[string]int)
	}
	return &user, nil
}

func (r *UserRepository) SaveUser(ctx context.Context, user *model.User) error {
	key := fmt.Sprintf("user:%d", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, key, data, 0).Err()
}

func (r *UserRepository) CreateUser(ctx context.Context, id int64, password string) (*model.User, error) {
	key := fmt.Sprintf("user:%d", id)

	// Check if user exists
	exists, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if exists > 0 {
		return nil, fmt.Errorf("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           id,
		PasswordHash: string(hashedPassword),
		Balance:      10000.0,
		Portfolio:    make(map[string]int),
		CreatedAt:    time.Now(),
	}

	if err := r.SaveUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
