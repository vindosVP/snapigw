package auth

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/vindosVP/snapigw/gen/go"
)

type Client struct {
	grpc authv1.AuthClient
}

func (c Client) SetBanned(ctx context.Context, userId int64, isBanned bool) (bool, error) {
	req := &authv1.SetBannedRequest{
		UserId:   userId,
		IsBanned: isBanned,
	}
	res, err := c.grpc.SetBanned(ctx, req)
	if err != nil {
		return false, err
	}
	return res.IsBanned, nil
}

func (c Client) SetDeleted(ctx context.Context, userId int64, isDeleted bool) (bool, error) {
	req := &authv1.SetDeletedRequest{
		UserId:    userId,
		IsDeleted: isDeleted,
	}
	res, err := c.grpc.SetDeleted(ctx, req)
	if err != nil {
		return false, err
	}
	return res.IsDeleted, nil
}

func (c Client) SetAdmin(ctx context.Context, userId int64, isAdmin bool) (bool, error) {
	req := &authv1.SetAdminRightsRequest{
		UserId:  userId,
		IsAdmin: isAdmin,
	}
	res, err := c.grpc.SetAdminRights(ctx, req)
	if err != nil {
		return false, err
	}
	return res.IsAdmin, nil
}

func (c Client) Register(ctx context.Context, email, password string) (int64, error) {
	req := &authv1.RegisterRequest{
		Email:    email,
		Password: password,
	}
	res, err := c.grpc.Register(ctx, req)
	if err != nil {
		return 0, err
	}
	return res.UserId, nil
}

func (c Client) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	req := &authv1.LoginRequest{
		Email:    email,
		Password: password,
	}
	res, err := c.grpc.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	tp := &TokenPair{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}
	return tp, nil
}

func (c Client) RefreshToken(ctx context.Context, token string) (*TokenPair, error) {
	req := &authv1.RefreshRequest{
		RefreshToken: token,
	}
	res, err := c.grpc.Refresh(ctx, req)
	if err != nil {
		return nil, err
	}
	tp := &TokenPair{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}
	return tp, nil
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to auth service")
	}
	client := authv1.NewAuthClient(conn)
	return &Client{client}, nil
}
