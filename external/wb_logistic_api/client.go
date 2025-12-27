package wb_logistic_api

import (
	"context"
	"fmt"
	"wb_logistic_assistant/external/wb_logistic_api/models"
	"wb_logistic_assistant/external/wb_logistic_api/request"
	"wb_logistic_assistant/external/wb_logistic_api/response"
	"wb_logistic_assistant/external/wb_logistic_api/session"
	"wb_logistic_assistant/external/wb_logistic_api/transport"
)

type Client struct {
	client transport.HTTPClient
	auth   *session.AuthService
}

func NewClient(client transport.HTTPClient) *Client {
	return &Client{
		client: client,
		auth:   session.NewAuthService(client),
	}
}

func (c *Client) getSessionToken(ctx context.Context, s *session.Session) (string, error) {
	if s == nil {
		return "", fmt.Errorf("session is nil")
	}
	if !s.IsAuth() {
		return "", fmt.Errorf("session is not auth")
	}

	if s.SessionTokenExpired() {
		err := c.RefreshSession(ctx, s)
		if err != nil {
			return "", err
		}
	}

	return s.SessionTokenString(), nil
}

// API Auth

func (c *Client) RequestAuthCode(ctx context.Context, login string) (*models.AuthCode, error) {
	code, err := c.auth.RequestAuthCode(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("failed request auth code: %w", err)
	}
	return code, nil
}

func (c *Client) ExchangeAuthCode(ctx context.Context, code int, sticker string) (*models.AuthAccessToken, error) {
	token, err := c.auth.ExchangeCode(ctx, code, sticker)
	if err != nil {
		return nil, fmt.Errorf("failed exchange code: %w", err)
	}
	return token, nil
}

func (c *Client) GetSessionToken(ctx context.Context, login, accessToken string) (*models.AuthSessionToken, *models.UserInfo, error) {
	if accessToken == "" {
		return nil, nil, fmt.Errorf("access token is empty")
	}
	if len(login) < 11 {
		return nil, nil, fmt.Errorf("login '%s' is too short, len: %d", login, len(login))
	}
	if login[0] == '+' {
		login = login[1:]
	}

	token, err := c.auth.GetSessionToken(ctx, login, accessToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed get session token: %w", err)
	}

	info, err := c.auth.GetUserInfo(ctx, token.Token.AccessToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed get user info: %w", err)
	}

	return token, info, nil
}

func (c *Client) RefreshAccessToken(ctx context.Context, refreshToken string) (*models.AuthAccessToken, error) {
	return c.auth.RefreshAccessToken(ctx, refreshToken)
}

func (c *Client) RefreshSession(ctx context.Context, session *session.Session) error {
	if session == nil {
		return fmt.Errorf("session is nil")
	}
	if session.Login() == "" || session.RefreshToken() == "" {
		return fmt.Errorf("invalid session")
	}
	accessToken, err := c.auth.RefreshAccessToken(ctx, session.RefreshToken())
	if err != nil {
		return fmt.Errorf("failed refresh session: %w", err)
	}

	sessionToken, userInfo, err := c.GetSessionToken(ctx, session.Login(), accessToken.AccessToken)
	if err != nil {
		return fmt.Errorf("failed get session token: %w", err)
	}

	session.SetAccessToken(accessToken)
	session.SetSessionToken(sessionToken)
	session.SetUserInfo(userInfo)

	return nil
}

// API Methods

func (c *Client) GetUserInfo(ctx context.Context, s *session.Session) (*response.UserGetInfoResponse, error) {
	token, err := c.getSessionToken(ctx, s)
	if err != nil {
		return nil, err
	}

	userInfo := s.UserInfo()
	if userInfo == nil {
		info, err := c.auth.GetUserInfo(ctx, s.AccessTokenString())
		if err != nil {
			return nil, fmt.Errorf("failed get user info: %w", err)
		}
		s.SetUserInfo(info)
	}
	if userInfo == nil || userInfo.UserDetails == nil {
		return nil, fmt.Errorf("user info is missing in session")
	}

	req := request.NewUserGetInfoRequest(c.client, token).
		ClientID(userInfo.ID)

	res, err := req.Do(ctx)
	if err != nil {
		if req.IsUnauthorized() {
			if err = c.RefreshSession(ctx, s); err != nil {
				return nil, fmt.Errorf("failed refresh unauthorized session %w", err)
			}
			res, err = req.Do(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed retry request after refresh session: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to do request: %w", err)
		}
	}
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get user info: %w", res.Error)
	}

	return &res, nil
}

func (c *Client) GetRemainsLastMileReports(ctx context.Context, s *session.Session) (*response.GetRemainsLastMileReportsResponse, error) {
	token, err := c.getSessionToken(ctx, s)
	if err != nil {
		return nil, err
	}

	req := request.NewGetRemainsLastMileReportsRequest(c.client, token)
	res, err := req.Do(ctx)
	if err != nil {
		if req.IsUnauthorized() {
			if err = c.RefreshSession(ctx, s); err != nil {
				return nil, fmt.Errorf("failed refresh unauthorized session %w", err)
			}
			res, err = req.Do(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed retry request after refresh session: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to do request: %w", err)
		}
	}
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get route reports: %s", res.Error.Error())
	}

	return &res, nil
}

