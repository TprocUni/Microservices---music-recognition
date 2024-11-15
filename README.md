# Microservices---music-recognition
This project implements a suite of **Go-based microservices**. The services enable music storage, search, and device integration using modern RESTful APIs, handling audio files encoded in Base64. Each microservice demonstrates modular design and effective inter-service communication, reflecting real-world enterprise computing practices.


### Microservice Implementation

#### Tracks Microservice

The **Tracks microservice** acts as a music track database and operates on port 3000. It allows for creating, listing, retrieving, and deleting music tracks. Each operation is implemented with RESTful endpoints and interacts with an SQL database for persistent storage.

1. **Creating Tracks**:
   - Endpoint: `PUT /tracks/{id}`
   - Accepts a JSON object containing:
     - `Id`: A string identifier for the track.
     - `Audio`: A WAV file encoded in Base64.
   - Responses:
     - `201 Created` or `204 No Content` for success.
     - `400 Bad Request` or `500 Internal Server Error` for failure.
   - Example usage:
     ```sh
     curl -v -X PUT -d @input localhost:3000/tracks/{id}
     ```

2. **Listing Tracks**:
   - Endpoint: `GET /tracks`
   - Returns a list of all track IDs.
   - Responses:
     - `200 OK` for success.
     - `500 Internal Server Error` for failure.

3. **Reading Tracks**:
   - Endpoint: `GET /tracks/{id}`
   - Fetches a specific track by ID.
   - Responses:
     - `200 OK` with track details on success.
     - `404 Not Found` or `500 Internal Server Error` for failure.

4. **Deleting Tracks**:
   - Endpoint: `DELETE /tracks/{id}`
   - Removes a track by ID.
   - Responses:
     - `204 No Content` for success.
     - `404 Not Found` or `500 Internal Server Error` for failure.

#### Search Microservice

The **Search microservice** replicates a "Hum-to-Search" feature and operates on port 3001. It identifies music tracks from audio fragments using the audd.io API for music recognition.

1. **Recognizing Tracks**:
   - Endpoint: `POST /search`
   - Accepts a JSON object with:
     - `Audio`: A WAV file fragment encoded in Base64.
   - Responses:
     - `200 OK` with the identified track ID.
     - `400 Bad Request`, `404 Not Found`, or `500 Internal Server Error` for failure.
   - Example usage:
     ```sh
     curl -v -X POST -d @input localhost:3001/search
     ```

2. **External Integration**:
   - Leverages the audd.io API to recognize music fragments, ensuring accurate and reliable identification.

#### CoolTown Microservice

The **CoolTown microservice** enables device integration by retrieving full music tracks based on fragments, similar to the Samsung Bixby feature. It operates on port 3002.

1. **Retrieving Tracks**:
   - Endpoint: `POST /cooltown`
   - Accepts a JSON object with:
     - `Audio`: A Base64-encoded WAV file fragment.
   - Responses:
     - `200 OK` with the full music track.
     - `400 Bad Request`, `404 Not Found`, or `500 Internal Server Error` for failure.
   - Example usage:
     ```sh
     curl -v -X POST -d @input localhost:3002/cooltown
     ```

2. **Functionality**:
   - Bridges the search and track storage microservices, enabling seamless integration for real-world device functionality.

---

### Key Features

1. **RESTful API Design**:
   - Standardized HTTP methods (`PUT`, `GET`, `POST`, `DELETE`) for CRUD operations.

2. **Audio Handling**:
   - Utilizes Base64 encoding for audio file storage and transfer.
   - WAV format ensures compatibility with most playback devices.

3. **Inter-Service Communication**:
   - Modular architecture enables the microservices to interact effectively, simulating a real-world system.

4. **Error Handling**:
   - Implements detailed status codes (`200`, `204`, `400`, `404`, `500`) for robust client-server communication.

5. **Extensibility**:
   - Designed to integrate with external APIs (e.g., audd.io) and additional microservices.

---

### Conclusion

This system provides a robust, modular framework for audio management and recognition, emulating functionalities showcased in Addison Rae’s promotional video. It demonstrates the practical application of microservices architecture, offering a scalable solution for music-related operations. By leveraging RESTful APIs, SQL databases, and external services, the system is well-suited for enterprise-level applications.
