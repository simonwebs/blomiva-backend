package tenant

type VerifySchoolRequest struct {
	TenantID string `json:"tenantId"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type TenantResponse struct {
	Tenant Tenant `json:"tenant"`
}

type TenantListResponse struct {
	Data  []Tenant `json:"data"`
	Total int64    `json:"total"`
	Page  int64    `json:"page"`
	Limit int64    `json:"limit"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
