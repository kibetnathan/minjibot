import type { CommandCategory } from "@/data/commands"
import { CommandCard } from "@/components/commands/CommandCard"

type CategoryGridProps = {
  categories: CommandCategory[]
}

export function CategoryGrid({ categories }: CategoryGridProps) {
  if (categories.length === 0) {
    return (
      <p className="py-12 text-center text-sm text-muted-foreground">
        No commands match your search.
      </p>
    )
  }

  return (
    <div className="grid gap-6 md:grid-cols-2">
      {categories.map((category) => (
        <CommandCard key={category.id} category={category} />
      ))}
    </div>
  )
}
