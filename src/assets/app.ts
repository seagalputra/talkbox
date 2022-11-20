type AppConfig = {
  databaseUrl?: string;
  databaseName?: string;
};

const config: AppConfig = {
  databaseUrl: import.meta.env.DATABASE_URL,
  databaseName: import.meta.env.DATABASE_NAME,
};

export { config };
