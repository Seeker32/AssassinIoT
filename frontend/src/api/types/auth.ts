export interface LoginRequest {
  username: string
  password: string
  tenant_key: string
}

export interface LoginResponse {
  token: string
  user: UserInfo
}

export interface UserInfo {
  id: number
  username: string
  display_name: string
  tenant_key: string
  role: string
}
