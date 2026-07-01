import http from './client'
import type { Device } from './types/device'
import type { PaginatedResponse, ListQueryParams } from './types/common'

export function fetchDevices(params: ListQueryParams) {
  return http.get<PaginatedResponse<Device>>('/devices', { params })
}

export function fetchDevice(id: string) {
  return http.get<Device>(`/devices/${id}`)
}

export function createDevice(data: Partial<Device>) {
  return http.post<Device>('/devices', data)
}

export function updateDevice(id: string, data: Partial<Device>) {
  return http.put<Device>(`/devices/${id}`, data)
}

export function deleteDevice(id: string) {
  return http.delete(`/devices/${id}`)
}
