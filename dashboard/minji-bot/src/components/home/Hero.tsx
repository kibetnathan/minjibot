import { ArrowRight } from "lucide-react"

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
  return (
    <section className="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8 lg:py-32">
      <div className="mx-auto max-w-3xl text-center">
        <h1 className="mb-6 font-heading text-4xl font-bold tracking-tight text-foreground sm:text-5xl lg:text-6xl">
          A Discord Bot with a Companion API and Dashboard
        </h1>
        <p className="mx-auto mb-10 max-w-2xl text-lg leading-relaxed text-muted-foreground sm:text-xl">
          Six feature categories. Dual slash (/) and prefix (-) support. Built
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
      </div>
    </section>
  )
}
