#!/usr/bin/env python3
"""Import the test ISO into Cobbler via XML-RPC.

Cobbler 4.0.0 no longer ships the Python `cobbler` CLI (it has been replaced
by the Go CLI at https://github.com/cobbler/cli), so the container startup can
no longer shell out to `cobbler import ...`. This performs the equivalent call
directly against cobblerd's XML-RPC API instead.
"""

import sys
import time
import xmlrpc.client

SERVER_URL = "http://localhost:25151/"
USERNAME = "cobbler"
PASSWORD = "cobbler"
POLL_INTERVAL_SECONDS = 5
TIMEOUT_SECONDS = 600


def main() -> int:
    server = xmlrpc.client.ServerProxy(SERVER_URL, allow_none=True)
    token = server.login(USERNAME, PASSWORD)

    event_id = server.background_import(
        {
            "name": "Ubuntu-20.04",
            "breed": "ubuntu",
            "path": "/extracted_iso_image/",
            # Cobbler's distro signatures have no default_autoinstall configured for the
            # "focal" (20.04) version, so the automatic template lookup fails unless one
            # is supplied explicitly. The built-in preseed template used by older Ubuntu
            # releases works fine here too.
            "autoinstall_file": "built-in-sample.seed",
        },
        token,
    )

    deadline = time.time() + TIMEOUT_SECONDS
    while time.time() < deadline:
        _statetime, _name, state, _read_by_who = server.get_task_status(event_id)
        if state == "complete":
            print(f"Import task {event_id} completed successfully.")
            return 0
        if state == "failed":
            print(f"Import task {event_id} failed:", file=sys.stderr)
            print(server.get_event_log(event_id), file=sys.stderr)
            return 1
        time.sleep(POLL_INTERVAL_SECONDS)

    print(f"Timed out after {TIMEOUT_SECONDS}s waiting for import task {event_id}.", file=sys.stderr)
    return 1


if __name__ == "__main__":
    sys.exit(main())
