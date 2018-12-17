# kongfig

Kongfig is a configuration management tool for the Kong API gateway.

## Usage

```bash
kongfig apply -f config.json --dry-run
```

### Available Commands

| Command   | Description                              |
| ---       | ---                                      |
| `apply`   | Apply a configuration to a Kong instance |
| `help`    | Help about any command                   |
| `version` | Print the version number of Kongfig      |

Use `kongfig [command] --help` for more information about a command.

## Contributing

1. Fork the project
2. Make your changes
3. `GOOS=darwin make`
4. The `kongfig` binary is now available

> **NOTE**: when building in OS X, you'll need to export the `GOOS` env, eg:

```bash
GOOS=darwin make
```

Dependencies are managed using [dep]. Please refer to its documentation if needed.

### Testing

Tests are run automatically on every build or via the `make test` target.
Additionaly you can run `make cover` to check your coverage.

[dep]: https://github.com/golang/dep
