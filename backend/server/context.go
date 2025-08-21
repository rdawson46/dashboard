package server

import (
	"context"
)

type contextKey string

const userContextKey contextKey = "user"

func contextWithUser(ctx context.Context, user *User_jwt) context.Context {
    return context.WithValue(ctx, userContextKey, user)
}

func userFromContext(ctx context.Context) (*User_jwt, bool) {
    user, ok := ctx.Value(userContextKey).(*User_jwt)
    return user, ok
}
