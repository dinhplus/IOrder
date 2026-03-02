interface BadgeProps {
  variant?: "default" | "success" | "warning" | "danger" | "outline";
  children: React.ReactNode;
  className?: string;
}

export function Badge({ variant = "default", children, className = "" }: BadgeProps): React.JSX.Element {
  const variantClasses: Record<NonNullable<BadgeProps["variant"]>, string> = {
    default: "bg-secondary text-secondary-foreground",
    success: "bg-success/10 text-success",
    warning: "bg-warning/10 text-warning",
    danger: "bg-danger/10 text-danger",
    outline: "border border-border bg-transparent text-foreground",
  };

  return (
    <span
      className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${variantClasses[variant]} ${className}`}
    >
      {children}
    </span>
  );
}
