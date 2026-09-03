import { Users, Shield, Zap, Server } from "lucide-react"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { DashboardHeader } from "@/components/dashboard/DashboardHeader"
import { StatCard } from "@/components/dashboard/StatCard"

export default function Dashboard() {
  return (
    <div className="min-h-screen bg-background font-sans antialiased">
      <DashboardHeader />
      <main className="mx-auto max-w-6xl px-4 py-10 sm:px-6">
        <div className="mb-8">
          <h1 className="mb-1 font-heading text-3xl font-bold tracking-tight text-foreground">
            Overview
          </h1>
          <p className="text-sm text-muted-foreground">
            Welcome to your MinjiBot dashboard.
          </p>
        </div>

        <div className="mb-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <StatCard icon={Server} label="Servers" value="—" />
          <StatCard icon={Users} label="Members" value="—" />
          <StatCard icon={Zap} label="Commands run" value="—" />
          <StatCard icon={Shield} label="Mod actions" value="—" />
        </div>

        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>Getting started</CardTitle>
              <Badge variant="outline">No backend yet</Badge>
            </div>
            <CardDescription>
              This dashboard is a static shell. Authentication and live data
              will be wired up later.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-xs leading-relaxed text-muted-foreground">
              Head to the{" "}
              <a
                href="/commands"
                className="font-medium text-primary underline-offset-4 hover:underline"
              >
                commands page
              </a>{" "}
              to explore everything MinjiBot can do.
            </p>
          </CardContent>
        </Card>
      </main>
    </div>
  )
}
