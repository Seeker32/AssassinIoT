import http from './client'
import type { LoginRequest, LoginResponse, UserInfo } from './types/auth'

export function loginApi(data: LoginRequest) {
  return http.post<LoginResponse>('/auth/login', data)
}

export function logoutApi() {
  return http.post('/auth/logout')
}

export function getUserInfo() {
  return http.get<UserInfo>('/auth/me')
}
