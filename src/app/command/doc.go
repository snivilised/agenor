// Package command provides CLI commands for the jay application.
package command

// Here is a tree visualisation of the command structure showing where
// flags are defined.
//
// Design note: cobra propagates persistent flags from parent to child
// commands. An earlier design used hidden ghost commands (nav, exec) as
// intermediaries so that nav-level flags (--subscribe, --action, etc.)
// were only inherited by navigation commands and not by utility commands
// (verify, theme). However, ghost commands do not appear in --help output,
// making walk, sprint, and query undiscoverable to the user. Since user
// experience takes priority, the ghost commands have been removed and each
// navigation leaf command now registers its own copy of the nav flags and
// families directly on its local flag set. Flag isolation is preserved
// because local flags (cmd.Flags()) do not propagate to siblings or parents.
//
// jay (root)
// │   PersistentFlags:
// │     --tui / -t
// │     --theme
// │
// ├── walk
// │     Flags:
// │       --subscribe / -s
// │       --action / -a
// │       --pipeline / -p
// │       --resume / -r
// │       [preview family]
// │         --dry-run
// │       [cascade family]
// │         --depth
// │         --no-recurse / -N
// │       [sampling family]
// │         --sample
// │         --num-files
// │         --num-folders
// │         --last
// │       [poly-filter family]
// │         --files-glob / -b
// │         --file-regex / -x
// │         --folders-glob / -g
// │         --folders-regex / -y
// │
// ├── sprint
// │     Flags:
// │       --subscribe / -s
// │       --action / -a
// │       --pipeline / -p
// │       --resume / -r
// │       [preview family]
// │         --dry-run
// │       [cascade family]
// │         --depth
// │         --no-recurse / -N
// │       [sampling family]
// │         --sample
// │         --num-files
// │         --num-folders
// │         --last
// │       [poly-filter family]
// │         --files-glob / -b
// │         --file-regex / -x
// │         --folders-glob / -g
// │         --folders-regex / -y
// │       [worker-pool family]
// │         --cpu
// │         --now
// │
// ├── query
// │     Flags:
// │       --subscribe / -s
// │       --action / -a
// │       --pipeline / -p
// │       [preview family]
// │         --dry-run
// │       [cascade family]
// │         --depth
// │         --no-recurse / -N
// │       [sampling family]
// │         --sample
// │         --num-files
// │         --num-folders
// │         --last
// │       [poly-filter family]
// │         --files-glob / -b
// │         --file-regex / -x
// │         --folders-glob / -g
// │         --folders-regex / -y
// │       (no --resume: query is read-only and cannot be resumed)
// │
// ├── verify
// │     Flags: (tbd)
// │
// └── theme
//       Flags: (tbd)
