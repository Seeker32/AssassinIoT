/** Property definition within a ThingModel */
export interface PropertyDef {
  identifier: string
  name: string
  data_type: 'int' | 'float' | 'string' | 'bool' | 'json'
  unit?: string
  mode?: 'rw' | 'r' | 'w'
  description?: string
  range?: { min?: number; max?: number }
}

/** Service (command) definition within a ThingModel */
export interface ServiceDef {
  identifier: string
  name: string
  input: Record<string, unknown>
  output: Record<string, unknown>
  description?: string
}

/** Event definition within a ThingModel */
export interface EventDef {
  identifier: string
  name: string
  data_type: string
  description?: string
}

export interface ThingModel {
  id: number
  model_key: string
  tenant_key: string
  name: string
  description: string
  category: string
  properties: Record<string, PropertyDef>
  services: Record<string, ServiceDef>
  events: Record<string, EventDef>
  version: string
  status: 'active' | 'deprecated' | 'disabled'
  created_at: string
  updated_at: string
}