func (c *Client) GetRemainsLastMileReportsRouteInfo(ctx context.Context, s *session.Session, routeID int) (*response.GetRemainsLastMileReportsRouteInfoResponse, error) {
	token, err := c.getSessionToken(ctx, s)
	if err != nil {
		return nil, err
	}

	req := request.NewGetRemainsLastMileReportsInfoRequest(c.client, token).
		RouteID(routeID)

	res, err := req.Do(ctx)
	if err != nil {
		if req.IsUnauthorized() {
			if err = c.RefreshSession(ctx, s); err != nil {
				return nil, fmt.Errorf("failed refresh unauthorized session %w", err)
			}
			res, err = req.Do(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed retry request after refresh session: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to do request: %w", err)
		}
	}
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get route reports: %s", res.Error.Error())
	}

	return &res, nil
}

func (c *Client) GetJobsScheduling(ctx context.Context, s *session.Session) (*response.GetJobsSchedulingResponse, error) {
	token, err := c.getSessionToken(ctx, s)
	if err != nil {
		return nil, err
	}

	userInfo := s.UserInfo()
	if userInfo == nil {
		info, err := c.auth.GetUserInfo(ctx, s.AccessTokenString())
		if err != nil {
			return nil, fmt.Errorf("failed get user info: %w", err)
		}
		s.SetUserInfo(info)
	}
	if userInfo == nil || userInfo.UserDetails == nil {
		return nil, fmt.Errorf("user info is missing in session")
	}

	req := request.NewGetJobsSchedulingRequest(c.client, token).
		SupplierID(userInfo.UserDetails.SupplierID)

	res, err := req.Do(ctx)
	if err != nil {
		if req.IsUnauthorized() {
			if err = c.RefreshSession(ctx, s); err != nil {
				return nil, fmt.Errorf("failed refresh unauthorized session %w", err)
			}
			res, err = req.Do(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed retry request after refresh session: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to do request: %w", err)
		}
	}
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get jobs scheduling: %s", res.Error.Error())
	}

	return &res, nil
}

func (c *Client) GetShipments(ctx context.Context, s *session.Session, params *models.GetShipmentParamsRequest) (*response.GetShipmentsResponse, error) {
	token, err := c.getSessionToken(ctx, s)
	if err != nil {
		return nil, err
	}
	if s.UserInfo() == nil || s.UserInfo().UserDetails == nil {
		return nil, fmt.Errorf("user info is missing in session")
	}

	req := request.NewGetShipmentsRequest(c.client, token).
		SupplierID(s.UserInfo().UserDetails.SupplierID).
		FromParams(params)

	res, err := req.Do(ctx)
	if err != nil {
		if req.IsUnauthorized() {
			if err = c.RefreshSession(ctx, s); err != nil {
				return nil, fmt.Errorf("failed refresh unauthorized session %w", err)
			}
			res, err = req.Do(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed retry request after refresh session: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to do request: %w", err)
		}
	}
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get shipments: %s", res.Error.Error())
	}

	return &res, nil
}

func (c *Client) GetShipmentInfo(ctx context.Context, s *session.Session, shipmentID int) (*response.GetShipmentInfoResponse, error) {
	token, err := c.getSessionToken(ctx, s)
	if err != nil {
		return nil, err
	}

	req := request.NewGetShipmentInfoRequest(c.client, token).
		ShipmentID(shipmentID)

	res, err := req.Do(ctx)
	if err != nil {
		if req.IsUnauthorized() {
			if err = c.RefreshSession(ctx, s); err != nil {
				return nil, fmt.Errorf("failed refresh unauthorized session %w", err)
			}
			res, err = req.Do(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed retry request after refresh session: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to do request: %w", err)
		}
	}
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get shipments: %s", res.Error.Error())
	}

	return &res, nil
}

func (c *Client) GetShipmentTransfers(ctx context.Context, s *session.Session, shipmentID int) (*response.GetShipmentTransfersResponse, error) {
	token, err := c.getSessionToken(ctx, s)
	if err != nil {
		return nil, err
	}

	req := request.NewGetShipmentTransfersRequest(c.client, token).
		ShipmentID(shipmentID)

	res, err := req.Do(ctx)
	if err != nil {
		if req.IsUnauthorized() {
			if err = c.RefreshSession(ctx, s); err != nil {
				return nil, fmt.Errorf("failed refresh unauthorized session %w", err)
			}
			res, err = req.Do(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed retry request after refresh session: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to do request: %w", err)
		}
	}
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get shipments: %s", res.Error.Error())
	}

	return &res, nil
}

