#! /bin/bash

# Requires xorriso (sudo apt-get install -y xorriso, sudo yum install xorriso -y, or sudo zypper install -y xorriso)
if [ -z "$1" ]
  then
    echo "No cobbler server url supplied"
fi

cobbler_commit=b3b6ee2391c5a1fb89f7796e4d9dc6538617485a # master as of 4/2/2022
cobbler_branch=master
iso_url=https://cdimage.ubuntu.com/ubuntu-legacy-server/releases/20.04/release/ubuntu-20.04.1-legacy-server-amd64.iso
iso_os=ubuntu
valid_iso_checksum=00a9d46306fbe9beb3581853a289490bc231c51f
iso_filename=$(echo ${iso_url##*/})
valid_extracted_iso_checksum=dd0b3148e1f071fb86aee4b0395fd63b
valid_git_checksum=6c9511b26946dd3f1f072b9f40eaeccf  # master as of 4/2/2022

[ -d "./docker/cobbler_source" ] && git_checksum=$(find ./docker/cobbler_source/ -type f -exec md5sum {} \; | sort -k 2 | md5sum | awk '{print $1}')
if [ -d "./docker/cobbler_source" ] && [ $git_checksum == $valid_git_checksum ]; then
  echo "Cobbler code already cloned and the correct version is checked out"
else
  rm -rf ./docker/cobbler_source
  git clone --shallow-since="2021-09-01" https://github.com/cobbler/cobbler.git -b $cobbler_branch docker/cobbler_source
  cd ./docker/cobbler_source
  printf "Changing to version of Cobbler being tested.\n\n"
  git checkout $cobbler_commit > /dev/null 2>&1
  rm -rf .git  # remove .git dir so the checksum is consistent
  cd -
fi

echo $(pwd)
if [ -f "$iso_filename" ] && [ $(sha1sum $iso_filename | awk '{print $1}') == "$valid_iso_checksum" ]; then
  echo "ISO already downloaded"
else
  rm $iso_filename
  wget $iso_url
fi

extracted_iso_checksum=$(find extracted_iso_image -type f -exec md5sum {} \; | sort -k 2 | md5sum | awk '{print $1}')
if [ -d "extracted_iso_image" ] && [ $extracted_iso_checksum == $valid_extracted_iso_checksum ]; then
   echo "ISO already extracted"
else
   xorriso -osirrox on -indev $iso_filename -extract / extracted_iso_image
fi

docker build -f ./docker/cobbler_source/docker/develop/develop.dockerfile -t cobbler-dev .
docker-compose -f docker/docker-compose.yml up -d

SERVER_URL=$1
printf "### Waiting for Cobbler to become available on ${SERVER_URL} \n\n"

attempt_counter=0
max_attempts=48

until $(curl --connect-timeout 1 --output /dev/null --silent ${SERVER_URL}); do
  if [ ${attempt_counter} -eq ${max_attempts} ];then
    echo "Max attempts reached"
    # Debug logs
    docker-compose -f ./docker/docker-compose.yml logs
    exit 1
  fi

  attempt_counter=$(($attempt_counter+1))
  sleep 5
done

docker-compose -f docker/docker-compose.yml logs
