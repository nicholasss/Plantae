# Plantae Server

A REST API server to track care and feeding of your house plants.

This is the server portion of the project. There is currently no front-end.

! This is currently a work in progress and may have bugs or malfunction.

## Style guide

JSON is in camelCase.
URL paths are in kebab-case.
Scripting variables (Bash, Hurl) and PostgreSQL table & column names are in snake_case.

## Goals

- Track multiple house plants.
- Surface information about its toxicity (to humans or pets).
- Organize plants into biomes or into rooms.
- Plan out your watering schedule.

## to do

### General efforts

- [x] Implement basic user account management.
- [x] Implement basic admin account management.
- [x] Implement plant species management.
- [x] Implement plant name management.
- [x] Implement plant type management.
- [x] Implement plant light need management.
- [x] Implement plant water need management.
- [ ] Implement basic end-user's plant tracking.
- [ ] Implement script to place 5 example plant's data into database.
- [ ] Create a backup scheme for the universal plant species data.

### Cleanup efforts

- [x] Fix respondWithError utility function to not reveal the error value to client
- [x] Replace all fmt.Errorf with errors.New where formatting or wrapping is not needed
- [ ] Transition JSON, SQL, and scripting variables to a single style. Either camelCase or snake_case.
- [x] Add 'application/json; charset=utf-8' Content-Type assert in all of the tests and fix where needed
- [x] Rename 'authorizeNormalAdmin' utility function or rewrite to actually authorize normal admins.
- [x] log -> log/slog package for logging level, and logging to file
- [x] Ensure all POST responses are 201 (return body of resource created makes sense to have a 201)
- [ ] test the following in sqlc.yaml config: `emit_pointers_for_null_types`
