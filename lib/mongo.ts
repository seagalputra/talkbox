import mongoose, { mongo } from "mongoose";

declare let global: { mongoose: any };

const DATABASE_URL = process.env.DATABASE_URL || "";

if (!DATABASE_URL) {
  throw new Error("Please define DATABASE_URL inside .env file or system env");
}

let mongooseCache = global.mongoose;

if (!mongooseCache) {
  mongooseCache = global.mongoose = {
    conn: null,
    promise: null,
  };
}

async function dbConnect() {
  if (mongooseCache.conn) {
    return mongooseCache.conn;
  }

  if (!mongooseCache.promise) {
    const opts = {
      bufferCommands: false,
    };

    mongooseCache.promise = mongoose
      .connect(DATABASE_URL, opts)
      .then((mongoose) => mongoose);
  }

  try {
    mongooseCache.conn = await mongooseCache.promise;
  } catch (err) {
    mongooseCache.promise = null;
    throw err;
  }

  return mongooseCache.conn;
}

export default dbConnect;
