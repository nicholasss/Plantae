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
- [ ] Implement plant type management.
- [ ] Implement plant light need management.
- [ ] Implement basic end-user's plant tracking.
- [ ] Begin entering plant information into plant table.
- [ ] Create a backup scheme for the universal plant species data.

### Cleanup efforts

- [ ] Fix respondWithError utility function to not reveal the error value to client
- [ ] Replace all fmt.Errorf with errors.New where formatting or wrapping is not needed
- [ ] Transition JSON, SQL, and scripting variables to a single style. Either camelCase or snake_case.
- [ ] Add 'application/json' Content-Type assert in all of the tests and fix where needed
- [ ] Rename 'authorizeNormalAdmin' utility function or rewrite to actually authorize normal admins.
