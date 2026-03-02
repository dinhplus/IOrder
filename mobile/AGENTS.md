# Mobile — Code Instructions (GitHub Copilot / Claude)

## Stack
- **React Native** + **Expo SDK 52**
- **expo-router v4** (file-based routing)
- **TypeScript** (strict mode)
- **pnpm** package manager
- **Jest** + **jest-expo** for testing

## Project Layout

```
mobile/
├── app/                        # expo-router pages (file-based)
│   ├── _layout.tsx             # Root Stack navigator
│   └── index.tsx               # Home screen
├── src/
│   ├── screens/                # Full-page screen components
│   ├── components/
│   │   ├── ui/                 # Base UI components
│   │   └── features/           # Feature-specific components
│   ├── lib/
│   │   └── api/
│   │       └── client.ts       # Typed fetch wrapper (apiClient)
│   ├── types/
│   │   └── api.ts              # APIError class + response interfaces
│   ├── hooks/                  # Custom React hooks
│   └── styles/                 # Theme and shared StyleSheet values
├── __tests__/                  # Jest test files
├── app.json                    # Expo config (bundle IDs, plugins)
├── babel.config.js
├── tsconfig.json               # strict: true, @/* → src/*
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

### API URL Configuration
Set the API URL in `app.json` under `expo.extra`:
```json
{
  "expo": {
    "extra": {
      "apiUrl": "http://localhost:8080"
    }
  }
}
```

### TypeScript
- Strict mode — no implicit `any`
- Path alias: `@/` maps to `src/`
- Shared types with frontend where possible

### Screens
```tsx
import { View, Text, StyleSheet } from "react-native";

export default function ExampleScreen() {
  return (
    <View style={styles.container}>
      <Text style={styles.title}>Title</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, alignItems: "center", justifyContent: "center" },
  title: { fontSize: 24, fontWeight: "bold" },
});
```

### Navigation (expo-router)
```typescript
import { router } from "expo-router";

// Navigate
router.push("/screens/detail");
router.replace("/");

// Link component
import { Link } from "expo-router";
<Link href="/screens/detail">Go to Detail</Link>
```

### Cross-Platform
- Test on both iOS and Android simulators
- Use `Platform.OS` only when behavior truly differs
- Prefer cross-platform components over platform-specific

## Quality Gates (before PR)
```bash
cd mobile
pnpm run typecheck    # tsc --noEmit — must pass (0 errors)
pnpm run lint         # eslint — must pass
pnpm test             # jest --passWithNoTests — must pass
```
