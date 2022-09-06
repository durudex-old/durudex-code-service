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

package v1

import (
	"context"

	"github.com/durudex/durudex-code-service/internal/service"
	v1 "github.com/durudex/durudex-code-service/pkg/pb/durudex/v1"
)

// User gRPC server handler.
type UserHandler struct {
	service service.User
	v1.UnimplementedUserCodeServiceServer
}

// Creating a new user gRPC handler.
func NewUserHandler(service service.User) *UserHandler {
	return &UserHandler{service: service}
}

// Create verify user email code handler.
func (h *UserHandler) CreateVerifyUserEmailCode(ctx context.Context, input *v1.CreateVerifyUserEmailCodeRequest) (*v1.CreateVerifyUserEmailCodeResponse, error) {
	if err := h.service.CreateVerifyEmailCode(ctx, input.Email); err != nil {
		return &v1.CreateVerifyUserEmailCodeResponse{}, err
	}

	return &v1.CreateVerifyUserEmailCodeResponse{}, nil
}

// Verify user email code handler.
func (h *UserHandler) VerifyUserEmailCode(ctx context.Context, input *v1.VerifyUserEmailCodeRequest) (*v1.VerifyUserEmailCodeResponse, error) {
	status, err := h.service.VerifyEmailCode(ctx, input.Email, input.Code)
	if err != nil {
		return &v1.VerifyUserEmailCodeResponse{Status: false}, err
	}

	return &v1.VerifyUserEmailCodeResponse{Status: status}, nil
}
