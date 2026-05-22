# Frontend — Supply Chain Resilience

Vue 3 frontend for the Supply Chain Resilience assessment tool.

## Requirements

- Node.js ≥ 20

## Install

```bash
npm install
```

## Run in development mode

The frontend calls the PocketBase backend using the `VITE_API_URL` environment variable. Set it before starting the dev server:

```bash
VITE_API_URL=http://127.0.0.1:8090 npm run dev
```

Or create a `web/.env.local` file:

```
VITE_API_URL=http://127.0.0.1:8090
```

Then run:

```bash
npm run dev
```

## Build for production

```bash
npm run build
```

Output is written to `dist/`.

## Survey content

Survey forms are defined in JSON files under `src/surveys/`:

- `form_dk.json` — Danish survey
- `form_en.json` — English survey

The schema follows the [FormKit schema](https://formkit.com/essentials/schema) format. Changes to these files are picked up automatically by the dev server.

## Supported languages

English (`en`) and Danish (`da`). Translation strings live in `src/i18n/`.
