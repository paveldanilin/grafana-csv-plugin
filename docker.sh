#!/usr/bin/env bash

function docker_install()
{
  echo "Install docker..."

  yum update -y
  curl http://repo.net.billing.ru/scripts/add-docker-ce-repo.sh | sudo sh

  yum install -y epel-release
  yum install -y yum-utils device-mapper-persistent-data lvm2 python36 docker-ce docker-ce-cli containerd.io docker-compose wget
  systemctl start docker
  systemctl enable docker
}

function container_build_and_up()
{
  echo "Build container and up..."

  docker build --no-cache --tag=grafana .
  docker run -d -p 3000:3000 --name=grafana grafana
}

function container_down()
{
  echo "Stop container..."

  docker stop grafana
}

main()
{
  cmd="$1"
  case $cmd in
    up)
      container_build_and_up
      shift # past argument
    ;;
    iup)
      docker_install
      container_build_and_up
      shift # past argument
    ;;
    down)
      container_down
      shift # past argument
    ;;
    *)    # unknown option
      echo "ERROR: unknown command \"$cmd\""
    ;;
  esac
}

main $1
