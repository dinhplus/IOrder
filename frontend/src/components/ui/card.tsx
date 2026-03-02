interface CardProps {
  children: React.ReactNode;
  className?: string;
}

interface CardHeaderProps {
  children: React.ReactNode;
  className?: string;
}

interface CardTitleProps {
  children: React.ReactNode;
  className?: string;
}

interface CardDescriptionProps {
  children: React.ReactNode;
  className?: string;
}

interface CardContentProps {
  children: React.ReactNode;
  className?: string;
}

export function Card({ children, className = "" }: CardProps): React.JSX.Element {
  return (
    <div className={`rounded-lg border border-border bg-background shadow-sm ${className}`}>
      {children}
    </div>
  );
}

export function CardHeader({ children, className = "" }: CardHeaderProps): React.JSX.Element {
  return <div className={`flex flex-col gap-1.5 p-6 ${className}`}>{children}</div>;
}

export function CardTitle({ children, className = "" }: CardTitleProps): React.JSX.Element {
  return <h3 className={`text-lg font-semibold leading-none ${className}`}>{children}</h3>;
}

export function CardDescription({
  children,
  className = "",
}: CardDescriptionProps): React.JSX.Element {
  return <p className={`text-sm text-muted-foreground ${className}`}>{children}</p>;
}

export function CardContent({ children, className = "" }: CardContentProps): React.JSX.Element {
  return <div className={`p-6 pt-0 ${className}`}>{children}</div>;
}
