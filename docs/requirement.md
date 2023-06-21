Create an application that could connect to any PostgreSQL database and provide an API to query slowest database queries.
Expect a high load of calls to come in and add a cache layer within the API to handle the load.

Add an additional API layer to be used for demo purposes that would CRUD dummy data entities into the DB so that the DB Query API would have something to work with.

You must use Fiber.

Requirements:
- Support pagination
- Support filtering by SELECT, INSERT,UPDATE, DELETE
- Support sorting by time spent
- Cache should be used for frequently accessed calls (with proper invalidation)
- Project & architecture documentation
- 80%+ test coverage
