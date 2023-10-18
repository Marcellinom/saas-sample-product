package responses

type User struct {
	ID         string `json:"id" example:"UUID"`
	ActiveRole string `json:"active_role" example:"super admin"`
	Roles      []Role `json:"roles"`
}

type Role struct {
	Name        string   `json:"name" example:"Yoga"`
	Permissions []string `json:"permissions" example:"[]string{bahagia,menangis}"`
	IsDefault   bool     `json:"is_default"`
}
