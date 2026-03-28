# lf-explore-skill

## What It Does
The `/lf:explore` skill is a thinking partner mode that helps users explore ideas, investigate problems, and clarify requirements before or during a change. It reads the codebase and existing change artifacts but never implements code. When insights crystallize, it offers a natural transition to `/lf:formalize`.

## Requirements

### Requirement: Explore skill is a thinking partner
The system SHALL provide an embedded `/lf:explore` skill that enters a thinking partner mode focused on exploration, not implementation.

#### Scenario: Explore prevents implementation
- **WHEN** the user is in explore mode and asks to implement code
- **THEN** the skill reminds the user to exit explore mode first (e.g., by formalizing a change with `/lf:formalize`)

#### Scenario: Explore reads existing changes
- **WHEN** the user invokes `/lf:explore` and active changes exist under `.littlefactory/changes/`
- **THEN** the skill reads available artifacts (proposal.md, design.md, tasks.json, specs) to ground the conversation in existing context

#### Scenario: Explore offers transition to formalize
- **WHEN** insights crystallize during exploration and the conversation reaches a natural decision point
- **THEN** the skill offers to transition to `/lf:formalize` to capture the discussion as a structured change

### Requirement: Explore skill uses visual communication
The system SHALL have the explore skill use ASCII diagrams and visual representations to communicate architecture, data flows, and comparisons.

#### Scenario: Explore visualizes architecture
- **WHEN** the conversation involves system architecture, data flows, or option comparisons
- **THEN** the skill uses ASCII diagrams, tables, or visual representations to clarify the discussion

## Boundaries

## Gotchas
