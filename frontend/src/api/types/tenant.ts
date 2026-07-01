export interface Tenant {
  id: number
  tenant_key: string
  name: string
  description: string
  status: 'active' | 'disabled'
  created_at: string
  updated_at: string
}
