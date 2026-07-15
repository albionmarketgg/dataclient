# Source this before running go/gcc/wails:  source scripts/env.sh
#
# Puts Go, your GOPATH bin, and a MinGW-w64 (UCRT) gcc on PATH — the cgo build
# (gopacket + the Npcap SDK) needs a C compiler. Override MINGW_BIN if your gcc
# lives elsewhere, or put machine-specific paths in scripts/env.local.sh (ignored).
export PATH="/c/Program Files/Go/bin:$HOME/go/bin:$PATH"

MINGW_BIN="${MINGW_BIN:-/mingw64/bin}"
[ -d "$MINGW_BIN" ] && export PATH="$MINGW_BIN:$PATH"

# Optional machine-local overrides (gitignored).
_here="$(cd "$(dirname "${BASH_SOURCE[0]:-$0}")" && pwd)"
[ -f "$_here/env.local.sh" ] && . "$_here/env.local.sh"
