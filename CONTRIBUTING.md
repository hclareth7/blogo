# Contributing to BLOGO

Thank you for your interest in contributing to BLOGO.

## Development Philosophy

BLOGO follows **Specification Driven Development (SDD)**. Every feature starts with a specification before implementation.

```
Specification -> Review -> Implementation -> Testing -> Release
```

## How to Contribute

### Reporting Issues

- Search existing issues before creating a new one.
- Include steps to reproduce, expected behavior, and actual behavior.
- Include your Go version and operating system.

### Proposing Features

1. Open an issue describing the feature and its use case.
2. If accepted, write a specification in `specs/` following the existing format.
3. Get the specification reviewed before starting implementation.

### Submitting Code

1. Fork the repository.
2. Create a feature branch from `main`.
3. Ensure your change has a corresponding spec in `specs/` (for new features).
4. Write tests for your changes.
5. Run the full quality pipeline before submitting:

```bash
make fmt
make lint
make test
```

6. Submit a pull request with a clear description of the change.

## Code Standards

- Follow standard Go conventions and idioms.
- Format code with `gofumpt`.
- Pass `golangci-lint` with no warnings.
- Use `slog` for structured logging.
- Write tests with `t.Parallel()` by default.
- Escape all HTML output via Go templates.
- Use parameterized statements for any queries.

## Commit Messages

- Use imperative mood: "Add feature" not "Added feature".
- Keep the first line under 72 characters.
- Reference issue numbers when applicable.

## Project Structure

See `CLAUDE.md` for the full project structure and conventions.

## Content vs Code

BLOGO's code is licensed under **Apache 2.0**. The content displayed by BLOGO belongs to its respective authors and is licensed under their original terms (e.g., **CC BY-NC-ND 4.0**). Contributions must not modify the original content.

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).
