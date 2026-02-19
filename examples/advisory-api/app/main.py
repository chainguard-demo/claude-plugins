from flask import Flask, jsonify

app = Flask(__name__)

ADVISORIES = [
    {
        "id": "CVE-2024-1234",
        "severity": "HIGH",
        "package": "openssl",
        "fixed_in": "3.0.8",
        "description": "Buffer overflow in TLS handshake processing",
    },
    {
        "id": "CVE-2024-5678",
        "severity": "MEDIUM",
        "package": "curl",
        "fixed_in": "8.4.0",
        "description": "HSTS bypass via redirect",
    },
    {
        "id": "CVE-2024-9012",
        "severity": "LOW",
        "package": "zlib",
        "fixed_in": "1.3.1",
        "description": "Integer overflow in inflate",
    },
]


@app.route("/")
def index():
    return jsonify({"service": "advisory-api", "status": "ok", "version": "1.0.0"})


@app.route("/advisories")
def get_advisories():
    return jsonify(ADVISORIES)


@app.route("/advisories/<cve_id>")
def get_advisory(cve_id):
    match = next((a for a in ADVISORIES if a["id"] == cve_id), None)
    if not match:
        return jsonify({"error": "not found"}), 404
    return jsonify(match)
