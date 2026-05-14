package invitations

import (
	"context"
)

func GetInvitationByTokenService(ctx context.Context, token string) (Invitation, error) {
	return GetInvitationByToken(ctx, token)
}