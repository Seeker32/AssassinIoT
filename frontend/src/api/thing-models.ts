import http from './client'
import type { ThingModel } from './types/thing-model'
import type { PaginatedResponse, ListQueryParams } from './types/common'

export function fetchThingModels(params: ListQueryParams) {
  return http.get<PaginatedResponse<ThingModel>>('/thing-models', { params })
}

export function fetchThingModel(id: number) {
  return http.get<ThingModel>(`/thing-models/${id}`)
}

export function createThingModel(data: Partial<ThingModel>) {
  return http.post<ThingModel>('/thing-models', data)
}

export function updateThingModel(id: number, data: Partial<ThingModel>) {
  return http.put<ThingModel>(`/thing-models/${id}`, data)
}

export function deleteThingModel(id: number) {
  return http.delete(`/thing-models/${id}`)
}
