## gRPC server-client app for anomaly detection in a stream of data

The project is an anomaly detection system for streaming frequency data. The system consists of a gRPC-based server that generates and streams frequency data, a client that processes this data to detect anomalies, and a PostgreSQL database for storing detected anomalies.
This project was developed as a part of School 21 curriculum.

1. **Transmitter (gRPC Server)**: 
   - Implements a gRPC server that streams frequency data in a specific format. Each message contains a `session_id`, a `frequency` value, and a `timestamp`.
   - The server generates the data based on a normal distribution, with random parameters for mean and standard deviation for each new client connection.
   - Logs generated values and metadata to ensure traceability and reproducibility.

2. **Anomaly Detection Client**:
   - Connects to the gRPC server and receives a continuous stream of frequency data.
   - Approximates the mean and standard deviation of the incoming frequency values in real-time. 
   - After determining accurate distribution parameters, the client transitions to an anomaly detection phase.
   - Uses a configurable anomaly detection threshold (coefficient `k`) to identify outliers in the data stream, based on deviations from the mean.

3. **PostgreSQL Integration**:
   - Anomalies detected by the client are stored in a PostgreSQL database, with fields for `session_id`, `frequency`, and `timestamp`.
   - The project uses an ORM to interact with the database, ensuring that SQL queries are secure and maintainable.

4. **System Integration**:
   - The components (server, client, and database) work together seamlessly to detect and log anomalies in the streamed data.
   - The system can be used for anomaly detection on any real-time data stream following similar statistical patterns.

This project demonstrates practical skills in real-time data processing, anomaly detection, and integration of various technologies (gRPC, Go, PostgreSQL, and bunORM), while adhering to best practices for secure and efficient coding.

## Instructions for Running the Project

1. **Prerequisites**:
   - Make sure you have Go installed (version 1.21 is recommended).
   - Install `protoc` (Protocol Buffers compiler) and the Go plugins for `protoc` (`protoc-gen-go` and `protoc-gen-go-grpc`).
   - Ensure a PostgreSQL database is running.
   - Adjust PostgreSQL connection details in the Taskfile.

2. **Build and Run the Server**:

   Use the `server` task to build and start the gRPC server:
   ```bash
   task server
   ```
   The server will start streaming frequency data. Make sure it runs without errors before proceeding.

3. **Build and Run the Client**:

   Start the anomaly detection client with the appropriate PostgreSQL connection configuration:
   ```bash
   task client
   ```
   Make sure to adjust the DSN in the command if your PostgreSQL setup differs.

4. **Run Tests**:

   To run the test suite for the anomaly detection logic, execute:
   ```bash
   task test
   ```

By following these steps, you should be able to generate code, build the server and client, and run the complete system, including testing the anomaly detection functionality. Make sure the PostgreSQL server is configured correctly, as the client will log detected anomalies to the database.
