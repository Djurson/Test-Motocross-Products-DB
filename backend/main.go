package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

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
	router.HandleFunc("/bike/{bike_brand}/{bike_model}/{bike_year}", getProductsFromBike(db)).Methods("GET")
	router.HandleFunc("/bike/{name}", getProductsFromBrand(db)).Methods("GET")
	router.HandleFunc("/bike/{name}/{category}", getProductCategoryFromBrand(db)).Methods("GET")
	router.HandleFunc("/upload", uploadFileHandler(db)).Methods("POST")

	// wrap the router with CORS and JSON content type middlewares
	enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))

	// start server
	log.Fatal(http.ListenAndServe(":8000", enhancedRouter))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// check if the request is for CORS preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass down the request to the next middleware (or final header)
		next.ServeHTTP(w,r)
	})
}

func createSchema(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS brands (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) UNIQUE NOT NULL
		)`,

		`CREATE TABLE IF NOT EXISTS engine_sizes (
			id SERIAL PRIMARY KEY,
			brand_id INTEGER NOT NULL REFERENCES brands(id),
			size_cc INTEGER NOT NULL
		)`,

		`CREATE UNIQUE INDEX IF NOT EXISTS idx_engine_sizes_brand_id_size_cc
			ON engine_sizes (brand_id, size_cc)`,

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
    		description TEXT DEFAULT,
    		brand VARCHAR(100),
    		is_universal BOOLEAN DEFAULT FALSE
		);`,

		`CREATE TABLE IF NOT EXISTS motorcycles (
			id SERIAL PRIMARY KEY,
			brand_id INTEGER NOT NULL REFERENCES brands(id),
			model_id INTEGER NOT NULL REFERENCES models(id),
			engine_size_id INTEGER NOT NULL REFERENCES engine_sizes(id),
			startyear INTEGER NOT NULL,
			endyear INTEGER NOT NULL,
			full_name VARCHAR(200)
		)`,

		`COMMENT ON COLUMN motorcycles.full_name IS 'Generated always as: brand name + model name + engine size + years'`,

		`CREATE UNIQUE INDEX IF NOT EXISTS idx_motorcycles_brand_model_engine_start_end
			ON motorcycles (brand_id, model_id, engine_size_id, startyear, endyear)`,

		`CREATE TABLE IF NOT EXISTS product_compatibility (
			product_id VARCHAR(50) NOT NULL REFERENCES products(id),
			motorcycle_id INTEGER NOT NULL REFERENCES motorcycles(id),
			PRIMARY KEY (product_id, motorcycle_id)
		)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// Set JSON Content-Type
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w,r)
	})
}

// get product for specific bike
func getProductsFromBike(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bike_brand := vars["bike_brand"]
		model_name := vars["bike_model"]
		bike_year := vars["bike_year"]

		query := `
		SELECT p.id, p.name, p.category_id, p.description, p.brand, p.is_universal
		FROM products p
		JOIN product_compatibility pc ON p.id = pc.product_id
		JOIN motorcycles m ON pc.motorcycle_id = m.id
		JOIN brands b ON m.brand_id = b.id
		JOIN models mo ON m.model_id = mo.id
		JOIN engine_sizes es ON m.engine_size_id = es.id
		WHERE b.name = $1 AND mo.name = $2 AND m.startyear <= $3 AND m.endyear >= $3`

		yearInt, err := strconv.Atoi(bike_year)
		if err != nil {
			http.Error(w, "Invalid year format", http.StatusBadRequest)
			return
		}
		
		rows, err := db.Query(query, bike_brand, model_name, yearInt)
		if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []Product
		for rows.Next() {
			var p Product
			err := rows.Scan(&p.ID, &p.Name, &p.CategoryID, &p.Description, &p.Brand, &p.IsUniversal)
			if err != nil {
				http.Error(w, "Error scanning row", http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}

		if err = rows.Err(); err != nil {
			http.Error(w, "Error reading rows", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func getProductsFromBrand(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bike_brand := vars["bike_brand"]

		query := `
		SELECT DISTINCT p.id, p.name, p.category_id, p.description, p.brand, p.is_universal
		FROM products p
		JOIN product_compatibility pc ON p.id = pc.product_id
		JOIN motorcycles m ON pc.motorcycle_id = m.id
		JOIN brands b ON m.brand_id = b.id
		WHERE b.name = $1`

		rows, err := db.Query(query, bike_brand)
		if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []Product
		for rows.Next() {
			var p Product
			err := rows.Scan(&p.ID, &p.Name, &p.CategoryID, &p.Description, &p.Brand, &p.IsUniversal)
			if err != nil {
				http.Error(w, "Error scanning row", http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}

		if err = rows.Err(); err != nil {
			http.Error(w, "Error reading rows", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func getProductCategoryFromBrand(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bike_brand := vars["bike_brand"]
		category := vars["category"]
		
		query := `
		SELECT DISTINCT p.id, p.name, p.category_id, p.description, p.brand, p.is_universal
		FROM products p
		JOIN categories c ON p.category_id = c.id
		JOIN product_compatibility pc ON p.id = pc.product_id
		JOIN motorcycles m ON pc.motorcycle_id = m.id
		JOIN brands b ON m.brand_id = b.id
		WHERE c.name = $1 AND b.name = $2`

		rows, err := db.Query(query, category, bike_brand)
		if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []Product
		for rows.Next() {
			var p Product
			err := rows.Scan(&p.ID, &p.Name, &p.CategoryID, &p.Description, &p.Brand, &p.IsUniversal)
			if err != nil {
				http.Error(w, "Error scanning row", http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}

		if err = rows.Err(); err != nil {
			http.Error(w, "Error reading rows", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func uploadFileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20) // Max 10MB
		if err != nil {
			http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		rootCategory := vars["category"]

		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Could not get file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		csvreader(file, db, rootCategory);
	}
}