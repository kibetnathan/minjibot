export function Navbar() {
  return (
    <header className="border-b border-border">
      <div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
        <div className="flex items-center gap-2">
          <img src="beans.svg" alt="Minji" className="h-[2rem]" />
          <span className="font-heading text-xl font-bold">MinjiBot</span>
        </div>
        <nav className="hidden items-center gap-6 text-sm text-muted-foreground md:flex">
          <a
            href="#features"
            className="transition-colors hover:text-foreground"
          >
            Features
          </a>
          <a href="#tech" className="transition-colors hover:text-foreground">
            Tech Stack
          </a>
        </nav>
      </div>
    </header>
  )
}
