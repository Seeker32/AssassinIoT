import http from './client'
import type { Tenant } from './types/tenant'
import type { PaginatedResponse, ListQueryParams } from './types/common'

export function fetchTenants(params: ListQueryParams) {
  return http.get<PaginatedResponse<Tenant>>('/tenants', { params })
}

export function fetchTenant(id: number) {
  return http.get<Tenant>(`/tenants/${id}`)
}

export function createTenant(data: Partial<Tenant>) {
  return http.post<Tenant>('/tenants', data)
}

export function updateTenant(id: number, data: Partial<Tenant>) {
  return http.put<Tenant>(`/tenants/${id}`, data)
}

export function deleteTenant(id: number) {
  return http.delete(`/tenants/${id}`)
}
