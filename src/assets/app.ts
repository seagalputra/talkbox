type AppConfig = {
  databaseUrl?: string;
  databaseName?: string;
  jwtSecret?: string;
};

const config: AppConfig = {
  databaseUrl: import.meta.env.DATABASE_URL,
  databaseName: import.meta.env.DATABASE_NAME,
  jwtSecret: import.meta.env.JWT_SECRET,
};

export { config };
