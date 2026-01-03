import { Activity, CheckCircle, XCircle, ShieldX } from 'lucide-react'

const StatsCards = ({ results }) => {
  if (!results) return null

  const stats = [
    {
      label: 'Total Requests',
      value: results.total_requests,
      icon: Activity,
      iconColor: 'text-blue-600'
    },
    {
      label: 'Success',
      value: results.success_count,
      icon: CheckCircle,
      iconColor: 'text-green-600'
    },
    {
      label: 'Errors',
      value: results.error_count,
      icon: XCircle,
      iconColor: 'text-red-600'
    },
    {
      label: 'Blocked',
      value: results.blocked_count,
      icon: ShieldX,
      iconColor: 'text-purple-600'
    }
  ]

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      {stats.map((stat, index) => {
        const Icon = stat.icon
        return (
          <div key={index} className="bg-white rounded-lg p-6 shadow-md border border-gray-200">
            <div className="flex items-center gap-2 mb-3">
              <Icon size={20} className={stat.iconColor} />
              <span className="text-sm font-medium text-gray-700">{stat.label}</span>
            </div>
            <div className="text-3xl font-bold text-gray-900">{stat.value}</div>
          </div>
        )
      })}
    </div>
  )
}

export default StatsCards
