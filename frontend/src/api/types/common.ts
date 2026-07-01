// Common API types used across all modules

/** Standard paginated list response */
export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  page_size: number
}

/** Standard API response wrapper */
export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

/** Generic list query params */
export interface ListQueryParams {
  page?: number
  page_size?: number
  keyword?: string
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}
