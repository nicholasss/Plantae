# Plantae Server

! This project is not currently under development !

This is the server portion of the project. There is currently no front-end.

A REST API to keep track of different types of plant species, their common names, and care.

Server admins (super admins) can promote/demote registered users to admin status.
Admins have the power to create, update, delete, and list:

- plant species
- common names of plants species
- plant types
- lighting and water needs
- and the setting/unsetting of plants as specific plant types, lighting needs, and water needs

Users are limited to specific endpoints that have to do with user registration, login, and managing their plants.

There is an accompanying OpenAPI spec within the `openapi` directory.

## Goals of the project

- Create a plant management REST API
- Use the `net/http` package instead of a framework
- Utilize a PostgreSQL database as the backend
- Complete an OpenAPI spec from scratch

## Shortcomings of the project

The web server was written from the "top down". This meant that I was worried about what information would need to be stored by the server instead of expanding functionality.

The web server was developed stand alone, without a front end. This meant that I was not able to perform usability testing. While I did write integration tests for the `Hurl` tool to perform, they are not the same as the insights and feedback from using the web server directly.

Due to the large numbers of different categories of information, it quickly became time intensive to implement the basic features. Without being able to use or benefit from the features implemented, it was a slog to work through.

## What was learned

1. Focus on a iterative feature-focused approach instead of a high level design focused approach.
2. Build a front end to integrate with. This will allow for better design of the back end and crucial feedback.
3. Integration testing with `Hurl` proved very useful to both learn more about HTTP servers and to correctly implement features on specific endpoints.
4. Type out (or write) terminology to be used across the API that should stay consistent.
5. Consider how the data will be used by the client during design phase.

## Style guide

JSON is in camelCase.
URL paths are in kebab-case.
Scripting variables (Bash, Hurl) and PostgreSQL table & column names are in snake_case.

## to do

### General efforts

- [x] Implement basic user account management.
- [x] Implement basic admin account management.
- [x] Implement plant species management.
- [x] Implement plant name management.
- [x] Implement plant type management.
- [x] Implement plant light need management.
- [x] Implement plant water need management.
- [x] Implement language preference with user registration.
- [ ] Implement basic end-user's plant tracking.
- [ ] Implement script to place 5 example plant's data into database.
- [ ] Create a backup scheme for the universal plant species data.

### Cleanup efforts

- [x] Fix respondWithError utility function to not reveal the error value to client
- [x] Replace all fmt.Errorf with errors.New where formatting or wrapping is not needed
- [x] Add 'application/json; charset=utf-8' Content-Type assert in all of the tests and fix where needed
- [x] Rename 'authorizeNormalAdmin' utility function or rewrite to actually authorize normal admins.
- [x] log -> log/slog package for logging level, and logging to file
- [x] Ensure all POST responses are 201 (return body of resource created makes sense to have a 201)
- [x] Ensure all GET requests do not give a 204 and send an empty table instead.
- [x] Fix delete queries to update 'updated_at' and 'updated_by'
- [ ] Return either [] or {} when database return is null/empty
