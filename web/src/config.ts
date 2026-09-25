export const config = {
  baseUrl: import.meta.env.BASE_URL,
  apiBaseUrl: import.meta.env.VITE_API_URL || 'http://127.0.0.1:8090',
  companyName: import.meta.env.VITE_COMPANY_NAME || '',
  logoPath: import.meta.env.VITE_LOGO_PATH || '',
  landingHtmlPath: import.meta.env.VITE_LANDING_HTML || ''
}
