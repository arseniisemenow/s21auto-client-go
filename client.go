package s21client

import (
	"github.com/arseniisemenow/s21auto-client-go/requests"

	"github.com/go-resty/resty/v2"
)

var S21GqlUrl = "https://platform.21-school.ru/services/graphql"

type Client struct {
	gqlUrl string

	resty *resty.Client

	authProvider AuthProvider
}

func (client *Client) UseAuth(authProvider AuthProvider) {
	client.authProvider = authProvider
}

func (client *Client) applyAuth(request *resty.Request) (err error) {
	credentials, err := client.authProvider.GetAuthCredentials(request.Context())

	if err != nil {
		return
	}

	request.SetAuthToken(credentials.Token).SetHeader("schoolid", credentials.SchoolId)

	if credentials.ContextHeaders != nil {
		request.SetHeader("X-EDU-SCHOOL-ID", credentials.ContextHeaders.XEDUSchoolID)
		request.SetHeader("X-EDU-PRODUCT-ID", credentials.ContextHeaders.XEDUProductID)
		request.SetHeader("X-EDU-ORG-UNIT-ID", credentials.ContextHeaders.XEDUOrgUnitID)
		request.SetHeader("X-EDU-ROUTE-INFO", credentials.ContextHeaders.XEDURouteInfo)
	}

	return
}

func (client *Client) R() *requests.RequestContext {
	request := client.resty.R()

	return requests.NewRequestContext(request, client.gqlUrl)
}

func New(authProvider AuthProvider) *Client {
	client := &Client{
		gqlUrl:       S21GqlUrl,
		authProvider: authProvider,
		resty:        resty.New(),
	}

	// Set default headers for GraphQL API
	client.resty.SetHeader("Origin", "https://platform.21-school.ru")
	client.resty.SetHeader("Referer", "https://platform.21-school.ru/")

	client.resty.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		return client.applyAuth(r)
	})

	return client
}
