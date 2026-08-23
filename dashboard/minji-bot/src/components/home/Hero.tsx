import {
  ArrowRight,
  Bot,
  Zap,
  Shield,
  Globe,
  Terminal,
  Sparkles,
} from "lucide-react"

const buttonBase =
  "group/button inline-flex shrink-0 items-center justify-center rounded-none border border-transparent bg-clip-padding text-xs font-medium whitespace-nowrap transition-all outline-none select-none focus-visible:border-ring focus-visible:ring-1 focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"

const buttonVariants = {
  default: "bg-primary text-primary-foreground hover:bg-primary/80",
  outline: "border-border bg-background hover:bg-muted hover:text-foreground",
}

const buttonSizes = {
  lg: "h-9 gap-1.5 px-2.5",
  default: "h-8 gap-1.5 px-2.5",
}

export function Hero() {
  const features = [
    {
      icon: Bot,
      title: "Discord Bot",
      desc: "Slash & prefix commands for moderation, utility, fun, roleplay, social lookups.",
    },
    {
      icon: Zap,
      title: "REST API",
      desc: "Echo-powered API with pgx/v5 and generated sqlc queries.",
    },
    {
      icon: Shield,
      title: "Moderation",
      desc: "Ban, timeout, purge, lockdown, roles, warnings, audit logs.",
    },
    {
      icon: Globe,
      title: "Social Lookups",
      desc: "GitHub, Twitter, YouTube, Twitch, TikTok, Reddit profiles.",
    },
    {
      icon: Terminal,
      title: "Type-Safe DB",
      desc: "PostgreSQL with goose migrations, sqlc codegen, pgtype mappings.",
    },
    {
      icon: Sparkles,
      title: "Dashboard",
      desc: "React + Vite + shadcn/ui with dark mode and type-safe hooks.",
    },
  ]

  return (
    <section className="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8 lg:py-32">
      <div className="mx-auto max-w-3xl text-center">
        <h1 className="mb-6 font-heading text-4xl font-bold tracking-tight text-foreground sm:text-5xl lg:text-6xl">
          A Discord Bot with a Companion API and Dashboard
        </h1>
        <p className="mx-auto mb-10 max-w-2xl text-lg leading-relaxed text-muted-foreground sm:text-xl">
          Six feature categories. Dual slash (/) and prefix (!) support. Built
          with Go and React.
        </p>
        <div className="mb-16 flex flex-col items-center justify-center gap-4 sm:flex-row">
          <a
            href="https://github.com/kibetnathan/minjibot"
            target="_blank"
            rel="noopener noreferrer"
            className={`${buttonBase} ${buttonVariants.default} ${buttonSizes.lg} flex items-center gap-2`}
          >
            View on GitHub
            <ArrowRight className="h-4 w-4" />
          </a>
          <a
            href="#features"
            className={`${buttonBase} ${buttonVariants.outline} ${buttonSizes.lg} flex items-center gap-2`}
          >
            Explore Features
          </a>
        </div>

        <div className="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-6">
          {features.map((feature, index) => (
            <div
              key={index}
              className="group rounded-xl border border-border bg-card p-4 text-left transition-all hover:border-primary/50 hover:shadow-lg"
            >
              <div className="mb-3 w-fit rounded-lg bg-primary/10 p-2">
                <feature.icon className="h-5 w-5 text-primary" />
              </div>
              <h3 className="font-heading text-sm font-semibold text-foreground">
                {feature.title}
              </h3>
              <p className="mt-1 text-xs leading-relaxed text-muted-foreground">
                {feature.desc}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
