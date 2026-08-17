# exitcodes

Generates exit code constants for various languages, based on [this list compiled in the systemd docs](https://www.freedesktop.org/software/systemd/man/latest/systemd.exec.html#Process%20Exit%20Codes).

This repo is the source of truth. It produces the contents of two downstream repos:

| Language | Repo | Package |
| --- | --- | --- |
| Go | [cdzombak/exitcode_go](https://github.com/cdzombak/exitcode_go) | `github.com/cdzombak/exitcode_go` |
| Python | [cdzombak/exitcode_py](https://github.com/cdzombak/exitcode_py) | `exitcode` |

## How it works

`exitcodes.tsv` lists every exit code: its numeric value, its canonical name, optional per-language name overrides, a description, and the standard it comes from. `main.go` embeds that TSV along with `golang.tmpl` and `python.tmpl`, and renders one file per language:

```shell
go run main.go -go ../exitcode_go/generated.go -py ../exitcode_py/exitcode/__init__.py
```

Either flag may be omitted to generate only one language.

To add or change an exit code, edit `exitcodes.tsv` — never the generated files. The name columns exist because a mechanical transformation of the C-style name isn't always right: `EX_NOPERM` and `EXIT_NOPERMISSION` would otherwise collide, and acronyms like `CPU`, `PAM`, and `SELinux` don't camel-case correctly on their own.

## Automation

The `Sync Generated Code` workflow runs on every push to `main`. It regenerates both languages, and for each downstream repo opens a pull request if the output changed. Branches are named after a hash of the generated file, so an unchanged run reuses the existing PR instead of opening a duplicate; superseded PRs are closed automatically.

Releases are cut in the downstream repos, by tagging them there — not here.

## License

MIT; see [LICENSE](LICENSE).

## Author

Chris Dzombak ([dzombak.com](https://www.dzombak.com); [GitHub @cdzombak](https://github.com/cdzombak)).
