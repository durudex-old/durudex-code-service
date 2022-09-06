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

package service

import (
	"context"

	"github.com/durudex/durudex-code-service/internal/config"
	"github.com/durudex/durudex-code-service/internal/domain"
	"github.com/durudex/durudex-code-service/internal/repository/redis"
	"github.com/durudex/durudex-code-service/pkg/code"
	v1 "github.com/durudex/durudex-code-service/pkg/pb/durudex/v1"
)

// User code service interface.
type User interface {
	// Create verify user email code.
	CreateVerifyEmailCode(ctx context.Context, email string) error
	// Verify user email code.
	VerifyEmailCode(ctx context.Context, email string, input uint64) (bool, error)
}

// User code service structure.
type UserService struct {
	repos redis.User
	email v1.EmailUserServiceClient
	cfg   config.CodeConfig
}

// Creating a new user code service.
func NewUserService(repos redis.User, email v1.EmailUserServiceClient, cfg config.CodeConfig) *UserService {
	return &UserService{repos: repos, email: email, cfg: cfg}
}

// Create verify user email code.
func (s *UserService) CreateVerifyEmailCode(ctx context.Context, email string) error {
	// Generating a random uint64 code.
	c, err := code.Generate(s.cfg.MaxLength, s.cfg.MinLength)
	if err != nil {
		return err
	}

	// Creating a new user verification email code.
	if err := s.repos.CreateByEmail(ctx, email, c, s.cfg.TTL); err != nil {
		return err
	}

	// Sending an email to a user with a verification code.
	if _, err := s.email.SendEmailUserCode(ctx, &v1.SendEmailUserCodeRequest{
		Email:    email,
		Username: "new user",
		Code:     c,
	}); err != nil {
		return err
	}

	return nil
}

// Verify user email code.
func (s *UserService) VerifyEmailCode(ctx context.Context, email string, input uint64) (bool, error) {
	// Getting code by email.
	c, err := s.repos.GetByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	// Check input code.
	if input != c {
		return false, &domain.Error{Code: domain.CodeInvalidArgument, Message: "Invalid Code"}
	}

	return true, nil
}
