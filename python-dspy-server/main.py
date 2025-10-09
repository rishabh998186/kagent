from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Optional, Dict, Any
import dspy
import os

app = FastAPI(title="DSPy Compilation Service")

# Initialize DSPy with environment variables
def init_dspy():
    """Initialize DSPy with the configured LLM"""
    api_key = os.getenv("OPENAI_API_KEY")
    model = os.getenv("DSPY_MODEL", "gpt-4")
    
    if api_key:
        lm = dspy.LM(model=f"openai/{model}", api_key=api_key)
        dspy.configure(lm=lm)

init_dspy()

class SignatureFieldModel(BaseModel):
    name: str
    type: str = "string"
    description: Optional[str] = None
    prefix: Optional[str] = None

class CompileRequest(BaseModel):
    inputs: List[SignatureFieldModel]
    outputs: List[SignatureFieldModel]
    instructions: Optional[str] = None
    module: str = "ChainOfThought"

class CompileResponse(BaseModel):
    compiled_prompt: str
    signature_dict: Dict[str, Any]
    module_type: str

@app.get("/health")
def health_check():
    """Health check endpoint"""
    return {"status": "healthy"}

@app.post("/compile", response_model=CompileResponse)
def compile_prompt(request: CompileRequest):
    """
    Compile a DSPy signature into a prompt
    """
    try:
        # Build signature dynamically
        signature_class = build_signature(request)
        
        # Get the appropriate DSPy module
        module_instance = get_module(request.module, signature_class)
        
        # Generate the compiled prompt
        compiled_prompt = generate_prompt_from_module(module_instance)
        
        return CompileResponse(
            compiled_prompt=compiled_prompt,
            signature_dict={
                "inputs": [f.dict() for f in request.inputs],
                "outputs": [f.dict() for f in request.outputs],
                "instructions": request.instructions
            },
            module_type=request.module
        )
    
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

def build_signature(request: CompileRequest):
    """Build a DSPy signature from the request"""
    
    # Create signature fields
    fields = {}
    
    # Add input fields
    for field in request.inputs:
        desc = field.description or f"Input field: {field.name}"
        fields[field.name] = dspy.InputField(desc=desc, prefix=field.prefix)
    
    # Add output fields
    for field in request.outputs:
        desc = field.description or f"Output field: {field.name}"
        fields[field.name] = dspy.OutputField(desc=desc, prefix=field.prefix)
    
    # Create the signature class dynamically
    signature = type(
        "DynamicSignature",
        (dspy.Signature,),
        {
            "__doc__": request.instructions or "Dynamic DSPy Signature",
            **fields
        }
    )
    
    return signature

def get_module(module_type: str, signature):
    """Get the appropriate DSPy module instance"""
    
    module_map = {
        "Predict": dspy.Predict,
        "ChainOfThought": dspy.ChainOfThought,
        "ReAct": dspy.ReAct
    }
    
    if module_type not in module_map:
        raise ValueError(f"Unknown module type: {module_type}")
    
    ModuleClass = module_map[module_type]
    return ModuleClass(signature)

def generate_prompt_from_module(module):
    """Generate a prompt string from a DSPy module"""
    
    # For basic compilation, we extract the signature instructions
    # In a real scenario, DSPy modules compile to prompts through their forward() calls
    
    signature = module.signature
    
    # Build prompt template
    prompt_parts = []
    
    # Add instructions
    if hasattr(signature, '__doc__') and signature.__doc__:
        prompt_parts.append(signature.__doc__)
    
    # Add input fields
    prompt_parts.append("\n--- Inputs ---")
    for name, field in signature.input_fields.items():
        desc = getattr(field, 'desc', name)
        prompt_parts.append(f"{name}: {desc}")
    
    # Add output fields
    prompt_parts.append("\n--- Outputs ---")
    for name, field in signature.output_fields.items():
        desc = getattr(field, 'desc', name)
        prompt_parts.append(f"{name}: {desc}")
    
    return "\n".join(prompt_parts)

if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("PORT", "8000"))
    uvicorn.run(app, host="0.0.0.0", port=port)
