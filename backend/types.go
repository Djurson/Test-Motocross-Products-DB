package api

type Brand struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Model struct {
	ID      int    `json:"id"`
	BrandID int    `json:"brand_id"`
	Name    string `json:"name"`
}

type EngineSize struct {
	ID      int `json:"id"`
	BrandID int `json:"brand_id"`
	SizeCC  int `json:"size_cc"`
}

type Motorcycle struct {
	ID           int     `json:"id"`
	BrandID      int     `json:"brand_id"`
	ModelID      int     `json:"model_id"`
	EngineSizeID int     `json:"engine_size_id"`
	StartYear    int     `json:"start_year"`
	EndYear      int     `json:"end_year"`
	FullName     *string `json:"full_name,omitempty"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	CategoryID  int     `json:"category_id"`
	Description *string `json:"description,omitempty"`
	Brand       *string `json:"brand,omitempty"`
	IsUniversal bool    `json:"is_universal"`
}