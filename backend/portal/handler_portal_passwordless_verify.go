package portal

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	pbportal "github.com/shank318/doota/pb/doota/portal/v1"
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
)

func (p *Portal) PasswordlessVerify(ctx context.Context, c *connect.Request[pbportal.PasswordlessStartVerify]) (*connect.Response[pbportal.JWT], error) {
	logger := logging.Logger(ctx, p.logger)
	email := strings.TrimSpace(c.Msg.Email)
	code := strings.TrimSpace(c.Msg.Code)

	logger.Info("passwordless verify", zap.String("email", email), zap.String("code", strings.Repeat("*", len(code))))

	if !strings.Contains(email, "@") && !strings.Contains(email, ".") {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid email address"))
	}

	if code == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid code"))
	}

	ip, err := getIP(c.Header(), c.Peer().Addr)
	if err != nil {
		return nil, fmt.Errorf("get client IP: %w", err)
	}

	jwt, err := p.authUsecase.VerifyPasswordless(ctx, email, code, ip)
	if err != nil {
		var connectErr *connect.Error
		if errors.As(err, &connectErr) {
			return nil, connectErr
		}
		return nil, fmt.Errorf("failed to verify passwordless flow: %w", err)
	}
	return connect.NewResponse(jwt), nil
}
