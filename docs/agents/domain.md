# Domain Docs

How the engineering skills should consume this repo's domain documentation when exploring the codebase.

## Layout

This repo uses a single-context layout:

- Root `CONTEXT.md`
- Root `docs/adr/`

## Before Exploring, Read These If Present

- `CONTEXT.md` at the repo root.
- ADRs under root `docs/adr/` that are relevant to the work.

If these files or directories do not exist, proceed silently. Do not flag their absence or suggest creating them upfront.

`CONTEXT.md` and ADRs are created lazily by domain and documentation skills when terms or decisions actually get resolved. Do not create them during this setup.

## Use The Glossary's Vocabulary

When your output names a domain concept, use the term as defined in `CONTEXT.md`. Do not drift to synonyms the glossary explicitly avoids.

If the concept you need is not in the glossary yet, either reconsider whether you are inventing language the project does not use, or note the gap for domain modeling.

## Flag ADR Conflicts

If your output contradicts an existing ADR, surface it explicitly rather than silently overriding it.
