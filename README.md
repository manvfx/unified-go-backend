# Unified Go Backend API

## Overview

This is a Go-based backend API that uses MongoDB, Redis, and Elasticsearch. The project includes endpoints for user authentication, profile management, and access control with roles and permissions. 

## Features

- User registration, login, and email verification
- JWT-based authentication
- Role-based access control
- Redis for caching and storing verification codes
- Elasticsearch for logging
- Swagger for API documentation

## Prerequisites

- Docker
- Docker Compose

## Getting Started

### Local Development

1. **Clone the repository:**

    ```sh
    git clone https://github.com/yourusername/yourproject.git
    cd yourproject
    ```

2. **Create a `.env` file:**

    ```sh
    touch .env
    ```

    Add the following environment variables to the `.env` file (update with your actual values):

    ```env
    MONGO_URI=mongodb://mongo:27017/mdmdb
    REDIS_ADDR=redis:6379
    SMTP_HOST=smtp.mailtrap.io
    SMTP_PORT=587
    SMTP_USER=your_mailtrap_user
    SMTP_PASSWORD=your_mailtrap_password
    JWT_SECRET=your_jwt_secret
    ELASTICSEARCH_URL=http://elasticsearch:9200
    ```

3. **Build and run the Docker containers:**

    ```sh
    docker-compose up --build
    ```

    This will start the application along with MongoDB, Redis, and Elasticsearch services.

4. **Access the application:**

    - API: `http://localhost:8080`
    - Swagger Documentation: `http://localhost:8080/swagger/index.html`

## Deploying to Ubuntu VPS

### Step 1: Prepare Your Ubuntu VPS

1. **Update the system:**

    ```sh
    sudo apt update && sudo apt upgrade -y
    ```

2. **Install required packages:**

    ```sh
    sudo apt install curl git -y
    ```

### Step 2: Install Docker

1. **Install Docker:**

    ```sh
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    ```

2. **Add user to Docker group:**

    ```sh
    sudo usermod -aG docker ${USER}
    ```

    Log out and log back in to apply the group changes.

3. **Verify Docker installation:**

    ```sh
    docker --version
    ```

### Step 3: Install Docker Compose

1. **Download Docker Compose:**

    ```sh
    sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    ```

2. **Apply executable permissions:**

    ```sh
    sudo chmod +x /usr/local/bin/docker-compose
    ```

3. **Verify Docker Compose installation:**

    ```sh
    docker-compose --version
    ```

### Step 4: Prepare Your Project

1. **Clone your project:**

    ```sh
    git clone https://github.com/manvfx/unified-go-backend.git
    cd yourproject
    ```

2. **Create a `.env` file:**

    ```sh
    nano .env
    ```

    Add the following environment variables to the `.env` file (update with your actual values):

    ```env
    MONGO_URI=mongodb://mongo:27017/mdmdb
    REDIS_ADDR=redis:6379
    SMTP_HOST=smtp.mailtrap.io
    SMTP_PORT=587
    SMTP_USER=your_mailtrap_user
    SMTP_PASSWORD=your_mailtrap_password
    JWT_SECRET=your_jwt_secret
    ELASTICSEARCH_URL=http://elasticsearch:9200
    ```

### Step 5: Build and Run the Containers

1. **Build and start the containers:**

    ```sh
    docker-compose up --build -d
    ```

    The `-d` flag runs the containers in detached mode.

2. **Verify the deployment:**

    ```sh
    docker-compose ps
    ```

    Ensure all containers are running. Check the logs if there are any issues:

    ```sh
    docker-compose logs
    ```

### Step 6: Access Your Application

Your application should now be running on `http://<your-vps-ip>:8080`. Verify by visiting the URL or using a tool like Postman to test the endpoints.

### Optional: Set Up a Reverse Proxy with Nginx

For production deployment, it's a good practice to use a reverse proxy like Nginx to manage SSL certificates and handle requests.

1. **Install Nginx:**

    ```sh
    sudo apt install nginx -y
    ```

2. **Configure Nginx:**

    Create a new Nginx configuration file:

    ```sh
    sudo nano /etc/nginx/sites-available/unified-go-backend
    ```

    Add the following configuration (update `server_name` with your domain or VPS IP):

    ```nginx
    server {
        listen 80;
        server_name your_domain_or_ip;

        location / {
            proxy_pass http://localhost:8080;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
    ```

    Enable the configuration and restart Nginx:

    ```sh
    sudo ln -s /etc/nginx/sites-available/unified-go-backend /etc/nginx/sites-enabled/
    sudo nginx -t
    sudo systemctl restart nginx
    ```

3. **Optional: Set Up SSL with Let's Encrypt**

    If you have a domain, secure your application with SSL using Let's Encrypt.

    1. **Install Certbot:**

        ```sh
        sudo apt install certbot python3-certbot-nginx -y
        ```

    2. **Obtain and install SSL certificate:**

        ```sh
        sudo certbot --nginx -d your_domain
        ```

        Follow the prompts to complete the installation. Certbot will automatically configure Nginx to use the SSL certificate.

## Authors

- **Mohammad Amin Najafi** - [manvfx](https://github.com/manvfx)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
