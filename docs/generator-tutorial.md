# Generator Tutorial

This guide walks through using the Kikplate generator to create a new project from a plate. It covers the full flow from understanding the plate format to running the generator and inspecting the output.

## What the Generator Does

The generator reads a `plate.yaml` manifest, validates and coerces the values you supply, renders each template file with those values, and writes the result as a directory on disk (or pushes it directly to a remote Git repository).

The difference between generation and scaffolding is intent. Scaffolding clones an existing repository as is. Generation renders a structured set of template files according to a schema you define, producing a tailored project every time.

## Plate Structure

A plate is a directory with this layout:

```
my-plate/
  plate.yaml
  values.yaml
  templates/
    ...template files...
```

`plate.yaml` is the manifest. It declares the schema, modules, and file list.

`values.yaml` is an optional defaults file you can pass with `-f`.

`templates/` contains all renderable files, each ending in `.tmpl`.

## The plate.yaml Manifest

A minimal manifest:

```yaml
name: my-plate

schema:
  projectName:
    type: string
    required: true
  port:
    type: string
    required: false
    default: "8080"

modules:
  docker:
    enabled: false

files:
  - path: main.go
    template: templates/main.go.tmpl

  - path: Dockerfile
    condition: modules.docker.enabled
    template: templates/Dockerfile.tmpl
```

### Schema Field Types

| Type | Notes |
|------|-------|
| string | Default type. Coerced from any scalar. |
| bool | Accepts true, false, yes, no, 1, 0. |
| int | Integer only. Rejects floats. |
| number | Any numeric value. |
| enum | Must match one of the declared `values`. |

### Modules

Modules are named feature toggles. They default to enabled or disabled in the manifest and can be overridden per generation. The condition field on a file entry controls whether that file is rendered.

```yaml
modules:
  docker:
    enabled: true
```

### File Entries

Each entry in `files` has:

`path` — The output file path. Supports template expressions so paths themselves can include values.

`template` — A relative path to the template file inside the plate, a raw inline template string, or a remote URL.

`condition` — An optional expression. The file is skipped when it evaluates to false.

## Template Syntax

Templates use Go text/template syntax with a set of helper functions.

### Variables

Every schema field is available as `.fieldName`:

```
module {{ .modulePath }}
port: {{ .port }}
```

### Conditional Blocks

```
{{ if eq .database "postgres" }}
url: jdbc:postgresql://{{ .postgresHost }}:{{ .postgresPort }}/{{ .postgresDatabase }}
{{ else }}
url: jdbc:h2:mem:{{ .artifactId }}
{{ end }}
```

### Module Checks in Templates

```
{{ if .modules.docker.enabled }}
image: {{ .projectName | slugify }}
{{ end }}
```

### Helper Functions

| Function | Example |
|----------|---------|
| `lower` | `{{ .projectName \| lower }}` |
| `upper` | `{{ .projectName \| upper }}` |
| `trim` | `{{ .description \| trim }}` |
| `replace` | `{{ replace .packageName "." "/" }}` |
| `slugify` | `{{ .projectName \| slugify }}` |
| `default` | `{{ .port \| default "8080" }}` |

### Dynamic File Paths

The `path` field in the manifest is also a template, so the output path can include values:

```yaml
- path: src/main/java/{{ replace .packageName "." "/" }}/Application.java
  template: templates/Application.java.tmpl
```

## Example Plates

The `examples/` directory contains two ready-to-run plates.

### go-http-server

A minimal Go HTTP server.

Location: `examples/go-http-server/`

Features demonstrated:

Schema with string and enum fields. Docker module toggle. Dynamic file paths. `lower` and `default` template helpers.

Generate it:

```bash
kikplate generate --template examples/go-http-server \
  --set projectName=my-api \
  --set modulePath=github.com/you/my-api \
  --set port=9090 \
  --output-dir my-api
```

Or supply all values from the included values file:

```bash
kikplate generate --template examples/go-http-server \
  -f examples/go-http-server/values.yaml \
  --output-dir my-api
```

Enabling or disabling Docker at generation time:

```bash
kikplate generate --template examples/go-http-server \
  -f examples/go-http-server/values.yaml \
  --set modules.docker.enabled=false \
  --output-dir my-api
```

When `modules.docker.enabled` is false the `Dockerfile` is not generated at all.

### spring-boot-rest-api

A production-ready Spring Boot 3 REST API with CRUD endpoints, validation, error handling, and optional modules.

Location: `examples/spring-boot-rest-api/`

Features demonstrated:

All scalar schema types including int and enum. Four independently toggled modules: `docker`, `swagger`, `actuator`, `seedData`. Dynamic Java package paths using `replace`. Template variables embedded inside Java source code. Test templates rendered with the same values as production code. Conditional blocks inside `pom.xml` and `application.yml` switching between Postgres and H2.

Generate it with all modules on and Postgres:

```bash
kikplate generate --template examples/spring-boot-rest-api \
  -f examples/spring-boot-rest-api/values.yaml \
  --output-dir task-service
```

Generate a lightweight H2 development version with no Docker or seed data:

```bash
kikplate generate --template examples/spring-boot-rest-api \
  -f examples/spring-boot-rest-api/values.yaml \
  --set database=h2 \
  --set modules.docker.enabled=false \
  --set modules.seedData.enabled=false \
  --output-dir task-service-dev
```

When `modules.swagger.enabled` is false, `OpenApiConfig.java` is not generated and the springdoc dependency is not added to `pom.xml`. When `modules.actuator.enabled` is false, the management endpoint block is omitted from `application.yml`. When `modules.seedData.enabled` is false, `SeedDataLoader.java` is not generated.

## Generate to a Remote Repository

After generating locally you can inspect the output, then push it to a remote:

```bash
kikplate generate --template examples/spring-boot-rest-api \
  -f examples/spring-boot-rest-api/values.yaml \
  --repo https://github.com/you/task-service.git
```

This generates the project in a temporary directory, creates an initial Git commit, and pushes to the remote.

## Condition Expression Reference

Conditions on file entries support a small expression language.

Simple truthy lookup:

```yaml
condition: modules.docker.enabled
```

Equality check:

```yaml
condition: database == postgres
```

Negation:

```yaml
condition: "!modules.docker.enabled"
```

Compound expression:

```yaml
condition: modules.auth.enabled && database != none
```

## Interactive Mode

When you run `generate` without `-f` or `--set`, the CLI prompts you for each required and optional field defined in the schema:

```bash
kikplate generate myorg/spring-boot-rest-api
```

Output:

```
Generating: spring-boot-rest-api

projectName [required]: task-service
groupId (default: com.example): com.acme
artifactId (default: spring-api): task-service
packageName (default: com.example.springapi): com.acme.taskservice
javaVersion (enum: 17|21) (default: 21):
serverPort (int) (default: 8080): 8082
database (enum: postgres|h2) (default: h2): postgres
...
Enable module "docker"? [y/n] (default: y):
Enable module "swagger"? [y/n] (default: y):
```

## Writing Your Own Plate

1. Create a directory for the plate.
2. Add a `plate.yaml` that declares the schema, modules, and files list.
3. Add template files under a `templates/` subdirectory.
4. Add a `values.yaml` with sensible defaults for quick testing.
5. Run the generator locally to verify the output.

```bash
kikplate generate --template my-plate -f my-plate/values.yaml --output-dir test-output
```
