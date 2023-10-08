#!/bin/bash

# THIS SCRIPT IS RUN AS THE ROOT USER
# ANY FILE CREATED WILL BE OWNED BY ROOT

INSTANCE_ID=$(curl -s http://169.254.169.254/latest/meta-data/instance-id)
REGION=$(curl -s http://169.254.169.254/latest/meta-data/placement/region)

get_tag() {
    TAG="$1"

    aws ec2 describe-tags \
        --filters "Name=resource-id,Values=$INSTANCE_ID" \
        --region $REGION \
        --query 'Tags[].{Key:Key,Value:Value}' \
        --output text \
        | awk "/${TAG}/ {print \$2}"
}

# install useful packages
apt update -y

apt install -y \
    awscli

# install and configure go
cd ~
curl -OL https://go.dev/dl/go1.21.2.linux-amd64.tar.gz
tar -C /usr/local -xvf go1.21.2.linux-amd64.tar.gz
rm go1.21.2.linux-amd64.tar.gz

# initialize data directory
mkdir /data

# set node type
get_tag NodeType > /var/node_type

cd $HOME

