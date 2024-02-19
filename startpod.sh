podman pod create --name avscan-pod -p 8080:8080 -p 3310:3310
podman run -d --pod avscan-pod --name avscan-server --rm clamav/clamav
podman build -t avscan-api:latest .
podman run -d --pod avscan-pod --name avscan-api --rm avscan-api:latest
