package tenant

type CreateTenantRequest struct {
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	LegalName  string `json:"legalName"`
	SchoolCode string `json:"schoolCode"`
	Domain     string `json:"domain"`

	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Website     string `json:"website"`
	Description string `json:"description"`

	Address Address     `json:"address"`
	Geo     GeoLocation `json:"geo"`
	Owner   OwnerInfo   `json:"owner"`

	Consent  TenantConsent          `json:"consent"`
	Metadata map[string]interface{} `json:"metadata"`
}

type UpdateTenantRequest struct {
	Name        string `json:"name"`
	LegalName   string `json:"legalName"`
	SchoolCode  string `json:"schoolCode"`
	Domain      string `json:"domain"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Website     string `json:"website"`
	Description string `json:"description"`

	Logo    Logo        `json:"logo"`
	Banner  Banner      `json:"banner"`
	Address Address     `json:"address"`
	Geo     GeoLocation `json:"geo"`

	Features FeatureFlags   `json:"features"`
	Settings TenantSettings `json:"settings"`

	Metadata map[string]interface{} `json:"metadata"`
}