func (c *Client) GetTaresForOffices(ctx context.Context, s *session.Session, officeID int, destinationOfficeIDs []int, isDrive bool) (*response.GetTaresForOfficesResponse, error) {
	token, err := c.getSessionToken(ctx, s)
	if err != nil {
		return nil, err
	}

	req := request.NewGetTaresForOffices(c.client, token).
		SourceOfficeID(officeID).
		DestinationOfficeIDs(destinationOfficeIDs).
		IsDrive(isDrive)

	res, err := req.Do(ctx)
	if err != nil {
		if req.IsUnauthorized() {
			if err = c.RefreshSession(ctx, s); err != nil {
				return nil, fmt.Errorf("failed refresh unauthorized session %w", err)
			}
			res, err = req.Do(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed retry request after refresh session: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to do request: %w", err)
		}
	}
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get shipments: %s", res.Error.Error())
	}

	return &res, nil
}

func (c *Client) GetAssociationOfficesInfoByName(ctx context.Context, s *session.Session, param string, isDc bool) (*response.GetAssociationOfficesInfoByNameResponse, error) {
	token, err := c.getSessionToken(ctx, s)
	if err != nil {
		return nil, err
	}

	req := request.NewGetAssociationOfficesInfoByNameRequest(c.client, token).
		Param(param).
		IsDc(isDc)

	res, err := req.Do(ctx)
	if err != nil {
		if req.IsUnauthorized() {
			if err = c.RefreshSession(ctx, s); err != nil {
				return nil, fmt.Errorf("failed refresh unauthorized session %w", err)
			}
			res, err = req.Do(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed retry request after refresh session: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to do request: %w", err)
		}
	}
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get association offices info by name: %s", res.Error.Error())
	}

	return &res, nil
}

func (c *Client) GetAssociationRoutesInfoByName(ctx context.Context, s *session.Session, param string) (*response.GetAssociationRoutesInfoByNameResponse, error) {
	token, err := c.getSessionToken(ctx, s)
	if err != nil {
		return nil, err
	}

	req := request.NewGetAssociationRoutesInfoByNameRequest(c.client, token).
		Param(param)

	res, err := req.Do(ctx)
	if err != nil {
		if req.IsUnauthorized() {
			if err = c.RefreshSession(ctx, s); err != nil {
				return nil, fmt.Errorf("failed refresh unauthorized session %w", err)
			}
			res, err = req.Do(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed retry request after refresh session: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to do request: %w", err)
		}
	}
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get association routes info by name: %s", res.Error.Error())
	}

	return &res, nil
}

func (c *Client) GetWaySheets(ctx context.Context, s *session.Session, params *models.GetWaySheetsParamsRequest) (*response.GetWaySheetsResponse, error) {
	token, err := c.getSessionToken(ctx, s)
	if err != nil {
		return nil, err
	}

	req := request.NewGetWaySheetsRequestRequest(c.client, token).
		FromParams(params)

	res, err := req.Do(ctx)
	if err != nil {
		if req.IsUnauthorized() {
			if err = c.RefreshSession(ctx, s); err != nil {
				return nil, fmt.Errorf("failed refresh unauthorized session %w", err)
			}
			res, err = req.Do(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed retry request after refresh session: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to do request: %w", err)
		}
	}
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get way sheets: %s", res.Error.Error())
	}

	return &res, nil
}

func (c *Client) GetWaySheetInfo(ctx context.Context, s *session.Session, waySheetID int) (*response.GetWaySheetInfoResponse, error) {
	token, err := c.getSessionToken(ctx, s)
	if err != nil {
		return nil, err
	}

	req := request.NewGetWaySheetInfoRequest(c.client, token).
		WaySheetID(waySheetID)

	res, err := req.Do(ctx)
	if err != nil {
		if req.IsUnauthorized() {
			if err = c.RefreshSession(ctx, s); err != nil {
				return nil, fmt.Errorf("failed refresh unauthorized session %w", err)
			}
			res, err = req.Do(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed retry request after refresh session: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to do request: %w", err)
		}
	}
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get way sheet info: %s", res.Error.Error())
	}
	if res.Data.WaySheet == nil {
		return nil, fmt.Errorf("failed to get way sheet info, data way sheet info is empty")
	}

	return &res, nil
}

func (c *Client) GetWaySheetFinanceDetails(ctx context.Context, s *session.Session, waySheetID int) (*response.GetWaySheetFinanceDetailsResponse, error) {
	token, err := c.getSessionToken(ctx, s)
	if err != nil {
		return nil, err
	}

	req := request.NewGetWaySheetFinanceDetailsRequest(c.client, token).
		WaySheetID(waySheetID)

	res, err := req.Do(ctx)
	if err != nil {
		if req.IsUnauthorized() {
			if err = c.RefreshSession(ctx, s); err != nil {
				return nil, fmt.Errorf("failed refresh unauthorized session %w", err)
			}
			res, err = req.Do(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed retry request after refresh session: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to do request: %w", err)
		}
	}
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get shipments: %s", res.Error.Error())
	}

	return &res, nil
}
