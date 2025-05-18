# Scratchpad - A document to keep my thoughts and ideas in a structured manner

## Brainstorming

1. [X] Add exact requirements for the project as its easy to miss details when project grows
2. [X] Check if there is anything needed from the frontend that needs to be integrated into the backend that its not mentioned in Christians backend notes
3. [X] To kick-off the project, lets just get up a webserver to run that answers PONG to a PING or something like that. Better start small so can debug along the way
4. [X] Add some sample config to `main.go` just to be able to set and read `.env`variables
5. [X] Add a factory pattern that generates a `Config`instance. Keep it in `main.go`for now. Later can move it somewhere else
6. [ ] Probably will need a custom multiplexer for route organisation and middleware
7. [X] Need to check the ListenAndServe() function if its blocking or not. Like do i need to wrap it inside a go-routine or can i just run it without go-routine?

## Random thoughts

1. [X] It was more easier than i expected with the webserver. Really smoooth. Result below:

```bash
curl http://localhost:8080/ping
PONG
```
## Findings and learnings

1. Go handles "" as a string and '' as a Rune. IMPORTANT) There is a difference between `"a"` and `'a'`. It type infers `"a"`as a string and `'a' as a Rune.
2. Need to have dependency to read from `.env`files. Something similar to `dotenv`. Lets not add that now and instead pass them into via `bash`
3. `ListenAndServe()`is blocking so this needs to be wrapped in a go routine. Seen this similar pattern as IFFE functions as in ts/js
4. There is also a `ListenAndServeTLS`function that needs to be used if you deploy so it can handle https


## Major decisions and why i took them

1. [ ] Lets try to implement all additional features - Christian mentioned **atleast 1** but lets be ambitious
2. [X] Implement small webserver just to get something up and running and take it from there
3. [ ] Will use custom multiplexer to avoid global state issues as well as organize routes modularly for scalability, and enable middleware like logging, metrics and ofcourse CORS in deployment