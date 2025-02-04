CQRS Architecture with Go Microservices and Docker
======

Overview
--------
This project implements a CQRS (Command Query Responsibility Segregation) architecture using multiple microservices written in Go. It includes services for handling feed events, querying data, and pushing updates. The system supports PostgreSQL for relational data, NATS for messaging, Elasticsearch for search capabilities, and NGINX as a reverse proxy.

Services and Components
-----------------------
• Feed Service – Handles creation of feed items and publishes events.
• Query Service – Processes queries and interacts with both PostgreSQL and Elasticsearch for data retrieval.
• Pusher Service – Listens for events (via NATS) and pushes updates to clients.
• PostgreSQL – Serves as the core relational database.
• NATS – Provides messaging/stored streaming capabilities across services.
• Elasticsearch – Provides search abilities for query-related operations.
• NGINX – Acts as a reverse proxy, directing incoming HTTP requests to the appropriate service.

Getting Started
---------------
Prerequisites:
• Docker and Docker Compose installed.
• (Optionally) GNU Make to use predefined commands in the Makefile.

Building and Running the Project:
1. Open a terminal and navigate to the project directory.
2. Run the command below to build and start all services in detached mode:

       make run

   Alternatively, use Docker Compose directly:
       docker-compose up -d --build

3. Docker Compose will build each service based on the configurations:
   - The Dockerfile in the project builds the Go binaries.
   - The docker-compose.yml file orchestrates all parts of the application (databases, message queues, services, and NGINX).

Service Configurations:
• Environment variables (e.g., POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, NATS_ADDRESS, and ELASTICSEARCH_ADDRESS) are defined in docker-compose.yml and passed to each service container.
• The services are interconnected using Docker Compose dependencies ensuring correct startup order.

Development and Contribution
----------------------------
• To modify the Go code, update the source files in their respective directories (e.g., models, feed-service, query-service, pusher-service).
• After making changes, run:
       make run
to rebuild the images and deploy updated containers.

Dependencies:
The project uses several external libraries including:
• Gorilla Mux – HTTP router for building web applications.
• Gorilla Websocket – For real-time communication.
• NATS Go – For messaging with the NATS server.
• PostgreSQL driver – For database connectivity.
• Elasticsearch Go client – For interfacing with Elasticsearch.

Contact
-------
For questions or further information, contact the project maintainer.

Happy Coding!
