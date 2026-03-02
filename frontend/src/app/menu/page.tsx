import { listCategories, listItems } from "@/lib/api/menu";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { MenuCategory, MenuItem } from "@/types/menu";

const TENANT_ID = process.env.NEXT_PUBLIC_TENANT_ID ?? "";

async function getData(): Promise<{ categories: MenuCategory[]; items: MenuItem[] }> {
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

export default async function MenuPage(): Promise<React.JSX.Element> {
  const { categories, items } = await getData();

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

  return (
    <div className="flex flex-col gap-8">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">Menu</h2>
        <span className="text-sm text-muted-foreground">{items.length} items · {categories.length} categories</span>
      </div>

      {categories.length === 0 ? (
        <p className="text-muted-foreground">No menu categories found.</p>
      ) : (
        categories.map((cat) => {
          const catItems = items.filter((item) => item.category_id === cat.id);
          return (
            <section key={cat.id} className="flex flex-col gap-4">
              <div className="flex items-center gap-3">
                <h3 className="text-lg font-semibold">{cat.name}</h3>
                <Badge variant={cat.is_active ? "success" : "outline"}>{cat.is_active ? "Active" : "Inactive"}</Badge>
                <span className="text-xs text-muted-foreground capitalize">{cat.type}</span>
              </div>
              {catItems.length === 0 ? (
                <p className="text-sm text-muted-foreground pl-4">No items in this category.</p>
              ) : (
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                  {catItems.map((item) => (
                    <Card key={item.id}>
                      <CardHeader>
                        <div className="flex items-start justify-between gap-2">
                          <CardTitle className="text-base">{item.name}</CardTitle>
                          <Badge variant={item.is_available ? "success" : "danger"}>
                            {item.is_available ? "Available" : "Unavailable"}
                          </Badge>
                        </div>
                        {item.description && (
                          <p className="text-sm text-muted-foreground">{item.description}</p>
                        )}
                      </CardHeader>
                      <CardContent>
                        <p className="text-sm font-semibold text-primary">
                          {item.price.toLocaleString("vi-VN")}₫
                        </p>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              )}
            </section>
          );
        })
      )}
    </div>
  );
}
