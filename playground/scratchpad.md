# Scratchpad - A document to keep my thoughts and ideas in a structured manner

## Brainstorming

1. [X] Add exact requirements for the project as its easy to miss details when project grows
2. [X] Check if there is anything needed from the frontend that needs to be integrated into the backend that its not mentioned in Christians backend notes
3. [X] To kick-off the project, lets just get up a webserver to run that answers PONG to a PING or something like that. Better start small so can debug along the way
4. [X] Add some sample config to `main.go` just to be able to set and read `.env`variables
5. [X] Add a factory pattern that generates a `Config`instance. Keep it in `main.go`for now. Later can move it somewhere else
6. [ ] Probably will need a custom multiplexer for route organisation and middleware
7. [X] Need to check the ListenAndServe() function if its blocking or not. Like do i need to wrap it inside a go-routine or can i just run it without go-routine?
8. [X] I think the best way to separate out logic is to like seperate out like handling of http requests and responses, then the handling of the routes, another separation is like business logic. Also a seperation for data storage management. Maybe a little overkill for our minimnal application but i think maybe its good to show usually how i usually work in larger projects and how i like to structure my applications.

## Random thoughts

1. [X] It was more easier than i expected with the webserver. Really smoooth. Result below:
2. [X] Testing was initially very confusing but once you understood the concept by reading `https://pkg.go.dev/testing` it became much more clear.

```bash
curl http://localhost:8080/ping
PONG
```
## Findings and learnings

1. Go handles "" as a string and '' as a Rune. IMPORTANT) There is a difference between `"a"` and `'a'`. It type infers `"a"`as a string and `'a' as a Rune.
2. Need to have dependency to read from `.env`files. Something similar to `dotenv`. Lets not add that now and instead pass them into via `bash`
3. `ListenAndServe()`is blocking so this needs to be wrapped in a go routine. Seen this similar pattern as IFFE functions as in ts/js
4. There is also a `ListenAndServeTLS`function that needs to be used if you deploy so it can handle https
5. %v in error messages: Formats values in a default, human-readable way (e.g., t.Errorf("Validate() error = %v, want %v", err, tt.wantErr) shows error details clearly).
6. %q for strings: Formats strings with quotes, escaping special characters (e.g., t.Errorf("Validate() Title = %q, want %q", tt.book.Title, got) ensures safe string output).
7. %w for error wrapping: Wraps errors to preserve the original error while adding context (e.g., fmt.Errorf("%w: %s", ErrInvalidStatus, status) allows errors.Is to check the wrapped error).


## Major decisions and why i took them

1. [ ] Lets try to implement all additional features - Christian mentioned **atleast 1** but lets be ambitious
2. [X] Implement small webserver just to get something up and running and take it from there
3. [X] Will use custom multiplexer to avoid global state issues as well as organize routes modularly for scalability, and enable middleware like logging, metrics and ofcourse CORS in deployment
4. [ ] As Book datastructure very simple. I choose not to add any regex checks or any third party validation libraries like typescripts Zod (but for go). Handled mostly with strings.TrimParse(). Can also later add validation to the frontend to be extra secure.
5. [ ] For separation of concerns and for scalability. I'll modularize the applications in:
   1. Handlers: Manages HTTP requests and responses
   2. Routes: Managed routing
   3. Services: Business logic
   4. Store: Data storage and retrieval