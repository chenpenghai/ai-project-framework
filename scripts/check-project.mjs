#!/usr/bin/env node
import { existsSync, readdirSync, readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");
const errors = [];

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
  say(state.current_work && typeof state.current_work.title === "string", "state.json missing current_work.title");
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
    else if (!existsSync(join(root, doc.path))) errors.push(`document ${doc.id} path does not exist: ${doc.path}`);
  }
  for (const id of ["product", "architecture", "engineering", "current-work"]) {
    say(docsById.has(id), `documents.json missing core id: ${id}`);
  }
}

if (ownership?.facts) {
  const factIds = new Set();
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

if (errors.length) {
  console.error(`check-project: ${errors.length} error(s)`);
  for (const err of errors) console.error(`- ${err}`);
  process.exit(1);
}

console.log("check-project: ok");
