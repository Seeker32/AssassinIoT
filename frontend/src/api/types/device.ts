export interface Device {
  id: string
  tenant_key: string
  model_key: string
  device_name: string
  access_key: string
  firmware_ver: string
  properties_cfg: Record<string, unknown>
  status: 'active' | 'inactive' | 'disabled'
  online: boolean
  metadata: Record<string, unknown>
  last_seen: string | null
  created_at: string
  updated_at: string
}
