package tests

import (
	"testing"

	"microdrive_auth/tests/suite"

	ssov1 "microdrive_auth/gen/go"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshToken_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	loginResp, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)

	oldRefresh := loginResp.GetRefreshToken()
	require.NotEmpty(t, oldRefresh)

	refreshResp, err := st.AuthClient.RefreshToken(ctx, &ssov1.RefreshTokenRequest{
		RefreshToken: oldRefresh,
		AppId:        appID,
	})
	require.NoError(t, err)

	require.NotEmpty(t, refreshResp.GetAccessToken())
	require.NotEmpty(t, refreshResp.GetRefreshToken())

	assert.NotEqual(t, oldRefresh, refreshResp.GetRefreshToken())
}
