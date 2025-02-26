Backend System for Order Processing.

**Problem Statement:**
Build a backend system to manage and process orders in an e-commerce platform. The system should:

**Core Functionality:**
Provide a RESTful API to accept orders with fields such as - user_id, order_id, item_ids, and total_amount.
Simulate asynchronous order processing using an in-memory queue (e.g., Python queue.Queue or equivalent).
Provide an API to check the status of orders - Pending, Processing and Completed

**Implement an API to fetch key metrics, including:**
Total number of orders processed.
Average processing time for orders.
Count of orders in each status -
Pending, Processing, Completed

**Constraints:**
Database: Use SQLite/PostgreSQL/MySQL for order storage.
Queue: Use an in-memory queue for asynchronous processing.

Scalability: 
Ensure the system can handle 1,000 concurrent orders (simulate load).

**Deliverables**
Functional Backend Code.
A fully functioning backend service with:
RESTful APIs for order management and metrics reporting.
Modular components for queuing, database operations, and metrics computation.

**Database Design:**
Schema to store orders with fields: order_id, user_id, item_ids, total_amount, and status.
SQL scripts for schema creation and sample data population.
Queue Processing:An asynchronous queue that processes orders and updates their status.

**Metrics API:**
Accurate reporting of metrics such as total orders processed, average processing time, and current order statuses.
