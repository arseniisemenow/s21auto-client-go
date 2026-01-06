package s21client

import "context"

type ContextHeaders struct {
	XEDUSchoolID    string
	XEDUProductID   string
	XEDUOrgUnitID   string
	XEDURouteInfo   string
}

type AuthCredentials struct {
	Token          string
	SchoolId       string
	ContextHeaders *ContextHeaders
}

type AuthProvider interface {
	GetAuthCredentials(ctx context.Context) (AuthCredentials, error)
}
