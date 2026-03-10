import { MenuView } from "@/components/features/menu/menu-view";
import { listCategories, listItems } from "@/lib/api/menu";

const TENANT_ID = process.env.NEXT_PUBLIC_TENANT_ID ?? "";

async function getData() {
  if (!TENANT_ID) return { categories: [], items: [] };
  try {
    const [categories, items] = await Promise.all([
      listCategories(TENANT_ID),
      listItems(TENANT_ID),
    ]);
    return { categories, items };
  } catch {
    return { categories: [], items: [] };
  }
}

export default async function MenuPage() {
  if (!TENANT_ID) {
    return (
      <div className="flex flex-col gap-4">
        <h2 className="text-2xl font-bold">Menu</h2>
        <p className="text-muted-foreground text-sm">
          Set <code className="font-mono bg-muted px-1 rounded">NEXT_PUBLIC_TENANT_ID</code> to view menu data.
        </p>
      </div>
    );
  }

  const { categories, items } = await getData();

  return (
    <div className="flex flex-col gap-8">
      <h2 className="text-2xl font-bold">Menu</h2>
      <MenuView initialCategories={categories} initialItems={items} tenantId={TENANT_ID} />
    </div>
  );
}
