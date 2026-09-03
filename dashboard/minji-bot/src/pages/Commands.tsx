import { useMemo, useState } from "react"
import { Navbar } from "@/components/layout/Navbar"
import { Footer } from "@/components/layout/Footer"
import { CommandsHero } from "@/components/commands/CommandsHero"
import { CommandSearch } from "@/components/commands/CommandSearch"
import { CategoryGrid } from "@/components/commands/CategoryGrid"
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { commandCategories, type CommandCategory } from "@/data/commands"

const ALL_TAB = "all"

function filterCategories(
  categories: CommandCategory[],
  query: string,
  activeTab: string
): CommandCategory[] {
  const trimmed = query.trim().toLowerCase()

  return categories.reduce<CommandCategory[]>((acc, category) => {
    if (activeTab !== ALL_TAB && activeTab !== category.id) {
      return acc
    }

    const commands =
      trimmed === ""
        ? category.commands
        : category.commands.filter(
            (command) =>
              command.name.toLowerCase().includes(trimmed) ||
              command.description.toLowerCase().includes(trimmed)
          )

    if (commands.length === 0) {
      return acc
    }

    acc.push({ ...category, commands })
    return acc
  }, [])
}

export default function Commands() {
  const [query, setQuery] = useState("")
  const [activeTab, setActiveTab] = useState(ALL_TAB)

  const totalCommands = useMemo(
    () =>
      commandCategories.reduce(
        (sum, category) => sum + category.commands.length,
        0
      ),
    []
  )

  const visibleCategories = useMemo(
    () => filterCategories(commandCategories, query, activeTab),
    [query, activeTab]
  )

  return (
    <div className="min-h-screen bg-background font-sans antialiased">
      <Navbar />
      <main>
        <CommandsHero total={totalCommands} />
        <div className="mx-auto max-w-7xl px-4 pb-20 sm:px-6 lg:px-8">
          <CommandSearch query={query} onQueryChange={setQuery} />
          <Tabs
            value={activeTab}
            onValueChange={(value) => {
              const next = Array.isArray(value) ? value[0] : value
              if (typeof next === "string") {
                setActiveTab(next)
              }
            }}
          >
            <ScrollArea className="mb-8 max-w-full">
              <TabsList className="h-9 w-max">
                <TabsTrigger value={ALL_TAB}>All</TabsTrigger>
                {commandCategories.map((category) => (
                  <TabsTrigger key={category.id} value={category.id}>
                    {category.name}
                  </TabsTrigger>
                ))}
              </TabsList>
              <ScrollBar orientation="horizontal" />
            </ScrollArea>
          </Tabs>
          <CategoryGrid categories={visibleCategories} />
        </div>
      </main>
      <Footer />
    </div>
  )
}
