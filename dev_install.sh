#!/usr/bin/env bash

yum update -y
yum install -y wget
yum install -y git

curl -sL https://rpm.nodesource.com/setup_12.x | bash -

echo "***************"
echo "Update node...."
echo "***************"
yum clean all && yum makecache fast
yum install -y gcc-c++ make
yum install -y nodejs
echo "*************"
echo "Node updated!"
echo "*************"
echo ""

echo "******************"
echo "Install golang...."
echo "******************"
wget https://dl.google.com/go/go1.14.linux-amd64.tar.gz
tar -xzf go1.14.linux-amd64.tar.gz
mv go /usr/local
export PATH=$PATH:/usr/local/go/bin
echo "******************************"
echo "Golang 1.14 has been installed!"
echo "******************************"
echo ""
