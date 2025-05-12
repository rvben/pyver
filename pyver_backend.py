#!/usr/bin/env python3
import sys
import json
from packaging.version import Version

def main():
    if len(sys.argv) < 3:
        print("Usage: pyver_backend.py <command> <version(s)>", file=sys.stderr)
        sys.exit(1)
    cmd = sys.argv[1]
    if cmd == "compare":
        v1 = Version(sys.argv[2])
        v2 = Version(sys.argv[3])
        if v1 < v2:
            print(-1)
        elif v1 > v2:
            print(1)
        else:
            print(0)
    elif cmd == "parse":
        v = Version(sys.argv[2])
        # Output as JSON for possible future use
        print(json.dumps({
            "public": v.public,
            "base_version": v.base_version,
            "is_prerelease": v.is_prerelease,
            "is_postrelease": v.is_postrelease,
            "is_devrelease": v.is_devrelease,
            "epoch": v.epoch,
            "release": v.release,
            "pre": v.pre,
            "post": v.post,
            "dev": v.dev,
            "local": v.local,
        }))
    else:
        print(f"Unknown command: {cmd}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()