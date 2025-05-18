# Scratchpad - A document to keep my thoughts and ideas in a structured manner

## Brainstorming

1. [X] Add exact requirements for the project as its easy to miss details when project grows
2. [X] Check if there is anything needed from the frontend that needs to be integrated into the backend that its not mentioned in Christians backend notes
3. [X] To kick-off the project, lets just get up a webserver to run that answers PONG to a PING or something like that. Better start small so can debug along the way
4. [ ] Add some sample config to `main.go` just to be able to set and read `.env`variables
5. [ ] Add a singleton pattern that generates a `Config`instance. Keep it in `main.go`for now. Later can move it somewhere else

## Random thoughts

1. [X] It was more easier than i expected with the webserver. Really smoooth. Result below:

## Findings and learnings

1. Go handles "" as a string and '' as a Rune. IMPORTANT) There is a difference between `"a"` and `'a'`. It type infers `"a"`as a string and `'a' as a Rune.

```bash
curl http://localhost:8080/ping
PONG
```

## Major decisions and why i took them

1. [ ] Lets try to implement all additional features - Christian mentioned **atleast 1** but lets be ambitious
2. [X] Implement small webserver just to get something up and running and take it from there