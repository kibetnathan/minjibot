import { useEffect, useState } from "react"
import { apiUrl } from "@/lib/api"

export type CurrentUser = {
  id: string
  email: string
  is_admin: boolean
}

type State =
  | { status: "loading" }
  | { status: "authenticated"; user: CurrentUser }
  | { status: "unauthenticated" }

// useCurrentUser fetches the currently authenticated user from /api/auth/me.
// It sends the session cookie with credentials: "include" so cross-origin calls
// (production dashboard vs API) carry the cookie. Refetches when `enabled`.
export function useCurrentUser(enabled = true): State {
  const [state, setState] = useState<State>({ status: "loading" })

  useEffect(() => {
    if (!enabled) {
      setState({ status: "unauthenticated" })
      return
    }

    let cancelled = false

    fetch(apiUrl("/api/auth/me"), {
      credentials: "include",
      headers: { Accept: "application/json" },
    })
      .then(async (res) => {
        if (res.status === 401) {
          if (!cancelled) setState({ status: "unauthenticated" })
          return
        }
        if (!res.ok) {
          throw new Error(`me request failed: ${res.status}`)
        }
        const user = (await res.json()) as CurrentUser
        if (!cancelled) setState({ status: "authenticated", user })
      })
      .catch((err) => {
        console.error("useCurrentUser:", err)
        if (!cancelled) setState({ status: "unauthenticated" })
      })

    return () => {
      cancelled = true
    }
  }, [enabled])

  return state
}