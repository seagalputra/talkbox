FROM node:16-alpine AS deps

RUN apk add --no-cache libc6-compat
WORKDIR /app

COPY package.json yarn.lock* ./
RUN yarn --frozen-lockfile

# Build Frontend web
FROM node:16-alpine AS frontend-builder
WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

ARG NEXT_PUBLIC_API_BASE_URL
ARG NEXT_PUBLIC_WS_BASE_URL

RUN yarn build

# Build Backend API
FROM golang:alpine AS backend-builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY .env* .
COPY api ./api
COPY main.go .
COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go build -o talkbox

FROM node:16-alpine AS runner

RUN npm install -g pm2

WORKDIR /app

ENV NODE_ENV production

RUN addgroup --system --gid 1001 app
RUN adduser --system --uid 1001 talkbox

COPY --from=frontend-builder /app/public ./public

COPY --from=frontend-builder --chown=talkbox:app /app/.next/standalone ./
COPY --from=frontend-builder --chown=talkbox:app /app/.next/static ./.next/static
COPY --from=frontend-builder --chown=talkbox:app /app/services.config.js ./

COPY --from=backend-builder --chown=talkbox:app /app/talkbox ./
COPY --from=backend-builder --chown=talkbox:app /app/env* ./

USER talkbox

EXPOSE 8080-8090

ENV PORT 8080

CMD ["pm2-runtime", "start", "services.config.js"]