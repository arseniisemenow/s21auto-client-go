package s21client

import (
	"context"

	"github.com/arseniisemenow/s21auto-client-go/auth"
)

type defaultAuthProvider struct {
	token          auth.Token
	schoolId       string
	contextHeaders *ContextHeaders
}

func (provider *defaultAuthProvider) refreshCredentials(ctx context.Context) (err error) {
	err = provider.token.Refresh(ctx)

	if err != nil {
		return err
	}

	if provider.schoolId == "" {
		user, err := auth.RequestUserData(provider.token, ctx)

		if err != nil {
			return err
		}

		provider.schoolId = user.Roles[0].SchoolID
	}

	if provider.contextHeaders == nil {
		headers, err := auth.RequestContextHeaders(provider.token, ctx)
		if err != nil {
			return err
		}
		provider.contextHeaders = &ContextHeaders{
			XEDUSchoolID:  headers.XEDUSchoolID,
			XEDUProductID: headers.XEDUProductID,
			XEDUOrgUnitID: headers.XEDUOrgUnitID,
			XEDURouteInfo: headers.XEDURouteInfo,
		}
	}

	return
}

func (provider *defaultAuthProvider) GetAuthCredentials(ctx context.Context) (credentials AuthCredentials, err error) {
	err = provider.refreshCredentials(ctx)

	if err != nil {
		return
	}

	credentials = AuthCredentials{
		Token:          provider.token.AccessToken,
		SchoolId:       provider.schoolId,
		ContextHeaders: provider.contextHeaders,
	}

	return
}

func DefaultAuth(username, password string) *defaultAuthProvider {
	return &defaultAuthProvider{
		token: auth.Token{
			Username: username,
			Password: password,
		},
	}
}
