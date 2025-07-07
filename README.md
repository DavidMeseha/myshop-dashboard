# Shop Dashboard 🛒

A vendor/admin dashboard backend for the TechShop e-commerce platform, built with Go, MongoDB, and Chi.

## 📦 Features

- Vendor authentication (JWT, middleware)
- Product management (CRUD, soft delete, republish)
- Category, tag, and query search
- Attribute and tag processing for product update and creation
- Pagination and filtering for products
- Vendor registration
- Unique seName/SKU generation
- MongoDB integration
- RESTful API design
- CORS support

## 🌐 Related Projects

- **Frontend (Next.js):** [Live Demo](https://techshop-commerce.vercel.app/) | [Code](https://github.com/DavidMeseha/TechShop-Ecommerce-Frontend)
- **Client Backend (Node.js/Express):** [Code](https://github.com/DavidMeseha/TechShop-Ecommerce-backend)

## 🛠️ Technology Stack

- **Language:** Go 1.21
- **Web Framework:** [Chi](https://github.com/go-chi/chi)
- **Database:** MongoDB
- **CORS:** [rs/cors](https://github.com/rs/cors)
- **Environment:** [godotenv](https://github.com/joho/godotenv)
- **Testing:** [Testify](https://github.com/stretchr/testify)

## 🚀 Getting Started

### Prerequisites

- Go >= 1.21
- MongoDB >= 5.0

### Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/DavidMeseha/myshop-dashboard.git
   cd shop-dashboard
   ```

2. Copy `.env.example` to `.env` and set your environment variables:

   ```
   PORT=8080
   MONGODB_URI=mongodb://localhost:27017
   MONGODB=shop_dashboard
   ORIGIN=http://localhost:3001
   CLIENT_SERVER=http://localhost:3000
   ```

3. Install dependencies:

   ```sh
   go mod tidy
   ```

4. Run the server (with [air](https://github.com/cosmtrek/air) for live reload):
   ```sh
   dev.bat | ./dev.bat
   ```
   Or build and run manually:
   ```sh
   go build -o main.exe ./cmd/api
   ./main.exe
   ```

## 📖 API Endpoints

- `GET /health` — Health check
- `POST /api/v1/admin/create/product` — Create product
- `POST /api/v1/admin/edit/product/{id}` — Edit product
- `DELETE /api/v1/admin/delete/product/{id}` — Soft delete product
- `POST /api/v1/admin/republish/product/{id}` — Republish product
- `GET /api/v1/admin/products` — List products (pagination, filtering)
- `GET /api/v1/admin/product/{id}` — Get product details
- `GET /api/v1/admin/find/vendors` — Search vendors
- `GET /api/v1/admin/find/categories` — Search categories
- `GET /api/v1/admin/find/tags` — Search tags
- `POST /api/v1/create/vendorSeName` — Generate unique vendor seName
- `POST /api/v1/register/vendor` — Register new vendor

## 🔒 Security

- JWT authentication (via client backend)
- CORS protection
- Input validation

## 📚 Notes

- This dashboard is intended for vendor/admin use.
- Requires the client backend to be running for full functionality as all auth gose through it.