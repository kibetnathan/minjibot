import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import type { CommandCategory } from "@/data/commands"
import { CommandRow } from "@/components/commands/CommandRow"

type CommandCardProps = {
  category: CommandCategory
}

export function CommandCard({ category }: CommandCardProps) {
  const Icon = category.icon

  return (
    <Card className="flex flex-col">
      <CardHeader>
        <div className="flex items-start justify-between gap-2">
          <div className="flex items-center gap-2">
            <span className="flex size-8 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
              <Icon className="size-4" />
            </span>
            <CardTitle>{category.name}</CardTitle>
          </div>
          <Badge variant="outline">{category.commands.length}</Badge>
        </div>
        <CardDescription>{category.description}</CardDescription>
      </CardHeader>
      <CardContent>
        <ul className="grid gap-3">
          {category.commands.map((command) => (
            <CommandRow key={command.name} command={command} />
          ))}
        </ul>
      </CardContent>
    </Card>
  )
}
