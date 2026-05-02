// Package command provides CLI commands for the jay application.
package command

// Here is a tree visualisation of the command structure showing where
// flags are defined.
//
// jay (root)
// │   PersistentFlags:
// │     --tui / -t
// │     --theme
// │     --language
// │
// ├── nav (ghost, hidden)
// │   │   PersistentFlags:
// │   │     --subscribe / -s
// │   │     --action / -a
// │   │     --pipeline / -p
// │   │     [cascade family]
// │   │       --depth
// │   │       --no-recurse / -N
// │   │     [sampling family]
// │   │       --sample
// │   │       --num-files
// │   │       --num-folders
// │   │       --last
// │   │     [preview family]
// │   │       --dry-run
// │   │     [poly-filter family]
// │   │       --files-glob / -b
// │   │       --file-regex / -x
// │   │       --folders-glob / -g
// │   │       --folders-regex / -y
// │   │
// │   ├── exec (ghost, hidden)
// │   │   │   PersistentFlags:
// │   │   │     --resume / -r
// │   │   │   Constraints:
// │   │   │     MarkFlagsOneRequired("action", "pipeline")
// │   │   │
// │   │   ├── walk
// │   │   │     LocalFlags: (none)
// │   │   │
// │   │   └── run
// │   │         LocalFlags:
// │   │           [worker-pool family]
// │   │             --cpu
// │   │             --now
// │   │
// │   └── query
// │         LocalFlags: (none)
// │
// ├── verify
// │     LocalFlags: (tbd)
// │
// └── theme
//       LocalFlags: (tbd)
