import os
import requests
from flask import Flask, request, jsonify
import io
import contextlib
import multiprocessing
import math
import json
import datetime

from loguru import logger

app = Flask(__name__)

BRAVE_API_KEY = os.getenv("BRAVE_API_KEY")

def execute_code_worker(code, result_queue):
    output = io.StringIO()
    error_output = io.StringIO()

    # Pre-defined globals for the execution environment
    safe_globals = {
        '__builtins__': __builtins__,
        'math': math,
        'json': json,
        'datetime': datetime,
        'os': None, # Disable direct OS access for some basic safety
    }

    try:
        with contextlib.redirect_stdout(output), contextlib.redirect_stderr(error_output):
            exec(code, safe_globals)
        result_queue.put({
            'result': output.getvalue(),
            'error': error_output.getvalue()
        })
    except Exception as e:
        result_queue.put({
            'result': output.getvalue(),
            'error': f"{error_output.getvalue()}\n{str(e)}"
        })

@app.route("/run", methods=["POST"])
def run_code():
    if request.json is None:
        logger.info('no json')
        return jsonify({'error': "no json provided"})

    code = request.json.get('code', "")

    if not code:
        return jsonify({'error': "no code provided"})

    result_queue = multiprocessing.Queue()
    process = multiprocessing.Process(target=execute_code_worker, args=(code, result_queue))

    try:
        process.start()
        # Wait for up to 10 seconds for the code to finish
        process.join(timeout=10)

        if process.is_alive():
            process.terminate()
            process.join()
            return jsonify({'error': "Execution timed out (10s limit)"})

        if result_queue.empty():
            return jsonify({'error': "No output received from execution process"})

        res = result_queue.get()

        # Combine stdout and stderr for the model to see everything
        combined_output = res['result']
        if res['error'].strip():
            combined_output += f"\n--- Errors/Warnings ---\n{res['error']}"

        return jsonify({'result': combined_output})

    except Exception as e:
        logger.error(f"Execution wrapper failed: {e}")
        return jsonify({'error': str(e)})


@app.route("/websearch", methods=["POST"])
def web_search():
    if request.json is None:
        return jsonify({'error': "no json provided"})

    query = request.json.get('query', "")

    logger.info(query)
    if not query:
        logger.info("no query provided")
        return jsonify({'error': "no query provided"})

    if not BRAVE_API_KEY:
        logger.error("BRAVE_API_KEY not set")
        return jsonify({'error': "BRAVE_API_KEY not set in environment"})

    try:
        headers = {
            "Accept": "application/json",
            "Accept-Encoding": "gzip",
            "X-Subscription-Token": BRAVE_API_KEY
        }
        params = {"q": query}
        response = requests.get(
            "https://api.search.brave.com/res/v1/web/search",
            headers=headers,
            params=params
        )
        response.raise_for_status()
        data = response.json()
        
        results = []
        if "web" in data and "results" in data["web"]:
            for res in data["web"]["results"][:5]: # top 5 results
                results.append(f"Title: {res.get('title')}\nLink: {res.get('url')}\nDescription: {res.get('description')}\n")
        
        formatted_results = "\n".join(results)
        if not formatted_results:
            formatted_results = "No results found."
            
        return jsonify({'result': formatted_results})
    except Exception as e:
        logger.error(f"Error during web search: {e}")
        return jsonify({'error': f"Error during web search: {e}"})


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5500, debug=True)
