#!/usr/bin/env node
import { spawnSync } from "node:child_process";
import { chmodSync, copyFileSync, cpSync, existsSync, mkdirSync, readFileSync, rmSync, statSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");

const installGitHook = () => {
  const gitDir = join(root, ".git");
  const hookSrc = join(root, ".githooks/pre-commit");
  if (!existsSync(gitDir) || !existsSync(hookSrc)) return;
  if (statSync(gitDir).isDirectory()) {
    const hooksDir = join(gitDir, "hooks");
    mkdirSync(hooksDir, { recursive: true });
    const dest = join(hooksDir, "pre-commit");
    copyFileSync(hookSrc, dest);
    try {
      chmodSync(dest, 0o755);
    } catch {
      // Windows may ignore the executable bit; Git still runs the hook.
    }
  }
  spawnSync("git", ["config", "core.hooksPath", ".githooks"], { cwd: root, stdio: "ignore" });
};
const specPath = join(root, ".project/template.json");
if (!existsSync(specPath)) {
  console.error("build-template: missing .project/template.json");
  process.exit(1);
}

const spec = JSON.parse(readFileSync(specPath, "utf8"));
const out = join(root, spec.out || "template");

if (existsSync(out)) rmSync(out, { recursive: true, force: true });
mkdirSync(out, { recursive: true });

for (const rel of spec.copy ?? []) {
  const from = join(root, rel);
  const to = join(out, rel);
  if (!existsSync(from)) {
    console.error(`build-template: missing ${rel}`);
    process.exit(1);
  }
  mkdirSync(dirname(to), { recursive: true });
  cpSync(from, to, { recursive: true });
}

for (const [dest, src] of Object.entries(spec.slots ?? {})) {
  const from = join(root, src);
  const to = join(out, dest);
  if (!existsSync(from)) {
    console.error(`build-template: missing slot ${src}`);
    process.exit(1);
  }
  mkdirSync(dirname(to), { recursive: true });
  cpSync(from, to);
}

installGitHook();
console.log(`build-template: wrote ${spec.out || "template"}`);
