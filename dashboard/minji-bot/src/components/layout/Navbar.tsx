import { Link } from "react-router-dom"
import { SiDiscord } from "react-icons/si"
import { buttonVariants } from "@/components/ui/button"
import { cn } from "@/lib/utils"

const INVITE_URL =
  "https://discord.com/oauth2/authorize?client_id=1542965379075281076"

export function Navbar() {
  return (
    <header className="border-b border-border">
      <div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
        <div className="flex items-center gap-2">
          <Link to="/" className="flex items-center gap-2">
            <img src="beans.svg" alt="Minji" className="h-[2rem]" />
            <span className="font-heading text-xl font-bold">MinjiBot</span>
          </Link>
        </div>
        <nav className="hidden items-center gap-6 text-sm text-muted-foreground md:flex">
          <Link to="/commands" className="transition-colors hover:text-foreground">
            Commands
          </Link>
          <a
            href="#features"
            className="transition-colors hover:text-foreground"
          >
            Features
          </a>
          <a href="#tech" className="transition-colors hover:text-foreground">
            Tech Stack
          </a>
          <a
            href={INVITE_URL}
            target="_blank"
            rel="noopener noreferrer"
            className={cn(buttonVariants({ variant: "default" }), "gap-2")}
          >
            <SiDiscord className="size-4" />
            Invite Bot
          </a>
        </nav>
      </div>
    </header>
  )
}
