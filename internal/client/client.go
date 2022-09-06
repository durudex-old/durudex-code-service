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

package client

import (
	"github.com/durudex/durudex-code-service/internal/config"
	v1 "github.com/durudex/durudex-code-service/pkg/pb/durudex/v1"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// Client structure.
type Client struct{ Email EmailClient }

// Email client structure.
type EmailClient struct {
	User v1.EmailUserServiceClient
	conn *grpc.ClientConn
}

// Creating a new client.
func NewClient(cfg config.ServiceConfig) *Client {
	log.Debug().Msg("Creating a new client...")

	emailServiceConn := ConnectToGRPCService(cfg.Email)

	return &Client{
		Email: EmailClient{
			User: v1.NewEmailUserServiceClient(emailServiceConn),
			conn: emailServiceConn,
		},
	}
}

// Closing a client connections.
func (c *Client) Close() {
	log.Info().Msg("Closing a client connections")

	if err := c.Email.conn.Close(); err != nil {
		log.Error().Err(err).Msg("failed to close email service connection")
	}
}
