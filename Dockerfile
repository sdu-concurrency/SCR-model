FROM golang:1.26.3 AS build-backend-stage

WORKDIR /app
ADD ./back/migrations ./migrations
COPY ./back/go.mod ./back/go.sum ./
RUN go mod download

COPY ./back/*.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /backend

# stage1 as builder
FROM node:lts-alpine3.23 AS build-front-stage

WORKDIR /app

ENV VITE_APP_TITLE="Process model - Supply Chain Resilience"
ENV VITE_API_URL="/"
# Self-hosted company branding (build-time, all optional):
#   VITE_COMPANY_NAME - prefills and hides the survey's step-1 company field
#   VITE_LOGO_PATH    - header logo path, file placed in web/public (e.g. /logo.png)
#   VITE_LANDING_HTML - landing page HTML fragment path, file placed in web/public
#                       (e.g. /landing.html; plain HTML, Tailwind not available)
ARG VITE_COMPANY_NAME=""
ARG VITE_LOGO_PATH=""
ARG VITE_LANDING_HTML=""
ENV VITE_COMPANY_NAME=${VITE_COMPANY_NAME} \
    VITE_LOGO_PATH=${VITE_LOGO_PATH} \
    VITE_LANDING_HTML=${VITE_LANDING_HTML}
# Copy the package.json and install dependencies
COPY web/package*.json ./
RUN npm install
# Set production mode after installing devDependencies (needed for vue-tsc etc.)
ENV NODE_ENV=production
# Copy rest of the files
COPY web ./

# Build the project
RUN npm run build

# Deploy the application binary into a lean image
FROM alpine:3.23 AS build-release-stage
RUN apk update && \
    apk add --no-cache postfix && \
    rm -rf /var/cache/apk/*

WORKDIR /

COPY --from=build-backend-stage /backend /backend
COPY --from=build-front-stage /app/dist /pb_public

EXPOSE 8080

COPY --from=build-backend-stage /app/migrations /migrations

ENTRYPOINT ["/backend","serve","--http=0.0.0.0:8080"]
