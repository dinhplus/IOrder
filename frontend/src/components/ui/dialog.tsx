"use client";

import { useEffect, useRef } from "react";

interface DialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  children: React.ReactNode;
}

interface DialogContentProps {
  children: React.ReactNode;
  className?: string;
}

interface DialogHeaderProps {
  children: React.ReactNode;
}

interface DialogTitleProps {
  children: React.ReactNode;
}

interface DialogFooterProps {
  children: React.ReactNode;
}

export function Dialog({ open, onOpenChange, children }: DialogProps): React.JSX.Element | null {
  const dialogRef = useRef<HTMLDialogElement>(null);

  useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) return;

    if (open && !dialog.open) {
      dialog.showModal();
    } else if (!open && dialog.open) {
      dialog.close();
    }
  }, [open]);

  useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) return;

    const handleClose = (): void => {
      onOpenChange(false);
    };

    const handleEscape = (e: Event): void => {
      e.preventDefault();
      onOpenChange(false);
    };

    dialog.addEventListener("close", handleClose);
    dialog.addEventListener("cancel", handleEscape);

    return () => {
      dialog.removeEventListener("close", handleClose);
      dialog.removeEventListener("cancel", handleEscape);
    };
  }, [onOpenChange]);

  if (!open) return null;

  return (
    <dialog
      ref={dialogRef}
      className="backdrop:bg-black/50 backdrop:backdrop-blur-sm bg-transparent p-0 max-w-lg w-full rounded-lg"
      onClick={(e) => {
        if (e.target === dialogRef.current) {
          onOpenChange(false);
        }
      }}
    >
      {children}
    </dialog>
  );
}

export function DialogContent({ children, className = "" }: DialogContentProps): React.JSX.Element {
  return (
    <div className={`bg-background border border-border rounded-lg shadow-lg ${className}`}>
      {children}
    </div>
  );
}

export function DialogHeader({ children }: DialogHeaderProps): React.JSX.Element {
  return <div className="flex flex-col gap-1.5 p-6">{children}</div>;
}

export function DialogTitle({ children }: DialogTitleProps): React.JSX.Element {
  return <h2 className="text-lg font-semibold leading-none tracking-tight">{children}</h2>;
}

export function DialogFooter({ children }: DialogFooterProps): React.JSX.Element {
  return (
    <div className="flex items-center justify-end gap-2 p-6 pt-0 border-t border-border/50">
      {children}
    </div>
  );
}
