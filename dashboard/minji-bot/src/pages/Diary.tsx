import { useCallback, useEffect, useState } from "react"
import { Loader2, Plus, Trash2 } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Button, buttonVariants } from "@/components/ui/button"
import { DashboardHeader } from "@/components/dashboard/DashboardHeader"
import { apiUrl } from "@/lib/api"
import { useCurrentUser } from "@/lib/useCurrentUser"

type DiaryEntry = {
  id: number
  content: string
  created_at: string
}

type DiaryResponse = {
  items: DiaryEntry[]
}

function formatTimestamp(ts: string): string {
  try {
    return new Date(ts).toLocaleString()
  } catch {
    return ts
  }
}

export default function Diary() {
  const me = useCurrentUser()
  const [entries, setEntries] = useState<DiaryEntry[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [draft, setDraft] = useState("")
  const [saving, setSaving] = useState(false)

  const load = useCallback(() => {
    if (me.status !== "authenticated") return
    setLoading(true)
    setError(null)
    fetch(apiUrl("/api/diary"), {
      credentials: "include",
      headers: { Accept: "application/json" },
    })
      .then(async (res) => {
        if (!res.ok) throw new Error(`${res.status}`)
        const data = (await res.json()) as DiaryResponse
        setEntries(data.items)
      })
      .catch((err) => {
        console.error("diary:", err)
        setError("Could not load your diary.")
      })
      .finally(() => setLoading(false))
  }, [me.status])

  useEffect(() => {
    load()
  }, [load])

  async function handleAdd() {
    const content = draft.trim()
    if (!content) return
    setSaving(true)
    setError(null)
    try {
      const res = await fetch(apiUrl("/api/diary"), {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json", Accept: "application/json" },
        body: JSON.stringify({ content }),
      })
      if (!res.ok) throw new Error(`${res.status}`)
      const created = (await res.json()) as DiaryEntry
      setEntries((prev) => [created, ...prev])
      setDraft("")
    } catch (err) {
      console.error("add entry:", err)
      setError("Could not save your entry.")
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete(id: number) {
    setError(null)
    try {
      const res = await fetch(apiUrl(`/api/diary/${id}`), {
        method: "DELETE",
        credentials: "include",
      })
      if (!res.ok && res.status !== 204) throw new Error(`${res.status}`)
      setEntries((prev) => prev.filter((e) => e.id !== id))
    } catch (err) {
      console.error("delete entry:", err)
      setError("Could not delete that entry.")
    }
  }

  return (
    <div className="min-h-screen bg-background font-sans antialiased">
      <DashboardHeader />
      <main className="mx-auto max-w-3xl px-4 py-10 sm:px-6">
        <h1 className="mb-1 font-heading text-3xl font-bold tracking-tight text-foreground">
          Diary
        </h1>
        <p className="mb-6 text-sm text-muted-foreground">
          Private journal entries — only visible to you.
        </p>

        {me.status !== "authenticated" ? (
          <Card>
            <CardContent className="py-8">
              <p className="text-sm text-muted-foreground">Log in to manage your diary.</p>
              <a
                href={apiUrl("/api/auth/discord")}
                className={buttonVariants({ variant: "default" }) + " mt-4"}
              >
                Log in with Discord
              </a>
            </CardContent>
          </Card>
        ) : (
          <>
            <Card className="mb-6">
              <CardHeader>
                <CardTitle>New entry</CardTitle>
                <CardDescription>Write down anything you want to remember.</CardDescription>
              </CardHeader>
              <CardContent className="space-y-3">
                <textarea
                  value={draft}
                  onChange={(e) => setDraft(e.target.value)}
                  rows={3}
                  placeholder="What happened today?"
                  className="w-full resize-none rounded-none border border-input bg-transparent px-2.5 py-1.5 text-xs outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-1 focus-visible:ring-ring/50"
                />
                <div className="flex items-center gap-3">
                  <Button onClick={handleAdd} disabled={saving || draft.trim() === ""}>
                    {saving ? (
                      <Loader2 className="size-4 animate-spin" />
                    ) : (
                      <Plus className="size-4" />
                    )}
                    {saving ? "Saving…" : "Add entry"}
                  </Button>
                  {error && <span className="text-xs text-destructive">{error}</span>}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Entries</CardTitle>
                <CardDescription>
                  {entries.length} total
                </CardDescription>
              </CardHeader>
              <CardContent>
                {loading ? (
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <Loader2 className="size-4 animate-spin" />
                    Loading diary…
                  </div>
                ) : entries.length === 0 ? (
                  <p className="text-xs text-muted-foreground">No entries yet.</p>
                ) : (
                  <div className="space-y-2">
                    {entries.map((entry) => (
                      <div
                        key={entry.id}
                        className="flex items-start justify-between gap-3 rounded-none border border-border p-3"
                      >
                        <div className="min-w-0">
                          <p className="text-xs whitespace-pre-wrap break-words text-foreground/90">
                            {entry.content}
                          </p>
                          <p className="mt-1 text-xs text-muted-foreground">
                            {formatTimestamp(entry.created_at)}
                          </p>
                        </div>
                        <Button
                          variant="ghost"
                          size="icon-sm"
                          aria-label="Delete entry"
                          onClick={() => handleDelete(entry.id)}
                        >
                          <Trash2 className="size-4" />
                        </Button>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          </>
        )}
      </main>
    </div>
  )
}