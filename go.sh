#!/bin/sh
envsubst '$BACKEND_SERVER' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.step1.conf
envsubst '$FRONTEND_SERVER' < /etc/nginx/nginx.step1.conf > /etc/nginx/nginx.conf
/usr/sbin/nginx -g 'daemon off;'