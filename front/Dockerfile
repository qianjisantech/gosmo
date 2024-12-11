# nginx镜像
FROM registry.cn-shanghai.aliyuncs.com/vue-gin-devops/nginx:latest

# 移除nginx容器的default.conf文件、nginx配置文件
RUN rm /etc/nginx/nginx.conf
#RUN rm /etc/nginx/conf.d/default.conf

# 把主机的nginx.conf文件复制到nginx容器的/etc/nginx文件夹下
COPY nginx.conf /etc/nginx/
# 拷贝前端vue项目打包后生成的文件到nginx下运行
COPY dist/  /usr/share/nginx/html/

#COPY certificate/     /etc/nginx/ssl/

# 暴露9090端口
EXPOSE 80
#EXPOSE 443
# 使用daemon off的方式将nginx运行在前台保证镜像不至于退出
CMD ["nginx", "-g", "daemon off;"]