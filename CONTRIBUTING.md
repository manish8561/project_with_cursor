# Contributing to Identity Hub Services

Thanks for your interest in contributing to the project. This repository contains both an Angular frontend and a Go-based microservices backend. We welcome contributions from developers working on either side of the stack.

## Project Structure

- `frontend/` — Angular application
- `backend/` — Go microservices and shared backend tooling
- `deploy/` — Docker and deployment configuration
- `README.md` — project overview and local setup guidance

## Ways to Contribute

- Report bugs and suggest improvements by opening an issue
- Fix documentation gaps or clarify setup instructions
- Improve frontend functionality, styling, or UX
- Add or update backend services, APIs, or shared logic
- Help with tests, security checks, and deployment reliability

## Before You Start

1. Fork the repository and create a branch for your work.
2. Make sure you have the required tools installed:
   - Frontend: Node.js and npm
   - Backend: Go 1.26+
   - Docker and Docker Compose for local integration testing
3. Read the project README for local environment setup and service layout.

## Frontend Contributions

### Requirements

- Use the Angular code conventions already present in the project
- Keep UI changes consistent with the existing design and patterns
- Add or update tests when behavior changes

### Local setup

```bash
cd frontend
npm install
npm run start
```

### Validation

Run frontend checks as appropriate for the app, and verify the UI still works with the local backend/services.

## Backend Contributions

### Requirements

- Keep Go code formatted and idiomatic
- Prefer small, focused changes with clear responsibilities
- Add or update tests for bug fixes and new behavioral logic
- Ensure service-level changes do not break other microservices

### Local setup

```bash
go test ./...
```

For backend formatting and validation checks, the repository includes project-level commands:

```bash
make -C backend lint-format
make -C backend security-check
```

## Pull Request Guidelines

- Keep pull requests focused on one topic or bug fix
- Include a clear description of the changes and why they are needed
- Mention any affected services or frontend/backend areas
- Reference related issues when applicable
- Ensure the change has been validated locally before submitting

## Commit Guidance

- Use clear, descriptive commit messages
- Keep commits logically grouped
- Avoid unrelated refactors in feature or bug-fix PRs

## Code Review Expectations

All contributions are reviewed by maintainers and/or collaborators. We expect:

- Respectful and constructive feedback
- Focus on the code and the problem being solved
- Patience and collaboration while discussing changes

## Community Standards

Please follow the repository's [CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md) in all project interactions.

Thank you for helping improve Identity Hub Services.
