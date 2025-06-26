# Test-Motocross-Products-DB 🏍️

![Go](https://img.shields.io/badge/Go-1.24-blue?logo=go)
![Next.js](https://img.shields.io/badge/Next.js-15-black?logo=next.js)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13-blue?logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Compose-blue?logo=docker)
![License](https://img.shields.io/badge/License-MIT-green)

A simple test project for motocross products, and matching them to brand, model and year specific model built with:

- Next.js _(Frontend)_

- Go _(Backend API with Gorilla Mux)_

- PostgreSQL _(Database)_

- ShadCN _(Components)_

- Docker for containerization

# 📁 Project Structure

```go
.
├── backend/
│   ├── csvhandler.go               -- Script to import fitment data from CSV
│   ├── go.dockerfile               -- Go service Dockerfile
│   ├── main.go                     -- Main API and router logic
│   ├── test.go                     -- Test for CSV parsing
│   └── types.go                    -- Defined types for DB
│
├── frontend/
│   ├── app/                        -- Main page and logic
│   ├── components/                 -- UI components (cards, tables, dropdown)
│   └── next.dockerfile             -- Next.js service Dockerfile
│
├── docker-compose.yaml
└── README.md                       -- This file

```

# 🧠 Features
* 🧩 Normalized PostgreSQL schema for:
  
  - Brands
 
  - Models
 
  - Motorcycles (year-specific)
 
  - Products and product categories
 
  - Product-to-motorcycle compatibility
    
* 🛠 CSV importer written in Go
  
* 🔍 Query-ready for complex filtering like:
  
  - Brand + Model + Year + Category
 
  - Universal fitment across motorcycles
 
# 🚀 Getting Started

✅ Prerequisites
Make sure you have the following installed:

- [Docker](https://www.docker.com/products/docker-desktop)

- [Node.js](https://nodejs.org/en)

# 📋 Setup Instructions

1. Clone the repository

```bash
git clone https://github.com/Djurson/Test-Motocross-Products-DB.git
```

2. Move into to the repository

```bash
cd Test-Motocross-Products-DB
```

3. Create the database

```bash
docker compose up -d db
```

4. Build the backend

```bash
cd backend
docker compose build goapp
```

5. Compose the backend

```bash
cd ..
docker compose up goapp
```

6. Build or run the frontend

```bash
cd frontend
npm install
npm run dev
```

or

```bash
docker compose up -d nextapp
```

# 🔧 API Endpoints

The backend exposes a REST API under /api/go/users:

| Method | Endpoint                               | Description                                   |
| ------ | -------------------------------------- | --------------------------------------------- |
| GET    | `/brands`                              | Get all brands                                |
| GET    | `/brands/{brand}/models`               | Get all models for a brand                    |
| GET    | `/brands/{brand}/models/{model}/years` | Get all years for a specific model            |
| GET    | `/categories`                          | Get all categories                            |
| GET    | `/products`                            | Get products depending on filters             |
| POST   | `/upload`                              | Upload a csv file of products to the database |

# 💾 Database

The database is automatically started in Docker with the following default values:

- User: `motocross_user`

- Password: `m0tocr0ss_450`

- Database: `motocross_db`

The database is initialized with the following schema:

```sql
CREATE TABLE IF NOT EXISTS brands (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS models (
  id SERIAL PRIMARY KEY,
  brand_id INTEGER NOT NULL REFERENCES brands(id),
  name VARCHAR(100) NOT NULL,
  UNIQUE (brand_id, name)
);

CREATE TABLE IF NOT EXISTS categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  parent_id INTEGER REFERENCES categories(id),
  path VARCHAR(500),
  level INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS products (
  id VARCHAR(50) PRIMARY KEY,
  name VARCHAR(200) NOT NULL,
  category_id INTEGER NOT NULL REFERENCES categories(id),
  description TEXT DEFAULT '',
  for_brand VARCHAR(100),
  is_universal BOOLEAN DEFAULT FALSE,
  importer_name VARCHAR(100)
);

CREATE INDEX IF NOT EXISTS idx_products_is_universal ON products(is_universal);

CREATE TABLE IF NOT EXISTS motorcycles (
  id SERIAL PRIMARY KEY,
  brand_id INTEGER NOT NULL REFERENCES brands(id),
  model_id INTEGER NOT NULL REFERENCES models(id),
  startyear INTEGER NOT NULL,
  endyear INTEGER NOT NULL,
  full_name TEXT,
  UNIQUE (brand_id, model_id, startyear, endyear)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_motorcycles_brand_model_start_end
  ON motorcycles (brand_id, model_id, startyear, endyear);

CREATE TABLE IF NOT EXISTS product_compatibility (
  product_id VARCHAR(50) NOT NULL REFERENCES products(id),
  motorcycle_id INTEGER NOT NULL REFERENCES motorcycles(id),
  PRIMARY KEY (product_id, motorcycle_id)
);
```

# 👤 Author

Created by [@Djurson](https://github.com/Djurson)
