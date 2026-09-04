import { useEffect, useState } from "react"
import { useParams, Link } from "react-router-dom"
import { ArrowLeft, Loader2 } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { buttonVariants } from "@/components/ui/button"
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs"
import { ScrollArea } from "@/components/ui/scroll-area"
import { DashboardHeader } from "@/components/dashboard/DashboardHeader"
import { apiUrl } from "@/lib/api"
import { useCurrentUser } from "@/lib/useCurrentUser"

type DeletedMessage = {
  id: number
  channel_id: string
  message_id: string
  author_id: string
  author_name: string
  content: string
  attachments: unknown[] | null
  deleted_by: string
  deleted_by_name: string
  created_at: string
}

type ModAction = {
  id: number
  action: string
  actor_id: string
  actor_name: string
  target_id: string
  target_name: string
  metadata: Record<string, unknown> | null
  created_at: string
}

type PageResponse<T> = {
  total: number
  limit: number
  offset: number
  items: T[]
}

function formatTimestamp(ts: string): string {
  try {
    return new Date(ts).toLocaleString()
  } catch {
    return ts
  }
}

function actionColor(action: string): string {
  if (action.includes("BAN") || action.includes("KICK")) return "destructive"
  if (action.includes("TIMEOUT") || action.includes("WARN")) return "secondary"
  if (action.includes("NUKE") || action.includes("PURGE")) return "outline"
  return "default"
}

export default function GuildLogs() {
  const { guildId } = useParams<{ guildId: string }>()
  const me = useCurrentUser()

  const [deleted, setDeleted] = useState<PageResponse<DeletedMessage> | null>(null)
  const [actions, setActions] = useState<PageResponse<ModAction> | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (me.status !== "authenticated" || !guildId) return
    let cancelled = false
    setLoading(true)
    setError(null)

    Promise.all([
      fetch(apiUrl(`/api/logs/deleted?guild_id=${guildId}`), {
        credentials: "include",
        headers: { Accept: "application/json" },
      }).then(async (res) => {
        if (!res.ok) throw new Error(`${res.status}`)
        return (await res.json()) as PageResponse<DeletedMessage>
      }),
      fetch(apiUrl(`/api/logs/actions?guild_id=${guildId}`), {
        credentials: "include",
        headers: { Accept: "application/json" },
      }).then(async (res) => {
        if (!res.ok) throw new Error(`${res.status}`)
        return (await res.json()) as PageResponse<ModAction>
      }),
    ])
      .then(([d, a]) => {
        if (!cancelled) {
          setDeleted(d)
          setActions(a)
          setLoading(false)
        }
      })
      .catch((err) => {
        console.error("guild logs:", err)
        if (!cancelled) {
          setError("Could not load logs.")
          setLoading(false)
        }
      })

    return () => {
      cancelled = true
    }
  }, [me.status, guildId])

  return (
    <div className="min-h-screen bg-background font-sans antialiased">
      <DashboardHeader />
      <main className="mx-auto max-w-6xl px-4 py-10 sm:px-6">
        <div className="mb-6">
          <Link
            to="/dashboard"
            className={buttonVariants({ variant: "ghost", size: "sm" }) + " mb-2 -ml-2"}
          >
            <ArrowLeft className="mr-1 size-4" />
            Back to guilds
          </Link>
          <h1 className="mb-1 font-heading text-3xl font-bold tracking-tight text-foreground">
            Logs
          </h1>
          <p className="text-sm text-muted-foreground">
            Deleted messages and moderation actions for <code className="font-mono text-xs">{guildId}</code>.
          </p>
        </div>

        {me.status !== "authenticated" ? (
          <Card>
            <CardContent className="py-8">
              <p className="text-sm text-muted-foreground">Log in to view logs.</p>
              <a
                href={apiUrl("/api/auth/discord")}
                className={buttonVariants({ variant: "default" }) + " mt-4"}
              >
                Log in with Discord
              </a>
            </CardContent>
          </Card>
        ) : loading ? (
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Loader2 className="size-4 animate-spin" />
            Loading logs…
          </div>
        ) : error ? (
          <p className="text-sm text-destructive">{error}</p>
        ) : (
          <Tabs defaultValue="deleted">
            <TabsList variant="line">
              <TabsTrigger value="deleted">
                Deleted messages ({deleted?.total ?? 0})
              </TabsTrigger>
              <TabsTrigger value="actions">
                Mod actions ({actions?.total ?? 0})
              </TabsTrigger>
            </TabsList>

            <TabsContent value="deleted">
              <Card className="mt-2">
                <CardHeader>
                  <CardTitle>Deleted messages</CardTitle>
                  <CardDescription>
                    {deleted?.total ?? 0} total — showing {deleted?.items.length ?? 0}
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  {deleted === null || deleted.items.length === 0 ? (
                    <p className="text-xs text-muted-foreground">No deleted messages recorded yet.</p>
                  ) : (
                    <ScrollArea className="max-h-[600px]">
                      <div className="space-y-2">
                        {deleted.items.map((msg) => (
                          <div
                            key={msg.id}
                            className="flex flex-col gap-1 rounded-none border border-border p-3 text-xs"
                          >
                            <div className="flex items-center justify-between gap-2">
                              <span className="font-medium text-foreground">
                                {msg.author_name || msg.author_id}
                              </span>
                              <span className="text-muted-foreground">
                                {formatTimestamp(msg.created_at)}
                              </span>
                            </div>
                            {msg.content && (
                              <p className="text-foreground/80 whitespace-pre-wrap break-words">
                                {msg.content}
                              </p>
                            )}
                            <div className="flex items-center gap-2 text-muted-foreground">
                              <span>channel: {msg.channel_id}</span>
                              {msg.deleted_by_name && (
                                <span>
                                  deleted by: {msg.deleted_by_name}
                                </span>
                              )}
                            </div>
                          </div>
                        ))}
                      </div>
                    </ScrollArea>
                  )}
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="actions">
              <Card className="mt-2">
                <CardHeader>
                  <CardTitle>Moderation actions</CardTitle>
                  <CardDescription>
                    {actions?.total ?? 0} total — showing {actions?.items.length ?? 0}
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  {actions === null || actions.items.length === 0 ? (
                    <p className="text-xs text-muted-foreground">No moderation actions recorded yet.</p>
                  ) : (
                    <ScrollArea className="max-h-[600px]">
                      <div className="space-y-2">
                        {actions.items.map((log) => (
                          <div
                            key={log.id}
                            className="flex flex-col gap-1 rounded-none border border-border p-3 text-xs"
                          >
                            <div className="flex items-center justify-between gap-2">
                              <div className="flex items-center gap-2">
                                <Badge variant={actionColor(log.action) as "destructive" | "secondary" | "outline" | "default"}>
                                  {log.action}
                                </Badge>
                                <span className="text-foreground">
                                  {log.actor_name || log.actor_id}
                                </span>
                                <span className="text-muted-foreground">→</span>
                                <span className="text-foreground">
                                  {log.target_name || log.target_id}
                                </span>
                              </div>
                              <span className="text-muted-foreground">
                                {formatTimestamp(log.created_at)}
                              </span>
                            </div>
                            {log.metadata && typeof log.metadata === "object" && "reason" in log.metadata && (
                              <p className="text-muted-foreground">
                                Reason: {String(log.metadata.reason)}
                              </p>
                            )}
                          </div>
                        ))}
                      </div>
                    </ScrollArea>
                  )}
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        )}
      </main>
    </div>
  )
}