package routes

import (
	"context"
	"net/http"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"

	"github.com/julianstephens/go-utils/httputil/response"
)

func HandleError(responder *response.Responder, w http.ResponseWriter, r *http.Request, err error) {
	switch err {
	case context.Canceled:
		responder.ErrorWithStatus(w, r, http.StatusRequestTimeout, err, nil)
	case context.DeadlineExceeded:
		responder.ErrorWithStatus(w, r, http.StatusRequestTimeout, err, nil)
	case rpctypes.ErrEmptyKey:
		responder.BadRequest(w, r, err, nil)
	default:
		responder.InternalServerError(w, r, err, nil)
	}
}
