import { useState } from "react"
import { Link, useLocation } from "react-router-dom"
import { Menu, X } from "lucide-react"
import { SiDiscord } from "react-icons/si"
import { Button, buttonVariants } from "@/components/ui/button"
import { cn } from "@/lib/utils"

const INVITE_URL =
  "https://discord.com/oauth2/authorize?client_id=1542965379075281076"

const NAV_LINKS = [
  { to: "/commands", label: "Commands" },
  { href: "#features", label: "Features" },
  { href: "#tech", label: "Tech Stack" },
]

export function Navbar() {
  const [open, setOpen] = useState(false)

  return (
    <header className="border-b border-border">
      <div className="relative mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Link to="/" className="flex items-center gap-2">
              <img src="beans.svg" alt="Minji" className="h-[2rem]" />
              <span className="font-heading text-xl font-bold">MinjiBot</span>
            </Link>
          </div>

          <nav className="hidden items-center gap-6 text-sm text-muted-foreground md:flex">
            {NAV_LINKS.map((link) =>
              "to" in link ? (
                <Link
                  key={link.label}
                  to={link.to}
                  className="transition-colors hover:text-foreground"
                >
                  {link.label}
                </Link>
              ) : (
                <a
                  key={link.label}
                  href={link.href}
                  className="transition-colors hover:text-foreground"
                >
                  {link.label}
                </a>
              )
            )}
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

          <Button
            variant="ghost"
            size="icon"
            className="md:hidden"
            aria-label={open ? "Close menu" : "Open menu"}
            aria-expanded={open}
            onClick={() => setOpen((value) => !value)}
          >
            {open ? <X className="size-5" /> : <Menu className="size-5" />}
          </Button>
        </div>

        {open && (
          <nav className="mt-4 flex flex-col gap-4 border-t border-border pt-4 text-sm text-muted-foreground md:hidden">
            {NAV_LINKS.map((link) =>
              "to" in link ? (
                <Link
                  key={link.label}
                  to={link.to}
                  onClick={() => setOpen(false)}
                  className="transition-colors hover:text-foreground"
                >
                  {link.label}
                </Link>
              ) : (
                <a
                  key={link.label}
                  href={link.href}
                  onClick={() => setOpen(false)}
                  className="transition-colors hover:text-foreground"
                >
                  {link.label}
                </a>
              )
            )}
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
        )}
      </div>
    </header>
  )
}
