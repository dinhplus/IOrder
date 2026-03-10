"use client";

import { useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select } from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import {
  createCategory,
  updateCategory,
  deleteCategory,
  createMenuItem,
  updateMenuItem,
  deleteMenuItem,
  toggleItemAvailability,
  type CreateCategoryRequest,
  type UpdateCategoryRequest,
  type CreateMenuItemRequest,
  type UpdateMenuItemRequest,
} from "@/lib/api/menu";
import { formatCurrency } from "@/lib/format";
import type { MenuCategory, MenuItem } from "@/types/menu";

interface MenuViewProps {
  initialCategories: MenuCategory[];
  initialItems: MenuItem[];
  tenantId: string;
}

export function MenuView({ initialCategories, initialItems, tenantId }: MenuViewProps) {
  const [categories, setCategories] = useState<MenuCategory[]>(initialCategories);
  const [items, setItems] = useState<MenuItem[]>(initialItems);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [dialogType, setDialogType] = useState<"category" | "item">("category");
  const [editingCategory, setEditingCategory] = useState<MenuCategory | null>(null);
  const [editingItem, setEditingItem] = useState<MenuItem | null>(null);
  const [confirmDialog, setConfirmDialog] = useState<{
    open: boolean;
    title: string;
    message: string;
    onConfirm: () => void;
  }>({ open: false, title: "", message: "", onConfirm: () => {} });

  const [categoryForm, setCategoryForm] = useState({
    name: "",
    type: "food",
    description: "",
    sort_order: 0,
    is_active: true,
  });

  const [itemForm, setItemForm] = useState({
    category_id: "",
    name: "",
    description: "",
    price: 0,
    image_url: "",
    is_available: true,
    sort_order: 0,
  });

  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [actionError, setActionError] = useState("");

  const handleCreateCategory = () => {
    setDialogType("category");
    setEditingCategory(null);
    setCategoryForm({
      name: "",
      type: "food",
      description: "",
      sort_order: 0,
      is_active: true,
    });
    setError("");
    setDialogOpen(true);
  };

  const handleEditCategory = (cat: MenuCategory) => {
    setDialogType("category");
    setEditingCategory(cat);
    setCategoryForm({
      name: cat.name,
      type: cat.type,
      description: cat.description ?? "",
      sort_order: cat.sort_order,
      is_active: cat.is_active,
    });
    setError("");
    setDialogOpen(true);
  };

  const handleDeleteCategory = (id: string) => {
    const catItems = items.filter((item) => item.category_id === id);
    if (catItems.length > 0) {
      setConfirmDialog({
        open: true,
        title: "Cannot Delete Category",
        message: `This category has ${catItems.length} item(s). Please move or delete them before removing the category.`,
        onConfirm: () => setConfirmDialog((prev) => ({ ...prev, open: false })),
      });
      return;
    }
    setConfirmDialog({
      open: true,
      title: "Delete Category",
      message: "Are you sure you want to delete this category?",
      onConfirm: async () => {
        try {
          await deleteCategory(id, tenantId);
          setCategories((prev) => prev.filter((c) => c.id !== id));
          setActionError("");
        } catch (err: unknown) {
          setActionError((err as Error)?.message || "Failed to delete category");
        }
        setConfirmDialog((prev) => ({ ...prev, open: false }));
      },
    });
  };

  const handleCreateItem = (cat?: MenuCategory) => {
    setDialogType("item");
    setEditingItem(null);
    setItemForm({
      category_id: cat?.id ?? categories[0]?.id ?? "",
      name: "",
      description: "",
      price: 0,
      image_url: "",
      is_available: true,
      sort_order: 0,
    });
    setError("");
    setDialogOpen(true);
  };

  const handleEditItem = (item: MenuItem) => {
    setDialogType("item");
    setEditingItem(item);
    setItemForm({
      category_id: item.category_id,
      name: item.name,
      description: item.description ?? "",
      price: item.price,
      image_url: item.image_url ?? "",
      is_available: item.is_available,
      sort_order: item.sort_order,
    });
    setError("");
    setDialogOpen(true);
  };

  const handleDeleteItem = (id: string) => {
    setConfirmDialog({
      open: true,
      title: "Delete Menu Item",
      message: "Are you sure you want to delete this menu item?",
      onConfirm: async () => {
        try {
          await deleteMenuItem(id, tenantId);
          setItems((prev) => prev.filter((i) => i.id !== id));
          setActionError("");
        } catch (err: unknown) {
          setActionError((err as Error)?.message || "Failed to delete item");
        }
        setConfirmDialog((prev) => ({ ...prev, open: false }));
      },
    });
  };

  const handleToggleAvailability = async (item: MenuItem) => {
    try {
      await toggleItemAvailability(item.id, tenantId, !item.is_available);
      setItems((prev) =>
        prev.map((i) => (i.id === item.id ? { ...i, is_available: !i.is_available } : i))
      );
      setActionError("");
    } catch (err: unknown) {
      setActionError((err as Error)?.message || "Failed to toggle availability");
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError("");

    try {
      if (dialogType === "category") {
        if (editingCategory) {
          const updated: UpdateCategoryRequest = {
            name: categoryForm.name,
            type: categoryForm.type,
            description: categoryForm.description,
            sort_order: categoryForm.sort_order,
            is_active: categoryForm.is_active,
          };
          const result = await updateCategory(editingCategory.id, tenantId, updated);
          setCategories((prev) => prev.map((c) => (c.id === result.id ? result : c)));
        } else {
          const newCat: CreateCategoryRequest = {
            name: categoryForm.name,
            type: categoryForm.type,
            description: categoryForm.description || undefined,
            sort_order: categoryForm.sort_order,
          };
          const result = await createCategory(tenantId, newCat);
          setCategories((prev) => [...prev, result]);
        }
      } else {
        if (editingItem) {
          const updated: UpdateMenuItemRequest = {
            category_id: itemForm.category_id,
            name: itemForm.name,
            description: itemForm.description,
            price: itemForm.price,
            image_url: itemForm.image_url,
            is_available: itemForm.is_available,
            sort_order: itemForm.sort_order,
          };
          const result = await updateMenuItem(editingItem.id, tenantId, updated);
          setItems((prev) => prev.map((i) => (i.id === result.id ? result : i)));
        } else {
          const newItem: CreateMenuItemRequest = {
            category_id: itemForm.category_id,
            name: itemForm.name,
            description: itemForm.description || undefined,
            price: itemForm.price,
            image_url: itemForm.image_url || undefined,
            sort_order: itemForm.sort_order,
          };
          const result = await createMenuItem(tenantId, newItem);
          setItems((prev) => [...prev, result]);
        }
      }
      setDialogOpen(false);
    } catch (err: unknown) {
      setError((err as Error)?.message || "Failed to save");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <div className="flex items-center justify-between">
        <span className="text-sm text-muted-foreground">
          {items.length} items · {categories.length} categories
        </span>
        <div className="flex gap-2">
          <Button variant="outline" onClick={handleCreateCategory}>
            Add Category
          </Button>
          <Button onClick={() => handleCreateItem()}>Add Menu Item</Button>
        </div>
      </div>

      {actionError && (
        <div className="p-3 text-sm bg-danger/10 text-danger rounded-md flex items-center justify-between">
          <span>{actionError}</span>
          <button
            type="button"
            className="text-danger hover:text-danger/80 font-medium text-xs"
            onClick={() => setActionError("")}
          >
            Dismiss
          </button>
        </div>
      )}

      {categories.length === 0 ? (
        <p className="text-muted-foreground">No categories. Create one to get started.</p>
      ) : (
        categories.map((cat) => {
          const catItems = items.filter((item) => item.category_id === cat.id);
          return (
            <section key={cat.id} className="flex flex-col gap-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <h3 className="text-lg font-semibold">{cat.name}</h3>
                  <Badge variant={cat.is_active ? "success" : "outline"}>
                    {cat.is_active ? "Active" : "Inactive"}
                  </Badge>
                  <span className="text-xs text-muted-foreground capitalize">{cat.type}</span>
                </div>
                <div className="flex gap-2">
                  <Button variant="ghost" size="sm" onClick={() => handleCreateItem(cat)}>
                    + Item
                  </Button>
                  <Button variant="ghost" size="sm" onClick={() => handleEditCategory(cat)}>
                    Edit
                  </Button>
                  <Button variant="ghost" size="sm" onClick={() => handleDeleteCategory(cat.id)}>
                    Delete
                  </Button>
                </div>
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
                      <CardContent className="flex flex-col gap-2">
                        <p className="text-sm font-semibold text-primary">
                          {formatCurrency(item.price)}
                        </p>
                        <div className="flex gap-1 flex-wrap">
                          {item.tags.map((tag) => (
                            <span
                              key={tag}
                              className="text-xs px-2 py-0.5 rounded-full bg-muted text-muted-foreground"
                            >
                              {tag}
                            </span>
                          ))}
                        </div>
                        <div className="flex gap-2 mt-2">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleToggleAvailability(item)}
                          >
                            Toggle
                          </Button>
                          <Button variant="ghost" size="sm" onClick={() => handleEditItem(item)}>
                            Edit
                          </Button>
                          <Button variant="ghost" size="sm" onClick={() => handleDeleteItem(item.id)}>
                            Delete
                          </Button>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              )}
            </section>
          );
        })
      )}

      {/* Create/Edit Form Dialog */}
      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent>
          <form onSubmit={handleSubmit}>
            <DialogHeader>
              <DialogTitle>
                {dialogType === "category"
                  ? editingCategory
                    ? "Edit Category"
                    : "Create Category"
                  : editingItem
                  ? "Edit Menu Item"
                  : "Create Menu Item"}
              </DialogTitle>
            </DialogHeader>

            <div className="flex flex-col gap-4 p-6">
              {error && (
                <div className="p-3 text-sm bg-danger/10 text-danger rounded-md">{error}</div>
              )}

              {dialogType === "category" ? (
                <>
                  <div>
                    <Label htmlFor="cat-name">Name *</Label>
                    <Input
                      id="cat-name"
                      value={categoryForm.name}
                      onChange={(e) => setCategoryForm({ ...categoryForm, name: e.target.value })}
                      required
                    />
                  </div>
                  <div>
                    <Label htmlFor="cat-type">Type</Label>
                    <Select
                      id="cat-type"
                      value={categoryForm.type}
                      onChange={(e) => setCategoryForm({ ...categoryForm, type: e.target.value })}
                    >
                      <option value="food">Food</option>
                      <option value="beverage">Beverage</option>
                    </Select>
                  </div>
                  <div>
                    <Label htmlFor="cat-desc">Description</Label>
                    <Textarea
                      id="cat-desc"
                      value={categoryForm.description}
                      onChange={(e) =>
                        setCategoryForm({ ...categoryForm, description: e.target.value })
                      }
                    />
                  </div>
                  {editingCategory && (
                    <div className="flex items-center gap-2">
                      <input
                        type="checkbox"
                        id="cat-active"
                        checked={categoryForm.is_active}
                        onChange={(e) =>
                          setCategoryForm({ ...categoryForm, is_active: e.target.checked })
                        }
                        className="h-4 w-4 rounded border-border"
                      />
                      <Label htmlFor="cat-active" className="cursor-pointer">
                        Active
                      </Label>
                    </div>
                  )}
                </>
              ) : (
                <>
                  <div>
                    <Label htmlFor="item-category">Category *</Label>
                    <Select
                      id="item-category"
                      value={itemForm.category_id}
                      onChange={(e) => setItemForm({ ...itemForm, category_id: e.target.value })}
                      required
                    >
                      {categories.map((cat) => (
                        <option key={cat.id} value={cat.id}>
                          {cat.name}
                        </option>
                      ))}
                    </Select>
                  </div>
                  <div>
                    <Label htmlFor="item-name">Name *</Label>
                    <Input
                      id="item-name"
                      value={itemForm.name}
                      onChange={(e) => setItemForm({ ...itemForm, name: e.target.value })}
                      required
                    />
                  </div>
                  <div>
                    <Label htmlFor="item-desc">Description</Label>
                    <Textarea
                      id="item-desc"
                      value={itemForm.description}
                      onChange={(e) => setItemForm({ ...itemForm, description: e.target.value })}
                    />
                  </div>
                  <div>
                    <Label htmlFor="item-price">Price (smallest unit) *</Label>
                    <Input
                      id="item-price"
                      type="number"
                      value={itemForm.price}
                      onChange={(e) => setItemForm({ ...itemForm, price: Number(e.target.value) })}
                      required
                      min="0"
                    />
                    <p className="text-xs text-muted-foreground mt-1">
                      e.g., 50000 for 50,000 VND
                    </p>
                  </div>
                  <div>
                    <Label htmlFor="item-image">Image URL</Label>
                    <Input
                      id="item-image"
                      type="url"
                      value={itemForm.image_url}
                      onChange={(e) => setItemForm({ ...itemForm, image_url: e.target.value })}
                    />
                  </div>
                  {editingItem && (
                    <div className="flex items-center gap-2">
                      <input
                        type="checkbox"
                        id="item-available"
                        checked={itemForm.is_available}
                        onChange={(e) =>
                          setItemForm({ ...itemForm, is_available: e.target.checked })
                        }
                        className="h-4 w-4 rounded border-border"
                      />
                      <Label htmlFor="item-available" className="cursor-pointer">
                        Available
                      </Label>
                    </div>
                  )}
                </>
              )}
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setDialogOpen(false)}
                disabled={submitting}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={submitting}>
                {submitting ? "Saving..." : "Save"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Confirmation Dialog */}
      <Dialog open={confirmDialog.open} onOpenChange={(open) =>
        setConfirmDialog((prev) => ({ ...prev, open }))
      }>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{confirmDialog.title}</DialogTitle>
          </DialogHeader>
          <div className="p-6">
            <p className="text-sm text-muted-foreground">{confirmDialog.message}</p>
          </div>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setConfirmDialog((prev) => ({ ...prev, open: false }))}
            >
              Cancel
            </Button>
            <Button
              type="button"
              variant="danger"
              onClick={confirmDialog.onConfirm}
            >
              Confirm
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
