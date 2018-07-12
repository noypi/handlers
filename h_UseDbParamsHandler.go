package handlers

import (
	"net/http"

	"context"
)

type DbParams struct {
	Kind      string
	Namespace string
}

func UseDbParams(namespace, kind string, nexth http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		params := new(DbParams)
		params.Kind = kind
		params.Namespace = namespace
		nexth.ServeHTTP(w, r.WithContext(context.WithValue(ctx, "dbparams", params)))
	})
}
