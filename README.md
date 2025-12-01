✅ Here is the README without any code fences, no backticks, no escaping
→ Just copy-paste EXACTLY this into README.md
→ It will render perfectly on GitHub, VS Code, Goland, IntelliJ, GitLab, etc.
🐾 Pet Adoption API
A backend REST API built with Go, Gin, PostgreSQL, JWT Auth, and Docker for managing pets, shelters, and adoption requests.
✨ Overview
This API allows:
Users to register/login, browse pets, and apply for adoption
Shelter owners to manage adoption requests
Admins to manage shelters and pets
A background worker to asynchronously process adoption events
This project follows a clean, production-style structure with migrations, Docker support, and role-based access.
🏗️ Tech Stack
Layer	Technology
Language	Go (Golang)
Framework	Gin
Database	PostgreSQL
Auth	JWT
ORM	GORM
Migrations	golang-migrate
Containerization	Docker + docker-compose
Concurrency	Goroutines + Channels (Background Worker)
📂 Project Structure
pet-adoption-api/
├── cmd/
│ └── main.go
├── internal/
│ ├── config/
│ ├── database/
│ ├── handlers/
│ ├── middleware/
│ ├── models/
│ └── worker/
├── migrations/
├── Dockerfile
├── docker-compose.yml
└── README.md
🔑 Authentication
Register
POST /auth/register
Body:
{
"name": "Aida",
"email": "aida@example.com",
"password": "1234"
}
Login
POST /auth/login
Returns:
{
"token": "<jwt>"
}
Use token:
Authorization: Bearer <token>
🐶 Pets API
Method	Endpoint	Access	Description
GET	/pets	Public	List all pets
POST	/pets	Admin	Create pet
DELETE	/pets/:id	Admin	Delete pet
🏡 Shelters API
Method	Endpoint	Access	Description
GET	/shelters	Public	List shelters
POST	/shelters	Admin	Create shelter
❤️ Adoption API
Method	Endpoint	Access	Description
POST	/adoptions/:petID/apply	User	Apply for adoption
GET	/adoptions/my	User	View my adoption requests
GET	/adoptions/shelter	Shelter Owner/Admin	Requests for their shelter
PATCH	/adoptions/:id/approve	Shelter Owner/Admin	Approve request
PATCH	/adoptions/:id/reject	Shelter Owner/Admin	Reject request
🧵 Background Worker
A goroutine worker processes adoption events asynchronously.
Example log:
[WORKER] Processing adoption event → requestID=5 status=pending
Supports graceful shutdown with context cancellation.
🐳 Running with Docker
Build
docker-compose build
Start
docker-compose up
API will be available at:
http://localhost:8080
🛠 Local Development
Run migrations:
migrate -path migrations -database "$DB_URL" up
Run the server:
go run cmd/main.go
✔️ Requirements Checklist
Requirement	Status
JWT Authentication	✅
Multi-table Database	✅
CRUD Operations	✅
Migrations	✅
Role-based Access	✅
Concurrency / Worker	✅
Context + Graceful Shutdown	✅
Docker + docker-compose	✅
Tests	⏳
Documentation	✅