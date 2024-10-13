
# Signing Service Challenge in Go

An API service built in Go that allows customers to create signature devices with which they can sign arbitrary transaction data. This application supports multiple signature algorithms and provides flexibility in data persistence.

## Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
3. [Architecture](#architecture)
4. [Services](#services)
5. [Installation and Setup](#installation-and-setup)
6. [Usage](#usage)
7. [Testing](#testing)
8. [Additional Notes](#additional-notes)

## Introduction

This signing service provides a robust platform for customers to create and manage signature devices that enable signing of arbitrary transaction data. The service is designed to be concurrent and extensible, supporting RSA and ECC algorithms. Additionally, it allows for easy switching between in-memory and SQLite persistence storage, ensuring flexibility in deployment.

## Features

- **Concurrency**: The application is built to handle multiple requests simultaneously, leveraging Go's goroutines for efficient concurrency management.
- **Signature Algorithm Support**: Currently supports RSA and ECC algorithms with the ability to add new algorithms without modifying the core domain logic.
- **Data Persistence**: Signature devices are stored in memory by default, with an optional implementation for SQLite. Switching to a relational database can be done easily through a flag in the environment file.
- **API Documentation**: Swagger is used for API documentation, making it easy for developers to understand and utilize the service.

## Architecture

The application follows a structured architecture pattern:

- **MVC Structure**: The project is organized into Models, Views (limited in this case, as it is backend only), and Controllers.
- **Data Persistence**: By default, the application uses an in-memory store for signature devices. There is also a SQLite implementation that can be toggled using an environment variable, allowing for easy switching to a relational database.
- **RESTful API**: All interactions happen through a RESTful JSON API over HTTP(s), promoting a clean separation of concerns.

## Services

The application provides the following services:

### Health Check Service:
- **`GET /api/v0/health`**: Check the health status of the service.

### Signature Device Services:
- **`POST /api/v0/create-signature-device`**: Create a new signature device.
- **`POST /api/v0/sign-transaction`**: Sign a transaction using a specified signature device.
- **`GET /api/v0/devices`**: Retrieve a list of all signature devices.
- **`GET /api/v0/device`**: Retrieve a specific signature device by its ID.

## Installation and Setup

To set up the project locally, follow these steps:

1. **Clone the repository** (if applicable):
   ```bash
   git clone https://github.com/IIMohammad-AsghariII/signing-service-challenge-go.git
   cd signing-service-challenge-go
   ```

2. **Set environment variable for data store**:
   To use the SQLite database, change the following line in your `.env` file:
   ```
   DATA_STORE=memory
   ```
   to 
   ```
   DATA_STORE=db
   ```

   Ensure `CGO_ENABLED=1` is set for SQLite usage.

3. **Install dependencies**:
   Run the following command to get necessary packages:
   ```bash
   go mod tidy
   ```

4. **Run the application**:
   To start the application, run:
   ```bash
   go run main.go
   ```

5. **Access the application**:
   The application will be running at [http://localhost:8080](http://localhost:8080). You can change the default port in the `.env` file.

6. **API Documentation**:
   Swagger is used for API documentation. You can access it at [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html).

## Usage

You can find the postman collection in the postman folder.

### Creating a Signature Device

- **Endpoint**: `POST /api/v0/create-signature-device`
- **Request Body**:
  ```json
  {
    "id": "079bfcfe-4dd1-45fa-bb5f-e91565271060",
    "algorithm": "RSA",
    "label": "Mohammad"
  }
  ```
- **Response**:
  ```json
  {
    "ID": "079bfcfe-4dd1-45fa-bb5f-e91565271060",
    "PublicKey": "-----BEGIN RSA_PUBLIC_KEY-----\nMEgCQQDLTGczkUs545pHTtZBeKlOddEzz9yxaW49Nd/wG1wR6fgTfGPTl298QpLL\nP4wwJ5ktOhJV7nlANrRx5B+/bsZfAgMBAAE=\n-----END RSA_PUBLIC_KEY-----\n",
    "Label": "Mohammad",
    "SignatureCount": 0
  }
  ```

### Signing a Transaction

- **Endpoint**: `POST /api/v0/sign-transaction`
- **Request Body**:
  ```json
  {
    "deviceId": "079bfcfe-4dd1-45fa-bb5f-e91565271060",
    "data": "my_transaction_data"
  }
  ```
- **Response**:
  ```json
  {
    "signature": "<signature_base64_encoded>",
    "signed_data": "<signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>"
  }
  ```

### Listing Signature Devices

- **Endpoint**: `GET /api/v0/devices`
- **Response**:
  ```json
  [
    {
        "ID": "079bfcfe-4dd1-45fa-bb5f-e91565271060",
        "PublicKey": "-----BEGIN RSA_PUBLIC_KEY-----\nMEgCQQDLTGczkUs545pHTtZBeKlOddEzz9yxaW49Nd/wG1wR6fgTfGPTl298QpLL\nP4wwJ5ktOhJV7nlANrRx5B+/bsZfAgMBAAE=\n-----END RSA_PUBLIC_KEY-----\n",
        "Label": "Mohammad",
        "SignatureCount": 1
    },
    {
        "ID": "123e4567-e89b-12d3-a456-426614174000",
        "PublicKey": "-----BEGIN PUBLIC_KEY-----\nMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAERTZfQ/NvbjnyMGkiACzTEM1GsGYyCeZJ\nCHD69yVXOoZRofiWhTCHDZGvOdkIMM3d62/oA0euquTwsqdkmvRWtR8JSah+MHno\nZ9lGzmgLh7r+ZNzpIEbC1tb0QYL81EvL\n-----END PUBLIC_KEY-----\n",
        "Label": "",
        "SignatureCount": 0
    }
  ]
  ```

### Get Signature Device by ID

- **Endpoint**: `GET /api/v0/device?id=079bfcfe-4dd1-45fa-bb5f-e91565271060`
- **Response**:
  ```json
  {
    "ID": "079bfcfe-4dd1-45fa-bb5f-e91565271060",
        "PublicKey": "-----BEGIN RSA_PUBLIC_KEY-----\nMEgCQQDLTGczkUs545pHTtZBeKlOddEzz9yxaW49Nd/wG1wR6fgTfGPTl298QpLL\nP4wwJ5ktOhJV7nlANrRx5B+/bsZfAgMBAAE=\n-----END RSA_PUBLIC_KEY-----\n",
        "Label": "Mohammad",
        "SignatureCount": 1
  }
  ```

### Health Check

- **Endpoint**: `GET /api/v0/health`
- **Response**:
  ```json
  {
    "data": {
        "status": "pass",
        "version": "v0"
    }
  }
  ```

## Testing

All models, repositories, services, and controllers have been thoroughly tested. You can run the tests by navigating to the relevant test folder and executing:

```bash
go test
```

### Test Coverage

The tests cover:
- Signature device creation and retrieval functionality
- Transaction signing functionality
- Concurrent processing of requests
- Data persistence methods (both in-memory and SQLite)
- Verification of the signature algorithm implementations

## Additional Notes

- The request and response structures are separated into payload classes for better organization and clarity.
- To create a signature device, the user must provide a UUID as an ID and the algorithm type. The label is optional.
- Upon creation, the device can be used to sign transactions with the same ID and provided data.

### API Method Signatures

- `CreateSignatureDevice(id: string, algorithm: 'ECC' | 'RSA', [optional]: label: string): CreateSignatureDeviceResponse`
- `SignTransaction(deviceId: string, data: string): SignatureResponse`

### Response Structure

```json
{
    "signature": "<signature_base64_encoded>",
    "signed_data": "<signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>"
}
```

---

# Signature Service - Coding Challenge

## Instructions

This challenge is part of the software engineering interview process at fiskaly.

If you see this challenge, you've passed the first round of interviews and are now at the second and last stage.

We would like you to attempt the challenge below. You will then be able to discuss your solution in the skill-fit interview with two of our colleagues from the development department.

The quality of your code is more important to us than the quantity.

### Project Setup

For the challenge, we provide you with:

- Go project containing the setup
- Basic API structure and functionality
- Encoding / decoding of different key types (only needed to serialize keys to a persistent storage)
- Key generation algorithms (ECC, RSA)
- Library to generate UUIDs, included in `go.mod`

You can use these things as a foundation, but you're also free to modify them as you see fit.

### Prerequisites & Tooling

- Golang (v1.20+)

### The Challenge

The goal is to implement an API service that allows customers to create `signature devices` with which they can sign arbitrary transaction data.

#### Domain Description

The `signature service` can manage multiple `signature devices`. Such a device is identified by a unique identifier (e.g. UUID). For now you can pretend there is only one user / organization using the system (e.g. a dedicated node for them), therefore you do not need to think about user management at all.

When creating the `signature device`, the client of the API has to choose the signature algorithm that the device will be using to sign transaction data. During the creation process, a new key pair (`public key` & `private key`) has to be generated and assigned to the device.

The `signature device` should also have a `label` that can be used to display it in the UI and a `signature_counter` that tracks how many signatures have been created with this device. The `label` is provided by the user. The `signature_counter` shall only be modified internally.

##### Signature Creation

For the signature creation, the client will have to provide `data_to_be_signed` through the API. In order to increase the security of the system, we will extend this raw data with the current `signature_counter` and the `last_signature`.

The resulting string (`secured_data_to_be_signed`) should follow this format: `<signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>`

In the base case there is no `last_signature` (= `signature_counter == 0`). Use the `base64`-encoded device ID (`last_signature = base64(device.id)`) instead of the `last_signature`.

This special string will be signed (`Signer.sign(secured_data_to_be_signed)`) and the resulting signature (`base64` encoded) will be returned to the client. The signature response could look like this:

```json
{ 
    "signature": <signature_base64_encoded>,
    "signed_data": "<signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>"
}
```

After the signature has been created, the signature counter's value has to be incremented (`signature_counter += 1`).

#### API

For now we need to provide two main operations to our customers:

- `CreateSignatureDevice(id: string, algorithm: 'ECC' | 'RSA', [optional]: label: string): CreateSignatureDeviceResponse`
- `SignTransaction(deviceId: string, data: string): SignatureResponse`

Think of how to expose these operations through a RESTful HTTP-based API.

In addition, `list / retrieval operations` for the resources generated in the previous operations should be made available to the customers.

#### QA / Testing

As we are in the business of compliance technology, we need to make sure that our implementation is verifiably correct. Think of an automatable way to assure the correctness (in this challenge: adherence to the specifications) of the system.

#### Technical Constraints & Considerations

- The system will be used by many concurrent clients accessing the same resources.
- The `signature_counter` has to be strictly monotonically increasing and ideally without any gaps.
- The system currently only supports `RSA` and `ECDSA` as signature algorithms. Try to design the signing mechanism in a way that allows easy extension to other algorithms without changing the core domain logic.
- For now it is enough to store signature devices in memory. Efficiency is not a priority for this. In the future we might want to scale out. As you design your storage logic, keep in mind that we may later want to switch to a relational database.

#### Credits

This challenge is heavily influenced by the regulations for `KassenSichV` (Germany) as well as the `RKSV` (Austria) and our solutions for them.
