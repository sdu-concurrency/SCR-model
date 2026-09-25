/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_APP_TITLE: string
  readonly VITE_API_URL: string
  /** Self-hosted branding: prefills+hides the survey company field */
  readonly VITE_COMPANY_NAME: string
  /** Self-hosted branding: header logo path (file in web/public, e.g. "/logo.png") */
  readonly VITE_LOGO_PATH: string
  /** Self-hosted branding: landing page HTML fragment path (file in web/public, e.g. "/landing.html") */
  readonly VITE_LANDING_HTML: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
