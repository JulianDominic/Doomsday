FROM python:3.12-alpine

WORKDIR /app

# Copy the dependencies file to the working directory
COPY requirements.txt .

RUN pip install --upgrade pip
RUN pip install -r requirements.txt

# Copy the content of the local directory to the working directory
COPY . .

CMD ["python", "main.py"]
