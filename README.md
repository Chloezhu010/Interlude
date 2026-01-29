# Interlude

A terminal companion for [Claude Code](https://claude.ai/code) that keeps you in flow while waiting for responses.

## Features

While Claude thinks, Interlude offers micro-activities to keep you engaged:

- **Dev Jokes** — Quick laughs to pass the time
- **CS Trivia** — Multiple choice quiz questions
- **CS Flashcards** — Review core concepts (Big O, data structures, algorithms)

## Requirements

- **macOS** or **Linux** (Windows via WSL2 only)
- Go 1.21+
- [Claude Code](https://claude.ai/code) CLI

## Installation

```bash
# With Go 1.21+
go install github.com/Chloezhu010/Interlude/cmd/interlude@latest

# Or build from source
git clone https://github.com/Chloezhu010/Interlude.git
cd Interlude
make build    # or: go build -o interlude ./cmd/interlude
```

## Quick Start

```bash
# One-time setup (configures Claude Code hooks)
interlude init

# Start listening
interlude start
```

That's it. Interlude will appear when Claude has been working for 5+ seconds.

## Controls

| Key | Action |
|-----|--------|
| `j` | Dev joke |
| `t` | Trivia quiz |
| `f` | CS flashcard |
| `n` | Next (joke/question/card) |
| `space` | Reveal flashcard answer |
| `1-4` | Answer trivia |
| `b` / `Esc` | Back to menu |
| `q` | Quit |

## How It Works

Interlude uses Claude Code hooks to detect activity:

1. `interlude init` adds hooks to `~/.claude/settings.json`
2. When you submit a prompt, the hook writes "active" to `~/.interlude/status`
3. When Claude finishes, another hook writes "idle"
4. `interlude start` polls this file and shows the TUI after 3 seconds

## Customization

### Add your own jokes or flashcards

These are embedded at build time. To customize:

1. Fork/clone the repo
2. Edit the JSON files:
   - `internal/fun/jokes_data.json` — array of strings
   - `internal/fun/flashcards_data.json` — array of `{question, answer, explanation}`
3. Rebuild: `make build`

### Modify trivia questions

Trivia is cached locally after first fetch. Edit directly:

```bash
# Edit the cache (no rebuild needed)
$EDITOR ~/.interlude/trivia.json
```

Format:
```json
[
  {
    "question": "Your question?",
    "correct_answer": "Right answer",
    "incorrect_answers": ["Wrong 1", "Wrong 2", "Wrong 3"]
  }
]
```

To refresh from API, delete the cache:
```bash
rm ~/.interlude/trivia.json
```

## Credits

- **Trivia**: [Open Trivia Database](https://opentdb.com/) — Science: Computers category (CC BY-SA 4.0)
- **Jokes**: Curated from [devjoke](https://github.com/shrutikapoor08/devjoke)
- **TUI**: Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)

## License

MIT
