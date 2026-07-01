import http from './client'
import type { ModelCategory } from './types/model-category'
import type { PaginatedResponse, ListQueryParams } from './types/common'

export function fetchModelCategories(params: ListQueryParams) {
  return http.get<PaginatedResponse<ModelCategory>>('/model-categories', { params })
}

export function fetchModelCategory(id: number) {
  return http.get<ModelCategory>(`/model-categories/${id}`)
}

export function createModelCategory(data: Partial<ModelCategory>) {
  return http.post<ModelCategory>('/model-categories', data)
}

export function updateModelCategory(id: number, data: Partial<ModelCategory>) {
  return http.put<ModelCategory>(`/model-categories/${id}`, data)
}

export function deleteModelCategory(id: number) {
  return http.delete(`/model-categories/${id}`)
}
