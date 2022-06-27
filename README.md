# Croupier

An API to handle playable card operations.

##### Table of contents

* :hammer_and_pick: [Build](#hammer_and_pick-build)
* :zap: [Execution](#zap-execution)
* :sparkles: [Testing](#sparkles-testing)

## :hammer_and_pick: Build

From the root directory, build the project:
```sh
make all
```
This will run the tests before building the actual binary, necessary to run the API.
By running the tests first, we also ensure no breaking changes have been introduced.


## :zap: Execution

Once built, the binary is exported to a `target` directory. From the root directory, start
the API:
```sh
./target/croupier
```
By default, the port `8080` will be used, although it can be changed by setting the
`PORT` environment variable.

The API exposes three endpoints:
- POST `/decks`
    - Creates a deck of cards.
    - If desired:
      - Provide a request body with:
        - `shuffled` (bool) to create a shuffled deck.
      - Provide `cards`, the card codes e.g. `AS` for `Ace of Spades`, as a query parameter to
        create a partial deck.
- GET `/decks/:id`
    - Retrieves the deck associated with the provided ID.
- POST `/decks/:id/cards/draw`
  - Draws a certain number cards from the deck associated with the provided ID.
  - The number of cards to draw `count` must be provided as a query
    parameter.


## :sparkles: Testing

From the root directory, run unit tests:
```sh
make test
```
This will simply run tests. On the other hand, if we do want to export the coverage:
```sh
make test_coverage
```
This will export the coverage in a file (`coverage.out`) in a `target` directory. It can
easily be visualised running this command, from the root directory:
```sh
go tool cover -html=./target/coverage.out
```
