from flask import Flask, request, jsonify
import io
import contextlib

app = Flask(__name__)

@app.route("/run", methods=["POST"])
def run_code():
    if request.json is None:
        print('no json')
        return jsonify({'error': "no json provided"})

    code = request.json.get('code', "")
    output = io.StringIO()

    # safe_globals = {'__builtins__': {print}}
    try:
        with contextlib.redirect_stdout(output):
            # exec(code, safe_globals)
            exec(code)
        print("Value from code", output.getvalue())
        return jsonify({'result': output.getvalue()})
    except Exception as e:
        print(e)
        return jsonify({'error': str(e)})


@app.route("/websearch", methods=["POST"])
def web_search():
    if request.json is None:
        return jsonify({'error': "no json provided"})

    query = request.json.get('query', "")
    return jsonify({'result': query})


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5500, debug=True)
