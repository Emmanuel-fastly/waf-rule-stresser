import { Shield } from 'lucide-react'

const Header = () => {
  return (
    <header className="max-w-7xl mx-auto px-4 py-8">
      <div className="flex items-center gap-3 mb-2">
        <Shield size={32} className="text-blue-600" />
        <h1 className="text-3xl font-bold text-gray-900">WAF Rate Limit Tester</h1>
      </div>
      <p className="text-gray-600">
        Test WAF block & rate-limit rules by sending controlled HTTP traffic
      </p>
    </header>
  )
}

export default Header
