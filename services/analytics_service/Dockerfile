FROM python:3.11-slim

WORKDIR /app

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

COPY services/analytics_service/requirements.txt /app/requirements.txt
RUN pip install --no-cache-dir -r /app/requirements.txt

# Copy service code
COPY services/analytics_service/app /app/app

# Copy shared proto from repo root and generate python stubs
COPY proto /app/proto
RUN python -m grpc_tools.protoc \
  -I/app/proto \
  --python_out=/app/app/gen \
  --grpc_python_out=/app/app/gen \
  /app/proto/analytics/v1/analytics.proto

EXPOSE 50051

CMD ["python", "-m", "app.server"]

