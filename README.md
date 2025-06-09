# Plantae Server

Project to track your house plants. This is the server portion of the project.

## Style

JSON is camelCase.
Endpoints are kebab-case.
Scripting variables (Bash, Hurl) are in snake_case.

## Goals

- Track each house plant.
- Surface information about its toxicity (to humans or pets).
- Organize plants into biomes or into rooms.
- Plan out your watering schedule.

## Todo

- [ ] Implement basic user management.
- [ ] Implement basic plant tracking.
- [ ] Begin entering plant information into plant table.
- [ ] Create a backup scheme for the universal plant data.

### Cleanup efforts

- [ ] Fix respondWithError utility function to not reveal the error value to client
- [ ] Replace all fmt.Errorf with errors.New where formatting or wrapping is not needed
