# A simple To-Do api written in Go

## Usage

Run the service locally either by running from source or by using the provided docker image.

## Notes

- ListItems with a redis backing is currently unsupported

## Technical decisions

### Generating docs from code

We have the option of writing the code or the documentation first, and generating the other.
A choice was made to have the documentation generated from the code.
The reasoning was two fold;
- It allowed me to easily the maintain the code from a single source of truth (the `*.go` files)
- It allowed me to iterate quickly writing the code first.

### Using the httprouter library

While you can parse fields like `items/{{UUID}}` in plain Go, it can be a bit trickier and an unpleasant dev experience.
A choice was made to just use a router than a whole framework because the standard library Go offers is incredibly good.

### Middlewares

While we could just use the standard library, alice offers a simple alternative to call multiple middleware handlers.

    - recover middleware was chosen to ensure our service doesn't crash when there's a panic (there shouldn't be any panics, but still...)
    - logging middleware was chosen so that we can log every request that goes through our service, Zap was chosen as the logging framework as it has a simple to use API that allows us to extend our logging capability

### Data store

We use a datastore interface for our datastore layer, this allows us to implement multiple backends and so long as they support the operations provided they can be used.
This allowed us to create an in memory datastore in the first instance, and replace it with redis later.
Obviously the more datastore types we support the more work we need to do to maintain them all.

### Testing level

A decision was made to test the server using the in memory datastore.
This allows minimal use of mocks (none) to test the whole flow of the system.

## Todo

Todo list:
- api tests
- unit tests
- better errors
- helpful error messages (notfound etc)
- Makefile
- configure
    - ports
    - redis address
- documentation
    - openapi
    - comments
    - README
- healthcheck
- metrics?
- change in mem key from UUID to string