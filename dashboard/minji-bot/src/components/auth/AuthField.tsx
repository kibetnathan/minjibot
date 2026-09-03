import type { ComponentProps } from "react"
import { Input } from "@/components/ui/input"
import { cn } from "@/lib/utils"

type AuthFieldProps = ComponentProps<typeof Input> & {
  label: string
  id: string
}

export function AuthField({ label, id, className, ...props }: AuthFieldProps) {
  return (
    <div className="grid gap-1.5">
      <label
        htmlFor={id}
        className="text-xs font-medium text-foreground"
      >
        {label}
      </label>
      <Input id={id} className={cn("h-9", className)} {...props} />
    </div>
  )
}
