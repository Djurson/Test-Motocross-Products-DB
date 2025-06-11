package api

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"regexp"
	"strconv"
)

func csvreader(file multipart.File, db *sql.DB, rootCategory string) {
    reader := csv.NewReader(file)
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

	insertFromCSV(records, db, rootCategory)
}

func insertFromCSV(records [][]string, db *sql.DB, rootCategory string) error {
	rootCatID, err := getOrCreateCategoryWithParent(db, rootCategory, nil)
    if err != nil {
        return fmt.Errorf("kunde inte skapa/hämta root kategori: %w", err)
    }

    for _, row := range records {
        subCategoryName := row[0] // "Off Road"
        brandName := row[4]    // t.ex. "KTM"
        modelName := row[5]    // t.ex. "250SX"
        modYears := row[6]     // t.ex. "2008-2008"
        productCode := "KT" + row[9]
        productName := row[10] // t.ex. "Front Fork Spring 3.6N"

        // 1. Skapa/hämta underkategori med parent = "fjädrar"
        subCatID, err := getOrCreateCategoryWithParent(db, subCategoryName, &rootCatID)
        if err != nil {
            return fmt.Errorf("kunde inte skapa/hämta underkategori %s: %w", subCategoryName, err)
        }

        // 2. Hämta eller skapa brand
        brandID, err := getOrCreateBrand(db, brandName)
        if err != nil {
            return fmt.Errorf("kunde inte skapa/hämta brand %s: %w", brandName, err)
        }

        // 3. Hämta eller skapa model
        modelID, err := getOrCreateModel(db, brandID, modelName)
        if err != nil {
            return fmt.Errorf("kunde inte skapa/hämta model %s: %w", modelName, err)
        }

        // 4. Parsar motorstorlek från modellnamn (exempel, här hårdkodat)
        engineSizeCC := parseEngineSizeFromModel(modelName) // Du skriver denna funktion själv
        engineSizeID, err := getOrCreateEngineSize(db, brandID, engineSizeCC)
        if err != nil {
            return fmt.Errorf("kunde inte skapa/hämta engine size: %w", err)
        }

        // 5. Parsar start och end year från modYears
        var startYear, endYear int
        n, err := fmt.Sscanf(modYears, "%d-%d", &startYear, &endYear)
        if err != nil || n != 2 {
            // fallback: sätt båda till samma år om parsning misslyckas
            startYear = 0
            endYear = 0
        }

        // 6. Hämta eller skapa motorcycle
        motorcycleID, err := getOrCreateMotorcycle(db, brandID, modelID, engineSizeID, startYear, endYear)
        if err != nil {
            return fmt.Errorf("kunde inte skapa/hämta motorcycle: %w", err)
        }

        // 8. Skapa produkt kopplat till underkategori
        productID, err := getOrCreateProduct(db, productCode, productName, subCatID, brandName)
        if err != nil {
            return fmt.Errorf("kunde inte skapa/hämta produkt: %w", err)
        }

        // 9. Koppla produkt till motorcykel
        err = insertProductCompatibility(db, productID, motorcycleID)
        if err != nil {
            return fmt.Errorf("kunde inte skapa produkt_compatibility: %w", err)
        }
    }

    return nil
}

func getOrCreateCategoryWithParent(db *sql.DB, name string, parentID *int) (int, error) {
    var id int
    if parentID == nil {
        err := db.QueryRow(`SELECT id FROM categories WHERE name = $1 AND parent_id IS NULL`, name).Scan(&id)
        if err == sql.ErrNoRows {
            err = db.QueryRow(`INSERT INTO categories(name, parent_id) VALUES($1, NULL) RETURNING id`, name).Scan(&id)
        }
        return id, err
    } else {
        err := db.QueryRow(`SELECT id FROM categories WHERE name = $1 AND parent_id = $2`, name, *parentID).Scan(&id)
        if err == sql.ErrNoRows {
            err = db.QueryRow(`INSERT INTO categories(name, parent_id) VALUES($1, $2) RETURNING id`, name, *parentID).Scan(&id)
        }
        return id, err
    }
}

func parseEngineSizeFromModel(model string) int {
    re := regexp.MustCompile(`\d+`)
    match := re.FindString(model)
    if match == "" {
        return 0
    }
    size, err := strconv.Atoi(match)
    if err != nil {
        return 0
    }
    return size
}

func getOrCreateBrand(db *sql.DB, brandName string) (int, error) {
	var id int
	err := db.QueryRow(`
		INSERT INTO brands(name)
		VALUES($1)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, brandName).Scan(&id)
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

func getOrCreateEngineSize(db *sql.DB, brandID int, engineSizeCC int) (int, error) {
    var id int
    err := db.QueryRow(`SELECT id FROM engine_sizes WHERE brand_id = $1 AND size_cc = $2`, brandID, engineSizeCC).Scan(&id)
    if err == sql.ErrNoRows {
        err = db.QueryRow(`INSERT INTO engine_sizes(brand_id, size_cc) VALUES($1, $2) RETURNING id`, brandID, engineSizeCC).Scan(&id)
    }
    return id, err
}

func getOrCreateMotorcycle(db *sql.DB, brandID int, modelID int, engineSizeID int, startYear int, endYear int) (int, error) {
    var id int

    query := `
        INSERT INTO motorcycles (brand_id, model_id, engine_size_id, startyear, endyear)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (brand_id, model_id, engine_size_id, startyear, endyear)
        DO UPDATE SET startyear = EXCLUDED.startyear
        RETURNING id
    `

    err := db.QueryRow(query, brandID, modelID, engineSizeID, startYear, endYear).Scan(&id)
    return id, err
}

func getOrCreateProduct(db *sql.DB, productCode string, productName string, subCatID int, brandName string) (string, error) {
	var id string

	query := `
		INSERT INTO products(id, name, category_id, description, brand)
		VALUES($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name
		RETURNING id;
	`
	err := db.QueryRow(query, productCode, productName, subCatID, "", brandName).Scan(&id)
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