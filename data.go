cd /data/sim

./venv/bin/python3 - <<'PY'
import messages_pb2

bundle_schema = messages_pb2.ProtoReaderBundle.DESCRIPTOR

print("PROTOREADERBUNDLE FIELDS:")
for field in bundle_schema.fields:
    print(f"Field {field.number}: {field.name}")

reads_schema = bundle_schema.fields_by_name["reads"].message_type

print("\nINDIVIDUAL RAW READ FIELDS:")
for field in reads_schema.fields:
    print(f"Field {field.number}: {field.name}")
PY