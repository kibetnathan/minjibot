import { Bot, Zap, Shield, Globe, Terminal, Sparkles } from "lucide-react"

const features = [
  {
    icon: Bot,
    title: "Discord Bot",
    description:
      "Slash and prefix commands for moderation, utility, fun, roleplay, and social lookups.",
  },
  {
    icon: Zap,
    title: "REST API",
    description:
      "Echo-powered API on port 8080 with pgx/v5 and generated sqlc queries.",
  },
  {
    icon: Shield,
    title: "Moderation Suite",
    description:
      "Ban, timeout, purge, lockdown, roles, warnings, and audit log inspection.",
  },
  {
    icon: Globe,
    title: "Social Lookups",
    description:
      "GitHub, Twitter, YouTube, Twitch, TikTok, Reddit, and more profile fetching.",
  },
  {
    icon: Terminal,
    title: "Type-Safe Database",
    description:
      "PostgreSQL with goose migrations, sqlc codegen, and pgtype mappings.",
  },
  {
    icon: Sparkles,
    title: "Dashboard",
    description:
      "React + Vite + shadcn/ui frontend with dark mode and type-safe hooks.",
  },
]

export function Features() {
  return (
    <section
      id="features"
      className="mx-auto max-w-7xl bg-muted/30 px-4 py-20 sm:px-6 lg:px-8"
    >
      <div className="mb-16 text-center">
        <h2 className="mb-4 font-heading text-3xl font-bold text-foreground sm:text-4xl">
          Features
        </h2>
        <p className="mx-auto max-w-2xl text-muted-foreground">
          Dual support for slash (/) and classic prefix (!) commands across
          every category.
        </p>
      </div>
      <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        {features.map((feature, index) => (
          <div
            key={index}
            className="group rounded-xl border border-border bg-card p-6 transition-all hover:border-primary/50 hover:shadow-lg"
          >
            <div className="mb-4 w-fit rounded-lg bg-primary/10 p-3">
              <feature.icon className="h-6 w-6 text-primary" />
            </div>
            <h3 className="mb-2 font-heading text-lg font-semibold text-foreground">
              {feature.title}
            </h3>
            <p className="text-sm leading-relaxed text-muted-foreground">
              {feature.description}
            </p>
          </div>
        ))}
      </div>
    </section>
  )
}
