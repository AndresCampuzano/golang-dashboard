<p align="center">
  <img src="https://cdn-icons-png.flaticon.com/512/6295/6295417.png" width="100"  alt="cloud-img"/>
</p>
<h1 align="center">GOLANG-DASHBOARD</h1>
<p align="center">
	<img src="https://img.shields.io/github/license/AndresCampuzano/golang-dashboard?style=flat&color=0080ff" alt="license">
	<img src="https://img.shields.io/github/last-commit/AndresCampuzano/golang-dashboard?style=flat&logo=git&logoColor=white&color=0080ff" alt="last-commit">
	<img src="https://img.shields.io/github/languages/top/AndresCampuzano/golang-dashboard?style=flat&color=0080ff" alt="repo-top-language">
	<img src="https://img.shields.io/github/languages/count/AndresCampuzano/golang-dashboard?style=flat&color=0080ff" alt="repo-language-count">
<p>

<p align="center">
	<img src="https://img.shields.io/badge/Go-00ADD8.svg?style=flat&logo=Go&logoColor=white" alt="Go">
</p>
<hr>

## 🔗 Quick Links

<!-- TOC -->
  * [🔗 Quick Links](#-quick-links)
  * [📍 Overview](#-overview)
  * [📦 Features](#-features)
    * [Additional Information](#additional-information)
      * [Database Storage](#database-storage)
      * [AWS S3 Integration](#aws-s3-integration)
      * [Helper Functions](#helper-functions)
  * [📂 Repository Structure](#-repository-structure)
  * [🛣️ Endpoints](#-endpoints)
  * [🧩 Dependencies](#-dependencies)
  * [🚀 Getting Started](#-getting-started)
    * [⚙️ Installation](#-installation)
    * [🤖 Running golang-dashboard](#-running-golang-dashboard)
    * [⚙️ Configuration](#-configuration)
      * [PostgreSQL](#postgresql)
      * [AWS](#aws)
  * [🧪 Tests](#-tests)
<!-- TOC -->

---

## 📍 Overview

This is a JSON API server written in Go that provides endpoints for managing users, customers, products, sales, and expenses.

---

## 📦 Features

- User authentication and authorization
- CRUD operations for users, customers, products, sales, and expenses
- Integration with AWS S3 for file storage

### Additional Information

#### Database Storage
The server uses PostgreSQL as its database backend. It provides a `NewPostgresStore` function to create a new instance of `PostgresStore`, which establishes a connection to the PostgreSQL database and initializes necessary extensions.

#### AWS S3 Integration
The server integrates with AWS S3 for file storage. The BucketBasics struct encapsulates Amazon S3 actions such as uploading and deleting files. It provides methods for uploading base64-encoded images to an S3 bucket and deleting files from the bucket.

#### Helper Functions
The server includes helper functions for handling JSON responses, HTTP request routing, and working with PostgreSQL array types.

---

## 📂 Repository Structure

```md
└── golang-dashboard/
    ├── Makefile
    ├── api.go
    ├── auth.go
    ├── customer.service.go
    ├── customer.storage.go
    ├── customer.types.go
    ├── database.go
    ├── expense.service.go
    ├── expense.storage.go
    ├── expense.types.go
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── product.service.go
    ├── product.storage.go
    ├── product.types.go
    ├── s3.go
    ├── sale.service.go
    ├── sale.storage.go
    ├── sale.types.go
    ├── storage.go
    ├── user.service.go
    ├── user.storage.go
    ├── user.types.go
    └── utils.go
```

---

## 🛣️ Endpoints

The server exposes the following endpoints:

- `POST /login`: User login
- `POST /signup`: User sign up
- `GET /users`: Get all users
- `GET /users/{id}`: Get user by ID
- `GET /customers`: Get all customers
- `GET /customers/{id}`: Get customer by ID
- `POST /customers`: Create a new customer
- `PUT /customers/{id}`: Update customer by ID
- `DELETE /customers/{id}`: Delete customer by ID
- `GET /products`: Get all products
- `GET /products/{id}`: Get product by ID
- `POST /products`: Create a new product
- `PUT /products/{id}`: Update product by ID
- `DELETE /products/{id}`: Delete product by ID
- `GET /sales`: Get all sales
- `GET /sales/{id}`: Get sale by ID
- `POST /sales`: Create a new sale
- `GET /expenses`: Get all expenses
- `GET /expenses/{id}`: Get expense by ID
- `POST /expenses`: Create a new expense
- `PUT /expenses/{id}`: Update expense by ID
- `DELETE /expenses/{id}`: Delete expense by ID

---

## 🧩 Dependencies

This project depends on the following external packages:

- [github.com/aws/aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2) - AWS SDK for Go
- [github.com/golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) - JWT implementation for Go
- [github.com/google/uuid](https://github.com/google/uuid) - UUID generation library for Go
- [github.com/gorilla/mux](https://github.com/gorilla/mux) - Powerful HTTP router and URL matcher for building Go web servers
- [github.com/joho/godotenv](https://github.com/joho/godotenv) - Go port of Ruby's dotenv library for loading environment variables from .env files
- [github.com/lib/pq](https://github.com/lib/pq) - PostgreSQL driver for Go's database/sql package
- [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) - Supplementary cryptography libraries for Go

---

## 🚀 Getting Started

***Requirements***

You may need to configure environment variables or configuration files for database connection, AWS credentials, and other settings.

Ensure you have the following dependencies installed on your system:

* **Go**: `version go1.22.2`

### ⚙️ Installation

1. Clone the golang-dashboard repository:

```sh
git clone https://github.com/AndresCampuzano/golang-dashboard
```

2. Change to the project directory:

```sh
cd golang-dashboard
```

3. Install the dependencies:

```sh
go mod download
```

### 🤖 Running golang-dashboard

Use the following command to run golang-dashboard:

```sh
make run
```


### ⚙️ Configuration

This project requires the following environment variables to be set for proper operation:

#### PostgreSQL

- `POSTGRES_USER`: PostgreSQL username
- `POSTGRES_DB_NAME`: PostgreSQL database name
- `POSTGRES_PASSWORD`: PostgreSQL password

#### AWS

- `AWS_ACCESS_KEY_ID`: AWS access key ID for accessing AWS services
- `AWS_SECRET_ACCESS_KEY`: AWS secret access key for accessing AWS services
- `AWS_REGION`: AWS region where resources are located
- `AWS_S3_BUCKET_NAME`: Name of the AWS S3 bucket for file storage
- `AWS_S3_BUCKET_URL`: URL of the AWS S3 bucket for accessing stored files

## 🧪 Tests

Embrace the chaos, this is a work in progress!