# product-reviews

## Overview

This project implements a RESTful API system for managing products and their reviews, including the calculation of
average product ratings. The system is designed to be performant, concurrency-safe, and extensible, with a focus on
meeting the assignment requirements while keeping the architecture simple and maintainable.

The tech stack includes:

- Golang: Chosen as the primary language per the preference in the assignment, leveraging its concurrency features and
  simplicity.
- PostgreSQL: A relational database for persistent storage of products and reviews.
- Redis: An in-memory key-value store for caching and concurrency-safe average rating calculations.
- NATS: A lightweight messaging broker for publishing review-related events to other services.

Below, I outline the thought process behind the implementation, key design decisions, and trade-offs made during
development.

## Design Decisions

### Technology Choices

#### PostgreSQL

- **Why?**: PostgreSQL is a robust, general-purpose relational database that supports ACID transactions, making it a
  reliable choice for storing product and review data. Its strong support for relational queries also simplifies
  retrieving product details and review lists.
- **Trade-offs**: While a NoSQL database like MongoDB could offer flexibility for unstructured data, PostgreSQL’s
  relational model provides better consistency and integrity for this use case, where relationships between products and
  reviews are straightforward.

#### Redis

- **Why?**: Redis was selected for its low-latency key-value storage, ideal for caching average product ratings, reviews
  and
  ensuring concurrency safety. The SETNX (Set If Not Exists) command provides a lightweight locking mechanism to prevent
  race conditions during rating updates.
- **Trade-offs**: Using Redis adds a second data store, increasing complexity compared to relying solely on PostgreSQL.
  However, the performance gain for frequent rating calculations justifies this choice over recalculating averages from
  the database on every request.

#### NATS

- **Why?**: NATS is a simple, fast, and lightweight messaging broker that effectively demonstrates event publishing (
  e.g., when a review is added, modified, or deleted). Its ease of setup and pub/sub model make it suitable for
  notifying other services.
- **Trade-offs**: More robust brokers like RabbitMQ or Kafka offer advanced features (e.g., message persistence,
  delivery guarantees), but NATS was chosen for its simplicity and speed, aligning with the assignment’s scope of
  showcasing basic event publishing.

### API Design

The RESTful design follows standard conventions for clarity and interoperability. Separating product and review details
in the API responses optimizes payload size and aligns with the requirement to return only the average rating with
products.

### Average Rating Calculation

- **Implementation**: The average rating is calculated on-demand when a GET product API request is received. A SQL query
  computes the average from PostgreSQL, and the result is cached in Redis with a 5-minute TTL (time-to-live). Subsequent
  requests within the TTL reuse the cached value, reducing database load. Redis’s SETNX (Set If Not Exists) command is
  used to ensure concurrency safety when writing the computed average to the cache, preventing race conditions if
  multiple requests trigger the calculation simultaneously.
- **Why On-Demand Calculation?**: Calculating the average only when requested avoids unnecessary updates to Redis on
  every review change, simplifying the system. The 5-minute TTL ensures the cache stays reasonably fresh while
  minimizing database queries.
- **Trade-offs**: This approach trades off real-time accuracy for simplicity and reduced write overhead. If a review is
  added or modified, the cached average may be stale for up to 5 minutes. For more frequent updates, a shorter TTL or
  incremental updates triggered by review changes could be considered, though this would increase complexity and Redis
  writes. The use of SETNX ensures safe caching under concurrency.

### Event Publishing Simplification

- **Implementation**: Review events (create, update, delete) are published to NATS immediately after the database
  operation succeeds, without transactional guarantees between the database write and message publish.
- **Why?**: This simplifies the implementation and improves performance by avoiding complex coordination between
  PostgreSQL and NATS.
- **Trade-offs**: This approach introduces a risk of inconsistency—if the NATS publish fails after a database write,
  other services might miss the update. For the scope of this assignment, this trade-off prioritizes performance and
  simplicity over reliability.

### Key Trade-offs and Potential Improvements

#### Simplified Event Publishing

- **Trade-off**: By publishing events to NATS without transactional coupling to PostgreSQL writes, there’s a small
  chance of missed notifications if the service crashes between the database update and the NATS publish.
- **Improvement**: Two options could enhance reliability:
    - **Transactional Approach**: Use a two-phase commit or an Outbox pattern to ensure database writes and NATS
      publishes occur atomically. This increases complexity and latency but guarantees consistency.
    - **Separate Event Table**: Store events in a PostgreSQL table and use a cron job or background worker to publish
      them to NATS. This ensures no events are lost, even if the service crashes, at the cost of delayed notifications.

## Prerequisites

Before you get started, ensure you have the following installed on your system:

- **Docker**: Containerization platform for running the application, used to create and manage containers for consistent
  development and deployment.
- **Tasks**: Tool for launching and automating commands in your development workflow, simplifies running repetitive
  tasks.
    - `brew install go-task/tap/go-task`
- **Linter**: A static analysis tool for identifying and fixing programming errors, bugs, stylistic errors, and
  suspicious
  constructs in Go code.
    - `brew install golangci-lint`
- **NATS CLI** (Optionally): Command-line interface for interacting with NATS servers, used for managing and testing
  messaging systems.
    - `go install github.com/nats-io/natscli/nats@latest`

## Getting Started

1. Clone the repository
2. Navigate to the project directory
3. Run `task start` to initialize the development environment

## Useful Commands

### Viewing NATS Messages

To view the messages published in the `reviews-stream`, you can use the following command:

```bash
nats stream view reviews-stream --server=localhost:4222
```