import axios from "axios";

const http = axios.create({
  baseURL: `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/v1`,
});

http.interceptors.request.use((config) => {
  const token = localStorage.getItem("talkbox");
  if (token && config.headers) {
    config.headers["Authorization"] = `Bearer ${token}`;
  }
  return config;
});

export default http;
