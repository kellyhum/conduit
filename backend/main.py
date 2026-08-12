import sqlite3

from fastapi import FastAPI, Form, Request

app = FastAPI()
DB_NAME = 'conduit.db'

# initialize users and files databases if not exists
def init_db():
    conn = sqlite3.connect(DB_NAME)
    c = conn.cursor()
    c.execute('''CREATE TABLE IF NOT EXISTS users (
                    username TEXT PRIMARY KEY,
                    user_id TEXT,
                    public_key TEXT,
                    encryption_public_key TEXT,
                    ip_address TEXT)''')
    c.execute('''CREATE TABLE IF NOT EXISTS files (
                    filename TEXT,
                    owner_username TEXT)''')
    conn.commit()
    conn.close()

init_db()

### FASTAPI ###
# /saveuser = save username / id / public key / ip address
# conduit 'setup' command
@app.post("/saveuser")
async def register(request: Request,
                   username: str = Form(...),
                   user_id: str = Form(...),
                   public_key: str = Form(...),
                   encryption_public_key: str = Form(...)):
    ip_addr = request.client.host

    conn = sqlite3.connect(DB_NAME)
    c = conn.cursor()
    c.execute("INSERT OR REPLACE INTO users VALUES (?, ?, ?, ?, ?)", (username, user_id, public_key, encryption_public_key, ip_addr))
    conn.commit()
    conn.close()
    return {"status": f"user {username} synced to database"}

# /getuser = return person b's ip / public key
# conduit 'transfer' command
@app.get("/getuser/{username}")
async def get_user_info(username: str):
    conn = sqlite3.connect(DB_NAME)
    c = conn.cursor()
    c.execute("SELECT public_key, encryption_public_key, ip_address FROM users WHERE username=?", (username,))
    user = c.fetchone()
    conn.close()

    if user:
        return {"username": username,
                "public_key": user[0],
                "encryption_public_key": user[1],
                "ip_address": user[2]}

    return {"error": "user not found"}

# /uplaod = record file exists in sql database
# conduit 'upload' command
@app.post("/upload")
async def upload(filename: str = Form(...),
                 owner_username: str = Form(...)):
    conn = sqlite3.connect(DB_NAME)
    c = conn.cursor()
    c.execute("INSERT INTO files VALUES (?, ?)", (filename, owner_username))
    conn.commit()
    conn.close()
    return {"status": f"file {filename} uploaded by {owner_username}"}
