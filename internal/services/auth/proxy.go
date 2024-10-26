package auth

import (
	"net/http"

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
