# How to setup Docker Swarm

1. Build the docker image for all the services

  For example
  ```bash
  docker build -f Dockerfile -t suhelkapadia/listener-service:1.0.0 .
  ```

2. Push the docker image to the docker hub

  For example
  ```bash
  docker push suhelkapadia/listener-service:1.0.0
  ```

  Make sure you are logged in to the docker hub using `docker login` command

3. Create a `swarm.yml` file similar to the `docker-compose.yml` file, but alter the ports and the images to the ones pushed on the docker hub.

4. Initialize the swarm

  ```bash
  docker swarm init
  ```

5. Deploy the stack

  ```bash
  docker stack deploy -c swarm.yml myapp
  ```

6. Check the services

  ```bash
  docker service ls
  ```

## Scale the services

1. Scale the services

  ```bash
  docker service scale myapp_listener-service=3
  ```

2. Check the services

  ```bash
  docker service ls
  ```

3. Reduce to 2

    ```bash
    docker service scale myapp_listener-service=2
    ```

## Update the services

1. Build the docker image for the updated service (logger service let's say)

    For example
    ```bash
    docker build -f Dockerfile -t suhelkapadia/logger-service:1.0.1 .
    ```

2. Push the docker image to the docker hub

    For example
    ```bash
    docker push suhelkapadia/logger-service:1.0.1
    ```

3. Scale the service to 2, this is done to ensure zero downtime

    ```bash
    docker service scale myapp_logger-service=2
    ```

3. Update the service

    ```bash
    docker service update --image suhelkapadia/logger-service:1.0.1 myapp_logger-service
    ```

4. You could also downgrade the service if the update fails

    ```bash
    docker service update --image suhelkapadia/logger-service:1.0.0 myapp_logger-service
    ```

5. Scale the service back to 1

    ```bash
    docker service scale myapp_logger-service=1
    ```

## Stopping swarm

First method is to scale all services to 0.

But to remove the entire swarm, run the following command

```bash
docker stack rm myapp
```

And then leave the swarm

```bash
docker swarm leave --force
```
