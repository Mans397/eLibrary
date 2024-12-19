# E-library

## Project Description
The Library Management System is designed to manage user data and facilitate CRUD (Create, Read, Update, Delete) operations through a user-friendly web interface. It is tailored for small to medium-sized organizations or individuals who need an efficient and scalable solution for managing library resources and user accounts.

## Team Members
- Mansur Serikov
- Rostislav Prettser
- Nartay Shamil

## Screenshot
<img width="1512" alt="Снимок экрана 2024-12-19 в 20 00 28" src="https://github.com/user-attachments/assets/b3810cc2-cea6-436d-9953-235a0c759ee4" />
<img width="1512" alt="Снимок экрана 2024-12-19 в 20 00 38" src="https://github.com/user-attachments/assets/1a7420e5-ac4c-4a7b-bbf9-6f4b8f046a4d" />


## How to Start the Project

### Prerequisites
1. Install [IntelliJ IDEA](https://www.jetbrains.com/idea/) (Community или Ultimate Edition). Make sure to enable the Go plugin if it's not already installed.
2. Install [Go](https://golang.org/dl/) (version 1.20 or higher).
3. Install [PostgreSQL](https://www.postgresql.org/) and create a database.


### Step-by-Step Instructions
1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/library-management-system.git
   cd library-management-system
   ```
2. **Set up the database:**
   - Create a PostgreSQL database and user.
   - Update the `Database/connection.go` file with your database credentials.
3. **Run migrations:**
   ```bash
   go run migrate.go
   ```
4. **Start the server:**
   ```bash
   go run main.go
   ```
5. **Open the frontend:**
   - Open your browser and go to `http://localhost:8080/`.

### API Endpoints
- `GET /db/readUser`: Fetch a user by email.
- `POST /db/createUser`: Add a new user.
- `PUT /db/updateUser`: Update an existing user.
- `DELETE /db/deleteUser`: Remove a user by email.

## Tools and Resources
- **Programming Language:** Go (Golang)
- **Database:** PostgreSQL
- **Frameworks/Libraries:**
  - `gorm.io/gorm` for ORM
  - `net/http` for HTTP server
  - `encoding/json` for JSON parsing
- **Frontend:** Basic HTML, CSS, and JavaScript
- **Testing Tools:** Postman
- **Version Control:** Git and GitHub

## Notes
- Ensure the database is running before starting the server.
- Use Postman or the web interface to test CRUD operations.

---
Feel free to contact the team for further assistance or issues.

> **Deadline Reminder:** All tasks must be completed by December 20, 2024, at 21:00.

