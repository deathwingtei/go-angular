# Angular Best Practice Folder Structure

While React is relatively unopinionated about folder structure, Angular has a more structured approach, especially for enterprise applications. The modern standard for Angular (particularly with **Standalone Components**) is the **Core / Shared / Features** architecture.

Here is how you can map your React knowledge to the best practice Angular structure:

## The React to Angular Mapping

- **React `pages/`** ➡️ **Angular `features/` or `pages/`**: These are your "smart" components that are tied to routes.
- **React `components/`** ➡️ **Angular `shared/components/` or `ui/`**: These are your "dumb" or presentational components reused across multiple pages.
- **React `routes.js`** ➡️ **Angular `app.routes.ts` + feature routes**: Routing is hierarchical in Angular.
- **React `hooks/` / `context/`** ➡️ **Angular `core/services/` or `shared/utils/`**: Reusable logic and state management are usually handled by Services and Signals in Angular.

---

## Recommended Folder Structure Template

Here is the ideal structure inside your `src/app/` folder for a scalable Angular application:

```text
src/
└── app/
    ├── core/                   # 1. Singleton services, interceptors, guards
    │   ├── auth/               # Auth guards, auth services
    │   ├── http/               # HTTP interceptors
    │   └── layout/             # Main layout components (Header, Sidebar, Footer)
    │
    ├── shared/                 # 2. Reusable UI components, pipes, directives
    │   ├── components/         # e.g., button, card, modal (dumb components)
    │   ├── directives/         # Custom directives
    │   ├── pipes/              # Custom pipes (data formatting)
    │   └── utils/              # Helper functions
    │
    ├── features/               # 3. Your "Pages" or Domains
    │   ├── dashboard/          # A feature module/domain
    │   │   ├── components/     # Components specific ONLY to the dashboard
    │   │   ├── services/       # State/data services for dashboard
    │   │   ├── dashboard.component.ts      # The main "Page" component
    │   │   └── dashboard.routes.ts         # Child routes for dashboard
    │   │
    │   └── user-profile/       # Another feature
    │       ├── components/
    │       ├── user-profile.component.ts
    │       └── user-profile.routes.ts
    │
    ├── app.component.ts        # The root component
    ├── app.config.ts           # Global app providers (replaces app.module.ts)
    └── app.routes.ts           # Main routing file (Lazy loading features here)
```

### Breakdown of the Folders:

#### 1. `core/` (App-wide singletons)

This folder contains things that should be instantiated exactly **once** in your application.

- **Examples**: Authentication services, HTTP interceptors (for adding JWT tokens), global error handlers, and shell layout components (like your main Navbar or Footer).
- **React Equivalent**: Your main App Context, global providers, and top-level layout wrappers.

#### 2. `shared/` (The UI Library)

This is where your presentational (dumb) components live. These components do not fetch data themselves; they only receive data via `@Input()` and emit events via `@Output()`.

- **Examples**: `<app-custom-button>`, `<app-data-table>`, custom date-formatting pipes.
- **React Equivalent**: The generic `components/` folder (like a custom UI library or design system).

#### 3. `features/` (The "Pages" and Domains)

Instead of a flat `pages/` directory, Angular groups by _feature_ or _domain_. A feature folder contains everything needed for that slice of the app. The root component of a feature (e.g., `dashboard.component.ts`) acts as your **Page**.

- **Why?** This makes **Lazy Loading** incredibly easy. In your `app.routes.ts`, you lazy load the entire feature route file.
- **React Equivalent**: The `pages/` folder, but grouped with their specific sub-components and logic.

### Routing Example (app.routes.ts)

In modern Angular (v15+ with standalone components), your main route file will lazy-load these features:

```typescript
import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: 'dashboard',
    // Lazy load the feature routes
    loadChildren: () =>
      import('./features/dashboard/dashboard.routes').then((m) => m.DASHBOARD_ROUTES),
  },
  {
    path: 'profile',
    loadChildren: () =>
      import('./features/user-profile/user-profile.routes').then((m) => m.PROFILE_ROUTES),
  },
  {
    path: '',
    redirectTo: 'dashboard',
    pathMatch: 'full',
  },
];
```

### Tips for Success

1. **Use Angular CLI**: Don't create these files manually. Use the CLI to generate them (e.g., `ng generate component shared/components/button`).
2. **Keep it Standalone**: Since Angular 15+, you no longer need `NgModule`. Generate all components, directives, and pipes as `standalone: true`.
3. **Smart vs. Dumb**: Keep your `features/` components "smart" (fetching data, handling state) and your `shared/` components "dumb" (just rendering UI).
