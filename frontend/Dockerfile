# 開発用 Dockerfile
FROM node:24-alpine

WORKDIR /app

# 依存関係をコピー
COPY package.json package-lock.json* ./
RUN npm ci

# ソースコードはボリュームマウントで共有
COPY . .

EXPOSE 3000

# 開発サーバー起動
CMD ["npm", "run", "dev"]