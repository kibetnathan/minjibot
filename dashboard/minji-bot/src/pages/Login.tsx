import { useState, type FormEvent } from "react"
import { Link } from "react-router-dom"
import { AuthLayout } from "@/components/auth/AuthLayout"
import { AuthField } from "@/components/auth/AuthField"
import { Button } from "@/components/ui/button"

export default function Login() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")

  function handleSubmit(event: FormEvent) {
    event.preventDefault()
    // No backend yet — submission is intentionally a no-op.
    void email
    void password
  }

  return (
    <AuthLayout
      title="Welcome back"
      description="Log in to the MinjiBot dashboard."
    >
      <form onSubmit={handleSubmit} className="grid gap-4">
        <AuthField
          label="Email"
          id="login-email"
          type="email"
          autoComplete="email"
          placeholder="you@example.com"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          required
        />
        <AuthField
          label="Password"
          id="login-password"
          type="password"
          autoComplete="current-password"
          placeholder="••••••••"
          value={password}
          onChange={(event) => setPassword(event.target.value)}
          required
        />
        <Button type="submit" className="h-9">
          Log in
        </Button>
      </form>

      <p className="mt-6 text-center text-xs text-muted-foreground">
        Don&apos;t have an account?{" "}
        <Link
          to="/signup"
          className="font-medium text-primary underline-offset-4 hover:underline"
        >
          Sign up
        </Link>
      </p>
    </AuthLayout>
  )
}
