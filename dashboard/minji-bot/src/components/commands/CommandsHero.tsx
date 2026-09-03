import { Command } from "lucide-react"

type CommandsHeroProps = {
  total: number
}

export function CommandsHero({ total }: CommandsHeroProps) {
  return (
    <section className="mx-auto max-w-7xl px-4 pt-16 pb-8 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-3xl text-center">
        <span className="mb-4 inline-flex size-12 items-center justify-center rounded-xl bg-primary/10 text-primary">
          <Command className="size-6" />
        </span>
        <h1 className="mb-3 font-heading text-4xl font-bold tracking-tight text-foreground sm:text-5xl">
          Commands
        </h1>
        <p className="mb-4 text-lg leading-relaxed text-muted-foreground">
          Every working MinjiBot command, organized by category. Every command
          works as both a prefix command (
          <code className="font-mono">-cmd</code>) and a slash command (
          <code className="font-mono">/cmd</code>).
        </p>
        <p className="text-sm text-muted-foreground">
          <span className="font-semibold text-foreground">{total}</span>{" "}
          commands across{" "}
          <span className="font-semibold text-foreground">10 categories</span>
        </p>
      </div>
    </section>
  )
}
