# IOrder Frontend

Next.js 15 web dashboard for the IOrder restaurant ordering platform.

## Stack

- **Next.js 15** - React framework with App Router
- **TypeScript** - Type-safe JavaScript
- **Tailwind CSS v4** - Utility-first CSS framework
- **pnpm** - Fast, disk space efficient package manager

## Project Structure

```
frontend/
├── src/
│   ├── app/                      # Next.js App Router pages
│   │   ├── page.tsx              # Home page
│   │   ├── tenants/              # Tenant management
│   │   ├── menu/                 # Menu management
│   │   ├── orders/               # Orders overview
│   │   ├── kitchen/              # Kitchen display system
│   │   ├── floor-plan/           # Floor plan & tables
│   │   └── staff/                # Staff management
│   ├── components/
│   │   ├── ui/                   # Reusable UI components
│   │   └── features/             # Feature-specific components
│   ├── lib/
│   │   ├── api/                  # API client functions
│   │   └── format.ts             # Utility functions
│   └── types/                    # TypeScript type definitions
├── package.json
├── tsconfig.json
└── tailwindcss.config.ts
```

## Quick Start

### Prerequisites

- Node.js 20+
- pnpm 9+
- Backend API running on `http://localhost:8080`

### Setup

1. **Navigate to frontend directory:**
   ```bash
   cd frontend
   ```

2. **Install dependencies:**
   ```bash
   pnpm install
   ```

3. **Create environment file:**
   ```bash
   cp .env.example .env.local
   ```

4. **Configure environment variables:**
   ```bash
   # .env.local
   NEXT_PUBLIC_API_URL=http://localhost:8080
   NEXT_PUBLIC_TENANT_ID=<your-tenant-uuid>
   ```

   ⚠️ **Important:** Most features require `NEXT_PUBLIC_TENANT_ID` to be set. Get a tenant ID by:
   - Running the backend and accessing `/api/v1/tenants`
   - Creating a tenant through the Tenants page
   - Using a pre-seeded tenant from your database

