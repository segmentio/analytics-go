# Releasing

1. **Update the version constant** in `analytics.go`:

   ```go
   const Version = "x.y.z"
   ```

2. **Update `History.md`** — change the `Unreleased` heading to `x.y.z / YYYY-MM-DD`.
3. Commit both files:

   ```
   git commit -m "Release vx.y.z"
   ```

4. Open a PR and merge it to `v3.0`.
5. Tag the merged commit:

   ```
   git tag vx.y.z
   git push origin vx.y.z
   ```

> **Note:** The `Version` constant in `analytics.go` is the single source of truth — it is
> sent in the `User-Agent` header on every API request and is automatically used by the
> test fixtures. No other files need to be updated.

> **Note:** The tag *is* the release. There is no publish workflow and no credentials;
> `proxy.golang.org` serves the module from the tag. Tags here carry a `v` prefix, unlike
> the other Segment SDKs.
