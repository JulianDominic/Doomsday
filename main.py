from flask import Flask, render_template, url_for, request, redirect
import random
import datetime
import json

app = Flask(__name__)
RANDOM_DATE = [None, None]
GUESS = [None]
PREVIOUS_DATE = [None]
with open("./data/doomsday.json", "r") as f:
    DOOSMDAY_JSON = json.load(f)

def get_random_date():
    """Generate a random date between the start & end dates."""
    start = datetime.date(2000, 1, 1)
    end = datetime.date(2038, 1, 19)
    random_date = start + (end - start) * random.random()
    RANDOM_DATE[0] = random_date
    RANDOM_DATE[1] = random_date.strftime('%A')
    return random_date.strftime(r'%d-%B-%Y')


@app.route("/", methods=["GET", "POST"])
def index():
    """
    This returns the main page where the user can guess the day.
    'POST' is used ONLY when the user submits the day,
    and displays the message that corresponds to the user's guess.
    """
    random_date = get_random_date()
    if request.method == "POST":
        return render_template("index.html",
                                random_date=random_date,
                                correct=GUESS[0],
                                day=PREVIOUS_DATE[0])
    elif request.method == "GET":
        return render_template("index.html",
                               random_date=random_date,
                               correct=None,
                               day=PREVIOUS_DATE[0])


@app.route("/check", methods=["POST"])
def check_guess():
    GUESS[0] = request.form["guess"] == RANDOM_DATE[1]
    PREVIOUS_DATE[0] = RANDOM_DATE[1]
    # We use code 307 to PRESERVE the request type that was originally sent
    # This is so that the '/' route only renders the "correct" message when
    # the user submits. Previously, there was an issue where the user didn't submit
    # but it maintained the previous' answer, and the message did not go away after a refresh.
    return redirect(url_for('index'), code=307)


@app.route("/cheatsheet")
def cheatsheet(doomsdays:dict=DOOSMDAY_JSON):
    return render_template("cheatsheet.html", doomsdays=doomsdays)

if __name__ == "__main__":
    app.run(host="127.0.0.1",
            port="8000",
            debug=True)
