import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { buttonVariants } from "@/components/ui/button"
import { DashboardHeader } from "@/components/dashboard/DashboardHeader"
import { apiUrl } from "@/lib/api"
import { useCurrentUser } from "@/lib/useCurrentUser"

export default function Profile() {
  const me = useCurrentUser()

  return (
    <div className="min-h-screen bg-background font-sans antialiased">
      <DashboardHeader />
      <main className="mx-auto max-w-3xl px-4 py-10 sm:px-6">
        <h1 className="mb-1 font-heading text-3xl font-bold tracking-tight text-foreground">
          Profile
        </h1>
        <p className="mb-6 text-sm text-muted-foreground">
          Your MinjiBot account and how it interacts with the dashboard.
        </p>

        {me.status === "loading" ? (
          <p className="text-sm text-muted-foreground">Loading…</p>
        ) : me.status !== "authenticated" ? (
          <Card>
            <CardContent className="py-8">
              <p className="text-sm text-muted-foreground">Log in to view your profile.</p>
              <a
                href={apiUrl("/api/auth/discord")}
                className={buttonVariants({ variant: "default" }) + " mt-4"}
              >
                Log in with Discord
              </a>
            </CardContent>
          </Card>
        ) : (
          <Card>
            <CardHeader>
              <CardTitle>Account</CardTitle>
              <CardDescription>Details linked to your Discord sign-in.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between gap-4">
                <span className="text-xs font-medium text-muted-foreground">Discord ID</span>
                <code className="font-mono text-xs">{me.user.id}</code>
              </div>
              <div className="flex items-center justify-between gap-4">
                <span className="text-xs font-medium text-muted-foreground">Email</span>
                <span className="text-xs text-foreground">{me.user.email || "—"}</span>
              </div>
              <div className="flex items-center justify-between gap-4">
                <span className="text-xs font-medium text-muted-foreground">Role</span>
                <Badge variant={me.user.is_admin ? "default" : "outline"}>
                  {me.user.is_admin ? "Administrator" : "Member"}
                </Badge>
              </div>
            </CardContent>
          </Card>
        )}
      </main>
    </div>
  )
}