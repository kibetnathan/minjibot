import { Link } from "react-router-dom"
import { buttonVariants } from "@/components/ui/button"
import { apiUrl } from "@/lib/api"
import { useCurrentUser } from "@/lib/useCurrentUser"
import { cn } from "@/lib/utils"

export function DashboardHeader() {
  const me = useCurrentUser()

  return (
    <header className="border-b border-border">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-4 sm:px-6">
        <div className="flex items-center gap-2">
          <img src="beans.svg" alt="Minji" className="h-[2rem]" />
          <span className="font-heading text-xl font-bold">Dashboard</span>
        </div>
        <nav className="flex items-center gap-3 text-sm text-muted-foreground">
          <Link to="/" className="transition-colors hover:text-foreground">
            Home
          </Link>
          <Link
            to="/commands"
            className="transition-colors hover:text-foreground"
          >
            Commands
          </Link>
          <span className="hidden truncate text-xs text-muted-foreground sm:inline">
            {me.status === "authenticated" ? me.user.email : ""}
          </span>
          <a
            href={apiUrl("/api/auth/logout")}
            className={cn(buttonVariants({ variant: "outline" }))}
          >
            Log out
          </a>
        </nav>
      </div>
    </header>
  )
}
