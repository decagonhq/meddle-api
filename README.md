
## Installation Guidlines

### Install meddle using docker
To get this app working using docker, make sure
you have docker installed on your pc and docker-compose.yml, use the command below to get
things working

```bash
  git clone url
  cd my-project
  docker compose up --build
```
### Install meddle locally with pg-sql

Step 1: Install and configure pgAdmin locally.

Step 2: Create a new database.

Step 3: Make sure the credentials in the .env matches your DB credentials and run the following command.

```bash
  go run main.go
```