package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"mime/multipart"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func csvreader(file multipart.File, db *sql.DB, rootCategory string) {
	decoded := transform.NewReader(file, charmap.Windows1252.NewDecoder())

	reader := csv.NewReader(decoded)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Fel vid läsning av CSV:", err)
		return
	}

	if len(records) < 1 {
		fmt.Println("CSV-filen är tom")
		return
	}

	err = insertFromCSV(records, db, rootCategory)
	if err != nil {
		log.Println("Error when inserting: ", err)
	}
}

func insertFromCSV(records [][]string, db *sql.DB, rootCategory string) error {
	rootCatID, err := getOrCreateCategoryWithParent(db, rootCategory, nil)
	if err != nil {
		return fmt.Errorf("kunde inte skapa/hämta root kategori: %w", err)
	}

	for i, row := range records {
		if i == 0 {
			continue
		}
		subCategoryName := row[0]
		brandName := row[4]
		modelName := row[5]
		modYears := row[6]
		productCode := "KT" + row[9]
		productName := row[10]
		companyName := row[18]

		// 1. Skapa/hämta underkategori med parent = "fjädrar"
		subCatID, err := getOrCreateCategoryWithParent(db, subCategoryName, &rootCatID)
		if err != nil {
			return fmt.Errorf("kunde inte skapa/hämta underkategori %s: %w", subCategoryName, err)
		}

		isUniversal := brandName == "" || modelName == "" || modYears == ""

		// 6. Skapa produkt
		productID, err := getOrCreateProduct(db, productCode, productName, subCatID, brandName, isUniversal, companyName)
		if err != nil {
			return fmt.Errorf("kunde inte skapa/hämta produkt: %w", err)
		}

		// 7. Koppla endast om det inte är en universal-produkt
		if !isUniversal {
			brandID, err := getOrCreateBrand(db, brandName)
			if err != nil {
				return fmt.Errorf("kunde inte skapa/hämta brand %s: %w", brandName, err)
			}

			modelID, err := getOrCreateModel(db, brandID, modelName)
			if err != nil {
				return fmt.Errorf("kunde inte skapa/hämta model %s: %w", modelName, err)
			}

			startYear, endYear := parseModelYears(modYears)

			fullname := brandName + " " + modelName + " " + strconv.Itoa(startYear) + "-" + strconv.Itoa(endYear)
			motorcycleID, err := getOrCreateMotorcycle(db, brandID, modelID, startYear, endYear, fullname)
			if err != nil {
				return fmt.Errorf("kunde inte skapa/hämta motorcycle: %w", err)
			}

			err = insertProductCompatibility(db, productID, motorcycleID)
			if err != nil {
				return fmt.Errorf("kunde inte skapa produkt_compatibility: %w", err)
			}
		}
	}

	return nil
}

func getOrCreateCategoryWithParent(db *sql.DB, name string, parentID *int) (int, error) {
	var id int

	// Försök hämta kategori först
	if parentID == nil {
		err := db.QueryRow(`SELECT id FROM categories WHERE name = $1 AND parent_id IS NULL`, name).Scan(&id)
		if err == sql.ErrNoRows {
			// Skapa kategori utan parent
			err = db.QueryRow(`INSERT INTO categories(name, parent_id) VALUES($1, NULL) RETURNING id`, name).Scan(&id)
			if err != nil {
				return 0, err
			}
			// Uppdatera path och level för rotkategori
			path := fmt.Sprintf("/%d/", id)
			level := 0
			_, err = db.Exec(`UPDATE categories SET path = $1, level = $2 WHERE id = $3`, path, level, id)
			if err != nil {
				return 0, err
			}
		}
		return id, err
	} else {
		// Hämta förälder path och level
		var parentPath string
		var parentLevel int
		err := db.QueryRow(`SELECT path, level FROM categories WHERE id = $1`, *parentID).Scan(&parentPath, &parentLevel)
		if err != nil {
			return 0, err
		}

		err = db.QueryRow(`SELECT id FROM categories WHERE name = $1 AND parent_id = $2`, name, *parentID).Scan(&id)
		if err == sql.ErrNoRows {
			// Skapa kategori med parent
			err = db.QueryRow(`INSERT INTO categories(name, parent_id) VALUES($1, $2) RETURNING id`, name, *parentID).Scan(&id)
			if err != nil {
				return 0, err
			}

			// Uppdatera path och level baserat på förälder
			path := fmt.Sprintf("%s%d/", parentPath, id)
			level := parentLevel + 1
			_, err = db.Exec(`UPDATE categories SET path = $1, level = $2 WHERE id = $3`, path, level, id)
			if err != nil {
				return 0, err
			}
		}
		return id, err
	}
}

func parseModelYears(modYearStr string) (int, int) {
	var startyear, endyear int = 0, 9999
	var err error

	if strings.Contains(modYearStr, ">") {
		trimmed := strings.Trim(modYearStr, ">")
		startyear, err = strconv.Atoi(trimmed)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	} else if strings.Contains(modYearStr, "-") {
		splitstrs := strings.Split(modYearStr, "-")
		startyear, err = strconv.Atoi(splitstrs[0])
		if err != nil {
			fmt.Println("Error: ", err)
		}

		endyear, err = strconv.Atoi(splitstrs[1])
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
	return startyear, endyear
}

func getOrCreateBrand(db *sql.DB, brandName string) (int, error) {
	var id int
	err := db.QueryRow(`INSERT INTO brands (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id`, brandName).Scan(&id)
	if err == sql.ErrNoRows {
		// raden fanns redan, hämta den manuellt
		err = db.QueryRow(`SELECT id FROM brands WHERE name = $1`, brandName).Scan(&id)
	}
	return id, err
}

func getOrCreateModel(db *sql.DB, brandID int, modelName string) (int, error) {
	var id int
	err := db.QueryRow(`
		INSERT INTO models(brand_id, name)
		VALUES($1, $2)
		ON CONFLICT (brand_id, name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, brandID, modelName).Scan(&id)
	return id, err
}

func getOrCreateMotorcycle(db *sql.DB, brandID int, modelID int, startYear int, endYear int, fullname string) (int, error) {
	var id int

	query := `
    	INSERT INTO motorcycles (brand_id, model_id, startyear, endyear, full_name)
    	VALUES ($1, $2, $3, $4, $5)
    	ON CONFLICT (brand_id, model_id, startyear, endyear)
    	DO UPDATE SET full_name = EXCLUDED.full_name
    	RETURNING id
	`

	err := db.QueryRow(query, brandID, modelID, startYear, endYear, fullname).Scan(&id)
	return id, err
}

func getOrCreateProduct(db *sql.DB, productCode, productName string, subCatID int, brandName string, isUniversal bool, companyName string) (string, error) {
	var id string

	query := `
	INSERT INTO products(id, name, category_id, description, for_brand, is_universal, importer_name)
	VALUES($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (id) DO UPDATE 
	SET name = EXCLUDED.name, is_universal = EXCLUDED.is_universal
	RETURNING id;
	`
	err := db.QueryRow(query, productCode, productName, subCatID, companyName+" - "+productName, brandName, isUniversal, companyName).Scan(&id)
	return id, err
}

func insertProductCompatibility(db *sql.DB, productID string, motorcycleID int) error {
	_, err := db.Exec(`
		INSERT INTO product_compatibility (product_id, motorcycle_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, productID, motorcycleID)
	return err
}
