# Dependency Management

Go dependencies are automatically monitored and kept up-to-date using Dependabot.

## Dependabot Configuration

Configuration is defined in [`.github/dependabot.yaml`](../.github/dependabot.yaml):

- **Ecosystem**: Go modules (`gomod`)
- **Monitored Directories**: Root directory (`/`) and `src/` subdirectory
- **Update Schedule**: Weekly

Dependabot automatically creates pull requests for dependency updates, allowing the team to review and merge them routinely.

## Manual Dependency Updates

To manually update dependencies:

```sh
go get -u ./...
go mod tidy
```

## Special Cases

### govips

The `github.com/davidbyttow/govips` dependency is **intentionally ignored by Dependabot** and must be upgraded manually. Monitor this library separately and test thoroughly before upgrading due to its critical role in image processing operations.

To upgrade govips manually:

```sh
go get -u github.com/davidbyttow/govips
go mod tidy
go test -v nuggan  # Test thoroughly after updating
```

## Workflow

1. Dependabot creates PRs for routine dependency updates (weekly schedule).
2. Team reviews and tests updates.
3. Merge approved PRs.
4. Manually check and upgrade `govips` as needed, with thorough testing.
