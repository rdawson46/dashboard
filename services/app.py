from flask import Flask, request, jsonify
import io
import contextlib

from loguru import logger

app = Flask(__name__)

@app.route("/run", methods=["POST"])
def run_code():
    if request.json is None:
        logger.info('no json')
        return jsonify({'error': "no json provided"})

    code = request.json.get('code', "")

    if not code:
        return jsonify({'error': "no code provided"})

    output = io.StringIO()

    # safe_globals = {'__builtins__': {'print': print}}
    try:
        with contextlib.redirect_stdout(output):
            # exec(code, safe_globals)
            exec(code)
        return jsonify({'result': output.getvalue()})
    except Exception as e:
        logger.error(e)
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

    try:
        raise Exception('Not implemented yet')
        return jsonify({'result': results})
    except Exception as e:
        return jsonify({'error': f"Error during web search: {e}"})


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5500, debug=True)
