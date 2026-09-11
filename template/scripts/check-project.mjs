#!/usr/bin/env node
import { spawnSync } from "node:child_process";
import { existsSync, readdirSync, readFileSync, statSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const repoRoot = join(dirname(fileURLToPath(import.meta.url)), "..");
const root = resolve(process.argv[2] ? join(repoRoot, process.argv[2]) : repoRoot);
const inCI = Boolean(process.env.CI || process.env.GITHUB_ACTIONS);
const errors = [];

if (existsSync(join(root, ".project/template.json")) && !inCI) {
  const built = spawnSync(process.execPath, [join(repoRoot, "scripts/build-template.mjs")], {
    cwd: repoRoot,
    stdio: "inherit",
  });
  if (built.status) process.exit(built.status ?? 1);
}

const say = (ok, msg) => {
  if (!ok) errors.push(msg);
};

const read = (rel) => {
  const path = join(root, rel);
  if (!existsSync(path)) {
    errors.push(`missing file: ${rel}`);
    return null;
  }
  return readFileSync(path, "utf8");
};

const json = (rel) => {
  const text = read(rel);
  if (text == null) return null;
  try {
    return JSON.parse(text);
  } catch (err) {
    errors.push(`${rel} is not valid JSON: ${err.message}`);
    return null;
  }
};

for (const rel of [
  "AGENTS.md",
  "README.md",
  "docs/engineering/ENGINEERING.md",
  "docs/product/PRODUCT.md",
  "docs/architecture/ARCHITECTURE.md",
  "docs/work/CURRENT.md",
]) {
  read(rel);
}

for (const rel of [".project/queue.json", ".project/rules.json"]) {
  if (existsSync(join(root, rel))) errors.push(`forbidden file: ${rel}`);
}

const state = json(".project/state.json");
const documents = json(".project/documents.json");
const ownership = json(".project/ownership.json");
const routes = json(".project/context-routes.json");

if (state) {
  say(state.os_version === 2, `state.json os_version must be 2, got ${state.os_version}`);
  say(typeof state.phase === "string" && state.phase.length > 0, "state.json missing phase");
  say(
    state.current_work && typeof state.current_work.title === "string",
    "state.json missing current_work.title",
  );
  say(typeof state.last_updated === "string", "state.json missing last_updated");
}

const docsById = new Map();
if (documents?.documents) {
  for (const doc of documents.documents) {
    if (!doc.id) {
      errors.push("documents.json entry missing id");
      continue;
    }
    if (docsById.has(doc.id)) errors.push(`duplicate document id: ${doc.id}`);
    docsById.set(doc.id, doc);
    if (!doc.path) errors.push(`document ${doc.id} missing path`);
    else if (!existsSync(join(root, doc.path))) {
      errors.push(`document ${doc.id} path does not exist: ${doc.path}`);
    }
    if (doc.path && doc.path.startsWith("modules/")) {
      errors.push(`document ${doc.id} path is under modules/ (not authority until copied to docs/)`);
    }
  }
  for (const id of ["product", "architecture", "engineering", "current-work"]) {
    say(docsById.has(id), `documents.json missing core id: ${id}`);
  }
}

const factIds = new Set();
if (ownership?.facts) {
  for (const fact of ownership.facts) {
    if (!fact.id) {
      errors.push("ownership.json fact missing id");
      continue;
    }
    if (factIds.has(fact.id)) errors.push(`duplicate fact id: ${fact.id}`);
    factIds.add(fact.id);
    if (!docsById.has(fact.knowledge_owner)) {
      errors.push(
        `fact ${fact.id} knowledge_owner "${fact.knowledge_owner}" is not a documents.json id`,
      );
    }
  }
}

if (documents?.documents) {
  for (const doc of documents.documents) {
    for (const auth of doc.authority ?? []) {
      if (!factIds.has(auth)) {
        errors.push(`document ${doc.id} authority "${auth}" is not an ownership fact id`);
      }
    }
  }
}

if (routes) {
  const checkIds = (list, where, required) => {
    if (!Array.isArray(list)) {
      errors.push(`${where} must be an array`);
      return;
    }
    for (const id of list) {
      if (required && !docsById.has(id)) errors.push(`${where} unknown document id: ${id}`);
    }
  };
  checkIds(routes.always_read, "always_read", true);
  if (Array.isArray(routes.always_read) && routes.always_read.length > 2) {
    errors.push("always_read should stay small (engineering only, or engineering + one more)");
  }
  if (!Array.isArray(routes.routes)) errors.push("context-routes.json missing routes");
  else {
    for (const route of routes.routes) {
      checkIds(route.required, `route ${route.task} required`, true);
      checkIds(route.optional ?? [], `route ${route.task} optional`, false);
    }
  }
}

const adrDir = join(root, "docs/decisions");
if (existsSync(adrDir)) {
  const files = readdirSync(adrDir).filter((name) => name.endsWith(".md") && name !== "README.md");
  const numbers = new Map();
  for (const name of files) {
    const match = name.match(/^(\d{4})-/);
    if (!match) {
      errors.push(`ADR filename must be NNNN-slug.md: ${name}`);
      continue;
    }
    if (numbers.has(match[1])) {
      errors.push(`duplicate ADR number ${match[1]}: ${numbers.get(match[1])} and ${name}`);
    }
    numbers.set(match[1], name);
  }
}

const productPath = join(root, "docs/product/PRODUCT.md");
const productText = existsSync(productPath) ? readFileSync(productPath, "utf8") : "";
if (productText && !/^## 当前未决事项\s*$/m.test(productText)) {
  errors.push("docs/product/PRODUCT.md must have heading ## 当前未决事项 (sole OPEN list)");
}

const walkMd = (dir, acc = []) => {
  if (!existsSync(dir)) return acc;
  for (const name of readdirSync(dir)) {
    const p = join(dir, name);
    if (statSync(p).isDirectory()) walkMd(p, acc);
    else if (name.endsWith(".md")) acc.push(p);
  }
  return acc;
};

const relFromRoot = (abs) => abs.slice(root.length + 1).replaceAll("\\", "/");
const allowedDocsMd = new Set(["docs/decisions/README.md"]);
if (documents?.documents) {
  for (const doc of documents.documents) {
    if (doc.path) allowedDocsMd.add(doc.path.replaceAll("\\", "/"));
  }
}

for (const file of walkMd(join(root, "docs"))) {
  const rel = relFromRoot(file);
  const inDecisions = rel.startsWith("docs/decisions/");
  const adrName = inDecisions ? rel.slice("docs/decisions/".length) : "";
  const allowedAdr = inDecisions && !adrName.includes("/") && /^\d{4}-.+\.md$/.test(adrName);
  if (!allowedDocsMd.has(rel) && !allowedAdr) {
    errors.push(`unregistered docs file: ${rel}`);
  }
  if (rel === "docs/product/PRODUCT.md") continue;
  const text = readFileSync(file, "utf8");
  if (/^## 当前未决事项\s*$/m.test(text)) {
    errors.push(`${rel} must not have ## 当前未决事项; product OPEN list lives only in docs/product/PRODUCT.md`);
  }
}

const currentPath = join(root, "docs/work/CURRENT.md");
if (existsSync(currentPath)) {
  const currentText = readFileSync(currentPath, "utf8");
  if (!currentText.includes("docs/product/PRODUCT.md")) {
    errors.push("docs/work/CURRENT.md must point OPEN list to docs/product/PRODUCT.md, not copy it");
  }
}

for (const name of ["CLAUDE.md", ".github/copilot-instructions.md"]) {
  const path = join(root, name);
  if (!existsSync(path) || !statSync(path).isFile()) continue;
  const text = readFileSync(path, "utf8");
  if (!/AGENTS\.md/.test(text)) {
    errors.push(`${name} must point to AGENTS.md (one-line pointer, not a second rulebook)`);
  }
  if (text.length > 400) {
    errors.push(`${name} is too long; keep a one-line pointer to AGENTS.md`);
  }
}

if (errors.length) {
  console.error(`check-project: ${errors.length} error(s)`);
  for (const err of errors) console.error(`- ${err}`);
  process.exit(1);
}

if (existsSync(join(root, ".project/template.json")) && !process.argv[2]) {
  const spec = JSON.parse(readFileSync(join(root, ".project/template.json"), "utf8"));
  const out = spec.out || "template";
  if (!existsSync(join(root, out))) {
    console.error("check-project: 1 error(s)");
    console.error(`- missing generated ${out}/; commit it with the source`);
    process.exit(1);
  }
  const nested = spawnSync(process.execPath, [fileURLToPath(import.meta.url), out], {
    cwd: repoRoot,
    stdio: "inherit",
  });
  if (nested.status) process.exit(nested.status);
}

console.log("check-project: ok");
