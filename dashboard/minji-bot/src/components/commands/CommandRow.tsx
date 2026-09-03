import { Badge } from "@/components/ui/badge"
import type { Command } from "@/data/commands"

type CommandRowProps = {
  command: Command
}

export function CommandRow({ command }: CommandRowProps) {
  const Icon = command.icon

  return (
    <li className="flex flex-col gap-0.5">
      <div className="flex items-center gap-2">
        <Badge
          variant="outline"
          className="shrink-0 px-2 py-0.5 font-mono text-[0.7rem]"
        >
          {Icon ? <Icon className="mr-1 inline size-3 shrink-0" /> : null}
          {command.name}
        </Badge>
      </div>
      <p className="pl-0.5 text-xs leading-relaxed text-muted-foreground">
        {command.description}
      </p>
    </li>
  )
}
