// This file is not shipped into generated services — the generator only walks
// templates/base and templates/types, never the templates root.
//
// It exists so the Go toolchain treats templates/ as a separate module and stops
// descending into it. Without it, `go build ./...` and `go test ./...` try to
// compile the template tree as real packages and fail on unresolvable imports.
module gomicrogen.local/templates

go 1.23.0
