# Instagram Clone

This project is a simplified, microservices-based clone of Instagram. It is built using Go and utilizes a variety of technologies to handle different aspects of the application, such as user management, posts, interactions, and notifications.

## Microservices

The application is divided into the following microservices:

*   **Core Services:** Handles user registration, authentication, and user profile management.
*   **Post Services:** Manages the creation, retrieval, and storage of posts, including image handling.
*   **Interaction Services:** Manages user interactions with posts, such as likes, comments, and timelines.
*   **Notification Services:** Sends notifications to users about various events, such as new followers, likes, and comments.

## Technologies Used

*   **Go:** The primary programming language for all microservices.
*   **Gin:** A popular Go web framework used for building the APIs.
*   **PostgreSQL:** The relational database used for storing data for each microservice.
*   **RabbitMQ:** A message broker used for asynchronous communication between microservices.
*   **MinIO:** An object storage service used for storing images.
*   **Docker:** Used for containerizing the application and its dependencies.
*   **GORM:** The ORM library for Go.

## Getting Started

To get the project up and running locally, follow these steps:

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/instagram-clone.git
    cd instagram-clone
    ```
2.  **Set up environment variables:**
    Each microservice requires a `.env` file with the necessary environment variables. You can create a `.env` file in the root of each microservice directory (e.g., `core-services/.env`) and add the required variables.
3.  **Run the infrastructure services using Docker Compose:**

    ```bash
    cd instagram-clone-docker
    docker-compose up -d
    ```

    This will start the PostgreSQL databases, RabbitMQ, and MinIO services.

4.  **Run the microservices:**

    Open a new terminal for each microservice and run the following commands:

    **Core Services:**

    ```bash
    cd core-services
    go run main.go
    ```

    **Post Services:**

    ```bash
    cd post-services
    go run main.go
    ```

    **Interaction Services:**

    ```bash
    cd interaction-services
    go run main.go
    ```

    **Notification Services:**

    ```bash
    cd notification-services
    go run main.go
    ```

## API Endpoints

Each microservice exposes a set of API endpoints for interacting with its resources. Here is a summary of the available endpoints:

### Core Services

*   `POST /register`: Register a new user.
*   `POST /login`: Log in a user.
*   `GET /users/:id`: Get user profile information.
*   `POST /follow`: Follow a user.

### Post Services

*   `POST /posts`: Create a new post.
*   `GET /posts/:id`: Get a specific post.
*   `GET /users/:id/posts`: Get all posts from a user.

### Interaction Services

*   `POST /posts/:id/like`: Like a post.
*   `POST /posts/:id/comment`: Comment on a post.
*   `GET /timeline`: Get the user's timeline.

### Notification Services

*   `GET /notifications`: Get all notifications for the user.

## Directory Structure

The project is organized into the following directories:

```
instagram-clone/
├── core-services/          # Core microservice for user management
├── interaction-services/   # Interaction microservice for likes, comments, etc.
├── notification-services/  # Notification microservice for user notifications
├── post-services/          # Post microservice for creating and managing posts
└── instagram-clone-docker/ # Docker-related files
```

