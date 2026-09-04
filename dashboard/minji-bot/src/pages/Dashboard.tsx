import { useEffect, useState } from "react"
import { Link } from "react-router-dom"
import { Server, ScrollText, Trash2 } from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { buttonVariants } from "@/components/ui/button"
import { DashboardHeader } from "@/components/dashboard/DashboardHeader"
import { useCurrentUser } from "@/lib/useCurrentUser"
import { apiUrl } from "@/lib/api"

type GuildSummary = {
  id: string
  name: string
  premium_tier: number
  deleted_messages: number
  mod_actions: number
}

export default function Dashboard() {
  const me = useCurrentUser()
  const [guilds, setGuilds] = useState<GuildSummary[] | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (me.status !== "authenticated") return
    let cancelled = false
    fetch(apiUrl("/api/guilds"), { credentials: "include", headers: { Accept: "application/json" } })
      .then(async (res) => {
        if (!res.ok) throw new Error(`guilds request failed: ${res.status}`)
        return (await res.json()) as GuildSummary[]
      })
      .then((data) => {
        if (!cancelled) setGuilds(data)
      })
      .catch((err) => {
        console.error("guilds:", err)
        if (!cancelled) setError("Could not load guilds.")
      })
    return () => {
      cancelled = true
    }
  }, [me.status])

  const totalDeleted = guilds?.reduce((n, g) => n + g.deleted_messages, 0) ?? 0
  const totalActions = guilds?.reduce((n, g) => n + g.mod_actions, 0) ?? 0

  return (
    <div className="min-h-screen bg-background font-sans antialiased">
      <DashboardHeader />
      <main className="mx-auto max-w-6xl px-4 py-10 sm:px-6">
        <div className="mb-8">
          <h1 className="mb-1 font-heading text-3xl font-bold tracking-tight text-foreground">
            Logging
          </h1>
          <p className="text-sm text-muted-foreground">
            Deleted messages and moderation actions, split by guild.
          </p>
        </div>

        {me.status === "loading" ? (
          <p className="text-sm text-muted-foreground">Loading…</p>
        ) : me.status !== "authenticated" ? (
          <Card>
            <CardContent className="py-8">
              <p className="text-sm text-muted-foreground">
                Log in to view your guild logging data.
              </p>
              <a
                href={apiUrl("/api/auth/discord")}
                className={buttonVariants({ variant: "default", size: "sm" })}
              >
                Log in with Discord
              </a>
            </CardContent>
          </Card>
        ) : error ? (
          <p className="text-sm text-destructive">{error}</p>
        ) : guilds === null ? (
          <p className="text-sm text-muted-foreground">Loading…</p>
        ) : guilds.length === 0 ? (
          <p className="text-sm text-muted-foreground">
            No guilds found. Invite the bot to a server to start logging.
          </p>
        ) : (
          <div className="space-y-3">
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
              <Card>
                <CardHeader>
                  <CardDescription>Servers</CardDescription>
                  <CardTitle className="text-2xl">{guilds.length}</CardTitle>
                </CardHeader>
              </Card>
              <Card>
                <CardHeader>
                  <CardDescription>Deleted messages</CardDescription>
                  <CardTitle className="text-2xl">{totalDeleted}</CardTitle>
                </CardHeader>
              </Card>
              <Card>
                <CardHeader>
                  <CardDescription>Moderation actions</CardDescription>
                  <CardTitle className="text-2xl">{totalActions}</CardTitle>
                </CardHeader>
              </Card>
            </div>
            {guilds.map((g) => (
              <Card key={g.id}>
                <CardContent className="flex items-center justify-between gap-4">
                  <div className="flex items-center gap-3">
                    <span className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
                      <Server className="size-4" />
                    </span>
                    <div>
                      <p className="font-medium text-foreground">{g.name}</p>
                      <p className="text-xs text-muted-foreground">{g.id}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-4 text-xs text-muted-foreground">
                    <span className="inline-flex items-center gap-1">
                      <Trash2 className="size-3.5" />
                      {g.deleted_messages}
                    </span>
                    <span className="inline-flex items-center gap-1">
                      <ScrollText className="size-3.5" />
                      {g.mod_actions}
                    </span>
                    <Link
                      to={`/dashboard/guild/${g.id}`}
                      className={buttonVariants({ variant: "outline", size: "sm" })}
                    >
                      View logs
                    </Link>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </main>
    </div>
  )
}