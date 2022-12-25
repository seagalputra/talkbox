module.exports = {
  apps: [
    {
      name: "talkbox",
      script: "./server.js",
    },
    {
      name: "talkbox-api",
      script: "./talkbox",
    },
  ],
};
