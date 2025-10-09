# python-dspy-server

A Python FastAPI microservice for compiling and optimizing DSPy agent prompt configurations.  
This backend is designed to power the DSPy functionality in the kagent project.

## Features

- Compile DSPy prompt/signature configurations to text templates
- Run DSPy module optimization (planned)
- Simple HTTP API for integration with kagent

## Quickstart

### Run Locally

python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
cp .env.example .env # Edit with your credentials
python main.py


### Run with Docker

docker build -t kagent-dspy-server .
docker run --env-file .env -p 8000:8000 kagent-dspy-server

## API Endpoints

- `POST /compile`: Compile a DSPy signature to a prompt template
- `POST /optimize`: Run prompt optimization (partial/mock)
- `GET /health`: Service health check

## Configuration

Set required configuration in `.env` (see `.env.example`):

- `DSPY_LM_PROVIDER` (e.g., `openai`)
- `DSPY_MODEL` (e.g., `gpt-3.5-turbo`)
- `OPENAI_API_KEY` (if using OpenAI backend)
- `PORT` (default: 8000)

## Contributing

PRs and suggestions are welcome. Please follow code style and document your changes.

## License

[MIT License](../LICENSE)
EOF