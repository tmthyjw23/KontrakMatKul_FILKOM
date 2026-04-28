"use client";

import axios from "axios";

const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

function readToken(): string | null {
  if (typeof window === "undefined") {
    return null;
  }

  const localStorageToken = window.localStorage.getItem("token");
  if (localStorageToken) {
    return localStorageToken;
  }

  const cookieToken = document.cookie
    .split("; ")
    .find((cookie) => cookie.startsWith("token="))
    ?.split("=")[1];

  return cookieToken ? decodeURIComponent(cookieToken) : null;
}

export const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use((config) => {
  const token = readToken();

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});
