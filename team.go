package slack

import (
	"context"
	"net/url"
	"strconv"
)

const (
	DEFAULT_LOGINS_COUNT = 100
	DEFAULT_LOGINS_PAGE  = 1
)

type TeamResponse struct {
	Team TeamInfo `json:"team"`
	SlackResponse
}

type DiscoveryEnterpriseInfoResponse struct {
	Enterprise DiscoveryEnterpriseInfo `json:"enterprise"`
	SlackResponse
}

type DiscoveryEnterpriseInfo struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Domain      string                 `json:"domain"`
	EmailDomain string                 `json:"email_domain"`
	Icon        map[string]interface{} `json:"icon"`
	Teams       []TeamInfo             `json:"teams"`
}

type TeamInfo struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Domain      string                 `json:"domain"`
	EmailDomain string                 `json:"email_domain"`
	Icon        map[string]interface{} `json:"icon"`
}

type TeamProfileResponse struct {
	Profile TeamProfile `json:"profile"`
	SlackResponse
}

type TeamProfile struct {
	Fields []TeamProfileField `json:"fields"`
}

type TeamProfileField struct {
	ID             string          `json:"id"`
	Ordering       int             `json:"ordering"`
	Label          string          `json:"label"`
	Hint           string          `json:"hint"`
	Type           string          `json:"type"`
	PossibleValues []string        `json:"possible_values"`
	IsHidden       bool            `json:"is_hidden"`
	Options        map[string]bool `json:"options"`
}

type LoginResponse struct {
	Logins []Login `json:"logins"`
	Paging `json:"paging"`
	SlackResponse
}

