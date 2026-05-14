"use client"

import { useAuthStore } from "@/store/auth"
import { useRouter } from "next/navigation"
import { useEffect } from "react"
import { Activity, Users, DollarSign, LogOut } from "lucide-react"

export default function Dashboard() {
  const { isAuthenticated, logout } = useAuthStore()
  const router = useRouter()

  useEffect(() => {
    if (!isAuthenticated) {
      router.push("/login")
    }
  }, [isAuthenticated, router])

  if (!isAuthenticated) return null

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col">
      <nav className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <span className="text-xl font-bold text-gray-900">ERP Engine</span>
            </div>
            <div className="flex items-center">
              <button
                onClick={() => {
                  logout()
                  router.push("/login")
                }}
                className="flex items-center text-sm font-medium text-gray-500 hover:text-gray-700"
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
          <p className="mt-1 text-sm text-gray-500">Overview of key performance indicators.</p>
        </div>

        <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
          {/* Card 1 */}
          <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="p-5">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <DollarSign className="h-6 w-6 text-gray-400" />
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">Total Revenue</dt>
                    <dd className="flex items-baseline">
                      <div className="text-2xl font-semibold text-gray-900">$2,405,100.00</div>
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          {/* Card 2 */}
          <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="p-5">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <Users className="h-6 w-6 text-gray-400" />
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">Active Employees</dt>
                    <dd className="flex items-baseline">
                      <div className="text-2xl font-semibold text-gray-900">142</div>
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          {/* Card 3 */}
          <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="p-5">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <Activity className="h-6 w-6 text-gray-400" />
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">Pending Approvals</dt>
                    <dd className="flex items-baseline">
                      <div className="text-2xl font-semibold text-gray-900">12</div>
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Recent Ledger Activity */}
        <div className="mt-8">
          <h2 className="text-lg leading-6 font-medium text-gray-900 mb-4">Recent Ledger Activity</h2>
          <div className="bg-white shadow overflow-hidden sm:rounded-md border border-gray-200">
            <ul className="divide-y divide-gray-200">
              <li className="px-6 py-4 flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-blue-600">JE-2023-0102</p>
                  <p className="text-sm text-gray-500">Q3 Cloud Infrastructure Payment</p>
                </div>
                <div className="text-right">
                  <p className="text-sm font-semibold text-gray-900">POSTED</p>
                  <p className="text-sm text-gray-500">2023-10-24</p>
                </div>
              </li>
              <li className="px-6 py-4 flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-blue-600">JE-2023-0103</p>
                  <p className="text-sm text-gray-500">Client Invoice #8843</p>
                </div>
                <div className="text-right">
                  <p className="text-sm font-semibold text-gray-900">DRAFT</p>
                  <p className="text-sm text-gray-500">2023-10-25</p>
                </div>
              </li>
            </ul>
          </div>
        </div>
      </main>
    </div>
  )
}
