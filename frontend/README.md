# Frontend Application

This is the frontend application built with Angular.

## Prerequisites

- Node.js (v14 or later)
- npm (v6 or later)
- Angular CLI (v15 or later)

## Installation

1. Install dependencies:
```bash
npm install
```

## Development

To run the application in development mode:

```bash
ng serve
```

This will:
- Start the development server
- Use development environment configuration
- Open the application at `http://localhost:4200`
- Enable hot-reloading for development

## Production Build

To build the application for production:

```bash
ng build --configuration=production
```

This will:
- Create a production build in the `dist/frontend` directory
- Use production environment configuration
- Optimize the build for production
- Enable output hashing for cache busting

## Environment Configuration

The application uses different environment configurations for development and production:

- Development: `src/environments/environment.ts`
- Production: `src/environments/environment.prod.ts`

Make sure to update the `apiUrl` in both files according to your backend server configuration.

## Project Structure

- `src/app/` - Main application code
  - `services/` - Services including authentication
  - `components/` - Reusable components
  - `guards/` - Route guards
  - `environments/` - Environment configurations

## Features

- User authentication (login/register)
- Protected routes
- Responsive design
- Form validation
- Error handling
