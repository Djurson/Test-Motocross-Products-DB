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
	ID        int    `json:"id"`
	Brand     string `json:"brand"`
	Model     string `json:"model"`
	StartYear int    `json:"start_year"`
	EndYear   int    `json:"end_year"`
}

type Product struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	CategoryID   int          `json:"category_id"`
	CategoryName string       `json:"category_name"`
	CategoryPath string       `json:"category_path"`
	Description  string       `json:"description"`
	ForBrand     string       `json:"for_brand"`
	IsUniversal  bool         `json:"is_universal"`
	ImporterName string       `json:"importer_name"`
	Motorcycles  []Motorcycle `json:"motorcycles"`
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
