package routes

import (
	"context"
	"net/http"

	"github.com/julianstephens/go-utils/httputil/response"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
)



func HandleError(responder *response.Responder, w http.ResponseWriter, r *http.Request, err error) {
	switch err {
	case context.Canceled:
		responder.ErrorWithStatus(w, r, http.StatusRequestTimeout, err)
	case context.DeadlineExceeded:
		responder.ErrorWithStatus(w, r, http.StatusRequestTimeout, err)
	case rpctypes.ErrEmptyKey:
		responder.BadRequest(w, r, err)
	default:
		responder.InternalServerError(w, r, err)
	}
}
