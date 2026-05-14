"use client"

import { useAuthStore } from "@/store/auth"
import { useRouter } from "next/navigation"
import { useEffect } from "react"
import { Activity, Users, DollarSign, LogOut } from "lucide-react"
import { useQuery } from "@tanstack/react-query"
import apiClient from "@/lib/api"

interface DashboardMetrics {
  total_revenue: number
  active_employees: number
  pending_approvals: number
}

const fetchMetrics = async (): Promise<DashboardMetrics> => {
  const { data } = await apiClient.get('/dashboard/metrics')
  return data
}

export default function Dashboard() {
  const { isAuthenticated, logout } = useAuthStore()
  const router = useRouter()

  useEffect(() => {
    if (!isAuthenticated) {
      router.push("/login")
    }
  }, [isAuthenticated, router])

  const { data: metrics, isLoading, isError } = useQuery({
    queryKey: ['dashboardMetrics'],
    queryFn: fetchMetrics,
    enabled: isAuthenticated,
    refetchInterval: 30000, // Enterprise dashboards refresh periodically
  })

  if (!isAuthenticated) return null

  // Formatting for enterprise locale
  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(amount)
  }

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col">
      <nav className="bg-white border-b border-gray-200 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <span className="text-xl font-bold text-gray-900 tracking-tight">Enterprise ERP <span className="text-sm font-medium text-blue-600 bg-blue-50 px-2 py-1 rounded-md ml-2">SaaS Edition</span></span>
            </div>
            <div className="flex items-center">
              <button
                onClick={() => {
                  logout()
                  router.push("/login")
                }}
                className="flex items-center text-sm font-medium text-gray-500 hover:text-gray-900 transition-colors"
              >
                <LogOut className="mr-2 h-4 w-4" />
                Sign out
              </button>
            </div>
          </div>
        </div>
      </nav>

      <main className="flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 w-full">
        <div className="mb-8">
          <h1 className="text-2xl font-semibold text-gray-900">Executive Dashboard</h1>
          <p className="mt-1 text-sm text-gray-500">Real-time financial and operational intelligence.</p>
        </div>

        {isError && (
          <div className="bg-red-50 border-l-4 border-red-400 p-4 mb-8">
            <p className="text-sm text-red-700">Failed to load real-time metrics. Please ensure the backend services are running and you are authenticated.</p>
          </div>
        )}

        <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
          {/* Revenue Card */}
          <div className="bg-white overflow-hidden shadow rounded-lg border border-gray-100">
            <div className="p-5">
              <div className="flex items-center">
                <div className="flex-shrink-0 bg-green-50 rounded-md p-3">
                  <DollarSign className="h-6 w-6 text-green-600" />
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">Total Revenue (Ledger)</dt>
                    <dd className="flex items-baseline mt-1">
                      {isLoading ? (
                        <div className="h-8 w-32 bg-gray-200 animate-pulse rounded"></div>
                      ) : (
                        <div className="text-2xl font-semibold text-gray-900">{formatCurrency(metrics?.total_revenue || 0)}</div>
                      )}
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          {/* Employees Card */}
          <div className="bg-white overflow-hidden shadow rounded-lg border border-gray-100">
            <div className="p-5">
              <div className="flex items-center">
                <div className="flex-shrink-0 bg-blue-50 rounded-md p-3">
                  <Users className="h-6 w-6 text-blue-600" />
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">Active Employees</dt>
                    <dd className="flex items-baseline mt-1">
                      {isLoading ? (
                        <div className="h-8 w-16 bg-gray-200 animate-pulse rounded"></div>
                      ) : (
                        <div className="text-2xl font-semibold text-gray-900">{metrics?.active_employees || 0}</div>
                      )}
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          {/* Approvals Card */}
          <div className="bg-white overflow-hidden shadow rounded-lg border border-gray-100">
            <div className="p-5">
              <div className="flex items-center">
                <div className="flex-shrink-0 bg-orange-50 rounded-md p-3">
                  <Activity className="h-6 w-6 text-orange-600" />
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">Pending Approvals</dt>
                    <dd className="flex items-baseline mt-1">
                      {isLoading ? (
                        <div className="h-8 w-16 bg-gray-200 animate-pulse rounded"></div>
                      ) : (
                        <div className="text-2xl font-semibold text-gray-900">{metrics?.pending_approvals || 0}</div>
                      )}
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Real-time Ledger Notice */}
        <div className="mt-8 bg-blue-50 border border-blue-100 rounded-md p-4">
          <div className="flex">
            <div className="flex-shrink-0">
              <Activity className="h-5 w-5 text-blue-400" aria-hidden="true" />
            </div>
            <div className="ml-3 flex-1 md:flex md:justify-between">
              <p className="text-sm text-blue-700">
                Dashboard metrics are now powered by real-time aggregations directly from the Go Core Ledger Microservice.
              </p>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
