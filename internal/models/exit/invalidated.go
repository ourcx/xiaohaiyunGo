package exit

type Invalidated struct {
	Jwt   string `json:"jwt"`
	Email string `json:"mail"`
}