type Login struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	DateFirst int    `json:"date_first"`
	DateLast  int    `json:"date_last"`
	Count     int    `json:"count"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	ISP       string `json:"isp"`
	Country   string `json:"country"`
	Region    string `json:"region"`
}

type IntegrationLogsResponse struct {
	Logs   []IntegrationLog `json:"logs"`
	Paging `json:"paging"`
	SlackResponse
}

type IntegrationLog struct{}

type BillableInfoResponse struct {
	BillableInfo map[string]BillingActive `json:"billable_info"`
	SlackResponse
}

type BillingActive struct {
	BillingActive bool `json:"billing_active"`
}

// AccessLogParameters contains all the parameters necessary (including the optional ones) for a GetAccessLogs() request
type AccessLogParameters struct {
	Count int
	Page  int
}

type IntegrationLogsParameters struct {
	AppID      string `json:"app_id"`
	ChangeType string `json:"change_type"`
	Count      int    `json:"count"`
	Page       int    `json:"page"`
	ServiceID  string `json:"service_id"`
	TeamID     string `json:"team_id"`
	User       string `json:"user"`
}

// NewAccessLogParameters provides an instance of AccessLogParameters with all the sane default values set
func NewAccessLogParameters() AccessLogParameters {
	return AccessLogParameters{
		Count: DEFAULT_LOGINS_COUNT,
		Page:  DEFAULT_LOGINS_PAGE,
	}
}

func (api *Client) teamRequest(ctx context.Context, path string, values url.Values) (*TeamResponse, error) {
	response := &TeamResponse{}
	err := api.postMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, response.Err()
}

func (api *Client) billableInfoRequest(ctx context.Context, path string, values url.Values) (map[string]BillingActive, error) {
	response := &BillableInfoResponse{}
	err := api.postMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response.BillableInfo, response.Err()
}

func (api *Client) accessLogsRequest(ctx context.Context, path string, values url.Values) (*LoginResponse, error) {
	response := &LoginResponse{}
	err := api.postMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}
	return response, response.Err()
}

func (api *Client) teamProfileRequest(ctx context.Context, client httpClient, path string, values url.Values) (*TeamProfileResponse, error) {
	response := &TeamProfileResponse{}
	err := api.postMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}
	return response, response.Err()
}

type GetDiscoveryEnterpriseInfoParameters struct {
	Cursor         string
	IncludeDeleted bool
	Limit          int
}

// GetDiscoveryEnterpriseInfoParameters gets the Informations about all teams inside an Enterprise grid
func (api *Client) GetDiscoveryEnterpriseInfo(params *GetDiscoveryEnterpriseInfoParameters) (enterprise *DiscoveryEnterpriseInfo, nextCursor string, err error) {
	return api.GetDiscoveryEnterpriseInfoContext(context.Background(), params)
}

func (api *Client) GetDiscoveryEnterpriseInfoContext(ctx context.Context, params *GetDiscoveryEnterpriseInfoParameters) (enterprise *DiscoveryEnterpriseInfo, nextCursor string, err error) {
	values := url.Values{
		"token": {api.token},
	}
	if params.Cursor != "" {
		values.Add("cursor", params.Cursor)
	}
	if params.Limit != 0 {
		values.Add("limit", strconv.Itoa(params.Limit))
	}
	if params.IncludeDeleted {
		values.Add("include_deleted", strconv.FormatBool(params.IncludeDeleted))
	}
	response := DiscoveryEnterpriseInfoResponse{}
	err = api.postMethod(ctx, "discovery.enterprise.info", values, &response)
	if err != nil {
		return nil, "", err
	}

	return &response.Enterprise, response.ResponseMetadata.Cursor, response.Err()
}

// GetTeamInfo gets the Team Information of the user
func (api *Client) GetTeamInfo() (*TeamInfo, error) {
	return api.GetTeamInfoContext(context.Background())
}

// GetOtherTeamInfoContext gets Team information for any team with a custom context
func (api *Client) GetOtherTeamInfoContext(ctx context.Context, team string) (*TeamInfo, error) {
	if team == "" {
		return api.GetTeamInfoContext(ctx)
	}
	values := url.Values{
		"token": {api.token},
	}
	values.Add("team", team)
	response, err := api.teamRequest(ctx, "team.info", values)
	if err != nil {
		return nil, err
	}
	return &response.Team, nil
}

// GetOtherTeamInfo gets Team information for any team
func (api *Client) GetOtherTeamInfo(team string) (*TeamInfo, error) {
	return api.GetOtherTeamInfoContext(context.Background(), team)
}

// GetTeamInfoContext gets the Team Information of the user with a custom context
func (api *Client) GetTeamInfoContext(ctx context.Context) (*TeamInfo, error) {
	values := url.Values{
		"token": {api.token},
	}

	response, err := api.teamRequest(ctx, "team.info", values)
	if err != nil {
		return nil, err
	}
	return &response.Team, nil
}

// GetTeamProfile gets the Team Profile settings of the user
func (api *Client) GetTeamProfile() (*TeamProfile, error) {
	return api.GetTeamProfileContext(context.Background())
}

// GetTeamProfileContext gets the Team Profile settings of the user with a custom context
func (api *Client) GetTeamProfileContext(ctx context.Context) (*TeamProfile, error) {
	values := url.Values{
		"token": {api.token},
	}

	response, err := api.teamProfileRequest(ctx, api.httpclient, "team.profile.get", values)
	if err != nil {
		return nil, err
	}
	return &response.Profile, nil

}

// GetAccessLogs retrieves a page of logins according to the parameters given
func (api *Client) GetAccessLogs(params AccessLogParameters) ([]Login, *Paging, error) {
	return api.GetAccessLogsContext(context.Background(), params)
}

// GetAccessLogsContext retrieves a page of logins according to the parameters given with a custom context
func (api *Client) GetAccessLogsContext(ctx context.Context, params AccessLogParameters) ([]Login, *Paging, error) {
	values := url.Values{
		"token": {api.token},
	}
	if params.Count != DEFAULT_LOGINS_COUNT {
		values.Add("count", strconv.Itoa(params.Count))
	}
	if params.Page != DEFAULT_LOGINS_PAGE {
		values.Add("page", strconv.Itoa(params.Page))
	}

	response, err := api.accessLogsRequest(ctx, "team.accessLogs", values)
	if err != nil {
		return nil, nil, err
	}
	return response.Logins, &response.Paging, nil
}

// GetBillableInfo ...
func (api *Client) GetBillableInfo(user string) (map[string]BillingActive, error) {
	return api.GetBillableInfoContext(context.Background(), user)
}

// GetBillableInfoContext ...
func (api *Client) GetBillableInfoContext(ctx context.Context, user string) (map[string]BillingActive, error) {
	values := url.Values{
		"token": {api.token},
		"user":  {user},
	}

	return api.billableInfoRequest(ctx, "team.billableInfo", values)
}

// GetBillableInfoForTeam returns the billing_active status of all users on the team.
func (api *Client) GetBillableInfoForTeam() (map[string]BillingActive, error) {
	return api.GetBillableInfoForTeamContext(context.Background())
}

// GetBillableInfoForTeamContext returns the billing_active status of all users on the team with a custom context
func (api *Client) GetBillableInfoForTeamContext(ctx context.Context) (map[string]BillingActive, error) {
	values := url.Values{
		"token": {api.token},
	}

	return api.billableInfoRequest(ctx, "team.billableInfo", values)
}

// GetAccessLogs retrieves a page of logins according to the parameters given
func (api *Client) GetIntegrationLogs(params IntegrationLogsParameters) (*IntegrationLogsResponse, error) {
	return api.GetIntegrationLogsContext(context.Background(), params)
}

// GetAccessLogsContext retrieves a page of logins according to the parameters given with a custom context
func (api *Client) GetIntegrationLogsContext(ctx context.Context, params IntegrationLogsParameters) (*IntegrationLogsResponse, error) {
	values := url.Values{}
	if params.AppID != "" {
		values.Add("app_id", params.AppID)
	}
	if params.ChangeType != "" {
		values.Add("change_type", params.ChangeType)
	}
	if params.ServiceID != "" {
		values.Add("service_id", params.ServiceID)
	}
	if params.TeamID != "" {
		values.Add("team_id", params.TeamID)
	}
	if params.User != "" {
		values.Add("user", params.User)

	}
	if params.Count != DEFAULT_LOGINS_COUNT {
		values.Add("count", strconv.Itoa(params.Count))
	}
	if params.Page != DEFAULT_LOGINS_PAGE {
		values.Add("page", strconv.Itoa(params.Page))
	}

	response := &IntegrationLogsResponse{}

	err := api.getMethod(ctx, "team.integrationLogs", api.token, values, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
