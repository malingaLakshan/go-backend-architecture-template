cd /data/sim

./venv/bin/python3 - <<'PY'
import messages_pb2

bundle = messages_pb2.ProtoReaderBundle.DESCRIPTOR

print("ALL BUNDLE FIELDS:")
for field in bundle.fields:
    print(f" - Field {field.number}: {field.name}")

reads_field = bundle.fields_by_name["reads"]

print("\nALL INDIVIDUAL READ FIELDS:")
for field in reads_field.message_type.fields:
    print(f" - Field {field.number}: {field.name}")
PY