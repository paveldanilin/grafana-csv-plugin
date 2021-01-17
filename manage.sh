#!/usr/bin/env bash

function docker_install()
{
  echo "Nothing to do"

  # echo "MANAGE: Install docker..."

  # yum update -y
  # curl http://repo.net.billing.ru/scripts/add-docker-ce-repo.sh | sudo sh

  # yum install -y epel-release
  # yum install -y yum-utils device-mapper-persistent-data lvm2 python36 docker-ce docker-ce-cli containerd.io docker-compose
  # systemctl start docker
  # systemctl enable docker

  # echo "MANAGE: docker has been installed and started!"
}

function container_build()
{
  echo "MANAGE: Build \"grafana\" container..."

  docker rm grafana_csv-plugin
  docker build --no-cache --tag=grafana:csv-plugin .
}

function container_up()
{
  echo "MANAGE: Up \"grafana\" container..."

  docker run -d -p 3000:3000 --name=grafana_csv-plugin grafana:csv-plugin
}

function container_down()
{
  echo "MANAGE: Stop \"grafana\" container..."

  docker stop grafana_csv-plugin
}

main()
{
  cmd="$1"
  case $cmd in
    build)
      container_build
      shift # past argument
    ;;
    up)
      container_up
      shift # past argument
    ;;
    install)
      docker_install
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
