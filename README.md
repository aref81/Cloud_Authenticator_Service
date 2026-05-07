# Cloud Authenticator Service

A cloud-native **face-based identity verification** system built in Go. Users register with two photos; the system asynchronously verifies their identity by comparing facial similarity using the Imagga API, then notifies them by email with the result.

The system is split into two decoupled microservices — an **API Server** and an **Authenticator** — coordinated via RabbitMQ.

---

## How It Works

```
                 Client
                   │
                   │  POST /register  (name, email, national_code, pic1, pic2)
                   ▼
┌──────────────────────────────────────┐
│           API Server (:8000)         │
│                                      │
│  1. Validate & save user → PostgreSQL│
│  2. Upload pic1, pic2 → S3           │
│  3. Enqueue national_code → RabbitMQ │
└──────────────────┬───────────────────┘
                   │  AMQP queue: "reqs"
                   ▼
┌──────────────────────────────────────┐
│           Authenticator              │
│                                      │
│  1. Consume national_code from queue │
│  2. Download pic1, pic2 from S3      │
│  3. Detect faces → Imagga API        │
│  4. Compare similarity → Imagga API  │
│  5. score ≥ 80 → accepted            │
│     score < 80 → rejected            │
│  6. Update status → PostgreSQL       │
│  7. Send result email → Mailgun      │
└──────────────────────────────────────┘
```

---

## API Endpoints

### `POST /register`

Registers a new user and initiates asynchronous face verification.

**Request** — `multipart/form-data`:

| Field  | Type | Required | Description |
|--------|------|----------|-------------|
| `info` | JSON string | ✓ | `{"name": "", "email": "", "national_code": ""}` |
| `pic1` | file | ✓ | First face photo |
| `pic2` | file | ✓ | Second face photo for comparison |

**Response** `200 OK`:
```json
{ "message": "Successfully registered" }
```

---

### `GET /status`

Returns the current authentication status for a user. Requires the request to originate from the same IP address used during registration.

**Request** — query params or JSON body:

| Field           | Type   | Required |
|-----------------|--------|----------|
| `national_code` | string | ✓        |

**Response** `200 OK`:
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "national_code": "...",
  "ip": "...",
  "status": "accepted | rejected | pending"
}
```

> Requests from a different IP than the one used at registration are rejected with `403 Forbidden`.

---

## Architecture

### `apiserver`

REST API built with [Echo](https://echo.labstack.com/). Responsibilities:

- Validate and persist user records (name, email, national code, IP) to **PostgreSQL**
- Upload two face images to **S3-compatible storage**, keyed as `base64(national_code)_IMAGE_1` and `_IMAGE_2`
- Publish the encoded national code to the `reqs` RabbitMQ queue

### `authenticator`

Background worker service. Responsibilities:

- Subscribe to the `reqs` queue on **RabbitMQ**
- Download both images from **S3**
- Submit each image to the **Imagga Face Detection API** to obtain a `face_id`
- Call the **Imagga Face Similarity API** with both `face_id`s to get a confidence score (0–100)
- If `score ≥ 80` → `accepted`; otherwise `rejected`
- Update the user's status in **PostgreSQL**
- Send a result notification email via **Mailgun**

---

## Tech Stack

| Layer            | Technology                              |
|------------------|-----------------------------------------|
| Language         | Go                                      |
| HTTP Framework   | Echo v4                                 |
| Database         | PostgreSQL                              |
| Message Queue    | RabbitMQ (AMQP)                         |
| Object Storage   | S3-compatible (e.g. ArvanCloud)         |
| Face Recognition | [Imagga API](https://imagga.com)        |
| Email            | Mailgun                                 |
| Logging          | Logrus                                  |
| Containerization | Docker / Docker Compose                 |

---

## Project Structure

```
Cloud_Authenticator_Service/
├── docker-compose.yaml
├── apiserver/
│   ├── main.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── api/
│   │   ├── server.go               # Echo router setup
│   │   └── handlers/
│   │       ├── common.go           # Tool initialization (PSQL, RabbitMQ)
│   │       ├── registerHandler.go  # POST /register
│   │       └── statusHandler.go    # GET /status
│   ├── internal/model/
│   │   ├── user.go                 # User struct
│   │   └── handlerInterfaces.go    # Request/response models
│   ├── sql/
│   │   └── users.sql               # Table schema
│   └── utils/
│       ├── encoder.go              # Base64 helper
│       ├── broker/rabbitMQ.go      # RabbitMQ publisher
│       └── datasource/
│           ├── psql.go             # PostgreSQL client
│           └── s3.go               # S3 upload
└── authenticator/
    ├── main.go
    ├── Dockerfile
    ├── go.mod
    ├── internal/
    │   ├── listenerService.go          # RabbitMQ consumer loop
    │   └── authenticationProcessor.go # Core verification logic
    ├── internal/model/
    │   ├── user.go
    │   └── restInterfaces.go           # Imagga API response models
    └── utils/
        ├── encoder.go
        ├── broker/rabbitMQ.go          # RabbitMQ subscriber
        ├── datasource/
        │   ├── psql.go
        │   └── s3.go                   # S3 download
        └── 3rdPartyService/
            ├── faceDetection.go        # Imagga face detection
            ├── faceSimilarity.go       # Imagga similarity scoring
            └── mailService.go          # Mailgun email sender
