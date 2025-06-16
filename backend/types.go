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
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	CategoryID   int          `json:"category_id"`
	Description  string       `json:"description"`
	ForBrand     string       `json:"for_brand"`
	IsUniversal  bool         `json:"is_universal"`
	Motorcycles  []Motorcycle `json:"motorcycles"`
	ImporterName string       `json:"importer_name"`
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
