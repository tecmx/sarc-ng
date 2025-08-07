docs/
├── intro.mdx                        # Project introduction with interactive demo
│
├── getting-started/                 # Getting started guides
│   ├── installation.md              # Installation instructions
│   ├── development.md               # Development setup
│   ├── docker.md                    # Docker setup
│   └── quick-start.mdx              # Quick start with interactive elements
│
├── architecture/                    # System architecture
│   ├── overview.mdx                 # System architecture with diagrams
│   └── components.md                # Detailed component descriptions
│
├── api/                             # API conceptual documentation
│   ├── overview.mdx                 # API introduction with example widget
│   ├── authentication.mdx           # Authentication flow with interactive diagrams
│   ├── pagination.md                # Pagination explanation
│   ├── error-handling.md            # Error handling guidelines
│   └── examples/                    # Integration examples
│       ├── javascript.mdx           # JavaScript examples with live coding
│       ├── python.mdx               # Python examples
│       └── curl.md                  # cURL examples
│
├── api-reference/                   # Interactive API reference (auto-generated)
│   ├── buildings/                   # Generated from OpenAPI spec
│   ├── classes/
│   ├── lessons/
│   ├── resources/
│   ├── reservations/
│   └── system/                      # System endpoints like health, metrics
│
├── deployment/                      # Deployment guides
│   ├── overview.md                  # Deployment options
│   ├── aws.mdx                      # AWS deployment with architecture diagrams
│   ├── containerization.md          # Container deployment
│   └── serverless.md                # Serverless deployment
│
├── infrastructure/                  # Infrastructure docs
│   ├── overview.md                  # Infrastructure introduction
│   ├── aws-modules.md               # AWS modules documentation
│   └── terraform-examples.mdx       # Terraform examples with syntax highlighting
│
├── testing/                         # Testing guides
│   ├── overview.md                  # Testing strategy
│   ├── unit-tests.md                # Unit testing guide
│   └── integration-tests.mdx        # Integration testing with example runners
│
├── contributing/                    # Contribution guidelines
│   ├── workflow.md                  # Development workflow
│   ├── standards.md                 # Code standards
│   └── pull-requests.md             # PR process
│
├── tutorials/                       # Step-by-step tutorials
│   ├── creating-api.mdx             # Build an API tutorial with progress tracker
│   ├── resource-management.mdx      # Resource management tutorial
│   └── reservation-system.mdx       # Building a reservation system
│
├── reference/                       # Reference materials
│   ├── configuration.md             # Configuration reference
│   ├── environment-variables.md     # Environment variables reference
│   ├── cli-commands.mdx             # CLI commands with interactive terminal
│   └── troubleshooting.md           # Common issues and solutions
│
├── _components/                     # Self-contained React components for docs
│   ├── ApiPlayground.jsx            # Interactive API testing component
│   ├── CodeBlocks.jsx               # Enhanced code blocks with copy/run
│   ├── Diagrams.jsx                 # Architecture diagram components
│   └── ExampleRunner.jsx            # Code example runner
│
├── _theme/                          # Theme customizations
│   ├── custom.css                   # Custom styling
│   └── prism-theme.js               # Code syntax highlighting theme
│
├── _static/                         # Static assets for docs
│   ├── openapi.yaml                 # OpenAPI specification
│   ├── diagrams/                    # Architecture diagrams
│   │   ├── system-overview.svg
│   │   └── deployment-arch.svg
│   ├── images/                      # Documentation images
│   │   ├── screenshots/
│   │   └── icons/
│   └── examples/                    # Example files to download
│       ├── docker-compose.yml
│       └── configuration-template.yaml
│
├── _config/                         # Documentation configuration
│   ├── docusaurus.config.js         # Docusaurus configuration
│   ├── sidebars.js                  # Sidebar navigation structure
│   ├── plugins.js                   # Plugin configurations
│   └── theme.js                     # Theme configuration
│
└── _scripts/                        # Documentation-related scripts
    ├── generate-api-docs.js         # Script to generate API docs from OpenAPI
    ├── validate-examples.js         # Script to validate code examples
    └── update-diagrams.js           # Script to update diagrams from source
