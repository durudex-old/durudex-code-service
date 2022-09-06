/*
 * Copyright © 2022 Durudex
 *
 * This file is part of Durudex: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Durudex is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Durudex. If not, see <https://www.gnu.org/licenses/>.
 */

package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/durudex/durudex-code-service/pkg/database/redis"
)

// Redis module name.
const UserEmailModule string = "user:email"

// User code repository interface.
type User interface {
	// Creating a new user verification email code.
	CreateByEmail(ctx context.Context, email string, code uint64, ttl time.Duration) error
	// Getting a user verification email code.
	GetByEmail(ctx context.Context, email string) (uint64, error)
}

// User code repository structure.
type UserRepository struct{ redis redis.Redis }

// Creating a new user code repository.
func NewUserRepository(redis redis.Redis) *UserRepository {
	return &UserRepository{}
}

// Creating a new user verification email code.
func (r *UserRepository) CreateByEmail(ctx context.Context, email string, code uint64, ttl time.Duration) error {
	key := fmt.Sprintf("%s:%s", UserEmailModule, email)

	return r.redis.SetEX(ctx, key, code, ttl).Err()
}

// Getting a user verification email code.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (uint64, error) {
	key := fmt.Sprintf("%s:%s", UserEmailModule, email)

	return r.redis.Get(ctx, key).Uint64()
}
