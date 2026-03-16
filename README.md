### ESDEP
Easy deployment system that runs as a background service
It utilizes ssh keys for pulling git repos

### Building
1. Clone the thing
2. Fetch deps
```
go mod tidy
```
3. Build it
```
go build -o whatever ./cli/main.go
```

### Config
You can find sample config in ```config.yaml``` file