package auth

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/vindosVP/snapigw/internal/utils/response"
)

type Proxy struct {
	client *Client
	l      zerolog.Logger
}

func (p *Proxy) SetAdminHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		reqId := c.GetString("requestId")
		lg := p.l.With().Str("requestId", reqId).Logger()
		req := &SetAdminRequest{}
		err := c.BindJSON(req)
		if err != nil {
			lg.Info().Msg("invalid request structure")
			response.Err(c, http.StatusBadRequest, "invalid request structure")
			return
		}
		err = validator.New().Struct(req)
		if err != nil {
			lg.Info().Msg("invalid request")
			response.Err(c, http.StatusBadRequest, "validation error")
			return
		}
		userIdParam := c.Param("id")
		if userIdParam == "" {
			lg.Info().Msg("invalid user id")
			response.Err(c, http.StatusBadRequest, "invalid user id")
			return
		}
		userId, err := strconv.Atoi(userIdParam)
		if err != nil {
			lg.Info().Msg("invalid user id")
			response.Err(c, http.StatusBadRequest, "invalid user id")
			return
		}
		if userId == c.GetInt("userId") {
			lg.Info().Msg("user can not set admin flag to himself")
			response.Err(c, http.StatusBadRequest, "user can not set admin flag to himself")
			return
		}

		meta := map[string]string{"requestId": reqId}
		ctx := metadata.NewOutgoingContext(c, metadata.New(meta))
		admin, err := p.client.SetAdmin(ctx, int64(userId), *req.IsAdmin)
		if err != nil {
			s, ok := status.FromError(err)
			if !ok {
				lg.Error().Stack().Msg("failed to create error from code")
				response.Err(c, http.StatusInternalServerError, "failed to set admin flag")
				return
			}
			switch s.Code() {
			case codes.InvalidArgument:
				lg.Error().Msg("no requestId specified")
				response.Err(c, http.StatusInternalServerError, "failed to set admin flag")
			case codes.FailedPrecondition:
				response.Err(c, http.StatusBadRequest, "user does not exist")
			default:
				response.Err(c, http.StatusInternalServerError, "failed to set admin flag")
			}
			return
		}
		response.OkMsg(c, http.StatusOK, &SetAdminResponse{IsAdmin: admin}, "admin flag set successfully")
	}
}

func (p *Proxy) SetDeletedHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		reqId := c.GetString("requestId")
		lg := p.l.With().Str("requestId", reqId).Logger()
		req := &SetDeletedRequest{}
		err := c.BindJSON(req)
		if err != nil {
			lg.Info().Msg("invalid request structure")
			response.Err(c, http.StatusBadRequest, "invalid request structure")
			return
		}
		err = validator.New().Struct(req)
		if err != nil {
			lg.Info().Msg("invalid request")
			response.Err(c, http.StatusBadRequest, "validation error")
			return
		}
		userIdParam := c.Param("id")
		if userIdParam == "" {
			lg.Info().Msg("invalid user id")
			response.Err(c, http.StatusBadRequest, "invalid user id")
			return
		}
		userId, err := strconv.Atoi(userIdParam)
		if err != nil {
			lg.Info().Msg("invalid user id")
			response.Err(c, http.StatusBadRequest, "invalid user id")
			return
		}
		if userId == c.GetInt("userId") {
			lg.Info().Msg("user can not set deleted flag to himself")
			response.Err(c, http.StatusBadRequest, "user can not set deleted flag to himself")
			return
		}

		meta := map[string]string{"requestId": reqId}
		ctx := metadata.NewOutgoingContext(c, metadata.New(meta))
		deleted, err := p.client.SetDeleted(ctx, int64(userId), *req.IsDeleted)
		if err != nil {
			s, ok := status.FromError(err)
			if !ok {
				lg.Error().Stack().Msg("failed to create error from code")
				response.Err(c, http.StatusInternalServerError, "refresh failed")
				return
			}
			switch s.Code() {
			case codes.InvalidArgument:
				lg.Error().Msg("no requestId specified")
				response.Err(c, http.StatusInternalServerError, "failed to set deleted flag")
			case codes.FailedPrecondition:
				response.Err(c, http.StatusBadRequest, "user does not exist")
			default:
				response.Err(c, http.StatusInternalServerError, "failed to set deleted flag")
			}
			return
		}
		response.OkMsg(c, http.StatusOK, &SetDeletedResponse{IsDeleted: deleted}, "set deleted flag successfully")
	}
}

func (p *Proxy) SetBannedHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		reqId := c.GetString("requestId")
		lg := p.l.With().Str("requestId", reqId).Logger()
		req := &SetBannedRequest{}
		err := c.BindJSON(req)
		if err != nil {
			lg.Info().Msg("invalid request structure")
			response.Err(c, http.StatusBadRequest, "invalid request structure")
			return
		}
		err = validator.New().Struct(req)
		if err != nil {
			lg.Info().Msg("invalid request")
			response.Err(c, http.StatusBadRequest, err.Error())
			return
		}
		userIdParam := c.Param("id")
		if userIdParam == "" {
			lg.Info().Msg("invalid user id")
			response.Err(c, http.StatusBadRequest, "invalid user id")
			return
		}
		userId, err := strconv.Atoi(userIdParam)
		if err != nil {
			lg.Info().Msg("invalid user id")
			response.Err(c, http.StatusBadRequest, "invalid user id")
			return
		}
		if userId == c.GetInt("userId") {
			lg.Info().Msg("user can not set ban flag to himself")
			response.Err(c, http.StatusBadRequest, "user can not set ban flag to himself")
			return
		}

		meta := map[string]string{"requestId": reqId}
		ctx := metadata.NewOutgoingContext(c, metadata.New(meta))
		banned, err := p.client.SetBanned(ctx, int64(userId), *req.IsBanned)
		if err != nil {
			s, ok := status.FromError(err)
			if !ok {
				lg.Error().Stack().Msg("failed to create error from code")
				response.Err(c, http.StatusInternalServerError, "failed to set banned flag")
				return
			}
			switch s.Code() {
			case codes.InvalidArgument:
				lg.Error().Msg("no requestId specified")
				response.Err(c, http.StatusInternalServerError, "failed to set banned flag")
			case codes.FailedPrecondition:
				response.Err(c, http.StatusBadRequest, "user does not exist")
			default:
				response.Err(c, http.StatusInternalServerError, "failed to set banned flag")
			}
			return
		}
		response.OkMsg(c, http.StatusOK, &SetBannedResponse{IsBanned: banned}, "set banned flag successfully")
	}
}

