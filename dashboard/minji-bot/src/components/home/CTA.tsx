export function CTA() {
  return (
    <section className="mx-auto max-w-7xl bg-muted/30 px-4 py-20 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-2xl text-center">
        <h2 className="mb-4 font-heading text-3xl font-bold text-foreground sm:text-4xl">
          Ready to get started?
        </h2>
        <p className="mb-8 text-muted-foreground">
          Clone the repo, configure your environment, and spin up the database.
        </p>
        <div className="overflow-x-auto rounded-xl border border-border bg-card p-6 text-left font-mono text-sm">
          <pre className="text-muted-foreground">
            <code>{`# 1. Clone and install
git clone https://github.com/kibetnathan/minjibot
cd minjibot

# 2. Configure environment (copy .env.example to .env)
# Set POSTGRES_*, DB_URL, GOOSE_*, TESTING_DB

# 3. Start database
make docker/up

# 4. Run migrations
make goose-migrate-up

# 5. Generate database queries
sqlc generate

# 6. Start API server
make run`}</code>
          </pre>
        </div>
      </div>
    </section>
  )
}
