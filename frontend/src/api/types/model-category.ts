export interface ModelCategory {
  id: number
  tenant_key: string
  category_key: string
  display_name: string
  description: string
  icon: string
  sort_order: number
  status: 'active' | 'disabled'
  created_at: string
  updated_at: string
}
