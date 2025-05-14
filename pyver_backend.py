#!/usr/bin/env python3
import sys
import json
from packaging.version import Version, InvalidVersion

def format_pre(pre):
    if pre is None:
        return ""
    return f"{pre[0]}{pre[1]}"

def format_post(post):
    if post is None:
        return ""
    return f"post{post}"

def format_dev(dev):
    if dev is None:
        return ""
    return f"dev{dev}"

def format_local(local):
    if local is None:
        return ""
    if isinstance(local, tuple):
        return ".".join(str(x) for x in local)
    return str(local)

def main():
    if len(sys.argv) < 3:
        print("Usage: pyver_backend.py <command> <version(s)>", file=sys.stderr)
        sys.exit(1)
    cmd = sys.argv[1]
    try:
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
            print(json.dumps({
                "normalized": str(v),
                "public": v.public,
                "base_version": v.base_version,
                "is_prerelease": v.is_prerelease,
                "is_postrelease": v.is_postrelease,
                "is_devrelease": v.is_devrelease,
                "epoch": v.epoch,
                "release": v.release,
                "pre": format_pre(v.pre),
                "post": format_post(v.post),
                "dev": format_dev(v.dev),
                "local": format_local(v.local),
            }))
        else:
            print(f"Unknown command: {cmd}", file=sys.stderr)
            sys.exit(1)
    except InvalidVersion as e:
        print(f"Invalid version: {e}", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()