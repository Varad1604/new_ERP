"use client"

import axios from "axios"
import { useAuthStore } from "@/store/auth"

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1",
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 10000, // Enterprise APIs should fail fast
})

// Request Interceptor: Inject JWT and distributed tracing ID
apiClient.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().token
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    // Inject correlation ID for distributed tracing (Datadog/Jaeger)
    config.headers["X-Request-ID"] = crypto.randomUUID()

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response Interceptor: Centralized error handling
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Auto-logout on unauthorized
      useAuthStore.getState().logout()
      if (typeof window !== "undefined") {
        window.location.href = "/login"
      }
    }

    // Standardize error message shape for UI consumption
    const customError = new Error(
      error.response?.data?.error || "An unexpected system error occurred."
    )
    return Promise.reject(customError)
  }
)

export default apiClient
