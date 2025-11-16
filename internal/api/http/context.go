package httpapi

import (
	"context"

	"github.com/leonidlivshits/Avito-MVP/internal/domain"
)

type ctxKey string

const (
	ctxKeyRole   ctxKey = "role"
	ctxKeyUserID ctxKey = "user_id"
	ctxKeyReqID  ctxKey = "request_id"
)

func WithRole(ctx context.Context, r domain.Role) context.Context {
	return context.WithValue(ctx, ctxKeyRole, r)
}

func GetRole(ctx context.Context) domain.Role {
	if v := ctx.Value(ctxKeyRole); v != nil {
		if r, ok := v.(domain.Role); ok {
			return r
		}
	}
	return domain.RoleObserver
}

func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKeyUserID, id)
}

func GetUserID(ctx context.Context) string {
	if v := ctx.Value(ctxKeyUserID); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKeyReqID, id)
}

func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(ctxKeyReqID); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