5. **Start development server:**
   ```bash
   pnpm dev
   ```

   Open [http://localhost:3000](http://localhost:3000)

## Pages

### Home (`/`)
- Platform overview
- System health check
- Feature cards

### Tenants (`/tenants`)
- **View:** List all restaurant tenants with status
- **Create:** Add new tenant with configuration
- **Edit:** Update tenant details (name, currency, timezone)
- **Multi-tenancy:** Platform-level, no tenant ID required

### Menu (`/menu`)
- **Categories:** Create, edit, delete menu categories (food/beverage)
- **Items:** Full CRUD for menu items with:
  - Pricing (in smallest currency unit)
  - Images, descriptions, tags
  - Availability toggle
  - Category assignment

### Orders (`/orders`)
- **List:** View all orders grouped by status
- **Statuses:** DRAFT → SUBMITTED → CONFIRMED → IN_PREPARATION → READY → SERVED → PAYMENT_REQUESTED → PAID → CLOSED
- **Details:** Order ID, table, total, timestamps, notes

### Kitchen Display (`/kitchen`)
- **Lanes:** Submitted, In Preparation, Ready
- **Real-time:** Order queue for kitchen staff
- **Minimal UI:** Focus on order IDs and timestamps

### Floor Plan (`/floor-plan`)
- **View:** Floor plans with tables
- **Tables:** Name, capacity, status (available/occupied/reserved)
- **Visual:** Position-based layout (future: drag-drop)

### Staff (`/staff`)
- **List:** Staff members by tenant
- **Roles:** Owner, Manager, Waiter, Kitchen, Bartender, Cashier, etc.
- **Status:** Active/inactive toggle

## Environment Variables

```bash
# Required
NEXT_PUBLIC_API_URL=http://localhost:8080

# Optional (but most features need it)
NEXT_PUBLIC_TENANT_ID=<uuid>
```

## Development Workflow

### Run Development Server
```bash
pnpm dev
```

### Type Check
```bash
pnpm typecheck
```

### Lint Code
```bash
pnpm lint
```

### Build for Production
```bash
pnpm build
pnpm start
```

## API Client

The API client (`src/lib/api/client.ts`) handles:
- ✅ Automatic JSON serialization
- ✅ Error handling with typed `APIError`
- ✅ Tenant ID headers (`X-Tenant-ID`)
- ✅ TypeScript type safety

### Example Usage

```typescript
import { listTenants, createTenant } from "@/lib/api/tenants";

// List all tenants
const tenants = await listTenants();

// Create a tenant
const newTenant = await createTenant({
  slug: "my-restaurant",
  name: "My Restaurant",
  currency: "VND",
  timezone: "Asia/Ho_Chi_Minh",
});
```

## UI Components

All UI components are in `src/components/ui/`:

- **Button** - Primary, outline, ghost variants
- **Badge** - Status indicators with color variants
- **Card** - Content containers with header/footer
- **Dialog** - Modal dialogs
- **Input** - Text input with label support
- **Select** - Dropdown select
- **Textarea** - Multi-line text input
- **Table** - Data table with header/body/rows
- **Label** - Form labels

### Example Component

```tsx
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

export function MyComponent() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Hello World</CardTitle>
      </CardHeader>
      <CardContent>
        <Button onClick={() => alert("Clicked!")}>Click Me</Button>
      </CardContent>
    </Card>
  );
}
```

## Type Definitions

All API types are defined in `src/types/`:
- `tenant.ts` - Tenant model
- `menu.ts` - MenuCategory, MenuItem
- `order.ts` - Order, OrderStatus
- `staff.ts` - Staff model
- `table.ts` - FloorPlan, RestaurantTable
- `health.ts` - Health check response
- `api.ts` - APIError, APIErrorResponse

## Code Conventions

### File Naming
- **Pages:** `page.tsx` (Next.js convention)
- **Components:** `kebab-case.tsx`
- **Types:** `kebab-case.ts`

### Component Pattern
```tsx
// Server Component (default)
export default async function PageName() {
  const data = await fetchData();
  return <View data={data} />;
}

// Client Component (when needed)
"use client";

export function InteractiveComponent() {
  const [state, setState] = useState();
  return ...;
}
```

### API Error Handling
```typescript
try {
  const result = await apiClient.get("/endpoint");
} catch (err) {
  if (err instanceof APIError) {
    console.error(err.code, err.message);
  }
}
```

## Troubleshooting

### "No orders/menu/staff found"
- ✅ Check `NEXT_PUBLIC_TENANT_ID` is set in `.env.local`
- ✅ Verify tenant exists in backend database
- ✅ Confirm backend is running and accessible

### API Connection Errors
- ✅ Check `NEXT_PUBLIC_API_URL` points to backend
- ✅ Verify backend is running (`http://localhost:8080/health`)
- ✅ Check CORS configuration in backend

### Build Errors
- Run `pnpm typecheck` to find TypeScript errors
- Run `pnpm lint` to find linting issues
- Clear `.next` folder and rebuild

### Styles Not Loading
- Tailwind v4 uses `@tailwindcss/postcss` - check `postcss.config.mjs`
- Ensure `globals.css` imports Tailwind directives

## Production Deployment

### Build
```bash
pnpm build
```

### Environment Variables
Set in production environment:
```
NEXT_PUBLIC_API_URL=https://api.yourproduction.com
NEXT_PUBLIC_TENANT_ID=<optional-default-tenant>
```

### Start Production Server
```bash
pnpm start
```

Or deploy to Vercel/Netlify/AWS Amplify - they auto-detect Next.js.

## Future Enhancements

- [ ] Real-time order updates (WebSockets/Server-Sent Events)
- [ ] Drag-and-drop floor plan editor
- [ ] Order item details modal
- [ ] Customer-facing mobile ordering UI
- [ ] Analytics dashboard
- [ ] Dark mode toggle
- [ ] Multi-language support (i18n)
- [ ] Image upload for menu items
- [ ] Payment integration UI
- [ ] Staff permissions & authentication

## Contributing

1. Create feature branch: `feat/<issue>-description`
2. Write TypeScript with strict types
3. Run `pnpm typecheck && pnpm lint && pnpm build`
4. Open PR with conventional commit message

## Additional Documentation

- [Backend README](../backend/README.md)
- [Architecture Overview](../docs/architecture/system-architecture.md)
- [Working Agreement](../docs/process/working-agreement.md)
- [Frontend Agent Instructions](AGENTS.md)

## License

Proprietary - All rights reserved
