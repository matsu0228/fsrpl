FROM google/cloud-sdk:alpine

WORKDIR /usr/src/app
RUN apk add --update --no-cache openjdk8 \
    && gcloud components update \
    && gcloud components install cloud-firestore-emulator beta --quiet