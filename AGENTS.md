# Project Agent Rules (AGENTS.md)

## 1. Core Philosophy

- **Keep It Simple & No Technical Debt**: Code must be straightforward, easy to
  use, and maintainable. Avoid complex design patterns or unmaintainable
  shortcuts.

## 2. File Architecture Constraints

- **Separation**: UI logic (`renderScreen`) belongs in `ui.go`. Core logic
  belongs in `engine.go`. Do not pollute `main.go` with business logic.
- **Testing**: All moved logic must be covered in respective `_test.go` files.

## 3. API Constraints

- **Native Only**: You must strictly use the native Go standard library APIs
  only.

## 4. Formatting

- **Line Length**: Every single line of code, comments, and markdown MUST be
  strictly under 80 characters.
