#!/usr/bin/env bash
# nolint-revive-ginkgo-tests
# Adds //nolint:revive // ok to Ginkgo/Gomega dot-imports in Go _test.go files.
# Compatible with bash and zsh on macOS and Linux.
#
# Usage:
#   source nolint-revive-ginkgo-tests.sh
#   nolint-revive-ginkgo-tests [--dry-run] [path]
#
# Arguments:
#   --dry-run   Preview changes without modifying any files.
#   path        Root directory to search (defaults to current directory).

nolint-revive-ginkgo-tests() {
    local dry_run=false
    local search_root="."

    # ── Argument parsing ────────────────────────────────────────────────────────
    for arg in "$@"; do
        case "$arg" in
            --dry-run) dry_run=true ;;
            -*) echo "Unknown flag: $arg" >&2; return 1 ;;
            *)  search_root="$arg" ;;
        esac
    done

    # ── Patterns to match ───────────────────────────────────────────────────────
    local ginkgo_pattern='. "github.com/onsi/ginkgo/v2"'
    local gomega_pattern='. "github.com/onsi/gomega"'
    local nolint_marker='//nolint:revive'
    local nolint_comment='//nolint:revive // ok'

    local files_affected=0
    local modifications=0

    # ── Walk all *_test.go files ────────────────────────────────────────────────
    # Use -print0 / read -d '' to handle filenames with spaces safely.
    while IFS= read -r -d '' file; do

        local file_modified=false
        local file_mod_count=0
        local line_num=0
        local output_lines=()

        # Read the file line-by-line into an array for processing.
        # 'mapfile' (bash 4+) is used; fall back to a while-read loop for zsh / older bash.
        local lines=()
        while IFS= read -r line || [[ -n "$line" ]]; do
            lines+=("$line")
        done < "$file"

        for line in "${lines[@]}"; do
            (( line_num++ ))

            # Check whether this line is a target import and lacks the marker.
            local is_target=false
            if [[ "$line" == *"$ginkgo_pattern"* || "$line" == *"$gomega_pattern"* ]]; then
                if [[ "$line" != *"$nolint_marker"* ]]; then
                    is_target=true
                fi
            fi

            if $is_target; then
                (( modifications++ ))
                (( file_mod_count++ ))
                file_modified=true

                if $dry_run; then
                    if ! $file_modified || [[ $file_mod_count -eq 1 ]]; then
                        # Print file header only on the first match in this file.
                        :
                    fi
                    printf '  📄 %s  (line %d)\n' "$file" "$line_num"
                    printf '     - %s\n' "$line"
                    printf '     + %s %s\n' "$line" "$nolint_comment"
                fi

                output_lines+=("${line} ${nolint_comment}")
            else
                output_lines+=("$line")
            fi
        done

        if $file_modified; then
            (( files_affected++ ))

            if $dry_run; then
                printf '📁 %s — %d line(s) would be modified\n\n' "$file" "$file_mod_count"
            else
                # Write the modified content back directly.
                # Check for trailing newline before overwriting.
                local has_trailing_newline=false
                [[ "$(tail -c1 "$file" | wc -c)" -eq 1 ]] && has_trailing_newline=true

                local write_first=true
                for out_line in "${output_lines[@]}"; do
                    if $write_first; then
                        printf '%s' "$out_line"
                        write_first=false
                    else
                        printf '\n%s' "$out_line"
                    fi
                done > "$file"
                $has_trailing_newline && printf '\n' >> "$file"
            fi
        fi

    done < <(find "$search_root" -type f -name '*_test.go' -print0 2>/dev/null)

    # ── Summary ─────────────────────────────────────────────────────────────────
    echo ""
    if $dry_run; then
        printf '⛔️  %d file(s) will be affected\n' "$files_affected"
        printf '♦️   %d modification(s) will be made\n' "$modifications"
    else
        printf '✅  %d file(s) affected\n' "$files_affected"
        printf '🟢  %d modification(s) made\n' "$modifications"
    fi
}
