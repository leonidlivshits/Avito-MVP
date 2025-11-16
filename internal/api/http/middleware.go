package httpapi

import (
	"net/http"
	"strings"

	"github.com/leonidlivshits/Avito-MVP/internal/domain"
)

func TokenAuthMiddleware(adminToken, userToken, authorToken, reviewerToken string) func(http.Handler) http.Handler {
	tokenRole := map[string]domain.Role{}
	if adminToken != "" {
		tokenRole[adminToken] = domain.RoleAdmin
	}
	if userToken != "" {
		tokenRole[userToken] = domain.RoleObserver
	}
	if authorToken != "" {
		tokenRole[authorToken] = domain.RoleAuthor
	}
	if reviewerToken != "" {
		tokenRole[reviewerToken] = domain.RoleReviewer
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if xr := r.Header.Get("X-Role"); xr != "" {
				ctx = WithRole(ctx, domain.ParseRole(xr))
			}
			if uid := r.Header.Get("X-User-ID"); uid != "" {
				ctx = WithUserID(ctx, uid)
			}

			if GetRole(ctx) == domain.RoleObserver && GetUserID(ctx) == "" {
				auth := r.Header.Get("Authorization")
				if strings.HasPrefix(auth, "Bearer ") {
					tok := strings.TrimPrefix(auth, "Bearer ")
					if role, ok := tokenRole[tok]; ok {
						ctx = WithRole(ctx, role)
					}
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
