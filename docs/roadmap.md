# Kikplate

> Kikplate is a declarative application architecture platform for building, sharing, and generating production-ready software systems.

Kikplate is not a template gallery.

It is an ecosystem for defining, distributing, and generating reusable application architectures called plates.

Each plate represents a versioned, configurable blueprint for building software systems.

---

# Vision

Modern software development repeatedly rebuilds the same foundations:

- authentication systems
- database setups
- project structure
- deployment pipelines
- backend/frontend architectures

Kikplate solves this by introducing a standardized architecture reuse layer for software systems.

Instead of copying boilerplates or relying on AI-generated code, developers use structured, versioned architecture definitions.

---

# Core System

Kikplate is built on five layers:

---

## 1. Plates (Architecture Definitions)

Plates are reusable application blueprints.

A plate defines:

- project structure
- stack (frameworks, runtime, DB)
- optional modules (auth, payments, admin, etc.)
- configuration schema
- generation rules

A plate is NOT just a repository.

It is a structured architecture definition.

---

## 2. Generator Engine (Core Runtime)

Transforms plates into real projects.

Plate + Configuration → Generated Project

Responsibilities:

- template rendering
- conditional generation
- feature toggling
- dependency resolution
- project scaffolding

---

## 3. Registry (Ecosystem Layer)

The registry manages distribution and ecosystem intelligence.

Responsibilities:

- plate discovery
- versioning
- metadata storage
- search and ranking
- health monitoring
- trust scoring
- validation
- community submissions
- private registries

---

## 4. CLI (Developer Interface)

The CLI is the primary developer interface.

Commands:

- search plates
- create projects
- update projects
- add modules
- inspect ecosystem state

Example:

kik search saas  
kik create my-app --plate saas-starter  
kik add auth  
kik update  
kik doctor  

---

## 5. AI Layer (Intelligence System)

AI assists but does not replace the system.

It:

- understands user intent
- selects plates
- generates configuration (values)
- composes architectures

Flow:

User → AI → Configuration → Kikplate → Generated Project

---

# Kikplate Architecture Layers

```text
                     ┌────────────────────┐
                     │      AI Layer      │
                     │  (optional brain)  │
                     │ intent → config    │
                     └─────────┬──────────┘
                               │
                               ▼
                     ┌────────────────────┐
                     │     CLI Layer      │
                     │ user interface     │
                     │ commands & UX      │
                     └─────────┬──────────┘
                               │
                               ▼
                     ┌────────────────────┐
                     │ Generator Engine   │
                     │ core runtime       │
                     │ plate execution    │
                     └─────────┬──────────┘
                               │
        ┌──────────────────────┴──────────────────────┐
        │                                             │
        ▼                                             ▼
┌────────────────────┐                   ┌────────────────────┐
│      Plates        │                   │     Registry       │
│ architecture defs  │                   │ ecosystem layer    │
│ templates + schema │                   │ versions + trust   │
└────────────────────┘                   └────────────────────┘
```

# Key Principles

## Architecture First
Kikplate focuses on full application architectures, not snippets.

## Declarative Design
Projects are defined through configuration.

## Reproducibility
Every output is deterministic and versioned.

## Ecosystem Trust
Registry ensures quality and health of plates.

## AI-Native but Structured
AI enhances Kikplate but does not replace deterministic generation.

---

# Roadmap

---

## Phase 1 — Foundation

- CLI stabilization
- registry hardening
- plate schema definition
- GitHub-based ingestion
- metadata standardization
- install reliability improvements

---

## Phase 2 — Ecosystem Structure

- versioning system
- health monitoring
- trust scoring
- dependency validation
- search improvements
- lifecycle states (active, deprecated, archived)
- community submission pipeline

---

## Phase 3 — Developer Workflow

- improved CLI UX
- doctor / update / search commands
- smarter scaffolding flow
- feature toggles during creation
- local caching

---

## Phase 4 — Generation System

Introduce configuration-based generation:

- values.yaml support
- schema-driven plates
- conditional file rendering
- modular architecture composition

Example:

auth:
  enabled: true
database:
  type: postgres
payments:
  stripe: true

---

## Phase 5 — AI Layer

- natural language → architecture
- AI-generated configs
- plate recommendation engine
- multi-plate composition assistance

---

## Phase 6 — Platform Expansion

- private registries
- enterprise workspaces
- registry federation
- governance and permissions
- API + SDK ecosystem
- CI validation system

---

# Open Source vs Premium

## Free (Core)
- CLI
- public registry
- generator engine
- plate standard
- public plates

## Premium
- private registries
- AI assistance layer
- advanced analytics
- enterprise governance
- federation
- advanced generation system

---

# Final Positioning

Kikplate is a declarative application architecture platform for defining, sharing, and generating production-ready software systems through reusable, versioned, configurable blueprints called plates.

It is not a template system.

It is an architecture generation ecosystem.