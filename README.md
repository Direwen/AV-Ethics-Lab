# Interactive AV Ethics Lab

> A web-based interactive simulator to study ethical decision-making in autonomous vehicle scenarios

**🌐 Live Application**: [https://av-ethics-lab.vercel.app/](https://av-ethics-lab.vercel.app/)

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.25.4-blue.svg)](https://golang.org/)
[![Nuxt Version](https://img.shields.io/badge/nuxt-4.2.1-green.svg)](https://nuxt.com/)

**Author**: Min Sik Hein (Direwen) | Final Year Computer Science Student | Software Engineer & Researcher  
**Institution**: University of Sunderland  
**Module**: CET300

## 📋 Table of Contents

- [Overview](#overview)
- [Background & Motivation](#background--motivation)
- [Features](#features)
- [Architecture](#architecture)
- [Technology Stack](#technology-stack)
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Research Methodology](#research-methodology)
- [System Design Diagrams](#system-design-diagrams)
- [Screenshots](#screenshots)
- [Contributing](#contributing)
- [License](#license)

## 🎯 Overview

This project develops a web-based interactive simulator to study how people make ethical decisions in autonomous vehicle (AV) scenarios. Unlike traditional trolley problem-style surveys that force binary choices, this system presents users with context-rich, time-pressured driving scenarios where they must decide how an autonomous vehicle should act when harm is unavoidable.

### Key Objectives

- **Capture Instinctive Decision-Making**: Present realistic scenarios with time pressure to observe natural moral reasoning
- **Move Beyond Binary Choices**: Use ranking-based responses instead of simple yes/no decisions
- **Ensure Statistical Rigor**: Implement balanced factorial design with controlled variables
- **Maintain Privacy**: Anonymous data collection with GDPR compliance
- **Generate Actionable Insights**: Provide data that can inform AV policy and algorithm design

## 🧠 Background & Motivation

Most existing research on autonomous vehicle (AV) ethics relies on trolley problem–style experiments, such as the MIT Moral Machine, which force participants to choose between two predefined harmful outcomes. While useful at scale, these approaches oversimplify real driving conditions by ignoring uncertainty, time pressure, cultural context, and the difference between action and inaction.

Real-world traffic accidents rarely present clean binary choices. Instead, drivers (human or autonomous) must make decisions under incomplete information, spatial constraints, and competing risks. As a result, data derived from abstract moral dilemmas may fail to translate into practical AV decision-making or policy.

This project was created to address that gap by providing a controlled but realistic simulation environment where users make ethical decisions *as the vehicle*, under time pressure, using action-based choices rather than entity-based valuation. The goal is to collect empirically meaningful data that better reflects how humans reason about ethical trade-offs in real traffic scenarios.

## ✨ Features

### Core Functionality

- **Interactive Scenario Simulation**: Visual grid-based representation of driving scenarios
- **Time-Pressured Decisions**: 20-second timer per scenario to capture instinctive responses
- **Ranking-Based Responses**: Drag-and-drop interface to rank three action options
- **Balanced Experiment Design**: Ensures equal exposure to critical variables (visibility, brake status, entity behavior)
- **Personalized Feedback**: LLM-generated ethical archetype classification after completion
- **Public Dashboard**: Real-time analytics and insights from collected data

### Research Features

- **Trident of Doom Mechanism**: Three-way dilemma system (maintain, swerve left, swerve right)
- **Surface-Aware Entity Placement**: Realistic placement respecting drivable/walkable surfaces
- **Behavior Control**: Strict control over primary entity behavior (violation vs. compliant)
- **Duplicate Detection**: Tracks returning participants while maintaining anonymity

## 🏗️ Architecture

### System Components

```
┌─────────────────┐         ┌─────────────────┐         ┌─────────────────┐
│   Frontend      │         │   API Server    │         │   LLM Services  │
│   (Nuxt.js)     │◄───────►│   (Go/Echo)     │◄───────►│   (Groq/OpenRouter)
│                 │  HTTP   │                 │  API    │                 │
└─────────────────┘         └─────────────────┘         └─────────────────┘
         │                          │
         │                          │
         └──────────────────────────┘
                    │
                    ▼
         ┌─────────────────┐
         │   PostgreSQL    │
         │   Database      │
         └─────────────────┘
```

### Key Design Patterns

- **Service Layer Architecture**: Separation of concerns with dedicated services for sessions, scenarios, responses, and dashboard
- **LLM Pool System**: Separate client pools for scenario generation and feedback generation
- **Template-Based Scenarios**: Pre-computed grid templates with in-memory caching for performance
- **Transaction Management**: Atomic operations for response submission and session completion

## 🛠️ Technology Stack

### Frontend
- **Framework**: Nuxt.js 4 (Vue 3)
- **Language**: TypeScript
- **Styling**: TailwindCSS 4
- **State Management**: Pinia
- **UI Components**: Maz UI
- **Charts**: Chart.js / Vue-ChartJS
- **Drag & Drop**: Vue Draggable Next
- **Fingerprinting**: FingerprintJS

### Backend
- **Language**: Go 1.25.4
- **Framework**: Echo v4
- **ORM**: GORM
- **Database**: PostgreSQL
- **LLM Integration**: LangChainGo
- **Authentication**: JWT (golang-jwt/jwt)
- **Validation**: go-playground/validator

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Environment**: dotenv configuration

## 🚀 Getting Started

### Prerequisites

- Go 1.25.4 or higher
- Node.js 18+ and npm/pnpm/yarn
- PostgreSQL 14+
- Docker & Docker Compose (optional)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd AV-Ethics-Lab
   ```

2. **Backend Setup**
   ```bash
   cd go-server
   cp .env.example .env  # Configure your environment variables
   go mod download
   ```

3. **Frontend Setup**
   ```bash
   cd frontend
   npm install  # or pnpm install / yarn install
   ```

4. **Database Setup**
   ```bash
   # Using Docker Compose
   docker-compose up -d postgres
   
   # Or configure your PostgreSQL connection in .env
   ```

5. **Environment Variables**

   **Backend (.env in go-server/)**
   ```env
   DATABASE_URL=postgres://user:password@localhost:5432/av_ethics_db
   JWT_SECRET=your-secret-key-here
   SERVER_PORT=:8080
   LOCAL_FRONTEND_PORT=3000
   SESSION_EXPIRATION=4h
   GROQ_API_KEY=your-groq-api-key
   OPENROUTER_API_KEY=your-openrouter-api-key
   TRIDENT_ZONE_DISTANCE=3
   TRIDENT_ZONE_DEPTH=3
   BACKGROUND_ENTITIES_MIN=2
   BACKGROUND_ENTITIES_MAX=4
   ```

   **Frontend (nuxt.config.ts)**
   ```typescript
   runtimeConfig: {
     public: {
       apiBase: "http://localhost:8080",
       tridentZoneDistance: 1,
       timerDuration: 20
     }
   }
   ```

### Running the Application

1. **Start the Backend**
   ```bash
   cd go-server
   go run cmd/api/main.go
   ```

2. **Start the Frontend**
   ```bash
   cd frontend
   npm run dev
   ```

3. **Access the Application**
   - Frontend: http://localhost:3000
   - API: http://localhost:8080

### Docker Deployment

```bash
docker-compose up -d
```

## 📁 Project Structure

```
AV-Ethics-Lab/
├── frontend/                 # Nuxt.js frontend application
│   ├── app/
│   │   ├── components/       # Vue components
│   │   ├── pages/           # Route pages
│   │   ├── stores/          # Pinia stores
│   │   ├── types/           # TypeScript types
│   │   └── composables/     # Vue composables
│   └── nuxt.config.ts       # Nuxt configuration
│
├── go-server/                # Go backend API
│   ├── cmd/api/             # Main application entry
│   ├── internal/
│   │   ├── config/          # Database configuration
│   │   ├── session/         # Session management
│   │   ├── scenario/        # Scenario generation
│   │   ├── response/        # Response handling
│   │   ├── dashboard/       # Analytics dashboard
│   │   ├── platform/llm/    # LLM integration
│   │   ├── template/        # Template management
│   │   └── shared/          # Shared domain logic
│   └── pkg/                 # Package utilities
│
├── system design/           # Architecture documentation
│   ├── diagrams/           # PlantUML diagrams
│   └── notes/              # Design notes
│
└── docs/                    # Documentation
    └── design/ui/          # UI screenshots
```

## 🔬 Research Methodology

### Experimental Design

The system uses a **Balanced Factorial Design** with block randomization:

- **Critical Factors** (Balanced):
  - Visibility: Clear, Fog, Night, Rain (cycled)
  - Brake Status: Active, Failed, Fade (cycled)
  - Primary Behavior: 50% Violation / 50% Compliant

- **Randomized Factors**:
  - Road Condition: Dry, Wet, Icy
  - Location: US, UK, CN, FR
  - Speed: Low, Medium, High
  - Tailgater: Present/Absent

### The Trident of Doom

Each scenario presents a three-way dilemma:

- **Zone A (Forward)**: The inaction path - must contain a threat
- **Zone B (Left)**: Alternative swerve path
- **Zone C (Right)**: Alternative swerve path

This ensures no completely safe option exists, forcing meaningful ethical choices.

### Data Collection

- **Ranking Order**: User ranks three action options
- **Response Time**: Measured in milliseconds
- **Interaction Tracking**: Detects if user actually interacted
- **Timeout Detection**: Flags responses that exceeded time limit

### Differences from MIT Moral Machine

1. **Action Ranking vs. Entity Ranking**: Ranks actions (which rule to break) rather than entities (who is worth more)
2. **Omission vs. Commission**: Explicitly separates inaction from action
3. **Policy Relevance**: Data translates to cost functions for AV path-planning algorithms

## 📊 System Design Diagrams

For detailed architecture documentation, see the [System Design Diagrams](system%20design/diagrams/):

- **[Entity Relationship Diagram (ERD)](system%20design/diagrams/erd.puml)**: Database schema and relationships between sessions, scenarios, responses, and templates
- **[Site Map](system%20design/diagrams/sitemap.puml)**: Frontend routing structure showing public, guest, protected, and optional routes
- **[Sequence Diagram](system%20design/diagrams/sequence.puml)**: Complete experiment flow from session creation through scenario generation, response submission, and feedback generation

> **Note**: These diagrams are written in PlantUML format. To view them, use a PlantUML renderer (VS Code extension, online server, or CLI tool).

## 📸 Screenshots

### Home Page
![Home Page - Light Mode](docs/design/ui/home_page_light_mode.png)
![Home Page - Dark Mode](docs/design/ui/home_page_dark_mode.png)

### Consent & Demographics
![Consent Page - Dark Mode](docs/design/ui/consent_page_dark_mode.png)

### Experiment Interface
![Experiment Scenario 1](docs/design/ui/experiment1.png)
![Experiment Scenario 2](docs/design/ui/experiment2.png)

### Guide & Tutorial
![Guide - Dark Mode](docs/design/ui/guide_dark_mode.png)
![Guide Info - Dark Mode](docs/design/ui/guide_info_dark_mode.png)

### Dashboard Analytics
![Dashboard - Light Mode](docs/design/ui/dashboard_light_mode.png)
![Dashboard - Dark Mode](docs/design/ui/dashboard_dark_mode.png)

### Feedback & Results
![Feedback - Dark Mode](docs/design/ui/feedback_dark_mode.png)

### Loading States
![Loading](docs/design/ui/loading.png)

### Error Pages
![404 Page - Light Mode](docs/design/ui/404_light.png)
![404 Page - Dark Mode](docs/design/ui/404_dark.png)

## 🔒 Privacy & Ethics

- **Anonymous Data Collection**: No email addresses or IP addresses stored
- **Fingerprint-Based Identification**: Uses browser fingerprinting for uniqueness while maintaining anonymity
- **GDPR Compliant**: Privacy-first design with data minimization
- **Session Expiration**: Automatic cleanup of expired sessions
- **Consent-Based**: Explicit consent required before participation

## 🤝 Contributing

This is a research project for academic purposes. Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built as a final semester project for BSc Computer Science at University of Sunderland
- Inspired by the need for more realistic ethical decision-making research in autonomous vehicles
- Uses open-source technologies and frameworks
- Special thanks to the research participants who contributed their time and insights

## 📧 Contact

For questions about this research project, please contact the researcher (direwen.work@gmail.com).

---

**Note**: This project is part of ongoing research into ethical decision-making frameworks for autonomous vehicles. Your participation helps us understand human moral reasoning in complex traffic scenarios, which will inform the development of more ethical AI systems.
