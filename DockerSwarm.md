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
