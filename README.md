## Instruction
### Run with Docker
1. Make sure you have Go 1.22 and Docker installed on your machine.
2. Clone the repository.
3. Build the docker image:
```bash
    docker build -t yadro-computer-club .
```

4. Run the docker container with input file mounted:
```bash
    chmod u+x ./run.sh && ./run.sh ./examples/1
```

### Run without Docker
1. Make sure you have Go 1.22 installed on your machine.
2. Clone the repository.
3. Run:
```bash
    go run ./cmd/main.go ./examples/1
```

### Examples are located at `./examples` directory.