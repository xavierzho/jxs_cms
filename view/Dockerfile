# 第一阶段：构建 Vue2 项目
FROM node:20.19.5-alpine3.22 AS build

WORKDIR /app

# 只复制依赖相关文件，先安装依赖，避免频繁重新安装
COPY package*.json ./
RUN npm install --legacy-peer-deps

# 再复制源码，执行打包
COPY . .
RUN npm run build   # Vue2 默认 build 命令 (可能是 npm run build)

# 第二阶段：用 nginx 部署
FROM nginx:stable-alpine3.17-slim

# 删除默认配置，替换为自定义配置
RUN rm /etc/nginx/conf.d/default.conf
COPY data-cms.conf /etc/nginx/conf.d/data-cms.conf

# 拷贝前端构建产物到 nginx 容器
COPY --from=build /app/dist_prod /app/dist

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
