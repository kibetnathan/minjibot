import type { ComponentType } from "react"
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

type StatCardProps = {
  icon: ComponentType<{ className?: string }>
  label: string
  value: string
}

export function StatCard({ icon: Icon, label, value }: StatCardProps) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center gap-3">
        <span className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
          <Icon className="size-4" />
        </span>
        <div className="grid gap-0.5">
          <CardDescription>{label}</CardDescription>
          <CardTitle className="text-base">{value}</CardTitle>
        </div>
      </CardHeader>
    </Card>
  )
}
