import {
  Bot,
  Box,
  Cpu,
  Database,
  FileCode,
  GitBranch,
  LayoutGrid,
  Wind,
  Zap,
} from "lucide-react"

const techStack = [
  { name: "Go 1.26", icon: Box },
  { name: "DiscordGo", icon: Bot },
  { name: "Echo v5", icon: Zap },
  { name: "pgx/v5", icon: Database },
  { name: "sqlc", icon: GitBranch },
  { name: "goose", icon: Cpu },
  { name: "React 19", icon: LayoutGrid },
  { name: "Tailwind CSS v4", icon: Wind },
  { name: "shadcn/ui", icon: LayoutGrid },
  { name: "TypeScript", icon: FileCode },
]

export function TechStack() {
  return (
    <section id="tech" className="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8">
      <div className="mb-16 text-center">
        <h2 className="mb-4 font-heading text-3xl font-bold text-foreground sm:text-4xl">
          Tech Stack
        </h2>
        <p className="mx-auto max-w-2xl text-muted-foreground">
          Modern, type-safe tooling across the entire stack.
        </p>
      </div>
      <div className="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-5">
        {techStack.map((tech, index) => (
          <div
            key={index}
            className="rounded-lg border border-border bg-card p-4 text-center transition-colors hover:border-primary/50"
          >
            <tech.icon className="mx-auto mb-2 h-6 w-6 text-primary" />
            <span className="text-sm font-medium text-foreground">
              {tech.name}
            </span>
          </div>
        ))}
      </div>
    </section>
  )
}
