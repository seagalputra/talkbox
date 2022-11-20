import { MongoClient, MongoClientOptions, Db } from "mongodb";
import { config } from "./app";
declare let global: { mongo: Db };

const connectionUrl = config.databaseUrl ?? "";
let cachedMongo: Db;

async function connectToDb(
  url: string,
  options: MongoClientOptions | undefined
) {
  const connection = await new MongoClient(url, options).connect();
  return connection.db("talkbox");
}

async function getDb() {
  if (import.meta.env.MODE === "production") {
    const mongo = await connectToDb(connectionUrl, {});
    return mongo;
  } else {
    if (!global.mongo) {
      global.mongo = await connectToDb(connectionUrl, {});
      cachedMongo = global.mongo;
    }
    return cachedMongo;
  }
}

export { getDb };
