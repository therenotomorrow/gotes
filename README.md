# gotes

**gotes** is the homework for gRPC production course, you could find more here: [link](https://education.easyp.tech).

## Development

### System Requirements

```shell
go version
# go version go1.25.0 or higher

task --version
# 3.46.3 or higher (https://taskfile.dev)

psql --version
# psql (PostgreSQL) 18.1

redis-server --version
# Redis server v=8.4.0
```

### Download sources

```shell
PROJECT_ROOT=gotes
git clone https://github.com/therenotomorrow/gotes.git "$PROJECT_ROOT"
cd "$PROJECT_ROOT"
```

### Setup dependencies

```shell
# install dependencies
go mod tidy

# check code integrity
task tools:install qa # see other recipes by calling `task`

# apply migrations
task services:postgres
task services:goose -- up
task services:redis

# setup safe development (optional)
git config --local core.hooksPath .githooks
```

### Testing

```shell
# run quick checks
task test:smoke

# run with coverage
task test:cover
```

## License

MIT License. See the [LICENSE](./LICENSE) file for details.
