package ghlib

/*  types for github API */

type GhLimitResourcesCore struct {
    Limit int `json:"limit"`
    Remaining int `json:"remaining"`
    reset int `json:"reset"`
}

type GhLimitResourcesSearch struct {
    Limit int `json:"limit"`
    Remaining int `json:"remaining"`
    reset int `json:"reset"`
}

type GhLimitResources struct {
    Core GhLimitResourcesCore `json:"core"`
    Search GhLimitResourcesSearch `json:"search"`
}

type GhLimitRate struct {
    Limit int `json:"limit"`
    Remaining int `json:"remaining"`
    reset int `json:"reset"`
}

type GhLimit struct {
    Resources GhLimitResources `json:"resources"`
    Rate GhLimitRate `json:"rate"`
}


