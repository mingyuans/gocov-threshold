# gocov-threshold
This action is used to calculate the percentage of new code that is covered by tests in a Go project.
It supports global patterns and annotations to ignore certain files or lines from the coverage calculation.


# Quick Start
```yaml
      - name: Go Test
        run: |
          go test ./example/... -coverprofile=coverage.out
      
      - name: Coverage Threshold
        id: coverage_threshold
        uses: mingyuans/gocov-threshold@main
        with:
          path: .
          coverprofile: coverage.out
          module: github.com/mingyuans/gocov-threshold
          token: ${{ secrets.GITHUB_TOKEN }}
          conf: gocov-conf.yaml
          threshold: 80

      - name: Comment
        uses: mshick/add-pr-comment@v2
        with:
          message: |
            Coverage on new code: ${{ steps.coverage_threshold.outputs.gocov }}%
```


# Action Parameters

## Inputs

| Name           | Description                             | Required | Default |
|----------------|-----------------------------------------|----------|---------|
| `path`         | path to git repo                        | No       | `.`     |
| `coverprofile` | path to coverage profile                | Yes      | ``      |
| `module`       | the Go module name                      | Yes      | ``      |
| `threshold`    | coverage threshold (0-100)              | No       | `80`    |
| `logger-level` | logger level (debug, info, warn, error) | No       | `info`  |
| `conf`         | the config file                         | Yes      | ``      |

## Outputs

| Name    | Description                           |
|---------|---------------------------------------|
| `gocov` | the coverage difference (0.00-100.00) |


## Example of config file

```yaml
files:
  include:
    dirs:
  #   Only the directories that contain Go files will be considered
      - example
    patterns:
      - "*.go"

  exclude:
    dirs:
  #   The files that will be excluded from the coverage calculation
      - testdata
    patterns:
      - "*_mock.go"

statements:
  exclude:
  # Exclude statements that match these patterns
    patterns:
      - ".*mu.Lock.*"
```

## Preset annotations
The file will be ignored if it contains the following annotations:
```go
// gocover:ignore
package main

func CalculateSum(a, b int) int {
	return a + b
}
```

`println("Not divisible by 3")` will be ignored in the coverage calculation
```go
func exampleFunc2(value int) {
	if value%3 == 0 {
		println("Divisible by 3")
	} else {
		// gocover:ignore
		println("Not divisible by 3")
	}
}
```
