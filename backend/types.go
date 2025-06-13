package main

import "database/sql"

type Brand struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Model struct {
	ID      int    `json:"id"`
	BrandID int    `json:"brand_id"`
	Name    string `json:"name"`
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
	Name        string  `json:"name" validate:"required"`
	CategoryID  int     `json:"category_id" validate:"required"`
	Description *string `json:"description,omitempty"`
	Brand       *string `json:"brand,omitempty"`
	IsUniversal bool    `json:"is_universal"`
}

type Category struct {
	ID    int            `json:"id"`
	Name  string         `json:"name"`
	Path  sql.NullString `json:"path"`
	Level sql.NullInt32  `json:"level"`
}

type ModelYearRange struct {
	StartYear int `json:"startyear"`
	EndYear   int `json:"endyear"`
}
