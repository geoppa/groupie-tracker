# 🎵 Groupie Tracker

Groupie Tracker is a data-visualization platform that consumes an external RESTful API to present tracking charts, historical information, tour schedules, and relationship networks of legendary musical bands. Built with **Go** on the backend and **Semantic HTML/CSS/JS** on the frontend, it focuses on high performance, accessibility (A11y), and zero-panic error handling.

---

## 🏗️ Architecture & Code Execution Order

When a user interacts with the platform, the application follows a strict sequence of execution across isolated packages:
```
[User Browser] ──> [main.go] ──> [handlers] ──> [api] ──> [Zone01 API]
 └──<─── [Rendered HTML] ───<─────┴──────<──────┴─────<─────┘
```
 1. **`main.go` (Entry Point)**
   * Initializes the application instance.
   * Compiles all static template layouts from the `templates/` folder into memory via `template.ParseGlob`.
   * Configures static asset routing (`/static/`) for assets like stylesheets and client scripts.
   * Binds specific path endpoints to their respective logical handlers and runs the server engine loop on port `8080`.

2. **`internal/models/` (Blueprints)**
   * Contains strongly typed schemas that explicitly define how data must look.
   * Maps incoming API camelCase payloads into Go-standard PascalCase structs using declarative JSON field tags.

3. **`internal/api/` (Network Client)**
   * Manages outbound data requests using optimized memory-efficient standard body data streams via `json.NewDecoder`.
   * Unpacks API elements into safe data slices and structured lookup dictionaries.

4. **`internal/handlers/` (Business Controllers & Error Fallbacks)**
   * Validates incoming query strings, path patterns, and integers safely (`strconv.Atoi`).
   * Fetches data arrays from the network client layer, converts raw states into JSON formats (`template.JS`) to power client-side elements, and protects stream headers using temporary buffer pipes (`bytes.Buffer`).
   * Includes custom error fallbacks (`renderError`) that cleanly isolate operational glitches from dropping active user connections.

---

## 🚀 Getting Started

### Prerequisites
* **Go** installed locally on your machine (Version 1.16+ recommended).

### Running the Application Local Server
1. Clone or navigate to the root directory of the project
2. Fire up the development engine by executing:
   ```bash
   go run ./cmd/tracker
   ```
3. Open your favorite web browser and navigate to: http://localhost:8080

---

## 🧪 Testing Suite Execution

The repository contains an automated suite of continuous unit testing configurations to verify package boundaries, template engines, API models, and response logic in absolute isolation.

### Run All Unit Tests
To fire the evaluation engine across all available code packages simultaneously, execute the following command in your terminal:
```bash
go test ./...
```

### Run Tests with Detailed Output Logs
To review the step-by-step validation tracking blocks, performance measurements, and function tracing data hooks live, execute:
```bash
go test -v ./...
```

### What Each Test Validates
* **`main_test.go`**: Dynamically climbs up your directory paths from your testing context workspace to locate your `templates/` pool and validates layout syntax composition rules before compilation.
* **`internal/models/models_test.go`**: Validates JSON string mapping transformations into structural models, guaranteeing mapping accuracy across multidimensional city relationship indices (`map[string][]string`).
* **`internal/api/api_test.go`**: Isolates networking behaviors using simulated mock routing endpoints to verify memory allocation boundaries under successful and non-200 OK HTTP statuses.
* **`internal/handlers/handlers_test.go`**: Generates memory requests and recorders (`httptest.NewRecorder`) to inspect server codes, catch route injection attempts, and prevent pointer dereference panics.

---

## 🌟 Implemented Core Features
* **Dynamic Search Architecture**: Multi-level keyword search processing names, band creation dates, track releases, members, and geographic show tracking markers simultaneously.
* **Accessible Component Engineering**: Interfaces styled with semantic HTML patterns using explicit `aria-*` parameters to facilitate clean assistive screen reader support.
* **Color Customization Engine**: Isolated component structures combined with flexible client data theme controllers supporting Neon, Inferno, and Modern interfaces seamlessly.

## ⚖️ Copyright & Credits

This project was developed as part of the Zone01 curriculum by:

* **Frontend Architecture & Engineering**
* **ELEANA GEORGIOU** #elgeorgiou
* **Backend Architecture & Engineering**
* **GEORGIOS PAPADAKIS** #gpapadaki

Copyright © 2026 All rights reserved.