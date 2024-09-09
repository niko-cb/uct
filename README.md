# uct (Upsider Coding Test)

# The API
To run the api, you just need to use docker to start the containers (api and mysql).

```bash
docker-compose up
```

## Test data creation
To create the test data, you can use the following command:

```bash
go run backend/testdata/create.go
```

- This will create 20 companies, 1000 users and 10000 invoices all with due dates within a year from 'now'.
- The invoices and users are randomly linked to companies in the database.

## JWT Authentication

The api uses Echo's JWT library for authentication. To get a token, you can use the following command:

```bash
go run backend/testauth/jwt.go
```

- This will create a token for the user "User 1", but feel free to play with it.
- Token settings are default for this test.
- The secret key is stored in an environment variable for the test, but in a real scenario, it should be stored in a more secure way like a secret manager.

## Endpoints

- The test asks for api/invoices (GET and POST), but I've added versioning to it (v1: so it's /api/v1/invoices) to keep in mind that the api can grow and change and we may need to keep older versions running.

## ORM

- Using sqlboiler to generate the ORM models and queries.
- The models are generated in the backend/models folder.

## Table design

- Using tbls to generate the table design in Markdown format.
- The table design is generated in the backend/docs folder.

## Wire

- Using wire to generate the dependency injection code.

# Some notes

1. The test asks for a simple API, so I kept it simple.
2. Code readability and maintainability --> I used Clean Architecture principles to separate concerns and make the code more maintainable.
3. Testability --> I added tests for the main parts of the code. I kept the tests simple and didn't add too many as I wanted to finish in time. Normally I would test everything.
4. Security --> I used JWT for authentication and kept the secret key in an environment variable for the test.
5. Some Documentation --> I added some comments to the code to make it a little more understandable and used tbls to generate the table design which would make it easier for someone to understand the database design.
6. Code generation --> I used sqlboiler to generate the ORM models and queries and wire to generate the dependency injection code.
7. Docker --> I used docker to make it easier to run the api and the database.
