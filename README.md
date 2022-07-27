
## Installation Guidelines

### Provision app using docker
To run this app using docker, you first of all need to make sure you have docker installed on your system. If you have done that already then go ahead and run the following commands to start the application and all needed services.

```bash
  git clone git@github.com:decagonhq/meddle-api.git
  cd meddle-api
  docker compose up --build
```
### Provision app locally with pg-sql
To run this locally, follow the steps below

Step 1: Install and configure pgAdmin locally.

Step 2: Create a new database.

Step 3: Ask someone in the team for the .env credentials.

Step 3: Run the command below.
```bash
  go run main.go
```
### Api documentation link

```http://localhost:8080/swagger
```
