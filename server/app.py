from flask import Flask, request, jsonify
import sqlite3
import os

DB_PATH = os.path.join(os.path.dirname(__file__), "database", "users.db")
app = Flask(__name__)

def connect_db():
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row
    return conn

def init_db():
    with open(os.path.join(os.path.dirname(__file__), "database", "setup.sql")) as f:
        conn = connect_db()
        conn.executescript(f.read())
        conn.commit()
        conn.close()

init_db()

@app.route("/register", methods=["POST"])
def register():
    data = request.json
    username = data["username"]
    pubkey = data["public_key"].encode("utf-8")
    conn = connect_db()

    try:
        conn.execute("INSERT INTO users (username, public_key) VALUES (?, ?)", (username, pubkey))
        conn.commit()
    except sqlite3.IntegrityError:
        return jsonify({"error": "Username already exists"}), 400
    finally:
        conn.close()

    return jsonify({"message": "User created"}), 201

@app.route("/verify", methods=["POST"])
def verify():
    data = request.json
    username = data["username"]

    conn = connect_db()
    try:
        # Fetch user info
        cur = conn.execute("SELECT * FROM users WHERE username = ?", (username,))
        user = cur.fetchone()
        if not user:
            return jsonify({"error": "Username not found"}), 404

        # Fetch user files
        cur = conn.execute("SELECT file_name FROM files WHERE user_id = ?", (user["id"],))
        files = [row["file_name"] for row in cur.fetchall()]
    finally:
        conn.close()

    return jsonify({
        "username": user["username"],
        "public_key": user["public_key"].decode("utf-8"),
        "files": files
    })
    
if __name__ == "__main__":
    app.run(debug=True)
