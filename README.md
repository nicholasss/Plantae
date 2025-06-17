# Plantae Server

A REST API server to track your house plants over time.
This is the server portion of the project. There is currently no front-end.

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

- [x] Implement basic user management.
- [ ] Implement basic admin management.
- [x] Implement plant species management.
- [ ] Implement basic plant tracking.
- [ ] Begin entering plant information into plant table.
- [ ] Create a backup scheme for the universal plant data.

### Cleanup efforts

- [ ] Fix respondWithError utility function to not reveal the error value to client
- [ ] Replace all fmt.Errorf with errors.New where formatting or wrapping is not needed
- [ ] Formalize JSON, SQL, and scripting variables to a single style. Either camelCase or snake_case.
- [ ] Add 'application/json' Content-Type assert in all of the tests and fix where needed
- [ ] Rename 'authorizeNormalAdmin' utility function or rewrite to actually authorize normal admins.
