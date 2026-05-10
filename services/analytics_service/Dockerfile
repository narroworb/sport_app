FROM python:3.11-slim

WORKDIR /app

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PYTHONPATH=/app

# Install dependencies
COPY services/analytics_service/requirements.txt /app/requirements.txt
RUN pip install --no-cache-dir -r /app/requirements.txt

# Copy service code
COPY services/analytics_service/app /app/app

# Copy shared proto and generate stubs
COPY proto /app/proto

RUN python -m grpc_tools.protoc \
    -I/app/proto \
    --python_out=/app/app/gen \
    --grpc_python_out=/app/app/gen \
    /app/proto/analytics/v1/analytics.proto

# Create __init__.py files
RUN touch /app/app/gen/__init__.py && \
    touch /app/app/gen/analytics/__init__.py && \
    touch /app/app/gen/analytics/v1/__init__.py

EXPOSE 50051

CMD ["python", "-m", "app.server"]