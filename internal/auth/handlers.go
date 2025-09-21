package auth

import (
	"net/http"

	"github.com/julianstephens/feature-flag-service/internal/rbac/users"
	"github.com/julianstephens/go-utils/httputil/request"
	"github.com/julianstephens/go-utils/httputil/response"
)


func LoginHandler(authSvc *AuthClient, rbacUserService *users.RbacUserService, responder *response.Responder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := request.DecodeJSON(r, &req); err != nil {
			responder.BadRequest(w, r, err)
			return
		}

		_, err := rbacUserService.GetUserByEmail(req.Email)
		if err != nil {
			responder.Unauthorized(w, r, err)
			return
		}

		// if !authutils.CheckPasswordHash(req.Password, rbacUser.Password) {
		// 	responder.Unauthorized(w, r, err)
		// 	return
		// }


		// token, err := authSvc.GenerateToken(req.UserID, req.Username, req.Email, req.Roles, &req.CustomClaims)
		// if err != nil {
		// 	responder.InternalServerError(w, r, err)
		// 	return
		// }

		// res := LoginResponse{
		// 	Token: token,
		// }
		responder.OK(w, r, "hi")
	}
}