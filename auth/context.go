package auth

import (
	"context"
	"fmt"

	"github.com/arseniisemenow/s21auto-client-go/util"

	"github.com/go-resty/resty/v2"
)

const contextInfoUrl = "https://platform.21-school.ru/services/rest/edu-context/context-info"

type contextInfoResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ContextHeaders map[string]string `json:"contextHeaders"`
	} `json:"data"`
	Error interface{} `json:"error"`
}

type ContextHeaders struct {
	XEDUSchoolID  string
	XEDUProductID string
	XEDUOrgUnitID string
	XEDURouteInfo string
}

func RequestContextHeaders(token Token, ctx context.Context) (ContextHeaders, error) {
	client := resty.New()

	res, err := client.R().
		SetContext(ctx).
		SetAuthToken(token.AccessToken).
		SetHeader("Origin", "https://platform.21-school.ru").
		SetHeader("Referer", "https://platform.21-school.ru/").
		Get(contextInfoUrl)

	if err != nil {
		return ContextHeaders{}, fmt.Errorf("context-info request failed: %w", err)
	}

	if res.StatusCode() != 200 {
		return ContextHeaders{}, fmt.Errorf("context-info request failed with status %d: %s", res.StatusCode(), res.String())
	}

	contextInfo, err := util.UnmarshalJson[contextInfoResponse](res.Body())
	if err != nil {
		return ContextHeaders{}, fmt.Errorf("failed to parse context-info response: %w", err)
	}

	if !contextInfo.Success {
		return ContextHeaders{}, fmt.Errorf("context-info request returned unsuccessful response")
	}

	headers := contextInfo.Data.ContextHeaders
	return ContextHeaders{
		XEDUSchoolID:  headers["X-EDU-SCHOOL-ID"],
		XEDUProductID: headers["X-EDU-PRODUCT-ID"],
		XEDUOrgUnitID: headers["X-EDU-ORG-UNIT-ID"],
		XEDURouteInfo: headers["X-EDU-ROUTE-INFO"],
	}, nil
}
