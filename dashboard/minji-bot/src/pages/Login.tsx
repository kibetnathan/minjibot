import { AuthLayout } from "@/components/auth/AuthLayout"
import { DiscordLoginButton } from "@/components/auth/DiscordLoginButton"
import { Link } from "react-router-dom"

export default function Login() {
  return (
    <AuthLayout
      title="Welcome back"
      description="Log in to the MinjiBot dashboard."
    >
      <DiscordLoginButton className="w-full" />

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
