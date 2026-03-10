interface TableProps {
  children: React.ReactNode;
  className?: string;
}

interface TableHeaderProps {
  children: React.ReactNode;
}

interface TableBodyProps {
  children: React.ReactNode;
}

interface TableRowProps {
  children: React.ReactNode;
  className?: string;
}

interface TableHeadProps {
  children: React.ReactNode;
  className?: string;
}

interface TableCellProps {
  children: React.ReactNode;
  className?: string;
}

export function Table({ children, className = "" }: TableProps): React.JSX.Element {
  return (
    <div className="relative w-full overflow-auto">
      <table className={`w-full caption-bottom text-sm ${className}`}>{children}</table>
    </div>
  );
}

export function TableHeader({ children }: TableHeaderProps): React.JSX.Element {
  return <thead className="border-b border-border">{children}</thead>;
}

export function TableBody({ children }: TableBodyProps): React.JSX.Element {
  return <tbody className="[&_tr:last-child]:border-0">{children}</tbody>;
}

export function TableRow({ children, className = "" }: TableRowProps): React.JSX.Element {
  return (
    <tr className={`border-b border-border transition-colors hover:bg-muted/50 ${className}`}>
      {children}
    </tr>
  );
}

export function TableHead({ children, className = "" }: TableHeadProps): React.JSX.Element {
  return (
    <th
      className={`h-12 px-4 text-left align-middle font-medium text-muted-foreground ${className}`}
    >
      {children}
    </th>
  );
}

export function TableCell({ children, className = "" }: TableCellProps): React.JSX.Element {
  return <td className={`p-4 align-middle ${className}`}>{children}</td>;
}
