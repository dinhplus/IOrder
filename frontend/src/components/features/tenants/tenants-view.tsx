"use client";

import { useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select } from "@/components/ui/select";
import { createTenant, updateTenant, type CreateTenantRequest, type UpdateTenantRequest } from "@/lib/api/tenants";
import type { Tenant } from "@/types/tenant";

interface TenantsViewProps {
  initialTenants: Tenant[];
}

export function TenantsView({ initialTenants }: TenantsViewProps) {
  const [tenants, setTenants] = useState<Tenant[]>(initialTenants);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingTenant, setEditingTenant] = useState<Tenant | null>(null);
  const [formData, setFormData] = useState({
    slug: "",
    name: "",
    logo_url: "",
    timezone: "Asia/Ho_Chi_Minh",
    currency: "VND",
    is_active: true,
  });
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  const handleCreate = () => {
    setEditingTenant(null);
    setFormData({
      slug: "",
      name: "",
      logo_url: "",
      timezone: "Asia/Ho_Chi_Minh",
      currency: "VND",
      is_active: true,
    });
    setError("");
    setDialogOpen(true);
  };

  const handleEdit = (tenant: Tenant) => {
    setEditingTenant(tenant);
    setFormData({
      slug: tenant.slug,
      name: tenant.name,
      logo_url: tenant.logo_url ?? "",
      timezone: tenant.timezone,
      currency: tenant.currency,
      is_active: tenant.is_active,
    });
    setError("");
    setDialogOpen(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError("");

    try {
      if (editingTenant) {
        const updated: UpdateTenantRequest = {
          slug: formData.slug || undefined,
          name: formData.name || undefined,
          logo_url: formData.logo_url || undefined,
          timezone: formData.timezone,
          currency: formData.currency,
          is_active: formData.is_active,
        };
        const result = await updateTenant(editingTenant.id, updated);
        setTenants((prev) => prev.map((t) => (t.id === result.id ? result : t)));
      } else {
        const newTenant: CreateTenantRequest = {
          slug: formData.slug,
          name: formData.name,
          logo_url: formData.logo_url || undefined,
          timezone: formData.timezone,
          currency: formData.currency,
          is_active: formData.is_active,
        };
        const result = await createTenant(newTenant);
        setTenants((prev) => [result, ...prev]);
      }
      setDialogOpen(false);
    } catch (err: unknown) {
      setError((err as Error)?.message || "Failed to save tenant");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <div className="flex justify-end">
        <Button onClick={handleCreate}>Create Tenant</Button>
      </div>

      {tenants.length === 0 ? (
        <p className="text-muted-foreground">No tenants found.</p>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {tenants.map((tenant) => (
            <Card key={tenant.id}>
              <CardHeader>
                <div className="flex items-start justify-between gap-2">
                  <div>
                    <CardTitle className="text-base">{tenant.name}</CardTitle>
                    <CardDescription className="font-mono text-xs">{tenant.slug}</CardDescription>
                  </div>
                  <Badge variant={tenant.is_active ? "success" : "outline"}>
                    {tenant.is_active ? "Active" : "Inactive"}
                  </Badge>
                </div>
              </CardHeader>
              <CardContent className="flex flex-col gap-3">
                <dl className="flex flex-col gap-1 text-sm">
                  <div className="flex justify-between">
                    <dt className="text-muted-foreground">Currency</dt>
                    <dd className="font-medium">{tenant.currency}</dd>
                  </div>
                  <div className="flex justify-between">
                    <dt className="text-muted-foreground">Timezone</dt>
                    <dd className="font-medium text-xs">{tenant.timezone}</dd>
                  </div>
                </dl>
                <Button variant="outline" size="sm" onClick={() => handleEdit(tenant)}>
                  Edit
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent>
          <form onSubmit={handleSubmit}>
            <DialogHeader>
              <DialogTitle>{editingTenant ? "Edit Tenant" : "Create Tenant"}</DialogTitle>
            </DialogHeader>

            <div className="flex flex-col gap-4 p-6">
              {error && (
                <div className="p-3 text-sm bg-danger/10 text-danger rounded-md">{error}</div>
              )}

              <div>
                <Label htmlFor="name">Name *</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  required
                  placeholder="My Restaurant"
                />
              </div>

              <div>
                <Label htmlFor="slug">Slug *</Label>
                <Input
                  id="slug"
                  value={formData.slug}
                  onChange={(e) => setFormData({ ...formData, slug: e.target.value })}
                  required
                  placeholder="my-restaurant"
                  pattern="[a-z0-9-]+"
                />
                <p className="text-xs text-muted-foreground mt-1">
                  Lowercase letters, numbers, and hyphens only
                </p>
              </div>

              <div>
                <Label htmlFor="logo_url">Logo URL</Label>
                <Input
                  id="logo_url"
                  type="url"
                  value={formData.logo_url}
                  onChange={(e) => setFormData({ ...formData, logo_url: e.target.value })}
                  placeholder="https://example.com/logo.png"
                />
              </div>

              <div>
                <Label htmlFor="currency">Currency</Label>
                <Select
                  id="currency"
                  value={formData.currency}
                  onChange={(e) => setFormData({ ...formData, currency: e.target.value })}
                >
                  <option value="VND">VND - Vietnamese Dong</option>
                  <option value="USD">USD - US Dollar</option>
                  <option value="EUR">EUR - Euro</option>
                  <option value="GBP">GBP - British Pound</option>
                  <option value="JPY">JPY - Japanese Yen</option>
                </Select>
              </div>

              <div>
                <Label htmlFor="timezone">Timezone</Label>
                <Select
                  id="timezone"
                  value={formData.timezone}
                  onChange={(e) => setFormData({ ...formData, timezone: e.target.value })}
                >
                  <option value="Asia/Ho_Chi_Minh">Asia/Ho Chi Minh</option>
                  <option value="Asia/Bangkok">Asia/Bangkok</option>
                  <option value="Asia/Singapore">Asia/Singapore</option>
                  <option value="Asia/Tokyo">Asia/Tokyo</option>
                  <option value="America/New_York">America/New York</option>
                  <option value="Europe/London">Europe/London</option>
                </Select>
              </div>

              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="is_active"
                  checked={formData.is_active}
                  onChange={(e) => setFormData({ ...formData, is_active: e.target.checked })}
                  className="h-4 w-4 rounded border-border"
                />
                <Label htmlFor="is_active" className="cursor-pointer">
                  Active
                </Label>
              </div>
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
    </>
  );
}
