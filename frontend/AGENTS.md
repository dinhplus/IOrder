# Frontend — Code Instructions (GitHub Copilot / Claude)

## Stack
- **Next.js 15** (App Router)
- **TypeScript** (strict mode)
- **Tailwind CSS v4**
- **pnpm** package manager
- **ESLint 9** (next/core-web-vitals)

## Project Layout

```
frontend/
├── src/
│   ├── app/                    # Next.js App Router
│   │   ├── layout.tsx          # Root layout (HTML, header, global CSS)
│   │   ├── page.tsx            # Home page
│   │   └── globals.css         # @import "tailwindcss"
│   ├── components/
│   │   ├── ui/                 # Base UI components (Button, Input, etc.)
│   │   └── features/           # Feature-specific components
│   ├── lib/
│   │   └── api/
│   │       └── client.ts       # Typed fetch wrapper (apiClient)
│   └── types/
│       └── api.ts              # APIError class + response interfaces
├── public/                     # Static assets
├── next.config.ts
├── tsconfig.json               # strict: true, @/* path alias
├── package.json
├── postcss.config.mjs          # @tailwindcss/postcss
└── .env.example
```

## Rules

### API Client Usage
```typescript
import { apiClient } from "@/lib/api/client";
import { APIError } from "@/types/api";

try {
  const data = await apiClient.get<{ status: string }>("/health");
} catch (err) {
  if (err instanceof APIError) {
    console.error(err.code, err.message);
  }
}
```

### TypeScript
- Always use explicit return types on functions
- Use `interface` for object shapes, `type` for unions/intersections
- No `any` — use `unknown` and type guards instead
- Path alias: `@/` maps to `src/`

### Components
```tsx
// Prefer named exports
export function MyComponent({ title }: { title: string }) {
  return <h1 className="text-xl font-bold">{title}</h1>;
}
```

### Tailwind CSS
- Use utility classes directly — no separate CSS files unless global
- Responsive: `sm:`, `md:`, `lg:` prefixes

### Environment Variables
```typescript
// Client-side (available in browser)
const apiUrl = process.env.NEXT_PUBLIC_API_URL;

// Server-side only (never use NEXT_PUBLIC_ for secrets)
const secret = process.env.API_SECRET;
```

### Error Boundaries
For server components, use `error.tsx` files in the route directory.
For client components, wrap with `<ErrorBoundary>`.

## Quality Gates (before PR)
```bash
cd frontend
pnpm run typecheck    # tsc --noEmit — must pass (0 errors)
pnpm run lint         # eslint — must pass (0 warnings with --max-warnings=0)
pnpm run build        # next build — must succeed
```
