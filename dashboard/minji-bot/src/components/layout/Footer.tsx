import { SiGithub } from "react-icons/si"

export function Footer() {
  return (
    <footer className="border-t border-border py-8">
      <div className="mx-auto max-w-7xl px-4 text-center text-sm text-muted-foreground sm:px-6 lg:px-8">
        <p>Built with Go, React, and a lot of coffee.</p>
        <p className="mt-2 flex items-center justify-center gap-1">
          <a
            href="https://github.com/kibetnathan/minjibot"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-1 underline transition-colors hover:text-foreground"
          >
            <SiGithub className="h-4 w-4" />
            GitHub Repository
          </a>
        </p>
      </div>
    </footer>
  )
}
