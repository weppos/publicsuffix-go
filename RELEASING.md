# Releasing

This document describes the steps to release a new version of PublicSuffix for Go.

## Prerequisites

- You have commit access to the repository
- You have push access to the repository
- You have a GPG key configured for signing tags

## Release process

1. **Determine the new version** using [Semantic Versioning](https://semver.org/)

   ```shell
   VERSION=X.Y.Z
   ```

   - **MAJOR** version for incompatible API changes
   - **MINOR** version for backwards-compatible functionality additions
   - **PATCH** version for backwards-compatible bug fixes

2. **Update the version constant** with the new version

   Edit `publicsuffix/publicsuffix.go` and update the `Version` constant:

   ```go
   const (
       // Version identifies the current library version.
       // This is a pro forma convention given that Go dependencies
       // tends to be fetched directly from the repo.
       Version = "X.Y.Z"
       ...
   )
   ```

3. **Update the changelog** with the new version

   Edit `CHANGELOG.md` and add a new section for the release:

   ```markdown
   ## X.Y.Z

   - Description of changes
   ```

4. **Update dependencies**

   ```shell
   go get -u ./...
   go mod tidy
   ```

5. **Run tests** and confirm they pass

   ```shell
   make test
   ```

   or

   ```shell
   go test ./... -v
   ```

6. **Commit the new version**

   ```shell
   git commit -am "Release $VERSION"
   ```

7. **Create a signed tag**

   ```shell
   git tag -a v$VERSION -s -m "Release $VERSION"
   ```

8. **Push the changes and tag**

   ```shell
   git push origin main
   git push origin v$VERSION
   ```

## Post-release

- Verify the new version appears on [GitHub releases](https://github.com/weppos/publicsuffix-go/releases)
- Verify the new version is available via `go get github.com/weppos/publicsuffix-go@v$VERSION`
- Update [pkg.go.dev](https://pkg.go.dev/github.com/weppos/publicsuffix-go) (usually updates automatically)
- Announce the release if necessary

## Notes

- Go modules use git tags for versioning, so pushing the tag is what makes the new version available
- Unlike Ruby gems, there is no separate "publish" step - the tag itself publishes the version
- Go uses the `v` prefix for version tags (e.g., `v1.0.0`)
