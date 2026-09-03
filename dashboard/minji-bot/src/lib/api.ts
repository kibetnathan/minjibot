// Base URL of the Echo API.
//
// In development the dashboard and API share the :5173 origin through Vite's
// /api proxy, so an empty value means "same origin". In production the frontend
// is hosted separately (e.g. Netlify) from the API (e.g. Railway), so build with
// VITE_API_URL set to the API origin, e.g. https://api.example.com.
export const API_BASE_URL = (
  import.meta.env.VITE_API_URL as string | undefined
)?.replace(/\/+$/, "") ?? ""

// apiUrl builds an absolute URL for an API path, honoring VITE_API_URL in
// production and falling back to the same-origin proxy in development.
export function apiUrl(path: string): string {
  const normalized = path.startsWith("/") ? path : `/${path}`
  return `${API_BASE_URL}${normalized}`
}