```

---

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)
- A PostgreSQL instance
- A RabbitMQ instance (AMQP)
- An S3-compatible object storage bucket
- An [Imagga](https://imagga.com) account (API key + secret)
- A [Mailgun](https://mailgun.com) account

### Database Setup

Run the following SQL on your PostgreSQL instance before starting the services:

```sql
CREATE TABLE users (
    national_code VARCHAR(100) PRIMARY KEY,
    name          VARCHAR(40),
    email         VARCHAR(50),
    ip            VARCHAR(20),
    status        VARCHAR(50)
);
```

### Configuration

Both services read configuration from environment variables. Edit `docker-compose.yaml` or supply them via a `.env` file.

**`apiserver` variables:**

| Variable        | Description                        |
|-----------------|------------------------------------|
| `PS_URI`        | PostgreSQL connection string       |
| `S3_BUCKET`     | S3 bucket name                     |
| `S3_REGION`     | S3 region                          |
| `S3_ENDPOINT`   | S3 endpoint URL                    |
| `S3_ACCESS_KEY` | S3 access key                      |
| `S3_SECRET_KEY` | S3 secret key                      |
| `RB_URL`        | RabbitMQ AMQP connection URL       |

**`authenticator` variables** (all of the above, plus):

| Variable                           | Description                                              |
|------------------------------------|----------------------------------------------------------|
| `IMAGGA_API_KEY`                   | Imagga API key                                           |
| `IMAGGA_API_SECRET`                | Imagga API secret                                        |
| `IMAGGA_FACE_DETECTION_URL`        | `https://api.imagga.com/v2/faces/detections`             |
| `IMAGGA_SIMILARITY_DETECTION_URL`  | `https://api.imagga.com/v2/faces/similarity`             |
| `MAILGUN_DOMAIN`                   | Your Mailgun sending domain                              |
| `MAILGUN_API_KEY`                  | Mailgun API key                                          |

> ⚠️ **Security:** Never commit real credentials to version control. Use a `.env` file (add it to `.gitignore`) or a secrets manager in production.

### Run

```bash
git clone https://github.com/aref81/Cloud_Authenticator_Service.git
cd Cloud_Authenticator_Service

# Configure credentials in docker-compose.yaml or a .env file
docker compose up --build
```

The API server will be available at `http://localhost:8000`.

```bash
# Run in background
docker compose up --build -d

# Tail logs
docker compose logs -f

# Stop
docker compose down
```

### Example Requests

**Register a user:**

```bash
curl -X POST http://localhost:8000/register \
  -F 'info={"name":"John Doe","email":"john@example.com","national_code":"1234567890"}' \
  -F 'pic1=@/path/to/photo1.jpg' \
  -F 'pic2=@/path/to/photo2.jpg'
```

**Check authentication status:**

```bash
curl "http://localhost:8000/status?national_code=1234567890"
```

---

## Face Verification Details

The authenticator uses the [Imagga](https://imagga.com) API for AI-powered face comparison:

1. **Face Detection** — each image is submitted and a `face_id` token is returned.
2. **Similarity Scoring** — the two `face_id` tokens are compared; a score from `0` to `100` is returned.
3. **Decision** — score `≥ 80` → status set to `accepted`; score `< 80` → `rejected`. The user receives an email either way.
