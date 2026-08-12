cd /data/sim
./venv/bin/python3 - <<'PY'
from google.protobuf import descriptor_pb2, text_format
import messages_pb2

schema = descriptor_pb2.FileDescriptorProto()
messages_pb2.DESCRIPTOR.CopyToProto(schema)
print(text_format.MessageToString(schema))
PY