package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/chishkin-afk/posted/http-gateway/internal/application/dtos"
	"github.com/chishkin-afk/posted/http-gateway/pkg/errs"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authService interface {
	Register(ctx context.Context, req *dtos.RegisterRequest) (*dtos.Token, error)
	Login(ctx context.Context, req *dtos.LoginRequest) (*dtos.Token, error)
	UpdateUser(ctx context.Context, req *dtos.UpdateUserRequest, token string) (*dtos.User, error)
	DeleteUser(ctx context.Context, token string) error
	GetUserSelf(ctx context.Context, token string) (*dtos.User, error)
	GetUserByID(ctx context.Context, id string) (*dtos.User, error)
}

type postsService interface {
	Create(ctx context.Context, req *dtos.CreatePostRequest, token string) (*dtos.Post, error)
	Update(ctx context.Context, req *dtos.UpdatePostRequest, token string) (*dtos.Post, error)
	Delete(ctx context.Context, id string, token string) error
	GetByID(ctx context.Context, id string) (*dtos.Post, error)
	GetSelfPosts(ctx context.Context, page, limit uint32, token string) (*dtos.Posts, error)
}

type handlers struct {
	authService  authService
	postsService postsService
}

func (h *handlers) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.RegisterRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		resp, err := h.authService.Register(ctx.Request.Context(), &req)
		if err != nil {
			code, cleanErr := h.getCode(err)
			ctx.JSON(code, cleanErr)
			return
		}

		ctx.JSON(http.StatusCreated, resp)
	}
}

func (h *handlers) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.LoginRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		resp, err := h.authService.Login(ctx.Request.Context(), &req)
		if err != nil {
			code, cleanErr := h.getCode(err)
			ctx.JSON(code, cleanErr)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (h *handlers) UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := h.getToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		var req dtos.UpdateUserRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		resp, err := h.authService.UpdateUser(ctx.Request.Context(), &req, token)
		if err != nil {
			code, cleanErr := h.getCode(err)
			ctx.JSON(code, cleanErr)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (h *handlers) DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := h.getToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		if err := h.authService.DeleteUser(ctx.Request.Context(), token); err != nil {
			code, cleanErr := h.getCode(err)
			ctx.JSON(code, cleanErr)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (h *handlers) GetUserSelf() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := h.getToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		resp, err := h.authService.GetUserSelf(ctx.Request.Context(), token)
		if err != nil {
			code, cleanErr := h.getCode(err)
			ctx.JSON(code, cleanErr)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (h *handlers) GetUserByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid id",
			})
			return
		}

		resp, err := h.authService.GetUserByID(ctx.Request.Context(), id)
		if err != nil {
			code, cleanErr := h.getCode(err)
			ctx.JSON(code, cleanErr)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (h *handlers) CreatePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := h.getToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		var req dtos.CreatePostRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		resp, err := h.postsService.Create(ctx.Request.Context(), &req, token)
		if err != nil {
			code, cleanErr := h.getCode(err)
			ctx.JSON(code, cleanErr)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (h *handlers) UpdatePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := h.getToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		var req dtos.UpdatePostRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		resp, err := h.postsService.Update(ctx.Request.Context(), &req, token)
		if err != nil {
			code, cleanErr := h.getCode(err)
			ctx.JSON(code, cleanErr)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (h *handlers) DeletePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := h.getToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid id",
			})
			return
		}

		if err := h.postsService.Delete(ctx.Request.Context(), id, token); err != nil {
			code, cleanErr := h.getCode(err)
			ctx.JSON(code, cleanErr)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (h *handlers) GetPostByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid id",
			})
			return
		}

		resp, err := h.postsService.GetByID(ctx.Request.Context(), id)
		if err != nil {
			code, cleanErr := h.getCode(err)
			ctx.JSON(code, cleanErr)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (h *handlers) GetSelfPosts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := h.getToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		page, err := strconv.Atoi(ctx.Query("page"))
		if err != nil || page < 1 {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid page",
			})
			return
		}

		limit, err := strconv.Atoi(ctx.Query("size"))
		if err != nil || limit > 100 || limit < 1 {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid page",
			})
			return
		}

		resp, err := h.postsService.GetSelfPosts(ctx.Request.Context(), uint32(page), uint32(limit), token)
		if err != nil {
			code, cleanErr := h.getCode(err)
			ctx.JSON(code, cleanErr)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (h *handlers) getCode(err error) (int, *dtos.ErrMsg) {
	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, &dtos.ErrMsg{
			Error: errs.ErrInternalServer.Error(),
		}
	}

	switch st.Code() {
	case codes.AlreadyExists:
		return http.StatusConflict, &dtos.ErrMsg{
			Error: st.Message(),
		}
	case codes.DeadlineExceeded,
		codes.Canceled:
		return http.StatusGatewayTimeout, &dtos.ErrMsg{
			Error: st.Message(),
		}
	case codes.InvalidArgument:
		return http.StatusBadRequest, &dtos.ErrMsg{
			Error: st.Message(),
		}
	case codes.PermissionDenied:
		return http.StatusForbidden, &dtos.ErrMsg{
			Error: st.Message(),
		}
	case codes.NotFound:
		return http.StatusNotFound, &dtos.ErrMsg{
			Error: st.Message(),
		}
	case codes.Unauthenticated:
		return http.StatusUnauthorized, &dtos.ErrMsg{
			Error: st.Message(),
		}
	}

	return http.StatusInternalServerError, &dtos.ErrMsg{
		Error: errs.ErrInternalServer.Error(),
	}
}

func (h *handlers) getToken(ctx *gin.Context) (string, error) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		return "", errs.ErrInvalidToken
	}

	return token, nil
}
