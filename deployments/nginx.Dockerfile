#nginx/Dockerfile
FROM nginx:latest
EXPOSE 80
COPY ./configs/nginx.conf /etc/nginx/nginx.conf
ENTRYPOINT ["nginx","-g","daemon off;"]
