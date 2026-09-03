import { AuthLayout } from "@/components/auth/AuthLayout"
import { DiscordLoginButton } from "@/components/auth/DiscordLoginButton"
import { Link } from "react-router-dom"

export default function SignUp() {
  return (
    <AuthLayout
      title="Create your account"
      description="Start managing MinjiBot with your own dashboard."
    >
      <DiscordLoginButton label="Sign up with Discord" className="w-full" />

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
