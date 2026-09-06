import { useEffect, useState } from "react"
import { useParams, Link } from "react-router-dom"
import { ArrowLeft, Loader2 } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Button, buttonVariants } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { DashboardHeader } from "@/components/dashboard/DashboardHeader"
import { apiUrl } from "@/lib/api"
import { useCurrentUser } from "@/lib/useCurrentUser"

type Settings = {
  guild_id: string
  prefix: string
  language: string
  auto_moderation_enabled: boolean
  logging_channel_id: string
}

export default function GuildSettings() {
  const { guildId } = useParams<{ guildId: string }>()
  const me = useCurrentUser()

  const [prefix, setPrefix] = useState("-")
  const [language, setLanguage] = useState("en")
  const [autoMod, setAutoMod] = useState(false)
  const [loggingChannel, setLoggingChannel] = useState("")
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)
  const [saved, setSaved] = useState(false)

  useEffect(() => {
    if (me.status !== "authenticated" || !guildId) return
    let cancelled = false
    setLoading(true)
    setError(null)

    fetch(apiUrl(`/api/guilds/${guildId}/settings`), {
      credentials: "include",
      headers: { Accept: "application/json" },
    })
      .then(async (res) => {
        if (!res.ok) throw new Error(`${res.status}`)
        return (await res.json()) as Settings
      })
      .then((s) => {
        if (!cancelled) {
          setPrefix(s.prefix)
          setLanguage(s.language)
          setAutoMod(s.auto_moderation_enabled)
          setLoggingChannel(s.logging_channel_id)
          setLoading(false)
        }
      })
      .catch((err) => {
        console.error("guild settings:", err)
        if (!cancelled) {
          setError("Could not load settings.")
          setLoading(false)
        }
      })

    return () => {
      cancelled = true
    }
  }, [me.status, guildId])

  async function handleSave() {
    if (!guildId) return
    setSaving(true)
    setSaved(false)
    setError(null)

    const body: Record<string, unknown> = {
      prefix,
      language,
      auto_moderation_enabled: autoMod,
      logging_channel_id: loggingChannel,
    }

    try {
      const res = await fetch(apiUrl(`/api/guilds/${guildId}/settings`), {
        method: "PUT",
        credentials: "include",
        headers: { "Content-Type": "application/json", Accept: "application/json" },
        body: JSON.stringify(body),
      })
      if (!res.ok) {
        const data = (await res.json().catch(() => null)) ?? {}
        throw new Error((data as { error?: string }).error ?? `${res.status}`)
      }
      const updated = (await res.json()) as Settings
      setPrefix(updated.prefix)
      setLanguage(updated.language)
      setAutoMod(updated.auto_moderation_enabled)
      setLoggingChannel(updated.logging_channel_id)
      setSaved(true)
      setTimeout(() => setSaved(false), 2000)
    } catch (err) {
      console.error("save settings:", err)
      setError(String(err))
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="min-h-screen bg-background font-sans antialiased">
      <DashboardHeader />
      <main className="mx-auto max-w-3xl px-4 py-10 sm:px-6">
        <div className="mb-6">
          <Link
            to={`/dashboard/guild/${guildId}`}
            className={buttonVariants({ variant: "ghost", size: "sm" }) + " mb-2 -ml-2"}
          >
            <ArrowLeft className="mr-1 size-4" />
            Back to logs
          </Link>
          <h1 className="mb-1 font-heading text-3xl font-bold tracking-tight text-foreground">
            Server settings
          </h1>
          <p className="text-sm text-muted-foreground">
            Bot configuration for <code className="font-mono text-xs">{guildId}</code>.
          </p>
        </div>

        {me.status !== "authenticated" ? (
          <Card>
            <CardContent className="py-8">
              <p className="text-sm text-muted-foreground">Log in to manage settings.</p>
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
            Loading settings…
          </div>
        ) : error ? (
          <p className="text-sm text-destructive">{error}</p>
        ) : (
          <Card>
            <CardHeader>
              <CardTitle>Bot configuration</CardTitle>
              <CardDescription>
                These apply to this server. Changes take effect immediately.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-5">
              <div className="space-y-1.5">
                <label className="text-xs font-medium" htmlFor="prefix">
                  Command prefix
                </label>
                <Input
                  id="prefix"
                  value={prefix}
                  onChange={(e) => setPrefix(e.target.value)}
                  maxLength={3}
                  placeholder="-"
                />
                <p className="text-xs text-muted-foreground">
                  Prefix used for text commands, e.g. "-ban". Max 3 characters.
                </p>
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-medium" htmlFor="language">
                  Language
                </label>
                <select
                  id="language"
                  value={language}
                  onChange={(e) => setLanguage(e.target.value)}
                  className="h-8 w-full rounded-none border border-input bg-transparent px-2.5 py-1 text-xs outline-none focus-visible:border-ring focus-visible:ring-1 focus-visible:ring-ring/50"
                >
                  <option value="en">English</option>
                  <option value="es">Español</option>
                  <option value="fr">Français</option>
                </select>
              </div>

              <div className="space-y-1.5">
                <div className="flex items-center justify-between gap-4">
                  <div>
                    <label className="text-xs font-medium" htmlFor="autoMod">
                      Auto-moderation
                    </label>
                    <p className="text-xs text-muted-foreground">
                      Automatically handle configured moderation rules.
                    </p>
                  </div>
                  <input
                    id="autoMod"
                    type="checkbox"
                    checked={autoMod}
                    onChange={(e) => setAutoMod(e.target.checked)}
                    className="size-4 accent-primary"
                  />
                </div>
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-medium" htmlFor="loggingChannel">
                  Logging channel ID
                </label>
                <Input
                  id="loggingChannel"
                  value={loggingChannel}
                  onChange={(e) => setLoggingChannel(e.target.value)}
                  placeholder="Channel ID (e.g. 1234567890)"
                />
                <p className="text-xs text-muted-foreground">
                  Where moderation events are logged. Leave blank to disable.
                </p>
              </div>

              <div className="flex items-center gap-3 pt-2">
                <Button onClick={handleSave} disabled={saving}>
                  {saving && <Loader2 className="size-4 animate-spin" />}
                  {saving ? "Saving…" : "Save settings"}
                </Button>
                {saved && <span className="text-xs text-foreground">Saved.</span>}
              </div>
            </CardContent>
          </Card>
        )}
      </main>
    </div>
  )
}
