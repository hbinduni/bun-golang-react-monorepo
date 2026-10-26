#!/usr/bin/env python3
"""Check declared Go compiler versions without running application code."""
import pathlib
import re
import subprocess

ROOT = pathlib.Path(__file__).resolve().parents[1]
MODULES = ['server']
errors = []
versions = {}
for module in MODULES:
    text = (ROOT / module / "go.mod").read_text()
    minimum = re.search(r"^go (\S+)", text, re.M).group(1)
    preferred = re.search(r"^toolchain go(\S+)", text, re.M)
    selected = preferred.group(1) if preferred else minimum
    versions[module] = selected
    if tuple(map(int, selected.split("."))) < tuple(map(int, minimum.split("."))):
        errors.append(f"{module}: toolchain {selected} is older than minimum {minimum}")
    print(f"{module}/go.mod: minimum {minimum}, scan compiler {selected}")
files = subprocess.check_output(["git", "ls-files"], cwd=ROOT, text=True).splitlines()
for filename in files:
    path = pathlib.PurePosixPath(filename)
    if path.name.startswith("Dockerfile"):
        candidates = [m for m in MODULES if m == "." or str(path.parent) == m or str(path.parent).startswith(m + "/")]
        if not candidates:
            continue
        module = max(candidates, key=len)
        for tag in re.findall(r"(?im)^FROM\s+(?:--platform=\S+\s+)?golang:([^\s]+)", (ROOT / filename).read_text()):
            version = tag.split("-")[0].split("@")[0]
            if version != versions[module]:
                errors.append(f"{filename}: Docker tag {tag} differs from explicit compiler {versions[module]} (floating tags cannot prove patch consistency)")
for path in (ROOT / ".github/workflows").glob("*.y*ml"):
    text = path.read_text()
    for version in re.findall(r"go-version:\s*['\"]?([0-9]+(?:\.[0-9]+){1,2})", text):
        if version not in versions.values():
            errors.append(f"{path.relative_to(ROOT)}: explicit CI version {version} differs from declared module toolchains")
    for manifest in re.findall(r"go-version-file:\s*['\"]?([^'\"\s]+)", text):
        if "${{" not in manifest and not (ROOT / manifest).is_file():
            errors.append(f"{path.relative_to(ROOT)}: missing Go version file {manifest}")
for error in errors:
    print("ERROR:", error)
raise SystemExit(bool(errors))
