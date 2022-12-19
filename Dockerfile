
# Verificar essa issue antes de migrar para Yarn.
# https://github.com/yarnpkg/yarn/issues/77

FROM node:lts-alpine as dependencies
WORKDIR /site
COPY package.json package-lock.json ./
RUN apk add --no-cache libc6-compat

# about legacy-peer-deps flag: This flag isn't doing something new,rather it's telling NPM
# not to do something new, since NPM v7 now installs peerDependencies by default.
# See more about this flag here: https://stackoverflow.com/questions/66239691/what-does-npm-install-legacy-peer-deps-do-exactly-when-is-it-recommended-wh

RUN npm ci  --omit=dev --omit=optional --ignore-scripts --prefer-offline --legacy-peer-deps

FROM node:lts-alpine as builder
WORKDIR /site
COPY . .
COPY --from=dependencies /site/node_modules ./node_modules
RUN npm run build --legacy-peer-deps

FROM node:lts-alpine as runner
WORKDIR /site
COPY --from=builder /site/next.config.js /site/package.json  ./
COPY --from=builder /site/public ./public
COPY --from=builder /site/.next ./.next
COPY --from=dependencies /site/node_modules ./node_modules

# The production environment overrides the EXPOST with -p option.
EXPOSE 8080

# Next.js collects completely anonymous telemetry data about general usage.
# Learn more here: https://nextjs.org/telemetry
# Uncomment the following line in case you want to disable telemetry.
ENV NEXT_TELEMETRY_DISABLED 1

ENV NODE_ENV=production \
    DEFAULT_API_BASE_URL=$DEFAULT_API_BASE_URL \
    API_BASE_URL=$API_BASE_URL \
    PACKAGE_REPO_URL=$PACKAGE_REPO_URL \
    S3_REPO_URL=$S3_REPO_URL \
    ID_ANALYTICS=$ID_ANALYTICS \
    NEW_RELIC_APP_NAME=$NEW_RELIC_APP_NAME \
    NEW_RELIC_LICENSE_KEY=$NEW_RELIC_LICENSE_KEY

CMD ["npm", "run", "start"]