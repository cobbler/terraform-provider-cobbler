#! /bin/bash

# Requires xorriso (sudo apt-get install -y xorriso, sudo yum install xorriso -y, or sudo zypper install -y xorriso)
if [ -z "$1" ]
  then
    echo "No cobbler server url supplied"
fi

cobbler_image_tag=v4.0.0b5
iso_url=https://cdimage.ubuntu.com/ubuntu-legacy-server/releases/20.04/release/ubuntu-20.04.1-legacy-server-amd64.iso
iso_os=ubuntu
valid_iso_checksum=00a9d46306fbe9beb3581853a289490bc231c51f
iso_filename=$(echo ${iso_url##*/})
valid_extracted_iso_checksum=dd0b3148e1f071fb86aee4b0395fd63b

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

COBBLER_IMAGE_TAG=$cobbler_image_tag docker compose -f docker/compose.yml up -d

SERVER_URL=$1
printf "### Waiting for Cobbler to become available on ${SERVER_URL} \n\n"

attempt_counter=0
max_attempts=48

until $(curl --connect-timeout 1 --output /dev/null --silent ${SERVER_URL}); do
  if [ ${attempt_counter} -eq ${max_attempts} ];then
    echo "Max attempts reached"
    # Debug logs
    docker compose -f ./docker/compose.yml logs
    exit 1
  fi

  attempt_counter=$(($attempt_counter+1))
  sleep 5
done

# Cobbler 4.0.0 no longer ships the Python `cobbler` CLI (superseded by
# https://github.com/cobbler/cli), so seed the test distro via a direct XML-RPC call from the
# host instead of `cobbler import` - see import.py.
python3 docker/import.py

# Cobbler's Template.uri.path validator (cobbler/items/options/uri.py) requires the file to
# already exist under autoinstall_templates_dir - it never creates one - so the internal/template
# acceptance tests' backing files are touched into existence here.
docker compose -f docker/compose.yml exec cobbler touch \
  /var/lib/cobbler/templates/foo-resource-template-basic.j2 \
  /var/lib/cobbler/templates/foo-resource-template-change.j2 \
  /var/lib/cobbler/templates/foo-ds-template.j2 \
  /var/lib/cobbler/templates/foo-resource-template-import-not-found.j2

docker compose -f docker/compose.yml logs
