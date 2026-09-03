import { Search } from "lucide-react"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"

type CommandSearchProps = {
  query: string
  onQueryChange: (query: string) => void
}

export function CommandSearch({ query, onQueryChange }: CommandSearchProps) {
  return (
    <Card className="mb-8">
      <CardContent className="py-3">
        <div className="relative">
          <Search className="pointer-events-none absolute top-1/2 left-2 size-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            type="search"
            value={query}
            onChange={(event) => onQueryChange(event.target.value)}
            placeholder="Search commands by name or description…"
            className="pl-8"
          />
        </div>
      </CardContent>
    </Card>
  )
}
