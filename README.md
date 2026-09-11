# lapip

`lapip` is a Go rewrite of [Pipal](https://github.com/digininja/pipal), the
password dump analysis tool. It reads one password per line and writes the most
common passwords as a clean wordlist, one password per line.

## Install

Linux and macOS:

```sh
curl -fsSL https://raw.githubusercontent.com/ESHAYAT102/lapip/main/scripts/install.sh | sh
```

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/ESHAYAT102/lapip/main/scripts/install.ps1 | iex
```

The binary is installed to `~/.local/bin/lapip` on Linux and macOS, or
`$HOME\.local\bin\lapip.exe` on Windows.

## Usage

```sh
go build -o lapip .
./lapip passwords.txt
```

Options:

```text
-t int              number of top results (default 10)
-o file             write the report to a file instead of stdout
-m                 emit Markdown output
-numbers           add leetspeak variants and every 3- and 4-digit suffix, with no separator or one of .@#$!%&*_-+=
```

Examples:

```sh
./lapip -t 20 -o report.txt passwords.txt
./lapip -m passwords.txt > report.md
./lapip -numbers -o candidates.txt words.txt
```

`-numbers` independently keeps, uppercases, lowercases, or substitutes each letter:
`a → @`, `e → 3`, `s → 5`, `i/l → 1`, `o → 0`, `t → 7`, `b → 8`, `g → 9`
(uppercase letters work too). Every variant receives the numeric suffixes.
For example, `PASS` produces `P@SS`, `P@5S`, `P@S5`, `P@55`, `PA5S`, `PAS5`,
`PA55`, and `PASS`, plus lowercase and mixed-case combinations.
Variants are streamed, but output grows exponentially: `PASS` has 54 variants,
producing 7,722,054 lines including all numeric suffixes and separators.

## Uninstall

Linux and macOS:

```sh
curl -fsSL https://raw.githubusercontent.com/ESHAYAT102/lapip/main/scripts/uninstall.sh | sh
```

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/ESHAYAT102/lapip/main/scripts/uninstall.ps1 | iex
```

Input is processed as a stream, with a scanner buffer sized for password lines
up to 16 MiB. The original input is not modified. Use this tool only on
password data you are authorized to analyze.

## Development

```sh
go test ./...
go build ./...
```

## Scope

The analyzer keeps Pipal's frequency-based ranking behavior. Pipal's Ruby
checker and splitter plugin system, external word-list comparison, and optional
geographic checkers are not loaded automatically.
