package types

type Org struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	RegionId   int    `json:"region_id"`
	RegionName string `json:"region_name"`
	RrgJoinDay int    `json:"org_join_day"`
	SiteLogo   string `json:"site_logo"`
}

type User struct {
	Name    string `json:"real_name"`
	Avatar  string `json:"avatar_url"`
	Token   string `json:"token"`
	OrgList []Org  `json:"org_list"`
}
