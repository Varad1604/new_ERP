from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Dict, Any
import os

app = FastAPI(title="Enterprise AI Engine", version="1.0.0")

class NLQueryRequest(BaseModel):
    query: str
    context: Dict[str, Any] | None = None

class AIResponse(BaseModel):
    response: str
    metadata: Dict[str, Any]

import openai
import anthropic

# Provider-agnostic AI interface
class AIProviderInterface:
    def execute_prompt(self, prompt: str) -> str:
        raise NotImplementedError

class OpenAIProvider(AIProviderInterface):
    def __init__(self, api_key: str):
        self.client = openai.Client(api_key=api_key)

    def execute_prompt(self, prompt: str) -> str:
        response = self.client.chat.completions.create(
            model="gpt-4",
            messages=[{"role": "user", "content": prompt}],
            temperature=0.0
        )
        return response.choices[0].message.content

class AnthropicProvider(AIProviderInterface):
    def __init__(self, api_key: str):
        self.client = anthropic.Anthropic(api_key=api_key)

    def execute_prompt(self, prompt: str) -> str:
        response = self.client.messages.create(
            model="claude-3-opus-20240229",
            max_tokens=1000,
            messages=[{"role": "user", "content": prompt}]
        )
        return response.content[0].text

class StubAIProvider(AIProviderInterface):
    def execute_prompt(self, prompt: str) -> str:
        return f"Enterprise ERP Stub Response: The AI module received your query: '{prompt}'."

# Dependency Injection for AI Provider
def get_ai_provider() -> AIProviderInterface:
    provider_name = os.getenv("AI_PROVIDER", "stub").lower()

    if provider_name == "openai":
        api_key = os.getenv("OPENAI_API_KEY")
        if not api_key:
            raise ValueError("OPENAI_API_KEY is missing")
        return OpenAIProvider(api_key)

    elif provider_name == "anthropic":
        api_key = os.getenv("ANTHROPIC_API_KEY")
        if not api_key:
            raise ValueError("ANTHROPIC_API_KEY is missing")
        return AnthropicProvider(api_key)

    return StubAIProvider()

@app.get("/health")
def health_check():
    return {"status": "up", "service": "ai-engine"}

@app.post("/api/v1/ai/query", response_model=AIResponse)
def natural_language_query(request: NLQueryRequest):
    provider = get_ai_provider()
    try:
        # In a real scenario, we inject ERP context (e.g., schema info) into the prompt here
        result = provider.execute_prompt(request.query)
        return AIResponse(
            response=result,
            metadata={"provider": "stub", "tokens_used": 0}
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("PORT", 8000))
    uvicorn.run(app, host="0.0.0.0", port=port)
