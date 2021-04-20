# A simple To-Do api written in Go

## Usage

Run the service locally either by running from source or by using the provided docker image.


## Technical decisions

### Generating docs from code

We have the option of writing the code or the documentation first, and generating the other.
A choice was made to have the documentation generated from the code.
The reasoning was two fold;
- It allowed me to easily the maintain the code from a single source of truth (the `'.go` files)
- It allowed me to iterate quickly writing the code first.

### Using the httprouter library

While you can parse fields like `items/{{UUID}}` in plain Go, it can be a bit trickier and an unpleasant dev experience.
A choice was made to just use a router than a whole framework because the standard library Go offers is incredibly good.

### Middlewares

While we could just use the standard library, alice offers a simple alternative to call multiple middleware handlers.

    - recover middleware was chosen to ensure our service doesn't crash when there's a panic (there shouldn't be any panics, but still...)
    - logging middleware was chosen so that we can log every request that goes through our service, Zap was chosen as the logging framework as it has a simple to use API that allows us to extend our logging capability

### API testing

### Data store

## Todo

Todo list:
- api tests
- helpful error messages (notfound etc)
- external database (postgres/redis?)
- unit tests
- Makefile
- configure ports
- documentation
    - openapi
    - comments
    - README