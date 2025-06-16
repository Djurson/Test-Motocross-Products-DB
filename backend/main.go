package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// connect to the database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// create table if it doesn't exist
	if err := createSchema(db); err != nil {
		log.Fatal("Failed to create schema:", err)
	}

	// create router
	router := mux.NewRouter()
	router.HandleFunc("/brands", getBrandsHandler(db)).Methods("GET")
	router.HandleFunc("/brands/{brand}/models", getModelsByBrandHandler(db)).Methods("GET")
	router.HandleFunc("/brands/{brand}/models/{model}/years", getYearsHandler(db)).Methods("GET")
	router.HandleFunc("/categories", getCategoriesHandler(db)).Methods("GET")

	router.HandleFunc("/products", getFilteredProductsHandler(db)).Methods("GET")

	router.HandleFunc("/upload", uploadFileHandler(db)).Methods("POST")

	// wrap the router with CORS and JSON content type middlewares
	enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))

	// start server
	log.Fatal(http.ListenAndServe(":8000", enhancedRouter))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// check if the request is for CORS preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass down the request to the next middleware (or final header)
		next.ServeHTTP(w, r)
	})
}

func createSchema(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS brands (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) UNIQUE NOT NULL
		)`,

		`CREATE TABLE IF NOT EXISTS models (
			id SERIAL PRIMARY KEY,
			brand_id INTEGER NOT NULL REFERENCES brands(id),
			name VARCHAR(100) NOT NULL,
			UNIQUE (brand_id, name)
		)`,

		`CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			parent_id INTEGER REFERENCES categories(id),
			path VARCHAR(500),
			level INTEGER DEFAULT 0
		)`,

		`CREATE TABLE IF NOT EXISTS products (
    		id VARCHAR(50) PRIMARY KEY,
    		name VARCHAR(200) NOT NULL,
    		category_id INTEGER NOT NULL REFERENCES categories(id),
    		description TEXT DEFAULT '',
    		for_brand VARCHAR(100),
    		is_universal BOOLEAN DEFAULT FALSE,
			importer_name VARCHAR(100)
		)`,

		`CREATE INDEX IF NOT EXISTS idx_products_is_universal ON products(is_universal)`,

		`CREATE TABLE IF NOT EXISTS motorcycles (
    		id SERIAL PRIMARY KEY,
    		brand_id INTEGER NOT NULL REFERENCES brands(id),
    		model_id INTEGER NOT NULL REFERENCES models(id),
    		startyear INTEGER NOT NULL,
    		endyear INTEGER NOT NULL,
			full_name TEXT,
    		UNIQUE (brand_id, model_id, startyear, endyear)
		)`,

		`CREATE UNIQUE INDEX IF NOT EXISTS idx_motorcycles_brand_model_start_end
			ON motorcycles (brand_id, model_id, startyear, endyear)`,

		`CREATE TABLE IF NOT EXISTS product_compatibility (
			product_id VARCHAR(50) NOT NULL REFERENCES products(id),
			motorcycle_id INTEGER NOT NULL REFERENCES motorcycles(id),
			PRIMARY KEY (product_id, motorcycle_id)
		)`,
	}

	for i, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Printf("Schema error on query %d: %v\nSQL: %s", i, err, q)
			return err
		}
	}
	return nil
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set JSON Content-Type
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getBrandsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM brands ORDER BY name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var brands []Brand
		for rows.Next() {
			var b Brand
			if err := rows.Scan(&b.ID, &b.Name); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			brands = append(brands, b)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(brands)
	}
}

func getModelsByBrandHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		brand := vars["brand"]

		query := `
			SELECT mo.id, mo.name 
			FROM models mo 
			JOIN brands b ON mo.brand_id = b.id 
			WHERE b.name = $1 
			ORDER BY mo.name
		`

		rows, err := db.Query(query, brand)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var models []Model
		for rows.Next() {
			var m Model
			if err := rows.Scan(&m.ID, &m.Name); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			models = append(models, m)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models)
	}
}

func getYearsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		brand := vars["brand"]
		model := vars["model"]

		query := `
			SELECT m.startyear, m.endyear
			FROM motorcycles m
			JOIN brands b ON m.brand_id = b.id
			JOIN models mo ON m.model_id = mo.id
			WHERE b.name ILIKE $1 AND mo.name ILIKE $2
			GROUP BY m.startyear, m.endyear
			ORDER BY m.startyear
		`

		rows, err := db.Query(query, brand, model)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var ranges []ModelYearRange
		for rows.Next() {
			var r ModelYearRange
			if err := rows.Scan(&r.StartYear, &r.EndYear); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ranges = append(ranges, r)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ranges)
	}
}

func getCategoriesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT id, name, path, level
			FROM categories
			WHERE parent_id IS NULL
			ORDER BY name
		`)
		if err != nil {
			log.Printf("Database query error: %v", err) // <-- logga felet
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var categories []Category
		for rows.Next() {
			var c Category
			if err := rows.Scan(&c.ID, &c.Name, &c.Path, &c.Level); err != nil {
				log.Printf("Error scanning row: %v", err) // <-- logga felet
				http.Error(w, "Error scanning row", http.StatusInternalServerError)
				return
			}
			categories = append(categories, c)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(categories)
	}
}

func getFilteredProductsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		brand := queryParams.Get("brand")
		model := queryParams.Get("model")
		yearStr := queryParams.Get("year")
		categoryIDStr := queryParams.Get("category_id")

		var startyear, endyear int
		var err error

		if yearStr != "" {
			startyear, endyear, err = extractYears(yearStr)
			if err != nil {
				log.Printf("Error when extracting years: %v", err)
				http.Error(w, "Error when extracting years", http.StatusBadRequest)
				return
			}
		}

		baseQuery := `
			WITH RECURSIVE category_tree AS (
				SELECT id, name, parent_id, '/' || id || '/' AS path
				FROM categories
				WHERE parent_id IS NULL
				UNION ALL
				SELECT c.id, c.name, c.parent_id, ct.path || c.id || '/'
				FROM categories c
				JOIN category_tree ct ON c.parent_id = ct.id
			),
			category_paths AS (
				SELECT c1.id AS category_id, string_agg(c2.name, '/' ORDER BY position) AS full_category_path
				FROM category_tree c1
				JOIN category_tree c2 ON c1.path LIKE '%/' || c2.id || '/%'
				JOIN LATERAL (
					SELECT position
					FROM regexp_split_to_table(c1.path, '/') WITH ORDINALITY AS t(part, position)
					WHERE part = c2.id::text
					LIMIT 1
				) pos ON true
				GROUP BY c1.id
			)

		SELECT DISTINCT 
			p.id, p.name, p.category_id, cp.full_category_path,
			p.description, p.for_brand, p.is_universal, p.importer_name
		FROM products p
		LEFT JOIN product_compatibility pc ON p.id = pc.product_id
		LEFT JOIN motorcycles m ON pc.motorcycle_id = m.id
		LEFT JOIN brands b ON m.brand_id = b.id
		LEFT JOIN models mo ON m.model_id = mo.id
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN category_paths cp ON p.category_id = cp.category_id
		`
		whereClauses := []string{}
		args := []interface{}{}
		argPos := 1

		if brand != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("b.name ILIKE $%d", argPos))
			args = append(args, brand)
			argPos++
		}
		if model != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("mo.name ILIKE $%d", argPos))
			args = append(args, model)
			argPos++
		}
		if yearStr != "" {
			whereClauses = append(whereClauses,
				fmt.Sprintf("m.startyear <= $%d AND m.endyear >= $%d", argPos, argPos+1))
			args = append(args, startyear, endyear)
			argPos += 2
		}

		if categoryIDStr != "" {
			catID, err := strconv.Atoi(categoryIDStr)
			if err != nil {
				http.Error(w, "Invalid category_id", http.StatusBadRequest)
				return
			}
			whereClauses = append(whereClauses,
				fmt.Sprintf("(c.id = $%d OR c.path LIKE '%%/' || $%d || '/%%')", argPos, argPos))
			args = append(args, catID)
			argPos++
		}

		// Om inga specifika filters är satta → visa bara universella produkter
		if len(whereClauses) == 0 {
			whereClauses = append(whereClauses, "p.is_universal = TRUE")
		}

		finalQuery := baseQuery + " WHERE " + strings.Join(whereClauses, " AND ")
		log.Printf("Executing SQL:\n%s\nWith args: %+v\n", finalQuery, args)

		rows, err := db.Query(finalQuery, args...)
		if err != nil {
			log.Printf("Database query error: %v\n", err)
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		products, err := extractProductsFromRows(rows, db)
		if err != nil {
			http.Error(w, "Error scanning product", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func extractYears(yearString string) (int, int, error) {
	yearsSplit := strings.Split(yearString, "-")
	if len(yearsSplit) != 2 {
		return 0, 0, nil
	}
	startyear, err := strconv.Atoi(yearsSplit[0])
	if err != nil {
		return 0, 0, err
	}
	endyear, err := strconv.Atoi(yearsSplit[1])
	if err != nil {
		return 0, 0, err
	}

	return startyear, endyear, nil
}

func extractProductsFromRows(rows *sql.Rows, db *sql.DB) ([]Product, error) {
	var err error

	var products []Product

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.CategoryID, &p.CategoryPath, &p.Description, &p.ForBrand, &p.IsUniversal, &p.ImporterName); err != nil {
			log.Printf("Error scanning product: %v\n", err)
			return nil, err
		}

		p.Motorcycles, err = extractMotorcyclesForProduct(p.ID, db)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func extractMotorcyclesForProduct(prodID string, db *sql.DB) ([]Motorcycle, error) {
	// Hämta motorcyklar för produkten
	mcRows, err := db.Query(`
				SELECT m.id, b.name, mo.name, m.startyear, m.endyear
				FROM product_compatibility pc
				JOIN motorcycles m ON pc.motorcycle_id = m.id
				JOIN brands b ON m.brand_id = b.id
				JOIN models mo ON m.model_id = mo.id
				WHERE pc.product_id = $1
				LIMIT 10
				`, prodID)
	if err != nil {
		log.Printf("Error querying motorcycles: %v\n", err)
		return nil, err
	}

	defer mcRows.Close()
	var motorcycles []Motorcycle
	for mcRows.Next() {
		var m Motorcycle
		if err := mcRows.Scan(&m.ID, &m.Brand, &m.Model, &m.StartYear, &m.EndYear); err != nil {
			log.Printf("Error scanning motorcycle: %v\n", err)
			continue
		}
		motorcycles = append(motorcycles, m)
	}

	return motorcycles, nil
}

func uploadFileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(30 << 20)
		if err != nil {
			http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
			return
		}

		rootCategory := r.FormValue("category")
		if rootCategory == "" {
			http.Error(w, "Category is required", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Could not get file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		if !strings.HasSuffix(strings.ToLower(handler.Filename), ".csv") {
			http.Error(w, "Only .csv files allowed", http.StatusBadRequest)
			return
		}

		csvreader(file, db, rootCategory)
	}
}
