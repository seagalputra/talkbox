#!/usr/bin/env node

const path = require("path");
const { Seeder } = require("mongo-seeding");
const { loadEnvConfig } = require("@next/env");

const projectDir = process.cwd();
loadEnvConfig(projectDir);

const config = {
  database: process.env.DATABASE_URL,
};

const seeder = new Seeder(config);

const collections = seeder.readCollectionsFromPath(path.resolve("./data"), {
  transformers: [Seeder.Transformers.replaceDocumentIdWithUnderscoreId],
});

seeder
  .import(collections)
  .then(() => {
    console.log("Success");
  })
  .catch((err) => {
    console.log("Error", err);
  });
