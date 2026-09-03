import { useState, type FormEvent } from "react"
import { Link } from "react-router-dom"
import { AuthLayout } from "@/components/auth/AuthLayout"
import { AuthField } from "@/components/auth/AuthField"
import { Button } from "@/components/ui/button"

export default function SignUp() {
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")

  function handleSubmit(event: FormEvent) {
    event.preventDefault()
    // No backend yet — submission is intentionally a no-op.
    void name
    void email
    void password
  }

  return (
    <AuthLayout
      title="Create your account"
      description="Start managing MinjiBot with your own dashboard."
    >
      <form onSubmit={handleSubmit} className="grid gap-4">
        <AuthField
          label="Name"
          id="signup-name"
          type="text"
          autoComplete="name"
          placeholder="Minji"
          value={name}
          onChange={(event) => setName(event.target.value)}
          required
        />
        <AuthField
          label="Email"
          id="signup-email"
          type="email"
          autoComplete="email"
          placeholder="you@example.com"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          required
        />
        <AuthField
          label="Password"
          id="signup-password"
          type="password"
          autoComplete="new-password"
          placeholder="Create a password"
          value={password}
          onChange={(event) => setPassword(event.target.value)}
          minLength={8}
          required
        />
        <Button type="submit" className="h-9">
          Sign up
        </Button>
      </form>

      <p className="mt-6 text-center text-xs text-muted-foreground">
        Already have an account?{" "}
        <Link
          to="/login"
          className="font-medium text-primary underline-offset-4 hover:underline"
        >
          Log in
        </Link>
      </p>
    </AuthLayout>
  )
}