func (p *Proxy) RefreshHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		reqId := c.GetString("requestId")
		req := &RefreshRequest{}
		err := c.BindJSON(req)
		if err != nil {
			log.Info().Str("requestId", reqId).Msg("invalid request structure")
			response.Err(c, http.StatusBadRequest, "invalid request structure")
			return
		}
		err = validator.New().Struct(req)
		if err != nil {
			log.Info().Str("requestId", reqId).Msg("invalid request")
			response.Err(c, http.StatusBadRequest, err.Error())
			return
		}

		meta := map[string]string{"requestId": reqId}
		ctx := metadata.NewOutgoingContext(c, metadata.New(meta))
		tp, err := p.client.RefreshToken(ctx, req.RefreshToken)
		if err != nil {
			s, ok := status.FromError(err)
			if !ok {
				log.Error().Stack().Msg("failed to create error from code")
				response.Err(c, http.StatusInternalServerError, "refresh failed")
				return
			}
			switch s.Code() {
			case codes.InvalidArgument:
				log.Error().Msg("no requestId specified")
				response.Err(c, http.StatusInternalServerError, "refresh failed")
			case codes.FailedPrecondition:
				response.Err(c, http.StatusBadRequest, "user is unable to log in")
			default:
				response.Err(c, http.StatusInternalServerError, "refresh failed")
			}
			return
		}
		response.OkMsg(c, http.StatusOK, &RefreshResponse{AccessToken: tp.AccessToken, RefreshToken: tp.RefreshToken}, "refresh success")
	}
}

func (p *Proxy) LoginHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		reqId := c.GetString("requestId")
		req := &LoginRequest{}
		err := c.BindJSON(req)
		if err != nil {
			log.Info().Str("requestId", reqId).Msg("invalid request structure")
			response.Err(c, http.StatusBadRequest, "invalid request structure")
			return
		}
		err = validator.New().Struct(req)
		if err != nil {
			log.Info().Str("requestId", reqId).Msg("invalid request")
			response.Err(c, http.StatusBadRequest, err.Error())
			return
		}

		meta := map[string]string{"requestId": reqId}
		ctx := metadata.NewOutgoingContext(c, metadata.New(meta))
		tp, err := p.client.Login(ctx, req.Email, req.Password)
		if err != nil {
			s, ok := status.FromError(err)
			if !ok {
				log.Error().Stack().Msg("failed to create error from code")
				response.Err(c, http.StatusInternalServerError, "login failed")
				return
			}
			switch s.Code() {
			case codes.InvalidArgument:
				response.Err(c, http.StatusBadRequest, "invalid login or password")
			case codes.FailedPrecondition:
				response.Err(c, http.StatusBadRequest, "user is unable to log in")
			default:
				response.Err(c, http.StatusInternalServerError, "login failed")
			}
			return
		}
		response.OkMsg(c, http.StatusOK, &LoginResponse{AccessToken: tp.AccessToken, RefreshToken: tp.RefreshToken}, "login success")
	}
}

func (p *Proxy) RegisterHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		reqId := c.GetString("requestId")
		req := &RegisterRequest{}
		err := c.BindJSON(req)
		if err != nil {
			log.Info().Str("requestId", reqId).Msg("invalid request structure")
			response.Err(c, http.StatusBadRequest, "invalid request structure")
			return
		}
		err = validator.New().Struct(req)
		if err != nil {
			log.Info().Str("requestId", reqId).Msg("invalid request")
			response.Err(c, http.StatusBadRequest, err.Error())
			return
		}

		meta := map[string]string{"requestId": reqId}
		ctx := metadata.NewOutgoingContext(c, metadata.New(meta))
		id, err := p.client.Register(ctx, req.Email, req.Password)
		if err != nil {
			s, ok := status.FromError(err)
			if !ok {
				log.Error().Stack().Msg("failed to create error from code")
				response.Err(c, http.StatusInternalServerError, "register failed")
				return
			}
			switch s.Code() {
			case codes.InvalidArgument:
				log.Error().Msg("no requestId specified")
				response.Err(c, http.StatusInternalServerError, "register failed")
			case codes.FailedPrecondition:
				response.Err(c, http.StatusBadRequest, "user already exists")
			default:
				response.Err(c, http.StatusInternalServerError, "register failed")
			}
			return
		}
		response.OkMsg(c, http.StatusOK, &RegisterResponse{UserId: id}, "register success")
	}
}

func NewProxy(serviceAddr string, l zerolog.Logger) (*Proxy, error) {
	c, err := NewClient(serviceAddr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize auth client")
	}
	return &Proxy{client: c, l: l}, nil
}
