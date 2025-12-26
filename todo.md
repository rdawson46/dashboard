# To Do
## Immediate
* server send new desc for new chat with new id
    - chat transmits a new event out to the main page
    - main page will then pass the event to sidebar
    - sidebar consumes
         * trigger some kind of animation
* server end:
    - messages will now have to pull file content data for referenced file IDs

## Delayed
* queue around ollama to prevent overloading
* job registration
    - code job creation
    - handle functions jobs
        - logging within handle functions
* notification/background job style
* queue system for async tasks
    - tool execution
    - notification
* loading animation gets stuck on tool calls
* display:
    - reasoning/thinking
* tools (**web search** and code execution)
    - improve code execution prompting
* MCP
* monitoring system
* image support
* vector DB/RAG
* chat page styling 
    - long messages in input
    - sizing
    - background
    - scrolling while streaming
* db
    - tokens usage
* openai/anthropic/gemini integration
    - openai will require more work with various api - GPT5 - response api
