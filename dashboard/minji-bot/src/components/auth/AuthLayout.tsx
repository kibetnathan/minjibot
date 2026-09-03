import type { ReactNode } from "react"
import { Link } from "react-router-dom"

type AuthLayoutProps = {
  title: string
  description: string
  children: ReactNode
}

export function AuthLayout({
  title,
  description,
  children,
}: AuthLayoutProps) {
  return (
    <div className="grid min-h-screen bg-background font-sans antialiased lg:grid-cols-2">
      <aside className="hidden flex-col justify-between border-r border-border bg-card p-10 lg:flex">
        <div className="flex items-center gap-2">
          <img src="beans.svg" alt="Minji" className="h-[2rem]" />
          <span className="font-heading text-xl font-bold">MinjiBot</span>
        </div>
        <div className="max-w-md">
          <p className="text-2xl leading-relaxed font-medium text-foreground">
            Your server, your commands, your dashboard — all in one place.
          </p>
        </div>
        <p className="text-xs leading-relaxed text-muted-foreground">
          MinjiBot helps you moderate, roleplay, search, and have fun across
          your Discord server with a companion dashboard.
        </p>
      </aside>

      <main className="flex flex-col items-center justify-center px-4 py-12 sm:px-6 lg:px-8">
        <div className="mb-6 flex items-center gap-2 lg:hidden">
          <Link to="/" className="flex items-center gap-2">
            <img src="beans.svg" alt="Minji" className="h-[2rem]" />
            <span className="font-heading text-xl font-bold">MinjiBot</span>
          </Link>
        </div>

        <div className="w-full max-w-sm">
          <h1 className="mb-1 font-heading text-2xl font-bold tracking-tight text-foreground">
            {title}
          </h1>
          <p className="mb-8 text-sm leading-relaxed text-muted-foreground">
            {description}
          </p>
          {children}
        </div>
      </main>
    </div>
  )
}
