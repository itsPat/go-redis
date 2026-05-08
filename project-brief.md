# Mini Redis with Pub/Sub — Learning Project (Go)

## How to Help Me — Read This First

**I am writing this code myself. Do not write it for me.**

I'm using this project to learn Go. If you build it for me, the project loses all its value. Your job is to assist, not implement.

**Please do:**
- Answer questions when I ask them
- Explain concepts and idioms when I'm fuzzy on something
- Review code I've already written and point out bugs, anti-patterns, or non-idiomatic Go
- Suggest *approaches* (not implementations) when I'm fully stuck
- Show tiny illustrative snippets of language features in isolation when teaching is needed (e.g. "here's how `select` works in 4 lines") — but not the real project code

**Please don't:**
- Write functions for me unprompted
- Refactor my code without being asked
- Provide complete implementations of milestones, even if I ask, unless I've made a real attempt and clearly explain where I'm stuck
- Skip ahead and write later milestones because they seem trivial
- Volunteer "here, I'll just show you" when I'm only asking a conceptual question

If I ask you to write code outright, push back once and check that's really what I want. I might be tired or taking a shortcut and need a nudge to stay in the driver's seat.

When I ask "how do I X?" — explain the concept, show a minimal *illustrative* snippet if it genuinely helps, then let me write the actual project code myself.

---

## What I'm Building

A minimal Redis-like server in Go that supports:

- Basic key-value operations: `SET`, `GET`, `DEL`
- Pub/sub messaging: `SUBSCRIBE`, `PUBLISH`
- Multiple concurrent client connections over TCP

I'll test it with `netcat` / `telnet` initially, and as a stretch goal, with the real `redis-cli` once I implement the RESP protocol.

---

## Why I'm Building This

I'm learning Go, and specifically I want to grok **goroutines and channels**.

A plain key-value store would only exercise goroutines + mutexes. Adding pub/sub forces idiomatic channel use — specifically the pattern where one goroutine owns shared state and other goroutines communicate with it via channels. This is the "share memory by communicating" / actor model that defines Go's concurrency philosophy. I want to feel why this pattern exists by building something that genuinely benefits from it.

---

## Architecture

Three components:

1. **TCP server** — Listens on a port, accepts connections, spawns one goroutine per client connection.
2. **Store** — A `map[string]string` protected by `sync.RWMutex`. Handles `SET` / `GET` / `DEL`. Standard shared-memory-with-mutex pattern.
3. **Hub** — The pub/sub coordinator. Runs in its own goroutine and owns a `map[topic][]chan Message`. Other goroutines never touch this map directly — they communicate with the hub via channels (`subscribe`, `unsubscribe`, `publish`).

### The Key Pattern (The Thing I'm Here to Learn)

The Hub does **not** use a mutex. Instead:

- One goroutine has exclusive ownership of the subscriber map
- That goroutine sits in a `select` statement reading from channels
- When a client subscribes, they hand the hub their personal message channel
- When someone publishes, the hub fans the message out to each subscriber's channel
- Each connection goroutine reads from its own channel and writes to its TCP socket

A PUBLISH flow looks like:

```
client A's goroutine
  → hub's publish channel
    → hub goroutine
      → subscriber B's message channel
        → client B's goroutine
          → TCP write
```

This is the actor model in Go. The whole point of this project is to internalize this pattern.

---

## Milestones

Work through these in order. Don't skip ahead unless I explicitly ask.

1. **TCP echo server.** `net.Listen`, accept loop, `go handleConn` that echoes lines back. Confirms I understand the goroutine-per-connection model.
2. **Line-based command parser.** Read lines, split by spaces. No RESP yet — I'll test with `nc localhost 6379`.
3. **SET / GET / DEL** backed by the mutex-protected store.
4. **The Hub.** Build it as a standalone goroutine with its channel plumbing in place. No client wiring yet — just the structure.
5. **SUBSCRIBE / PUBLISH** wired through the hub end-to-end.
6. **(Stretch) RESP protocol** so the real `redis-cli` works against my server.
7. **(Stretch) Graceful shutdown** using `os/signal` and `context.Context`.

---

## Protocol

I'm starting with a simple line-based protocol for fast feedback:

```
SET foo bar
GET foo
DEL foo
SUBSCRIBE news
PUBLISH news hello
```

Testing flow: run my server, then `nc localhost 6379` and type commands. RESP is a stretch goal — I'll only tackle it once the architecture is solid, since swapping protocols is a parser change, not an architecture change.

---

## What I Already Know

- Basic Go syntax: types, functions, structs, methods, interfaces, packages
- I've read about goroutines and channels but haven't really used them in anger
- Comfortable with programming generally, networking concepts, and concurrency primitives in other languages

So you don't need to explain what a goroutine is at a one-sentence level, or how `func` works. But the deeper patterns — when to use buffered vs unbuffered channels, when `select` needs a `default` case, how to avoid goroutine leaks, when context cancellation belongs vs. doesn't — explain those when they come up. I'd rather hear "here's the principle and why it matters" than "here's the code, copy it."
