import { SiDiscord } from "react-icons/si"
import { Button } from "@/components/ui/button"
import { apiUrl } from "@/lib/api"
import { cn } from "@/lib/utils"

type DiscordLoginButtonProps = {
  className?: string
  label?: string
}

export function DiscordLoginButton({
  className,
  label = "Continue with Discord",
}: DiscordLoginButtonProps) {
  return (
    <Button
      type="button"
      className={cn("gap-2", className)}
      onClick={() => {
        window.location.href = apiUrl("/api/auth/discord")
      }}
    >
      <SiDiscord className="size-4" />
      {label}
    </Button>
  )
}
